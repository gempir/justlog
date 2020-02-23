package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/justlog/helix"
	log "github.com/sirupsen/logrus"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/justlog/filelog"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
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

type userRequestContext struct {
	echo.Context
	channelType string
	userType    string
}

// @title justlog API
// @version 1.0
// @description API for twitch logs

// @contact.name gempir
// @contact.url https://gempir.com
// @contact.email gempir.dev@gmail.com

// @license.name MIT
// @license.url https://github.com/gempir/justlog/blob/master/LICENSE

// @host logs.ivr.fi
// @BasePath /

// Init start the server
func (s *Server) Init() {
	http.Handle("/", corsHandler(http.HandlerFunc(s.route)))

	// e.GET("/:channelType/:channel/:userType/:user/:time", s.getUserLogsExact)

	// e.GET("/channel/:channel/user/:username/range", s.getUserLogsRangeByName)
	// e.GET("/channelid/:channelid/userid/:userid/range", s.getUserLogsRange)

	// e.GET("/channel/:channel/user/:username", s.getLastUserLogsByName)
	// e.GET("/channel/:channel/user/:username/random", s.getRandomQuoteByName)

	// e.GET("/channelid/:channelid/userid/:userid", s.getLastUserLogs)
	// e.GET("/channelid/:channelid/userid/:userid/random", s.getRandomQuote)

	// e.GET("/channelid/:channelid/range", s.getChannelLogsRange)
	// e.GET("/channel/:channel/range", s.getChannelLogsRangeByName)

	// e.GET("/channel/:channel", s.getCurrentChannelLogsByName)
	// e.GET("/channel/:channel/:year/:month/:day", s.getChannelLogsByName)
	// e.GET("/channelid/:channelid", s.getCurrentChannelLogs)
	// e.GET("/channelid/:channelid/:year/:month/:day", s.getChannelLogs)

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

	fmt.Fprint(w, r.URL.EscapedPath())
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

func reverse(input []string) []string {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
	return input
}

// getAllChannels godoc
// @Summary Get all joined channels
// @tags bot
// @Produce  json
// @Success 200 {object} api.RandomQuoteJSON json
// @Failure 500 {object} api.ErrorResponse json
// @Router /channels [get]
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

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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

func writeTextResponse(c echo.Context, cLog *chatLog) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	for _, cMessage := range cLog.Messages {
		switch cMessage.Type {
		case twitch.PRIVMSG:
			c.Response().Write([]byte(fmt.Sprintf("[%s] #%s %s: %s\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Username, cMessage.Text)))
		case twitch.CLEARCHAT:
			c.Response().Write([]byte(fmt.Sprintf("[%s] #%s %s\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Text)))
		case twitch.USERNOTICE:
			c.Response().Write([]byte(fmt.Sprintf("[%s] #%s %s\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Text)))
		}
	}

	return nil
}

func writeRawResponse(c echo.Context, cLog *chatLog) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	for _, cMessage := range cLog.Messages {
		c.Response().Write([]byte(cMessage.Raw + "\n"))
	}

	return nil
}

func writeJSONResponse(c echo.Context, logResult *chatLog) error {
	_, stream := c.QueryParams()["stream"]
	if stream {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		return json.NewEncoder(c.Response()).Encode(logResult)
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(logResult)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, echo.MIMEApplicationJSONCharsetUTF8, data)
}

func parseTimestamp(timestamp string) (time.Time, error) {

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil
}

func shouldReverse(c echo.Context) bool {
	_, ok := c.QueryParams()["reverse"]

	return c.QueryParam("order") == "reverse" || ok
}

func shouldRespondWithJSON(c echo.Context) bool {
	_, ok := c.QueryParams()["json"]

	return c.Request().Header.Get("Content-Type") == "application/json" || c.Request().Header.Get("accept") == "application/json" || c.QueryParam("type") == "json" || ok
}

func shouldRespondWithRaw(c echo.Context) bool {
	_, ok := c.QueryParams()["raw"]

	return c.QueryParam("type") == "raw" || ok
}
