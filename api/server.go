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
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/gempir/justlog/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	listenAddress string
	logPath       string
	fileLogger    *filelog.Logger
	helixClient   *helix.Client
	channels      []string
}

func NewServer(logPath string, listenAddress string, fileLogger *filelog.Logger, helixClient *helix.Client, channels []string) Server {
	return Server{
		listenAddress: listenAddress,
		logPath:       logPath,
		fileLogger:    fileLogger,
		helixClient:   helixClient,
		channels:      channels,
	}
}

func (s *Server) AddChannel(channel string) {
	s.channels = append(s.channels, channel)
}

func (s *Server) Init() {

	e := echo.New()
	e.HideBanner = true

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "", // disabled
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "", // disabled
		HSTSMaxAge:            0,  // disabled
		ContentSecurityPolicy: "", // disabled
	}))
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

	e.Static("/", "web/public/index.html")
	e.Static("/bundle.js", "web/public/bundle.js")

	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/index.html")
	})
	e.GET("/*", echoSwagger.WrapHandler)
	e.GET("/channels", s.getAllChannels)

	e.GET("/channel/:channel/user/:username/range", s.getUserLogsRangeByName)
	e.GET("/channelid/:channelid/userid/:userid/range", s.getUserLogsRange)

	e.GET("/channel/:channel/user/:username", s.getLastUserLogsByName)
	e.GET("/channel/:channel/user/:username/:year/:month", s.getUserLogsByName)
	e.GET("/channel/:channel/user/:username/random", s.getRandomQuoteByName)

	e.GET("/channelid/:channelid/userid/:userid", s.getLastUserLogs)
	e.GET("/channelid/:channelid/userid/:userid/:year/:month", s.getUserLogs)
	e.GET("/channelid/:channelid/userid/:userid/random", s.getRandomQuote)

	e.GET("/channelid/:channelid/range", s.getChannelLogsRange)
	e.GET("/channel/:channel/range", s.getChannelLogsRangeByName)

	e.GET("/channel/:channel", s.getCurrentChannelLogsByName)
	e.GET("/channel/:channel/:year/:month/:day", s.getChannelLogsByName)
	e.GET("/channelid/:channelid", s.getCurrentChannelLogs)
	e.GET("/channelid/:channelid/:year/:month/:day", s.getChannelLogs)

	e.Logger.Fatal(e.Start(s.listenAddress))
}

var (
	userHourLimit    = 744.0
	channelHourLimit = 24.0
)

type channel struct {
	UserID string `json:"userID"`
	Name   string `json:"name"`
}

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

type ErrorResponse struct {
	Message string `json:"message"`
}

type timestamp struct {
	time.Time
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
func (s *Server) getAllChannels(c echo.Context) error {
	response := new(AllChannelsJSON)
	response.Channels = []channel{}
	users, err := s.helixClient.GetUsersByUserIds(s.channels)

	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
	}

	for _, user := range users {
		response.Channels = append(response.Channels, channel{UserID: user.ID, Name: user.Login})
	}

	return c.JSON(http.StatusOK, response)
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

func shouldRespondWithJson(c echo.Context) bool {
	_, ok := c.QueryParams()["json"]

	return c.Request().Header.Get("Content-Type") == "application/json" || c.Request().Header.Get("accept") == "application/json" || c.QueryParam("type") == "json" || ok
}

func shouldRespondWithRaw(c echo.Context) bool {
	_, ok := c.QueryParams()["raw"]

	return c.QueryParam("type") == "raw" || ok
}
