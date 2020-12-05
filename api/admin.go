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

// swagger:model
type channelConfigsRequest struct {
	config.ChannelConfig
}

// swagger:route POST /admin/channelConfigs/{channelID} admin channelConfigs
//
// Will set the messageTypes logged for a channel
// https://github.com/gempir/go-twitch-irc/blob/master/message.go#L17
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Security:
//     - api_key:
//
//     Schemes: https
//
//     Responses:
//       200:
//		 400:
//       405:

// swagger:route DELETE /admin/channelConfigs/{channelID} admin deleteChannelConfigs
//
// Will reset the messageTypes logged for a channel
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Security:
//     - api_key:
//
//     Schemes: https
//
//     Responses:
//       200:
//		 400:
//       405:
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

		http.Error(w, "Uhhhhhh...", http.StatusBadRequest)
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

type channelsDeleteRequest struct {
	Channels []string `json:"channels"`
}

type channelConfigsJoinRequest struct {
	Channels []string `json:"channels"`
}

// swagger:route POST /admin/channels admin addChannels
//
// Will add the channels to log
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Security:
//     - api_key:
//
//     Schemes: https
//
//     Responses:
//       200:
//		 400:
//       405:
//       500:

// swagger:route DELETE /admin/channels admin deleteChannels
//
// Will remove the channels to log
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Security:
//     - api_key:
//
//     Schemes: https
//
//     Responses:
//       200:
//		 400:
//       405:
//       500:
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
