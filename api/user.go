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

	var chatMsg chatMessage
	switch parsedMessage.(type) {
	case *twitch.PrivateMessage:
		message := *parsedMessage.(*twitch.PrivateMessage)

		chatMsg = chatMessage{
			Timestamp:   timestamp{message.Time},
			Username:    message.User.Name,
			DisplayName: message.User.DisplayName,
			Text:        message.Message,
			Type:        message.Type,
			Channel:     message.Channel,
			Raw:         message.Raw,
		}
	case *twitch.ClearChatMessage:
		message := *parsedMessage.(*twitch.ClearChatMessage)

		chatMsg = chatMessage{
			Timestamp:   timestamp{message.Time},
			Username:    message.TargetUsername,
			DisplayName: message.TargetUsername,
			Text:        buildClearChatMessageText(message),
			Type:        message.Type,
			Channel:     message.Channel,
			Raw:         message.Raw,
		}
	case *twitch.UserNoticeMessage:
		message := *parsedMessage.(*twitch.UserNoticeMessage)

		chatMsg = chatMessage{
			Timestamp:   timestamp{message.Time},
			Username:    message.User.Name,
			DisplayName: message.User.DisplayName,
			Text:        message.SystemMsg + " " + message.Message,
			Type:        message.Type,
			Channel:     message.Channel,
			Raw:         message.Raw,
		}
	}

	return &chatLog{Messages: []chatMessage{chatMsg}}, nil
}

func (s *Server) writeAvailableLogs(w http.ResponseWriter, r *http.Request, q url.Values) {
	logs, err := s.fileLogger.GetAvailableLogsForUser(q.Get("channelid"), q.Get("userid"))
	if err != nil {
		http.Error(w, "failed to get available logs: "+err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(&logList{logs}, http.StatusOK, w, r)
}

func (s *Server) getUserLogs(request logRequest) (*chatLog, error) {
	logMessages, err := s.fileLogger.ReadLogForUser(request.channelid, request.userid, request.time.year, request.time.month)
	if err != nil {
		return &chatLog{}, err
	}

	if request.reverse {
		reverseSlice(logMessages)
	}

	var logResult chatLog

	for _, rawMessage := range logMessages {
		parsedMessage := twitch.ParseMessage(rawMessage)

		var chatMsg chatMessage

		switch parsedMessage.(type) {
		case *twitch.PrivateMessage:
			message := *parsedMessage.(*twitch.PrivateMessage)

			chatMsg = chatMessage{
				Timestamp:   timestamp{message.Time},
				Username:    message.User.Name,
				DisplayName: message.User.DisplayName,
				Text:        message.Message,
				Type:        message.Type,
				Channel:     message.Channel,
				Raw:         message.Raw,
			}
		case *twitch.ClearChatMessage:
			message := *parsedMessage.(*twitch.ClearChatMessage)

			chatMsg = chatMessage{
				Timestamp:   timestamp{message.Time},
				Username:    message.TargetUsername,
				DisplayName: message.TargetUsername,
				Text:        buildClearChatMessageText(message),
				Type:        message.Type,
				Channel:     message.Channel,
				Raw:         message.Raw,
			}
		case *twitch.UserNoticeMessage:
			message := *parsedMessage.(*twitch.UserNoticeMessage)

			chatMsg = chatMessage{
				Timestamp:   timestamp{message.Time},
				Username:    message.User.Name,
				DisplayName: message.User.DisplayName,
				Text:        message.SystemMsg + " " + message.Message,
				Type:        message.Type,
				Channel:     message.Channel,
				Raw:         message.Raw,
			}
		}

		logResult.Messages = append(logResult.Messages, chatMsg)
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

	var logResult chatLog

	for _, rawMessage := range logMessages {
		parsedMessage := twitch.ParseMessage(rawMessage)

		var chatMsg chatMessage

		switch parsedMessage.(type) {
		case *twitch.PrivateMessage:
			message := *parsedMessage.(*twitch.PrivateMessage)

			if message.Time.Unix() < fromTime.Unix() || message.Time.Unix() > toTime.Unix() {
				continue
			}

			chatMsg = chatMessage{
				Timestamp:   timestamp{message.Time},
				Username:    message.User.Name,
				DisplayName: message.User.DisplayName,
				Text:        message.Message,
				Type:        message.Type,
				Channel:     message.Channel,
				Raw:         message.Raw,
			}
		case *twitch.ClearChatMessage:
			message := *parsedMessage.(*twitch.ClearChatMessage)

			if message.Time.Unix() < fromTime.Unix() || message.Time.Unix() > toTime.Unix() {
				continue
			}

			chatMsg = chatMessage{
				Timestamp:   timestamp{message.Time},
				Username:    message.TargetUsername,
				DisplayName: message.TargetUsername,
				Text:        buildClearChatMessageText(message),
				Type:        message.Type,
				Channel:     message.Channel,
				Raw:         message.Raw,
			}
		case *twitch.UserNoticeMessage:
			message := *parsedMessage.(*twitch.UserNoticeMessage)

			chatMsg = chatMessage{
				Timestamp:   timestamp{message.Time},
				Username:    message.User.Name,
				DisplayName: message.User.DisplayName,
				Text:        message.SystemMsg + " " + message.Message,
				Type:        message.Type,
				Channel:     message.Channel,
				Raw:         message.Raw,
			}
		}

		logResult.Messages = append(logResult.Messages, chatMsg)
	}

	return &logResult, nil
}
