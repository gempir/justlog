package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (s *Server) getCurrentChannelLogs(c echo.Context) error {
	channelID := c.Param("channelid")
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()

	redirectURL := fmt.Sprintf("/channelid/%s/%d/%d/%d", channelID, year, month, day)
	if len(c.QueryString()) > 0 {
		redirectURL += "?" + c.QueryString()
	}
	return c.Redirect(http.StatusSeeOther, redirectURL)
}

func (s *Server) getCurrentChannelLogsByName(c echo.Context) error {
	channel := c.Param("channel")
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()

	redirectURL := fmt.Sprintf("/channel/%s/%d/%d/%d", channel, year, month, day)
	if len(c.QueryString()) > 0 {
		redirectURL += "?" + c.QueryString()
	}
	return c.Redirect(http.StatusSeeOther, redirectURL)
}

func (s *Server) getChannelLogsByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Failure fetching data from twitch")
	}

	names := c.ParamNames()
	names = append(names, "channelid")

	values := c.ParamValues()
	values = append(values, userMap[channel].ID)

	c.SetParamNames(names...)
	c.SetParamValues(values...)

	return s.getChannelLogs(c)
}

func (s *Server) getChannelLogsRangeByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Failure fetching data from twitch")
	}

	names := c.ParamNames()
	names = append(names, "channelid")

	values := c.ParamValues()
	values = append(values, userMap[channel].ID)

	c.SetParamNames(names...)
	c.SetParamValues(values...)

	return s.getChannelLogsRange(c)
}

func (s *Server) getChannelLogs(c echo.Context) error {
	channelID := c.Param("channelid")

	yearStr := c.Param("year")
	monthStr := c.Param("month")
	dayStr := c.Param("day")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Invalid year")
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Invalid month")
	}
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Invalid day")
	}

	logMessages, err := s.fileLogger.ReadLogForChannel(channelID, year, month, day)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Failure reading log")
	}

	if shouldReverse(c) {
		reverse(logMessages)
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
			}
		}

		logResult.Messages = append(logResult.Messages, chatMsg)
	}

	if shouldRespondWithJson(c) {
		return writeJSONResponse(c, &logResult)
	}

	if shouldRespondWithRaw(c) {
		return writeRawResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}

func (s *Server) getChannelLogsRange(c echo.Context) error {
	channelID := c.Param("channelid")

	fromTime, toTime, err := parseFromTo(c.QueryParam("from"), c.QueryParam("to"), channelHourLimit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var logMessages []string

	logMessages, _ = s.fileLogger.ReadLogForChannel(channelID, fromTime.Year(), int(fromTime.Month()), fromTime.Day())

	if fromTime.Month() != toTime.Month() {
		additionalMessages, _ := s.fileLogger.ReadLogForChannel(channelID, toTime.Year(), int(toTime.Month()), toTime.Day())

		logMessages = append(logMessages, additionalMessages...)
	}

	if len(logMessages) == 0 {
		return c.JSON(http.StatusNotFound, "No logs found")
	}

	if shouldReverse(c) {
		reverse(logMessages)
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
			}
		}

		logResult.Messages = append(logResult.Messages, chatMsg)
	}

	if shouldRespondWithJson(c) {
		return writeJSONResponse(c, &logResult)
	}

	if shouldRespondWithRaw(c) {
		return writeRawResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}
