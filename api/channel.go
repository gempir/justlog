package api

import (
	"errors"
	"strconv"

	"github.com/gempir/go-twitch-irc/v2"
)

// swagger:route GET /channel/{channel} logs channelLogs
//
// Get entire channel logs of current day
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Responses:
//       200: chatLog

// swagger:route GET /channel/{channel}/{year}/{month}/{day} logs channelLogsYearMonthDay
//
// Get entire channel logs of given day
//
//     Produces:
//     - application/json
//     - text/plain
//
//     Responses:
//       200: chatLog
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

	logResult := createLogResult()

	for _, rawMessage := range logMessages {
		logResult.Messages = append(logResult.Messages, createChatMessage(twitch.ParseMessage(rawMessage)))
	}

	return &logResult, nil
}

func (s *Server) getChannelLogsRange(request logRequest) (*chatLog, error) {
	fromTime, toTime, err := parseFromTo(request.time.from, request.time.to, channelHourLimit)
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
