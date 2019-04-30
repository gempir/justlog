package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type RandomQuoteJSON struct {
	Channel     string    `json:"channel"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	Message     string    `json:"message"`
	Timestamp   timestamp `json:"timestamp"`
}

func (s *Server) getLastUserLogs(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	year, month, err := s.fileLogger.GetLastLogYearAndMonthForUser(channelID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse{"No logs found"})
	}

	redirectURL := fmt.Sprintf("/channelid/%s/userid/%s/%d/%d", channelID, userID, year, month)
	if len(c.QueryString()) > 0 {
		redirectURL += "?" + c.QueryString()
	}
	return c.Redirect(303, redirectURL)
}

// getLastUserLogsByName godoc
// @Summary Redirect to last logs of user
// @tags user
// @Produce  json
// @Produce  plain
// @Param channel path string true "channelname"
// @Param username path string true "username"
// @Param json query any false "response as json"
// @Param type query string false "define response type only json supported currently, rest defaults to plain"
// @Success 303
// @Router /channel/{channel}/user/{username} [get]
func (s *Server) getLastUserLogsByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	username := strings.ToLower(c.Param("username"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse{"Failure fetching userIDs"})
	}
	var year int
	var month int
	year, month, err = s.fileLogger.GetLastLogYearAndMonthForUser(userMap[channel].ID, userMap[username].ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse{"No logs found"})
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
		return c.JSON(http.StatusInternalServerError, errorResponse{"Failure fetching userIDs"})
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

func (s *Server) getUserLogsByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	username := strings.ToLower(c.Param("username"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse{"Failure fetching userIDs"})
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

// getRandomQuoteByName godoc
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
		return c.JSON(http.StatusInternalServerError, errorResponse{"Failure fetching userIDs"})
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
				Message:     message.Message,
				Timestamp:   timestamp{message.Time},
			}

			return c.JSON(http.StatusOK, randomQ)
		}

		return c.String(http.StatusOK, message.Message)
	}

	return c.String(http.StatusNotFound, "No quote found")
}

func (s *Server) getUserLogs(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse{"Invalid year"})
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse{"Invalid month"})
	}

	logMessages, err := s.fileLogger.ReadLogForUser(channelID, userID, year, month)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse{"Failure reading log"})
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
			}
		}

		logResult.Messages = append(logResult.Messages, chatMsg)
	}

	if shouldRespondWithJson(c) {
		return writeJSONResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}

func (s *Server) getUserLogsRange(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	fromTime, toTime, err := parseFromTo(c.QueryParam("from"), c.QueryParam("to"), userHourLimit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{err.Error()})
	}

	var logMessages []string

	logMessages, _ = s.fileLogger.ReadLogForUser(channelID, userID, fromTime.Year(), int(fromTime.Month()))

	if fromTime.Month() != toTime.Month() {
		additionalMessages, _ := s.fileLogger.ReadLogForUser(channelID, userID, toTime.Year(), int(toTime.Month()))

		logMessages = append(logMessages, additionalMessages...)
	}

	if len(logMessages) == 0 {
		return c.JSON(http.StatusNotFound, errorResponse{"No logs found"})
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
				Text:        message.Message,
				Type:        message.Type,
				Channel:     message.Channel,
			}
		}

		logResult.Messages = append(logResult.Messages, chatMsg)
	}

	if shouldRespondWithJson(c) {
		return writeJSONResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}
