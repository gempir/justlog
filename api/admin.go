package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) authenticateAdmin(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.Header.Get("X-Api-Key")

	if apiKey == "" || apiKey != s.cfg.AdminAPIKey {
		http.Error(w, "No I don't think so.", http.StatusForbidden)
		return false
	}

	return true
}

type channelsDeleteRequest struct {
	// list of userIds
	Channels []string `json:"channels"`
}

type channelConfigsJoinRequest struct {
	// list of userIds
	Channels []string `json:"channels"`
}

// swagger:route POST /admin/channels admin addChannels
//
// Will add the channels to log, only works with userIds
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
// Will remove the channels to log, only works with userIds
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
			s.bot.Part(userData.Login)
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

	writeJSON(fmt.Sprintf("Doubters? Joined channels or already in: %v", request.Channels), http.StatusOK, w, r)
}
