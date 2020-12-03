package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gempir/go-twitch-irc/v2"
)

// RandomQuoteJSON response when request a random message
type RandomQuoteJSON struct {
	Channel     string    `json:"channel"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	Message     string    `json:"message"`
	Timestamp   timestamp `json:"timestamp"`
}

func (s *Server) getRandomQuote(request logRequest) (*chatLog, error) {
	rawMessage, err := s.fileLogger.ReadRandomMessageForUser(request.channelid, request.userid)
	if err != nil {
		return &chatLog{}, err
	}
	parsedMessage := twitch.ParseMessage(rawMessage)

	chatMsg := createChatMessage(parsedMessage)

	return &chatLog{Messages: []chatMessage{chatMsg}}, nil
}

// swagger:route GET /list logs list
//
// Lists available logs of a user
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Schemes: https
//
//     Responses:
//       200: logList
func (s *Server) writeAvailableLogs(w http.ResponseWriter, r *http.Request, q url.Values) {
	logs, err := s.fileLogger.GetAvailableLogsForUser(q.Get("channelid"), q.Get("userid"))
	if err != nil {
		http.Error(w, "failed to get available logs: "+err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(&logList{logs}, http.StatusOK, w, r)
}

// swagger:route GET /channel/{channel}/user/{username} logs userLogs
//
// Get user logs in channel of current month
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Responses:
//       200: chatLog

// swagger:route GET /channel/{channel}/user/{username}/{year}/{month} logs userLogsYearMonth
//
// Get user logs in channel of given year month
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Responses:
//       200: chatLog
func (s *Server) getUserLogs(request logRequest) (*chatLog, error) {
	logMessages, err := s.fileLogger.ReadLogForUser(request.channelid, request.userid, request.time.year, request.time.month)
	if err != nil {
		return &chatLog{}, err
	}

	if request.reverse {
		reverseSlice(logMessages)
	}

	logResult := createLogResult()

	for _, rawMessage := range logMessages {
		logResult.Messages = append(logResult.Messages, createChatMessage(twitch.ParseMessage(rawMessage)))
	}

	return &logResult, nil
}

func (s *Server) getUserLogsRange(request logRequest) (*chatLog, error) {

	fromTime, toTime, err := parseFromTo(request.time.from, request.time.to, userHourLimit)
	if err != nil {
		return &chatLog{}, err
	}

	var logMessages []string

	logMessages, _ = s.fileLogger.ReadLogForUser(request.channelid, request.userid, fmt.Sprintf("%d", fromTime.Year()), fmt.Sprintf("%d", int(fromTime.Month())))

	if fromTime.Month() != toTime.Month() {
		additionalMessages, _ := s.fileLogger.ReadLogForUser(request.channelid, request.userid, string(toTime.Year()), string(toTime.Month()))

		logMessages = append(logMessages, additionalMessages...)
	}

	if request.reverse {
		reverseSlice(logMessages)
	}

	logResult := createLogResult()

	for _, rawMessage := range logMessages {
		parsedMessage := twitch.ParseMessage(rawMessage)

		switch message := parsedMessage.(type) {
		case *twitch.PrivateMessage:
			if message.Time.Unix() < fromTime.Unix() || message.Time.Unix() > toTime.Unix() {
				continue
			}
		case *twitch.ClearChatMessage:
			if message.Time.Unix() < fromTime.Unix() || message.Time.Unix() > toTime.Unix() {
				continue
			}
		case *twitch.UserNoticeMessage:
			if message.Time.Unix() < fromTime.Unix() || message.Time.Unix() > toTime.Unix() {
				continue
			}
		}

		logResult.Messages = append(logResult.Messages, createChatMessage(parsedMessage))
	}

	return &logResult, nil
}
