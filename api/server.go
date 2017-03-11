package api

import (
	"github.com/labstack/echo"
	"net/http"
)

type Server struct {
	port    string
	logPath string
}

func NewServer(port string, logPath string) Server {
	return Server{
		port:    port,
		logPath: logPath,
	}
}

func (s *Server) Init() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/channel/:channel/user/:username", s.getCurrentUserLogs)
	e.GET("/channel/:channel", s.getCurrentChannelLogs)
	e.GET("/channel/:channel/:year/:month/:day", s.getDatedChannelLogs)
	e.GET("/channel/:channel/user/:username/:year/:month", s.getDatedUserLogs)
	e.GET("/channel/:channel/user/:username/random", s.getRandomQuote)

	e.Logger.Fatal(e.Start("127.0.0.1:" + s.port))
}
