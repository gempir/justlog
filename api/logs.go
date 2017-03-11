package api

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ErrorJSON simple json for default error response
type ErrorJSON struct {
	Error string `json:"Error"`
}

// Msg struct to define a simple message
type Msg struct {
	Channel       string `json:"channel"`
	Username      string `json:"username"`
	Message       string `json:"message"`
	Timestamp     string `json:"timestamp"`
	UnixTimestamp string `json:"unix_timestamp"`
	Duration      string `json:"duration"`
}

func (s *Server) getCurrentUserLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)
	year := time.Now().Year()
	month := time.Now().Month().String()
	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))

	redirectURL := fmt.Sprintf("/channel/%s/user/%s/%d/%s", channel, username, year, month)
	return c.Redirect(303, redirectURL)
}

func (s *Server) getCurrentChannelLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)
	year := time.Now().Year()
	month := time.Now().Month().String()
	day := time.Now().Day()

	redirectURL := fmt.Sprintf("/channel/%s/%d/%s/%d", channel, year, month, day)
	return c.Redirect(303, redirectURL)
}

func (s *Server) getDatedChannelLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)
	year := c.Param("year")
	month := strings.Title(c.Param("month"))
	day := c.Param("day")

	if year == "" || month == "" {
		year = strconv.Itoa(time.Now().Year())
		month = time.Now().Month().String()
	}

	content := ""

	file := fmt.Sprintf(s.logPath+"%s/%s/%s/%s/channel.txt", channel, year, month, day)
	if _, err := os.Stat(file + ".gz"); err == nil {
		file = file + ".gz"
		f, err := os.Open(file)
		if err != nil {
			errJSON := new(ErrorJSON)
			errJSON.Error = "error finding logs"
			return c.JSON(http.StatusNotFound, errJSON)
		}
		gz, err := gzip.NewReader(f)
		scanner := bufio.NewScanner(gz)
		if err != nil {
			errJSON := new(ErrorJSON)
			errJSON.Error = "error finding logs"
			return c.JSON(http.StatusNotFound, errJSON)
		}

		for scanner.Scan() {
			line := scanner.Text()
			content += line + "\r\n"
		}
		return c.String(http.StatusOK, content)
	} else {
		return c.File(file)
	}

}

func (s *Server) getDatedUserLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)
	year := c.Param("year")
	month := strings.Title(c.Param("month"))
	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))

	if year == "" || month == "" {
		year = strconv.Itoa(time.Now().Year())
		month = time.Now().Month().String()
	}

	content := ""

	file := fmt.Sprintf(s.logPath+"%s/%s/%s/%s.txt", channel, year, month, username)
	if _, err := os.Stat(file + ".gz"); err == nil {
		file = file + ".gz"
		f, err := os.Open(file)
		if err != nil {
			errJSON := new(ErrorJSON)
			errJSON.Error = "error finding logs"
			return c.JSON(http.StatusNotFound, errJSON)
		}
		gz, err := gzip.NewReader(f)
		scanner := bufio.NewScanner(gz)
		if err != nil {
			errJSON := new(ErrorJSON)
			errJSON.Error = "error finding logs"
			return c.JSON(http.StatusNotFound, errJSON)
		}

		for scanner.Scan() {
			line := scanner.Text()
			content += line + "\r\n"
		}
		return c.String(http.StatusOK, content)
	} else {
		return c.File(file)
	}

}

func (s *Server) getRandomQuote(c echo.Context) error {
	errJSON := new(ErrorJSON)
	errJSON.Error = "error finding logs"

	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)

	var userLogs []string
	var lines []string

	years, _ := ioutil.ReadDir(s.logPath + channel)
	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(s.logPath + channel + "/" + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()
			path := fmt.Sprintf("%s%s/%s/%s/%s.txt", s.logPath, channel, year, month, username)
			if _, err := os.Stat(path); err == nil {
				userLogs = append(userLogs, path)
			} else if _, err := os.Stat(path + ".gz"); err == nil {
				userLogs = append(userLogs, path+".gz")
			}
		}
	}
	if len(userLogs) < 1 {
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	file := userLogs[rand.Intn(len(userLogs))]

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return c.JSON(http.StatusNotFound, errJSON)
	}
	scanner := bufio.NewScanner(f)

	if strings.HasSuffix(file, ".gz") {
		gz, _ := gzip.NewReader(f)
		scanner = bufio.NewScanner(gz)
	}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}
	if len(lines) < 1 {
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	ranNum := rand.Intn(len(lines))
	line := lines[ranNum]
	lineSplit := strings.SplitN(line, "]", 2)
	return c.String(http.StatusOK, lineSplit[1])
}
