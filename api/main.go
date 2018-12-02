package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo"
)

type order string

var (
	orderDesc order = "DESC"
	orderAsc  order = "ASC"
)

type chatLog struct {
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Text      string             `json:"text"`
	Username  string             `json:"username"`
	Channel   string             `json:"channel"`
	Timestamp timestamp          `json:"timestamp"`
	Type      twitch.MessageType `json:"type"`
}

type timestamp struct {
	time.Time
}

func (t timestamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.UTC().Format(time.RFC3339) + "\""), nil
}

func (t *timestamp) UnmarshalJSON(data []byte) error {
	goTime, err := time.Parse(time.RFC3339, strings.TrimSuffix(strings.TrimPrefix(string(data[:]), "\""), "\""))
	if err != nil {
		return err
	}
	*t = timestamp{
		goTime,
	}
	return nil
}

func parseFromTo(from, to string, limit float64) (time.Time, time.Time, error) {
	var fromTime time.Time
	var toTime time.Time

	if from == "" && to == "" {
		fromTime = time.Now().AddDate(0, -1, 0)
		toTime = time.Now()
	} else if from == "" && to != "" {
		var err error
		toTime, err = parseTimestamp(to)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("Can't parse to timestamp: %s", err)
		}
		fromTime = toTime.AddDate(0, -1, 0)
	} else if from != "" && to == "" {
		var err error
		fromTime, err = parseTimestamp(from)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("Can't parse from timestamp: %s", err)
		}
		toTime = fromTime.AddDate(0, 1, 0)
	} else {
		var err error

		fromTime, err = parseTimestamp(from)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("Can't parse from timestamp: %s", err)
		}
		toTime, err = parseTimestamp(to)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("Can't parse to timestamp: %s", err)
		}

		if toTime.Sub(fromTime).Hours() > limit {
			return fromTime, toTime, errors.New("Timespan too big")
		}
	}

	return fromTime, toTime, nil
}

func writeTextResponse(c echo.Context, cLog *chatLog) error {
	c.Response().WriteHeader(http.StatusOK)

	for _, cMessage := range cLog.Messages {
		switch cMessage.Type {
		case twitch.PRIVMSG:
			c.Response().Write([]byte(fmt.Sprintf("[%s] #%s %s: %s\r\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Username, cMessage.Text)))
		case twitch.CLEARCHAT:
			c.Response().Write([]byte(fmt.Sprintf("[%s] #%s %s\r\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Text)))
		}
	}

	return nil
}

func writeJSONResponse(c echo.Context, logResult *chatLog) error {
	_, stream := c.QueryParams()["stream"]
	if stream {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		return json.NewEncoder(c.Response()).Encode(logResult)
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(logResult)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, echo.MIMEApplicationJSONCharsetUTF8, data)
}

func parseTimestamp(timestamp string) (time.Time, error) {

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil
}

func buildOrder(c echo.Context) order {
	dataOrder := orderAsc
	_, reverse := c.QueryParams()["reverse"]
	if reverse {
		dataOrder = orderDesc
	}

	return dataOrder
}
