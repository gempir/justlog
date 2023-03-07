package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/justlog/bot"

	"github.com/gempir/justlog/config"

	"github.com/gempir/justlog/helix"
	log "github.com/sirupsen/logrus"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/gempir/justlog/filelog"
)

// Server api server
type Server struct {
	listenAddress string
	logPath       string
	bot           *bot.Bot
	cfg           *config.Config
	fileLogger    filelog.Logger
	helixClient   helix.TwitchApiClient
	assetsHandler http.Handler
}

// NewServer create api Server
func NewServer(cfg *config.Config, bot *bot.Bot, fileLogger filelog.Logger, helixClient helix.TwitchApiClient, assets fs.FS) Server {
	build, err := fs.Sub(assets, "web/dist")
	if err != nil {
		log.Fatal("failed to read public assets")
	}

	return Server{
		listenAddress: cfg.ListenAddress,
		bot:           bot,
		logPath:       cfg.LogsDirectory,
		cfg:           cfg,
		fileLogger:    fileLogger,
		helixClient:   helixClient,
		assetsHandler: http.FileServer(http.FS(build)),
	}
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

// swagger:model
type AllChannelsJSON struct {
	Channels []channel `json:"channels"`
}

// swagger:model
type chatLog struct {
	Messages []chatMessage `json:"messages"`
}

// swagger:model
type logList struct {
	AvailableLogs []filelog.UserLogFile `json:"availableLogs"`
}

// swagger:model
type channelLogList struct {
	AvailableLogs []filelog.ChannelLogFile `json:"availableLogs"`
}

type chatMessage struct {
	Text        string             `json:"text"`
	Username    string             `json:"username"`
	DisplayName string             `json:"displayName"`
	Channel     string             `json:"channel"`
	Timestamp   timestamp          `json:"timestamp"`
	ID          string             `json:"id"`
	Type        twitch.MessageType `json:"type"`
	Raw         string             `json:"raw"`
	Tags        map[string]string  `json:"tags"`
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

	query, err := s.fillUserids(w, r)
	if err != nil {
		return
	}

	if url == "/list" {
		if s.cfg.IsOptedOut(query.Get("userid")) || s.cfg.IsOptedOut(query.Get("channelid")) {
			http.Error(w, "User or channel has opted out", http.StatusForbidden)
			return
		}

		s.writeAvailableLogs(w, r, query)
		return
	}

	if url == "/channels" {
		s.writeAllChannels(w, r)
		return
	}

	if url == "/optout" && r.Method == http.MethodPost {
		s.writeOptOutCode(w, r)
		return
	}

	if strings.HasPrefix(url, "/admin/channels") {
		success := s.authenticateAdmin(w, r)
		if success {
			s.writeChannels(w, r)
		}
		return
	}

	routedLogs := s.routeLogs(w, r)

	if !routedLogs {
		s.assetsHandler.ServeHTTP(w, r)
		return
	}
}

func (s *Server) fillUserids(w http.ResponseWriter, r *http.Request) (url.Values, error) {
	query := r.URL.Query()

	if query.Get("userid") == "" && query.Get("user") != "" {
		username := strings.ToLower(query.Get("user"))

		users, err := s.helixClient.GetUsersByUsernames([]string{username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		if len(users) == 0 {
			err := fmt.Errorf("could not find users")
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return nil, err
		}

		query.Set("userid", users[username].ID)
	}

	if query.Get("channelid") == "" && query.Get("channel") != "" {
		channelName := strings.ToLower(query.Get("channel"))

		users, err := s.helixClient.GetUsersByUsernames([]string{channelName})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		if len(users) == 0 {
			err := fmt.Errorf("could not find users")
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return nil, err
		}

		query.Set("channelid", users[channelName].ID)
	}

	return query, nil
}

func (s *Server) routeLogs(w http.ResponseWriter, r *http.Request) bool {

	request, err := s.newLogRequestFromURL(r)
	if err != nil {
		return false
	}
	if request.redirectPath != "" {
		http.Redirect(w, r, request.redirectPath, http.StatusFound)
		return true
	}

	if s.cfg.IsOptedOut(request.channelid) || s.cfg.IsOptedOut(request.userid) {
		http.Error(w, "User or channel has opted out", http.StatusForbidden)
		return true
	}

	var logs *chatLog
	if request.time.random {
		if request.isUserRequest {
			logs, err = s.getRandomQuote(request)
		} else {
			logs, err = s.getChannelRandomQuote(request)
		}
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
		http.Error(w, "could not load logs", http.StatusNotFound)
		return true
	}

	// Disable content type sniffing for log output
	w.Header().Set("X-Content-Type-Options", "nosniff")

	currentYear := fmt.Sprintf("%d", int(time.Now().Year()))
	currentMonth := fmt.Sprintf("%d", int(time.Now().Month()))

	if (request.time.year != "" && request.time.month != "") && (request.time.year < currentYear || (request.time.year == currentYear && request.time.month < currentMonth)) {
		writeCacheControl(w, r, time.Hour*8760)
	} else {
		writeCacheControlNoCache(w, r)
	}

	if request.responseType == responseTypeJSON {
		writeJSON(logs, http.StatusOK, w, r)
		return true
	}

	if request.responseType == responseTypeRaw {
		writeRaw(logs, http.StatusOK, w, r)
		return true
	}

	if request.responseType == responseTypeText {
		writeText(logs, http.StatusOK, w, r)
		return true
	}

	return false
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

// swagger:route GET /channels justlog channels
//
// List currently logged channels
//
//	Produces:
//	- application/json
//	- text/plain
//
//	Schemes: https
//
//	Responses:
//	  200: AllChannelsJSON
func (s *Server) writeAllChannels(w http.ResponseWriter, r *http.Request) {
	response := new(AllChannelsJSON)
	response.Channels = []channel{}
	users, err := s.helixClient.GetUsersByUserIds(s.cfg.Channels)

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

func writeCacheControl(w http.ResponseWriter, r *http.Request, cacheDuration time.Duration) {
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%.0f", cacheDuration.Seconds()))
}

func writeCacheControlNoCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
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

func createLogResult() chatLog {
	return chatLog{Messages: []chatMessage{}}
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

func createChatMessage(parsedMessage twitch.Message) chatMessage {
	switch message := parsedMessage.(type) {
	case *twitch.PrivateMessage:
		return chatMessage{
			Timestamp:   timestamp{message.Time},
			Username:    message.User.Name,
			DisplayName: message.User.DisplayName,
			Text:        message.Message,
			Type:        message.Type,
			Channel:     message.Channel,
			Raw:         message.Raw,
			ID:          message.ID,
			Tags:        message.Tags,
		}
	case *twitch.ClearChatMessage:
		return chatMessage{
			Timestamp:   timestamp{message.Time},
			Username:    message.TargetUsername,
			DisplayName: message.TargetUsername,
			Text:        buildClearChatMessageText(*message),
			Type:        message.Type,
			Channel:     message.Channel,
			Raw:         message.Raw,
			Tags:        message.Tags,
		}
	case *twitch.UserNoticeMessage:
		return chatMessage{
			Timestamp:   timestamp{message.Time},
			Username:    message.User.Name,
			DisplayName: message.User.DisplayName,
			Text:        message.SystemMsg + " " + message.Message,
			Type:        message.Type,
			Channel:     message.Channel,
			Raw:         message.Raw,
			ID:          message.ID,
			Tags:        message.Tags,
		}
	}

	return chatMessage{}
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
