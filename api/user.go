package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type RandomQuoteJSON struct {
	Channel     string    `json:"channel"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	Message     string    `json:"message"`
	Timestamp   timestamp `json:"timestamp"`
}

// @Summary Redirect to last logs of user
// @tags user
// @Produce  json
// @Produce  plain
// @Param channelid path string true "twitch userid"
// @Param userid path string true "twitch userid"
// @Param from query int false "unix timestamp, limit logs by timestamps from this point"
// @Param to query int false "unix timestamp, limit logs by timestamps to this point"
// @Param json query any false "response as json"
// @Param type query string false "define response type only json supported currently, rest defaults to plain"
// @Success 303
// @Success 200
// @Failure 404
// @Router /channelid/{channelid}/user/{userid} [get]
func (s *Server) getLastUserLogs(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	year, month, err := s.fileLogger.GetLastLogYearAndMonthForUser(channelID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{"No logs found"})
	}

	redirectURL := fmt.Sprintf("/channelid/%s/userid/%s/%d/%d", channelID, userID, year, month)
	if len(c.QueryString()) > 0 {
		redirectURL += "?" + c.QueryString()
	}
	return c.Redirect(303, redirectURL)
}

// @Summary Redirect to last logs of user
// @tags user
// @Produce  json
// @Produce  plain
// @Param channel path string true "channelname"
// @Param username path string true "username"
// @Param from query int false "unix timestamp, limit logs by timestamps from this point"
// @Param to query int false "unix timestamp, limit logs by timestamps to this point"
// @Param json query any false "response as json"
// @Param type query string false "define response type only json supported currently, rest defaults to plain"
// @Success 303
// @Success 200
// @Failure 500
// @Failure 404
// @Router /channel/{channel}/user/{username} [get]
func (s *Server) getLastUserLogsByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	username := strings.ToLower(c.Param("username"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
	}
	var year int
	var month int
	year, month, err = s.fileLogger.GetLastLogYearAndMonthForUser(userMap[channel].ID, userMap[username].ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{"No logs found"})
	}

	redirectURL := fmt.Sprintf("/channel/%s/user/%s/%d/%d", channel, username, year, month)
	if len(c.QueryString()) > 0 {
		redirectURL += "?" + c.QueryString()
	}
	return c.Redirect(303, redirectURL)
}

func (s *Server) getUserLogsRangeByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	username := strings.ToLower(c.Param("username"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
	}

	names := c.ParamNames()
	names = append(names, "channelid")
	names = append(names, "userid")

	values := c.ParamValues()
	values = append(values, userMap[channel].ID)
	values = append(values, userMap[username].ID)

	c.SetParamNames(names...)
	c.SetParamValues(values...)

	return s.getUserLogsRange(c)
}

// @Summary Get logs for user by year and month
// @tags user
// @Produce  json
// @Produce  plain
// @Param channel path string true "channelname"
// @Param username path string true "username"
// @Param year path string true "year of logs"
// @Param month path string true "month of logs"
// @Param from query int false "unix timestamp, limit logs by timestamps from this point"
// @Param to query int false "unix timestamp, limit logs by timestamps to this point"
// @Param json query any false "response as json"
// @Param type query string false "define response type only json supported currently, rest defaults to plain"
// @Success 200
// @Failure 500
// @Router /channel/{channel}/user/{username}/{year}/{month} [get]
func (s *Server) getUserLogsByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	username := strings.ToLower(c.Param("username"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
	}

	names := c.ParamNames()
	names = append(names, "channelid")
	names = append(names, "userid")

	values := c.ParamValues()
	values = append(values, userMap[channel].ID)
	values = append(values, userMap[username].ID)

	c.SetParamNames(names...)
	c.SetParamValues(values...)

	return s.getUserLogs(c)
}

// @Summary Get a random chat message from a user
// @tags user
// @Produce  json
// @Produce  plain
// @Param channel path string true "channelname"
// @Param username path string true "username"
// @Param json query any false "response as json"
// @Param type query string false "define response type only json supported currently, rest defaults to plain"
// @Success 200 {object} api.RandomQuoteJSON json
// @Router /channel/{channel}/user/{username}/random [get]
func (s *Server) getRandomQuoteByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	username := strings.ToLower(c.Param("username"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
	}

	names := c.ParamNames()
	names = append(names, "channelid")
	names = append(names, "userid")

	values := c.ParamValues()
	values = append(values, userMap[channel].ID)
	values = append(values, userMap[username].ID)

	c.SetParamNames(names...)
	c.SetParamValues(values...)

	return s.getRandomQuote(c)
}

// @Summary Get a random chat message from a user
// @tags user
// @Produce  json
// @Produce  plain
// @Param channelid path string true "twitch userid"
// @Param userid path string true "twitch userid"
// @Param json query any false "response as json"
// @Param type query string false "define response type only json supported currently, rest defaults to plain"
// @Success 200 {object} api.RandomQuoteJSON json
// @Router /channelid/{channelid}/userid/{userid}/random [get]
func (s *Server) getRandomQuote(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	rawMessage, err := s.fileLogger.ReadRandomMessageForUser(channelID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	parsedMessage := twitch.ParseMessage(rawMessage)

	switch parsedMessage.(type) {
	case *twitch.PrivateMessage:
		message := *parsedMessage.(*twitch.PrivateMessage)

		if shouldRespondWithJson(c) {

			randomQ := RandomQuoteJSON{
				Channel:     message.Channel,
				Username:    message.User.Name,
				DisplayName: message.User.DisplayName,
				Message:     message.Message,
				Timestamp:   timestamp{message.Time},
			}

			return c.JSON(http.StatusOK, randomQ)
		}

		return c.String(http.StatusOK, message.Message)
	case *twitch.ClearChatMessage:
		message := *parsedMessage.(*twitch.ClearChatMessage)

		if shouldRespondWithJson(c) {

			randomQ := RandomQuoteJSON{
				Channel:     message.Channel,
				Username:    message.TargetUsername,
				DisplayName: message.TargetUsername,
				Message:     buildClearChatMessageText(message),
				Timestamp:   timestamp{message.Time},
			}

			return c.JSON(http.StatusOK, randomQ)
		}

		return c.String(http.StatusOK, buildClearChatMessageText(message))
	}

	return c.String(http.StatusNotFound, "No quote found")
}

// @Summary Get logs for user by year and month
// @tags user
// @Produce  json
// @Produce  plain
// @Param channelid path string true "twitch userid"
// @Param userid path string true "twitch userid"
// @Param year path string true "year of logs"
// @Param month path string true "month of logs"
// @Param from query int false "unix timestamp, limit logs by timestamps from this point"
// @Param to query int false "unix timestamp, limit logs by timestamps to this point"
// @Param json query any false "response as json"
// @Param type query string false "define response type only json supported currently, rest defaults to plain"
// @Success 200
// @Failure 500
// @Router /channelid/{channelid}/userid/{userid}/{year}/{month} [get]
func (s *Server) getUserLogs(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Invalid year"})
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Invalid month"})
	}

	logMessages, err := s.fileLogger.ReadLogForUser(channelID, userID, year, month)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure reading log"})
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

	if shouldRespondWithJson(c) {
		return writeJSONResponse(c, &logResult)
	}

	if shouldRespondWithRaw(c) {
		return writeRawResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}

func (s *Server) getUserLogsRange(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	fromTime, toTime, err := parseFromTo(c.QueryParam("from"), c.QueryParam("to"), userHourLimit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
	}

	var logMessages []string

	logMessages, _ = s.fileLogger.ReadLogForUser(channelID, userID, fromTime.Year(), int(fromTime.Month()))

	if fromTime.Month() != toTime.Month() {
		additionalMessages, _ := s.fileLogger.ReadLogForUser(channelID, userID, toTime.Year(), int(toTime.Month()))

		logMessages = append(logMessages, additionalMessages...)
	}

	if len(logMessages) == 0 {
		return c.JSON(http.StatusNotFound, ErrorResponse{"No logs found"})
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

	if shouldRespondWithJson(c) {
		return writeJSONResponse(c, &logResult)
	}

	if shouldRespondWithRaw(c) {
		return writeRawResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}

func buildClearChatMessageText(message twitch.ClearChatMessage) string {
	if message.BanDuration == 0 {
		return fmt.Sprintf("%s has been banned", message.TargetUsername)
	} else {
		return fmt.Sprintf("%s has been timed out for %d seconds", message.TargetUsername, message.BanDuration)
	}
}
