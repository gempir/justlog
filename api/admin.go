package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gempir/justlog/config"
)

func (s *Server) authenticateAdmin(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.Header.Get("X-Api-Key")

	if apiKey == "" {
		http.Error(w, "No I don't think so.", http.StatusForbidden)
		return false
	}

	return true
}

type channelConfigsRequest struct {
	config.ChannelConfig
}

func (s *Server) writeConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "We'll see, we'll see. The winner gets tea.", http.StatusMethodNotAllowed)
		return
	}

	channelID := strings.TrimPrefix(r.URL.String(), "/admin/channelConfigs/")

	if _, ok := s.cfg.ChannelConfigs[channelID]; !ok {
		http.Error(w, "Uhhhhhh... unkown channel", http.StatusBadRequest)
		return
	}

	var request channelConfigsRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "ANYWAYS: "+err.Error(), http.StatusBadRequest)
		return
	}

	s.cfg.SetMessageTypes(channelID, request.MessageTypes)
	s.bot.UpdateMessageTypesToLog()

	writeJSON("Doubters? Updated "+channelID+" messageTypes to "+fmt.Sprintf("%v", request.MessageTypes), http.StatusOK, w, r)
}
