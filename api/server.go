package api

import (
	"github.com/op/go-logging"
	"github.com/labstack/echo"
	"net/http"
)

type Server struct {
	port string
	logPath string
	log logging.Logger
}

func NewServer(port string, logPath string, logger logging.Logger) Server {
	return Server{
		port: port,
		logPath: logPath,
		log: logger,
	}
}

func (s *Server) Init() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/channel/:channel/user/:username", s.getCurrentChannelLogs)
	e.GET("/channel/:channel/user/:username/:year/:month", s.getDatedChannelLogs)
	e.GET("/channel/:channel/user/:username/random", s.getRandomQuote)

	s.log.Error(e.Start("127.0.0.1:" + s.port).Error())
}