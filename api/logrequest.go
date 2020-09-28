package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type logRequest struct {
	channel          string
	user             string
	channelid        string
	userid           string
	time             logTime
	reverse          bool
	responseType     string
	redirectPath     string
	isUserRequest    bool
	isChannelRequest bool
}

// userRandomMessageRequest /channel/pajlada/user/gempir/random

type logTime struct {
	from   string
	to     string
	year   string
	month  string
	day    string
	random bool
}

var (
	pathRegex = regexp.MustCompile(`\/(channel|channelid)\/(\w+)(?:\/(user|userid)\/(\w+))?(?:(?:\/(\d{4})\/(\d{1,2})(?:\/(\d{1,2}))?)|(?:\/(range|random)))?`)
)

func (s *Server) newLogRequestFromURL(r *http.Request) (logRequest, error) {
	path := r.URL.EscapedPath()

	if path != strings.ToLower(path) {
		return logRequest{redirectPath: fmt.Sprintf("%s?%s", strings.ToLower(path), r.URL.Query().Encode())}, nil
	}

	if !strings.HasPrefix(path, "/channel") && !strings.HasPrefix(path, "/channelid") {
		return logRequest{}, errors.New("route not found")
	}

	url := strings.TrimRight(path, "/")

	matches := pathRegex.FindAllStringSubmatch(url, -1)
	if len(matches) == 0 || len(matches[0]) < 5 {
		return logRequest{}, errors.New("route not found")
	}

	request := logRequest{
		time: logTime{},
	}

	params := []string{}
	for _, match := range matches[0] {
		if match != "" {
			params = append(params, match)
		}
	}

	request.isUserRequest = len(params) > 4 && (params[3] == "user" || params[3] == "userid")
	request.isChannelRequest = len(params) < 4 || (len(params) >= 4 && params[3] != "user" && params[3] != "userid")

	if params[1] == "channel" {
		request.channel = params[2]
	}
	if params[1] == "channelid" {
		request.channelid = params[2]
	}
	if request.isUserRequest && params[3] == "user" {
		request.user = params[4]
	}
	if request.isUserRequest && params[3] == "userid" {
		request.userid = params[4]
	}

	var err error
	request, err = s.fillIds(request)
	if err != nil {
		log.Error(err)
		return logRequest{}, nil
	}

	if request.isUserRequest && len(params) == 7 {
		request.time.year = params[5]
		request.time.month = params[6]
	} else if request.isUserRequest && len(params) == 8 {
		return logRequest{}, errors.New("route not found")
	} else if request.isChannelRequest && len(params) == 6 {
		request.time.year = params[3]
		request.time.month = params[4]
		request.time.day = params[5]
	} else if request.isUserRequest && len(params) == 6 && params[5] == "random" {
		request.time.random = true
	} else if (request.isUserRequest && len(params) == 6 && params[5] == "range") || (request.isChannelRequest && len(params) == 4 && params[3] == "range") {
		request.time.from = r.URL.Query().Get("from")
		request.time.to = r.URL.Query().Get("to")
	} else {
		if request.isChannelRequest {
			request.time.year = fmt.Sprintf("%d", time.Now().Year())
			request.time.month = fmt.Sprintf("%d", time.Now().Month())
		} else {
			year, month, err := s.fileLogger.GetLastLogYearAndMonthForUser(request.channelid, request.userid)
			if err == nil {
				request.time.year = fmt.Sprintf("%d", year)
				request.time.month = fmt.Sprintf("%d", month)
			} else {
				request.time.year = fmt.Sprintf("%d", time.Now().Year())
				request.time.month = fmt.Sprintf("%d", time.Now().Month())
			}
		}

		timePath := request.time.year + "/" + request.time.month

		if request.isChannelRequest {
			request.time.day = fmt.Sprintf("%d", time.Now().Day())
			timePath += "/" + request.time.day
		}

		query := r.URL.Query()
		query.Del("from")
		query.Del("to")

		encodedQuery := ""
		if query.Encode() != "" {
			encodedQuery = "?" + query.Encode()
		}

		return logRequest{redirectPath: fmt.Sprintf("%s/%s%s", params[0], timePath, encodedQuery)}, nil
	}

	if _, ok := r.URL.Query()["reverse"]; ok {
		request.reverse = true
	} else {
		request.reverse = false
	}

	if _, ok := r.URL.Query()["json"]; ok || r.URL.Query().Get("type") == "json" || r.Header.Get("Content-Type") == "application/json" {
		request.responseType = responseTypeJSON
	} else if _, ok := r.URL.Query()["raw"]; ok || r.URL.Query().Get("type") == "raw" {
		request.responseType = responseTypeRaw
	} else {
		request.responseType = responseTypeText
	}

	return request, nil
}

func (s *Server) fillIds(request logRequest) (logRequest, error) {
	usernames := []string{}
	if request.channelid == "" && request.channel != "" {
		usernames = append(usernames, request.channel)
	}
	if request.userid == "" && request.user != "" {
		usernames = append(usernames, request.user)
	}

	ids, err := s.helixClient.GetUsersByUsernames(usernames)
	if err != nil {
		return request, err
	}

	if request.channelid == "" {
		request.channelid = ids[strings.ToLower(request.channel)].ID
	}
	if request.userid == "" {
		request.userid = ids[strings.ToLower(request.user)].ID
	}

	return request, nil
}
