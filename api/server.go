package api

import (
	"net/http"

	"github.com/gempir/justlog/filelog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	listenAddress string
	logPath       string
	fileLogger    *filelog.Logger
	channels      []string
}

func NewServer(logPath string, listenAddress string, fileLogger *filelog.Logger) Server {
	return Server{
		listenAddress: listenAddress,
		logPath:       logPath,
		fileLogger:    fileLogger,
		channels:      []string{},
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
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/channelid/:channelid/user/:userid", s.getCurrentUserLogs)
	e.GET("/channelid", s.getAllChannels)
	e.GET("/channelid/:channelid", s.getCurrentChannelLogs)
	e.GET("/channelid/:channelid/:year/:month/:day", s.getChannelLogs)
	e.GET("/channelid/:channelid/userid/:userid/:year/:month", s.getUserLogs)
	e.GET("/channelid/:channelid/userid/:userid/random", s.getRandomQuote)

	e.Logger.Fatal(e.Start(s.listenAddress))
}
