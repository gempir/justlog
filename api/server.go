package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// Server api server
type Server struct {
	port    string
	logPath string
}

// NewServer create Server
func NewServer() Server {
	return Server{
		logPath: "/var/twitch_logs",
	}
}

// Init api server
func (s *Server) Init() {

	e := echo.New()
	e.HideBanner = true

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/channel/:channel/user/:username", s.getCurrentUserLogs)
	e.GET("/channel/:channel", s.getCurrentChannelLogs)
	e.GET("/channel/:channel/:year/:month/:day", s.getDatedChannelLogs)
	e.GET("/channel/:channel/user/:username/:year/:month", s.getDatedUserLogs)
	e.GET("/channel/:channel/user/:username/random", s.getRandomQuote)

	fmt.Println("starting API on port :8025")
	e.Logger.Fatal(e.Start(":8025"))
}
