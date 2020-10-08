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

func (s *Server) writeChannelConfigs(w http.ResponseWriter, r *http.Request) {
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

type channelsDeleteRequest struct {
	Channels []string `json:"channels"`
}

type channelConfigsJoinRequest struct {
	Channels []string `json:"channels"`
}

func (s *Server) writeChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "We'll see, we'll see. The winner gets tea.", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodDelete {
		var request channelsDeleteRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "ANYWAYS: "+err.Error(), http.StatusBadRequest)
			return
		}

		s.cfg.RemoveChannels(request.Channels...)
		data, err := s.helixClient.GetUsersByUserIds(request.Channels)
		if err != nil {
			http.Error(w, "Failed to get channel names to leave, config might be already updated", http.StatusInternalServerError)
			return
		}
		for _, userData := range data {
			s.bot.Depart(userData.Login)
		}

		writeJSON(fmt.Sprintf("Doubters? Removed channels %v", request.Channels), http.StatusOK, w, r)
		return
	}

	var request channelConfigsJoinRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "ANYWAYS: "+err.Error(), http.StatusBadRequest)
		return
	}

	s.cfg.AddChannels(request.Channels...)
	data, err := s.helixClient.GetUsersByUserIds(request.Channels)
	if err != nil {
		http.Error(w, "Failed to get channel names to join, config might be already updated", http.StatusInternalServerError)
		return
	}
	for _, userData := range data {
		s.bot.Join(userData.Login)
	}

	writeJSON(fmt.Sprintf("Doubters? Joined channels %v", request.Channels), http.StatusOK, w, r)
}
