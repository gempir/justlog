package api

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gempir/go-twitch-irc/v2"
)

func (s *Server) getChannelLogs(request logRequest) (*chatLog, error) {
	yearStr := request.time.year
	monthStr := request.time.month
	dayStr := request.time.day

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return &chatLog{}, errors.New("invalid year")
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return &chatLog{}, errors.New("invalid month")
	}
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return &chatLog{}, errors.New("invalid day")
	}

	logMessages, err := s.fileLogger.ReadLogForChannel(request.channelid, year, month, day)
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
				ID:          message.ID,
			}
		case *twitch.ClearChatMessage:
			message := *parsedMessage.(*twitch.ClearChatMessage)

			var text string
			if message.BanDuration == 0 {
				text = fmt.Sprintf("%s has been banned", message.TargetUsername)
			} else {
				text = fmt.Sprintf("%s has been timed out for %d seconds", message.TargetUsername, message.BanDuration)
			}

			chatMsg = chatMessage{
				Timestamp:   timestamp{message.Time},
				Username:    message.TargetUsername,
				DisplayName: message.TargetUsername,
				Text:        text,
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
				ID:          message.ID,
			}
		}

		logResult.Messages = append(logResult.Messages, chatMsg)
	}

	return &logResult, nil
}

func (s *Server) getChannelLogsRange(request logRequest) (*chatLog, error) {
	fromTime, toTime, err := parseFromTo(request.time.from, request.time.to, userHourLimit)
	if err != nil {
		return &chatLog{}, err
	}

	var logMessages []string

	logMessages, _ = s.fileLogger.ReadLogForChannel(request.channelid, fromTime.Year(), int(fromTime.Month()), fromTime.Day())

	if fromTime.Month() != toTime.Month() {
		additionalMessages, _ := s.fileLogger.ReadLogForChannel(request.channelid, toTime.Year(), int(toTime.Month()), toTime.Day())

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
				ID:          message.ID,
			}
		case *twitch.ClearChatMessage:
			message := *parsedMessage.(*twitch.ClearChatMessage)

			if message.Time.Unix() < fromTime.Unix() || message.Time.Unix() > toTime.Unix() {
				continue
			}

			var text string
			if message.BanDuration == 0 {
				text = fmt.Sprintf("%s has been banned", message.TargetUsername)
			} else {
				text = fmt.Sprintf("%s has been timed out for %d seconds", message.TargetUsername, message.BanDuration)
			}

			chatMsg = chatMessage{
				Timestamp:   timestamp{message.Time},
				Username:    message.TargetUsername,
				DisplayName: message.TargetUsername,
				Text:        text,
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
				ID:          message.ID,
			}
		}

		logResult.Messages = append(logResult.Messages, chatMsg)
	}

	return &logResult, nil
}
