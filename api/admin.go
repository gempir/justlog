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

	if apiKey == "" || apiKey != s.cfg.AdminAPIKey {
		http.Error(w, "No I don't think so.", http.StatusForbidden)
		return false
	}

	return true
}

type channelConfigsDeleteRequest struct {
	MessageTypes bool `json:"messageTypes,omitempty"`
}

type channelConfigsRequest struct {
	config.ChannelConfig
}

func (s *Server) writeConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "We'll see, we'll see. The winner gets tea.", http.StatusMethodNotAllowed)
		return
	}

	channelID := strings.TrimPrefix(r.URL.String(), "/admin/channelConfigs/")

	if _, ok := s.cfg.ChannelConfigs[channelID]; !ok {
		http.Error(w, "Uhhhhhh... unkown channel", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		var request channelConfigsDeleteRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "ANYWAYS: "+err.Error(), http.StatusBadRequest)
			return
		}

		if request.MessageTypes {
			s.cfg.ResetMessageTypes(channelID)
			s.bot.UpdateMessageTypesToLog()
			writeJSON("Doubters? Reset "+channelID+" messageTypes", http.StatusOK, w, r)
			return
		}

	} else {
		var request channelConfigsRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "ANYWAYS: "+err.Error(), http.StatusBadRequest)
			return
		}

		s.cfg.SetMessageTypes(channelID, request.MessageTypes)
		s.bot.UpdateMessageTypesToLog()

		writeJSON("Doubters? Updated "+channelID+" messageTypes to "+fmt.Sprintf("%v", request.MessageTypes), http.StatusOK, w, r)
		return
	}

	http.Error(w, "Nothing here", http.StatusBadRequest)
	return
}
