package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func (s *Server) getCurrentChannelLogs(c echo.Context) error {
	channelID := c.Param("channelid")
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()

	redirectURL := fmt.Sprintf("/channelid/%s/%d/%d/%d", channelID, year, month, day)
	return c.Redirect(http.StatusSeeOther, redirectURL)
}

func (s *Server) getCurrentChannelLogsByName(c echo.Context) error {
	channel := c.Param("channel")
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()

	redirectURL := fmt.Sprintf("/channel/%s/%d/%d/%d", channel, year, month, day)
	return c.Redirect(http.StatusSeeOther, redirectURL)
}

func (s *Server) getChannelLogsByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Failure fetching userID")
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
		return c.JSON(http.StatusInternalServerError, "Failure fetching userID")
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
		channel, user, parsedMessage := twitch.ParseMessage(rawMessage)

		message := chatMessage{
			Timestamp:   timestamp{parsedMessage.Time},
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Text:        parsedMessage.Text,
			Type:        parsedMessage.Type,
			Channel:     channel,
		}

		logResult.Messages = append(logResult.Messages, message)
	}

	if shouldRespondWithJson(c) {
		return writeJSONResponse(c, &logResult)
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
		channel, user, parsedMessage := twitch.ParseMessage(rawMessage)

		if parsedMessage.Time.Unix() < fromTime.Unix() || parsedMessage.Time.Unix() > toTime.Unix() {
			continue
		}

		message := chatMessage{
			Timestamp:   timestamp{parsedMessage.Time},
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Text:        parsedMessage.Text,
			Type:        parsedMessage.Type,
			Channel:     channel,
		}

		logResult.Messages = append(logResult.Messages, message)
	}

	if shouldRespondWithJson(c) {
		return writeJSONResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}
