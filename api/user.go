package api

import (
	"fmt"

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

// func (s *Server) getLastUserLogs(c echo.Context) error {
// 	channelID := c.Param("channelid")
// 	userID := c.Param("userid")

// 	year, month, err := s.fileLogger.GetLastLogYearAndMonthForUser(channelID, userID)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, ErrorResponse{"No logs found"})
// 	}

// 	redirectURL := fmt.Sprintf("/channelid/%s/userid/%s/%d/%d", channelID, userID, year, month)
// 	if len(c.QueryString()) > 0 {
// 		redirectURL += "?" + c.QueryString()
// 	}
// 	return c.Redirect(303, redirectURL)
// }

// func (s *Server) getLastUserLogsByName(c echo.Context) error {
// 	channel := strings.ToLower(c.Param("channel"))
// 	username := strings.ToLower(c.Param("username"))

// 	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
// 	if err != nil {
// 		log.Error(err)
// 		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
// 	}
// 	var year int
// 	var month int
// 	year, month, err = s.fileLogger.GetLastLogYearAndMonthForUser(userMap[channel].ID, userMap[username].ID)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, ErrorResponse{"No logs found"})
// 	}

// 	redirectURL := fmt.Sprintf("/channel/%s/user/%s/%d/%d", channel, username, year, month)
// 	if len(c.QueryString()) > 0 {
// 		redirectURL += "?" + c.QueryString()
// 	}
// 	return c.Redirect(303, redirectURL)
// }

// func (s *Server) getUserLogsExact(c echo.Context) error {
// 	channel := strings.ToLower(c.Param("channel"))
// 	user := strings.ToLower(c.Param("user"))

// 	userMap := map[string]helix.UserData{}
// 	if c.Param("channelType") == "channel" || c.Param("userType") == "user" {
// 		var err error
// 		userMap, err = s.helixClient.GetUsersByUsernames([]string{channel, user})
// 		if err != nil {
// 			log.Error(err)
// 			return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
// 		}
// 	}

// 	names := c.ParamNames()
// 	values := c.ParamValues()
// 	names = append(names, "channelid")

// 	if c.Param("channelType") == "channel" {
// 		values = append(values, userMap[channel].ID)
// 	} else {
// 		values = append(values, c.Param("channel"))
// 	}

// 	names = append(names, "userid")
// 	if c.Param("userType") == "user" {
// 		values = append(values, userMap[user].ID)
// 	} else {
// 		values = append(values, c.Param("user"))
// 	}

// 	c.SetParamNames(names...)
// 	c.SetParamValues(values...)

// 	return s.getUserLogs(c)
// }

// func (s *Server) getUserLogsRangeByName(c echo.Context) error {
// 	channel := strings.ToLower(c.Param("channel"))
// 	username := strings.ToLower(c.Param("username"))

// 	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
// 	if err != nil {
// 		log.Error(err)
// 		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
// 	}

// 	names := c.ParamNames()
// 	names = append(names, "channelid")
// 	names = append(names, "userid")

// 	values := c.ParamValues()
// 	values = append(values, userMap[channel].ID)
// 	values = append(values, userMap[username].ID)

// 	c.SetParamNames(names...)
// 	c.SetParamValues(values...)

// 	return s.getUserLogsRange(c)
// }

// func (s *Server) getUserLogsByName(c echo.Context) error {
// 	channel := strings.ToLower(c.Param("channel"))
// 	username := strings.ToLower(c.Param("username"))

// 	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
// 	if err != nil {
// 		log.Error(err)
// 		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
// 	}

// 	names := c.ParamNames()
// 	names = append(names, "channelid")
// 	names = append(names, "userid")

// 	values := c.ParamValues()
// 	values = append(values, userMap[channel].ID)
// 	values = append(values, userMap[username].ID)

// 	c.SetParamNames(names...)
// 	c.SetParamValues(values...)

// 	if c.Param("time") == "range" {
// 		return s.getUserLogsRange(c)
// 	}

// 	return s.getUserLogs(c)
// }

// func (s *Server) getRandomQuoteByName(c echo.Context) error {
// 	channel := strings.ToLower(c.Param("channel"))
// 	username := strings.ToLower(c.Param("username"))

// 	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
// 	if err != nil {
// 		log.Error(err)
// 		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failure fetching data from twitch"})
// 	}

// 	names := c.ParamNames()
// 	names = append(names, "channelid")
// 	names = append(names, "userid")

// 	values := c.ParamValues()
// 	values = append(values, userMap[channel].ID)
// 	values = append(values, userMap[username].ID)

// 	c.SetParamNames(names...)
// 	c.SetParamValues(values...)

// 	return s.getRandomQuote(c)
// }

// func (s *Server) getRandomQuote(c echo.Context) error {
// 	channelID := c.Param("channelid")
// 	userID := c.Param("userid")

// 	rawMessage, err := s.fileLogger.ReadRandomMessageForUser(channelID, userID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}
// 	parsedMessage := twitch.ParseMessage(rawMessage)

// 	switch parsedMessage.(type) {
// 	case *twitch.PrivateMessage:
// 		message := *parsedMessage.(*twitch.PrivateMessage)

// 		if shouldRespondWithJSON(c) {

// 			randomQ := RandomQuoteJSON{
// 				Channel:     message.Channel,
// 				Username:    message.User.Name,
// 				DisplayName: message.User.DisplayName,
// 				Message:     message.Message,
// 				Timestamp:   timestamp{message.Time},
// 			}

// 			return c.JSON(http.StatusOK, randomQ)
// 		}

// 		return c.String(http.StatusOK, message.Message)
// 	case *twitch.ClearChatMessage:
// 		message := *parsedMessage.(*twitch.ClearChatMessage)

// 		if shouldRespondWithJSON(c) {

// 			randomQ := RandomQuoteJSON{
// 				Channel:     message.Channel,
// 				Username:    message.TargetUsername,
// 				DisplayName: message.TargetUsername,
// 				Message:     buildClearChatMessageText(message),
// 				Timestamp:   timestamp{message.Time},
// 			}

// 			return c.JSON(http.StatusOK, randomQ)
// 		}

// 		return c.String(http.StatusOK, buildClearChatMessageText(message))
// 	}

// 	return c.String(http.StatusNotFound, "No quote found")
// }

func (s *Server) getUserLogs(request userRequest) (*chatLog, error) {
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

func (s *Server) getUserLogsRange(request userRequest) (*chatLog, error) {

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

func buildClearChatMessageText(message twitch.ClearChatMessage) string {
	if message.BanDuration == 0 {
		return fmt.Sprintf("%s has been banned", message.TargetUsername)
	}

	return fmt.Sprintf("%s has been timed out for %d seconds", message.TargetUsername, message.BanDuration)
}
