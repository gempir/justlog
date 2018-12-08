package api

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

func (s *Server) getCurrentUserLogs(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	year := time.Now().Year()
	month := int(time.Now().Month())

	redirectURL := fmt.Sprintf("/channelid/%s/userid/%s/%d/%d", channelID, userID, year, month)
	return c.Redirect(303, redirectURL)
}

func (s *Server) getCurrentUserLogsByName(c echo.Context) error {
	channel := c.Param("channel")
	username := c.Param("username")

	year := time.Now().Year()
	month := int(time.Now().Month())

	redirectURL := fmt.Sprintf("/channel/%s/user/%s/%d/%d", channel, username, year, month)
	return c.Redirect(303, redirectURL)
}

func (s *Server) getUserLogsRangeByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	username := strings.ToLower(c.Param("username"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Failure fetching userIDs")
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
		return c.JSON(http.StatusInternalServerError, "Failure fetching userIDs")
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

func (s *Server) getRandomQuoteByName(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	username := strings.ToLower(c.Param("username"))

	userMap, err := s.helixClient.GetUsersByUsernames([]string{channel, username})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Failure fetching userIDs")
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
	userID := c.Param("userid")
	channelID := c.Param("channelid")

	var userLogs []string
	var lines []string

	if channelID == "" || userID == "" {
		return c.JSON(http.StatusNotFound, "error finding logs")
	}

	years, _ := ioutil.ReadDir(s.logPath + "/" + channelID)
	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(s.logPath + "/" + channelID + "/" + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()
			path := fmt.Sprintf("%s/%s/%s/%s/%s.txt", s.logPath, channelID, year, month, userID)
			if _, err := os.Stat(path); err == nil {
				userLogs = append(userLogs, path)
			} else if _, err := os.Stat(path + ".gz"); err == nil {
				userLogs = append(userLogs, path+".gz")
			}
		}
	}
	if len(userLogs) < 1 {
		return c.JSON(http.StatusNotFound, "error finding logs")
	}

	for _, logFile := range userLogs {
		f, _ := os.Open(logFile)

		scanner := bufio.NewScanner(f)

		if strings.HasSuffix(logFile, ".gz") {
			gz, _ := gzip.NewReader(f)
			scanner = bufio.NewScanner(gz)
		}

		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		f.Close()
	}

	ranNum := rand.Intn(len(lines))
	channel, user, message := twitch.ParseMessage(lines[ranNum])

	if shouldRespondWithJson(c) {

		randomQ := RandomQuoteJSON{
			Channel:     channel,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Message:     message.Text,
			Timestamp:   timestamp{message.Time},
		}

		return c.JSON(http.StatusOK, randomQ)
	}

	return c.String(http.StatusOK, fmt.Sprintf("%s: %s", user.DisplayName, message.Text))
}

func (s *Server) getUserLogs(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	yearStr := c.Param("year")
	monthStr := c.Param("month")

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

	logMessages, err := s.fileLogger.ReadLogForUser(channelID, userID, year, month)
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

func (s *Server) getUserLogsRange(c echo.Context) error {
	channelID := c.Param("channelid")
	userID := c.Param("userid")

	fromTime, toTime, err := parseFromTo(c.QueryParam("from"), c.QueryParam("to"), userHourLimit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var logMessages []string

	logMessages, _ = s.fileLogger.ReadLogForUser(channelID, userID, fromTime.Year(), int(fromTime.Month()))

	if fromTime.Month() != toTime.Month() {
		additionalMessages, _ := s.fileLogger.ReadLogForUser(channelID, userID, toTime.Year(), int(toTime.Month()))

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
