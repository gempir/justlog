package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/justlog/helix"
	log "github.com/sirupsen/logrus"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/justlog/filelog"
)

// Server api server
type Server struct {
	listenAddress string
	logPath       string
	fileLogger    *filelog.Logger
	helixClient   *helix.Client
	channels      []string
	assets        []string
	assetHandler  http.Handler
}

// NewServer create api Server
func NewServer(logPath string, listenAddress string, fileLogger *filelog.Logger, helixClient *helix.Client, channels []string) Server {
	return Server{
		listenAddress: listenAddress,
		logPath:       logPath,
		fileLogger:    fileLogger,
		helixClient:   helixClient,
		channels:      channels,
		assets:        []string{"/", "/bundle.js", "/swagger.json", "/swagger.html"},
		assetHandler:  http.FileServer(assets),
	}
}

// AddChannel adds a channel to the collection to output on the channels endpoint
func (s *Server) AddChannel(channel string) {
	s.channels = append(s.channels, channel)
}

const (
	responseTypeJSON = "json"
	responseTypeText = "text"
	responseTypeRaw  = "raw"
)

type userRequest struct {
	channel      string
	user         string
	channelid    string
	userid       string
	time         timeRequest
	reverse      bool
	responseType string
}

type timeRequest struct {
	from  string
	to    string
	year  string
	month string
	day   string
}

// Init start the server
func (s *Server) Init() {
	http.Handle("/", corsHandler(http.HandlerFunc(s.route)))

	log.Fatal(http.ListenAndServe(s.listenAddress, nil))
}

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	url := r.URL.EscapedPath()

	if contains(s.assets, url) {
		s.assetHandler.ServeHTTP(w, r)
		return
	}

	if url == "/channels" {
		s.getAllChannels(w, r)
		return
	}

	s.routeLogs(w, r)
}

func (s *Server) routeLogs(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimRight(r.URL.EscapedPath(), "/")

	matches := pathRegex.FindAllStringSubmatch(url, -1)
	if len(matches) == 0 || len(matches[0]) < 5 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	request := userRequest{
		time: timeRequest{},
	}

	params := []string{}
	for _, match := range matches[0] {
		if match != "" {
			params = append(params, match)
		}
	}

	isUserRequest := len(params) > 3 && (params[3] == "user" || params[3] == "userid")
	isChannelRequest := len(params) < 3 || (len(params) > 3 && params[3] != "user" && params[3] != "userid")

	if params[1] == "channel" {
		request.channel = params[2]
	}
	if params[1] == "channelid" {
		request.channelid = params[2]
	}
	if isUserRequest && params[3] == "user" {
		request.user = params[4]
	}
	if isUserRequest && params[3] == "userid" {
		request.userid = params[4]
	}

	if isUserRequest && len(params) == 7 {
		request.time.year = params[5]
		request.time.month = params[6]
	} else if isChannelRequest && len(params) == 6 {
		request.time.year = params[3]
		request.time.month = params[4]
		request.time.day = params[5]
	} else if (isUserRequest && len(params) == 6 && params[5] == "range") || (isChannelRequest && len(params) == 4 && params[3] == "range") {
		request.time.from = r.URL.Query().Get("from")
		request.time.to = r.URL.Query().Get("to")
	} else {
		request.time.year = string(time.Now().Year())
		request.time.month = string(time.Now().Month())
		query := r.URL.Query()
		query.Del("from")
		query.Del("to")

		encodedQuery := ""
		if query.Encode() != "" {
			encodedQuery = "?" + query.Encode()
		}
		http.Redirect(w, r, fmt.Sprintf("%s/%d/%d%s", params[0], time.Now().Year(), time.Now().Month(), encodedQuery), http.StatusPermanentRedirect)
		return
	}

	if _, ok := r.URL.Query()["reverse"]; ok {
		request.reverse = true
	} else {
		request.reverse = false
	}

	if _, ok := r.URL.Query()["json"]; ok || r.URL.Query().Get("type") == "json" {
		request.responseType = responseTypeJSON
	} else if _, ok := r.URL.Query()["raw"]; ok || r.URL.Query().Get("type") == "raw" {
		request.responseType = responseTypeRaw
	} else {
		request.responseType = responseTypeText
	}

	var err error
	request, err = s.fillIds(request)
	if err != nil {
		log.Error(err)
		http.Error(w, "could not fetch userids", http.StatusInternalServerError)
		return
	}

	var logs *chatLog
	if request.time.from != "" && request.time.to != "" {
		if isUserRequest {
			logs, err = s.getUserLogsRange(request)
		} else {
			logs, err = s.getChannelLogsRange(request)
		}

	} else {
		if isUserRequest {
			logs, err = s.getUserLogs(request)
		} else {
			logs, err = s.getChannelLogs(request)
		}
	}

	if err != nil {
		log.Error(err)
		http.Error(w, "could not load logs", http.StatusInternalServerError)
		return
	}

	if request.responseType == responseTypeJSON {
		writeJSON(logs, http.StatusOK, w, r)
		return
	}

	if request.responseType == responseTypeRaw {
		writeRaw(logs, http.StatusOK, w, r)
		return
	}

	if request.responseType == responseTypeText {
		writeText(logs, http.StatusOK, w, r)
		return
	}

	http.Error(w, "unkown response type", http.StatusBadRequest)
}

