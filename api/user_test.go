package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gempir/justlog/bot"
	"github.com/gempir/justlog/filelog"
	"github.com/gempir/justlog/helix"
	"github.com/stretchr/testify/mock"

	"github.com/gempir/justlog/config"
)

type helixClientMock struct {
	mock.Mock
}

func (m *helixClientMock) GetUsersByUserIds(channels []string) (map[string]helix.UserData, error) {
	data := make(map[string]helix.UserData)
	data["test"] = helix.UserData{Login: "gempir", ID: "77829817"}

	return data, nil
}

func (m *helixClientMock) GetUsersByUsernames(channels []string) (map[string]helix.UserData, error) {
	data := make(map[string]helix.UserData)
	data["test"] = helix.UserData{Login: "gempir", ID: "77829817"}

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

func TestChannels(t *testing.T) {
	t.Run("Returns channels", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/channels", nil)
		w := httptest.NewRecorder()

		server := createTestServer()
		server.route(w, r)

		assert.Contains(t, w.Body.String(), "gempir")
	})
}
