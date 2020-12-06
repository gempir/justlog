package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gempir/justlog/bot"
	"github.com/gempir/justlog/filelog"
	"github.com/gempir/justlog/helix"

	"github.com/gempir/justlog/config"
)

type helixClientMock struct {
}

func (m *helixClientMock) GetUsersByUserIds(channels []string) (map[string]helix.UserData, error) {
	data := make(map[string]helix.UserData)
	data["77829817"] = helix.UserData{Login: "gempir", ID: "77829817"}

	return data, nil
}

func (m *helixClientMock) GetUsersByUsernames(channels []string) (map[string]helix.UserData, error) {
	data := make(map[string]helix.UserData)
	data["gempir"] = helix.UserData{Login: "gempir", ID: "77829817"}

	return data, nil
}

func createTestServer() *Server {
	cfg := &config.Config{LogsDirectory: "./logs", Channels: []string{"gempir"}, ClientID: "123"}

	fileLogger := filelog.NewFileLogger(cfg.LogsDirectory)
	helixClient := new(helixClientMock)

	bot := bot.NewBot(cfg, helixClient, &fileLogger)

	apiServer := NewServer(cfg, bot, &fileLogger, helixClient, cfg.Channels)

	return &apiServer
}

func TestApiServer(t *testing.T) {
	server := createTestServer()

	t.Run("get channels", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/channels", nil)
		w := httptest.NewRecorder()

		server.route(w, r)
		assert.Contains(t, w.Body.String(), "gempir")
	})

	t.Run("get user logs", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/channel/gempir/user/gempir", nil)
		w := httptest.NewRecorder()

		server.route(w, r)
		assert.Equal(t, w.Code, 302)
	})

	t.Run("get channel logs", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/channel/gempir", nil)
		w := httptest.NewRecorder()

		server.route(w, r)
		assert.Equal(t, w.Code, 302)
	})
}