func clearHeaders(w http.ResponseWriter) {
	for key := range w.Header() {
		w.Header().Del(key)
	}
}

func (s *Server) fillIds(request userRequest) (userRequest, error) {
	usernames := []string{}
	if request.channelid == "" && request.channel != "" {
		usernames = append(usernames, request.channel)
	}
	if request.userid == "" && request.user != "" {
		usernames = append(usernames, request.user)
	}

	ids, err := s.helixClient.GetUsersByUsernames(usernames)
	if err != nil {
		return request, err
	}

	if request.channelid == "" {
		request.channelid = ids[request.channel].ID
	}
	if request.userid == "" {
		request.userid = ids[request.user].ID
	}

	return request, nil
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		}
	})
}

var (
	userHourLimit    = 744.0
	channelHourLimit = 24.0
	pathRegex        = regexp.MustCompile(`\/(channel|channelid)\/([a-zA-Z0-9]+)(?:\/(user|userid)\/([a-zA-Z0-9]+))?(?:(?:\/(\d{4})\/(\d{1,2})(?:\/(\d{1,2}))?)|(?:\/(range)))?`)
)

type channel struct {
	UserID string `json:"userID"`
	Name   string `json:"name"`
}

// AllChannelsJSON inlcudes all channels
type AllChannelsJSON struct {
	Channels []channel `json:"channels"`
}

type chatLog struct {
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Text        string             `json:"text"`
	Username    string             `json:"username"`
	DisplayName string             `json:"displayName"`
	Channel     string             `json:"channel"`
	Timestamp   timestamp          `json:"timestamp"`
	Type        twitch.MessageType `json:"type"`
	Raw         string             `json:"raw"`
}

// ErrorResponse a simple error response
type ErrorResponse struct {
	Message string `json:"message"`
}

type timestamp struct {
	time.Time
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func reverseSlice(input []string) []string {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
	return input
}

func (s *Server) getAllChannels(w http.ResponseWriter, r *http.Request) {
	response := new(AllChannelsJSON)
	response.Channels = []channel{}
	users, err := s.helixClient.GetUsersByUserIds(s.channels)

	if err != nil {
		log.Error(err)
		http.Error(w, "Failure fetching data from twitch", http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		response.Channels = append(response.Channels, channel{UserID: user.ID, Name: user.Login})
	}

	writeJSON(response, http.StatusOK, w, r)
}

func writeJSON(data interface{}, code int, w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

func writeRaw(cLog *chatLog, code int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)

	for _, cMessage := range cLog.Messages {
		w.Write([]byte(cMessage.Raw + "\n"))
	}
}

func writeText(cLog *chatLog, code int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)

	for _, cMessage := range cLog.Messages {
		switch cMessage.Type {
		case twitch.PRIVMSG:
			w.Write([]byte(fmt.Sprintf("[%s] #%s %s: %s\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Username, cMessage.Text)))
		case twitch.CLEARCHAT:
			w.Write([]byte(fmt.Sprintf("[%s] #%s %s\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Text)))
		case twitch.USERNOTICE:
			w.Write([]byte(fmt.Sprintf("[%s] #%s %s\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Text)))
		}
	}
}

func (t timestamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.UTC().Format(time.RFC3339) + "\""), nil
}

func (t *timestamp) UnmarshalJSON(data []byte) error {
	goTime, err := time.Parse(time.RFC3339, strings.TrimSuffix(strings.TrimPrefix(string(data[:]), "\""), "\""))
	if err != nil {
		return err
	}
	*t = timestamp{
		goTime,
	}
	return nil
}

func parseFromTo(from, to string, limit float64) (time.Time, time.Time, error) {
	var fromTime time.Time
	var toTime time.Time

	if from == "" && to == "" {
		fromTime = time.Now().AddDate(0, -1, 0)
		toTime = time.Now()
	} else if from == "" && to != "" {
		var err error
		toTime, err = parseTimestamp(to)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("Can't parse to timestamp: %s", err)
		}
		fromTime = toTime.AddDate(0, -1, 0)
	} else if from != "" && to == "" {
		var err error
		fromTime, err = parseTimestamp(from)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("Can't parse from timestamp: %s", err)
		}
		toTime = fromTime.AddDate(0, 1, 0)
	} else {
		var err error

		fromTime, err = parseTimestamp(from)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("Can't parse from timestamp: %s", err)
		}
		toTime, err = parseTimestamp(to)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("Can't parse to timestamp: %s", err)
		}

		if toTime.Sub(fromTime).Hours() > limit {
			return fromTime, toTime, errors.New("Timespan too big")
		}
	}

	return fromTime, toTime, nil
}

func parseTimestamp(timestamp string) (time.Time, error) {

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil
}
