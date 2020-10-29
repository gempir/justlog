package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/justlog/bot"

	"github.com/gempir/justlog/config"

	"github.com/gempir/justlog/helix"
	log "github.com/sirupsen/logrus"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/justlog/filelog"
)

// Server api server
type Server struct {
	listenAddress string
	logPath       string
	bot           *bot.Bot
	cfg           *config.Config
	fileLogger    *filelog.Logger
	helixClient   *helix.Client
	channels      []string
	assets        []string
	assetHandler  http.Handler
}

// NewServer create api Server
func NewServer(cfg *config.Config, bot *bot.Bot, fileLogger *filelog.Logger, helixClient *helix.Client, channels []string) Server {
	return Server{
		listenAddress: cfg.ListenAddress,
		bot:           bot,
		logPath:       cfg.LogsDirectory,
		cfg:           cfg,
		fileLogger:    fileLogger,
		helixClient:   helixClient,
		channels:      channels,
		assets:        []string{"/", "/favicon.ico", "/robots.txt"},
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

var (
	userHourLimit    = 744.0
	channelHourLimit = 24.0
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

type logList struct {
	AvailableLogs []filelog.UserLogFile `json:"availableLogs"`
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

// Init start the server
func (s *Server) Init() {
	http.Handle("/", corsHandler(http.HandlerFunc(s.route)))

	log.Infof("Listening on %s", s.listenAddress)
	log.Fatal(http.ListenAndServe(s.listenAddress, nil))
}

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	url := r.URL.EscapedPath()

	if contains(s.assets, url) || strings.HasPrefix(url, "/bundle") {
		s.assetHandler.ServeHTTP(w, r)
		return
	}

	query := s.fillUserids(w, r)

	if url == "/list" {
		s.writeAvailableLogs(w, r, query)
		return
	}

	if url == "/channels" {
		s.writeAllChannels(w, r)
		return
	}

	if strings.HasPrefix(url, "/admin/channelConfigs/") {
		success := s.authenticateAdmin(w, r)
		if success {
			s.writeChannelConfigs(w, r)
		}
		return
	}

	if strings.HasPrefix(url, "/admin/channels") {
		success := s.authenticateAdmin(w, r)
		if success {
			s.writeChannels(w, r)
		}
		return
	}

	s.routeLogs(w, r)
}

func (s *Server) fillUserids(w http.ResponseWriter, r *http.Request) url.Values {
	query := r.URL.Query()

	if query.Get("userid") == "" && query.Get("user") != "" {
		users, err := s.helixClient.GetUsersByUsernames([]string{query.Get("user")})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		query.Set("userid", users[query.Get("user")].ID)
	}

	if query.Get("channelid") == "" && query.Get("channel") != "" {
		users, err := s.helixClient.GetUsersByUsernames([]string{query.Get("channel")})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		query.Set("channelid", users[query.Get("channel")].ID)
	}

	return query
}

func (s *Server) routeLogs(w http.ResponseWriter, r *http.Request) {

	request, err := s.newLogRequestFromURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if request.redirectPath != "" {
		http.Redirect(w, r, request.redirectPath, http.StatusFound)
		return
	}

	var logs *chatLog
	if request.time.random {
		logs, err = s.getRandomQuote(request)
	} else if request.time.from != "" && request.time.to != "" {
		if request.isUserRequest {
			logs, err = s.getUserLogsRange(request)
		} else {
			logs, err = s.getChannelLogsRange(request)
		}

	} else {
		if request.isUserRequest {
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

	// Disable content type sniffing for log output
	w.Header().Set("X-Content-Type-Options", "nosniff")

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

func (s *Server) writeAllChannels(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(js)
}

func writeRaw(cLog *chatLog, code int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	for _, cMessage := range cLog.Messages {
		w.Write([]byte(cMessage.Raw + "\n"))
	}
}

func writeText(cLog *chatLog, code int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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

func buildClearChatMessageText(message twitch.ClearChatMessage) string {
	if message.BanDuration == 0 {
		return fmt.Sprintf("%s has been banned", message.TargetUsername)
	}

	return fmt.Sprintf("%s has been timed out for %d seconds", message.TargetUsername, message.BanDuration)
}
