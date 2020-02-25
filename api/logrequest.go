package api

import (
	"net/url"
	"regexp"
)

type logRequest struct {
	channel      string
	user         string
	channelid    string
	userid       string
	time         logTime
	reverse      bool
	responseType string
}

// userRandomMessageRequest /channel/pajlada/user/gempir/random

// userRangeRequst /channel/pajlada/user/gempir/range?from=123&to=124
// userDateRequest /channel/pajlada/user/gempir/2020/2
// userRequest /channel/pajlada/user/gempir --> redirect ^
// channelRangeRequest /channel/pajlada/range?from=123&to=124
// channelDateRequest /channel/pajlada/2020/2/25
// channelRequest /channel/pajlada -- redirect ^

type logTime struct {
	from  string
	to    string
	year  string
	month string
	day   string
}

var (
	pathRegex = regexp.MustCompile(`\/(channel|channelid)\/([a-zA-Z0-9]+)(?:\/(user|userid)\/([a-zA-Z0-9]+))?(?:(?:\/(\d{4})\/(\d{1,2})(?:\/(\d{1,2}))?)|(?:\/(range)))?`)
)

func newLogRequestFromURL(url string, query url.Values) logRequest {
	// if strings.StartsWith(url, "/channel") {
	// 	// do channel stuff
	// }

	return logRequest{}
}
