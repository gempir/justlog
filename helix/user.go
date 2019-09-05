package helix

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	clientID   string
	httpClient *http.Client
}

var (
	userCacheByID       map[string]UserData
	userCacheByUsername map[string]UserData
)

func init() {
	userCacheByID = map[string]UserData{}
	userCacheByUsername = map[string]UserData{}
}

func NewClient(clientID string) Client {
	return Client{
		clientID:   clientID,
		httpClient: &http.Client{},
	}
}

type userResponse struct {
	Data []UserData `json:"data"`
}

type UserData struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

func (c *Client) GetUsersByUserIds(userIDs []string) (map[string]UserData, error) {
	var filteredUserIDs []string

	result := make(map[string]UserData)

	for _, id := range userIDs {
		if strings.HasPrefix(id, "chatrooms:") {
			result[id] = UserData{
				ID:    id,
				Login: id,
			}
			continue
		}

		if _, ok := userCacheByID[id]; ok {
			result[id] = userCacheByID[id]
		} else {
			filteredUserIDs = append(filteredUserIDs, id)
		}
	}

	if len(filteredUserIDs) == 1 {
		params := "?id=" + filteredUserIDs[0]

		err := c.makeRequest(params)
		if err != nil {
			return nil, err
		}

	} else if len(filteredUserIDs) > 1 {
		var params string

		for index, id := range filteredUserIDs {
			if index == 0 {
				params += "?id=" + id
			} else {
				params += "&id=" + id
			}
		}

		err := c.makeRequest(params)
		if err != nil {
			return nil, err
		}
	}

	for _, id := range filteredUserIDs {
		result[id] = userCacheByID[id]
	}

	return result, nil
}

func (c *Client) GetUsersByUsernames(usernames []string) (map[string]UserData, error) {
	var filteredUsernames []string

	result := make(map[string]UserData)

	for _, username := range usernames {
		if strings.HasPrefix(username, "chatrooms:") {
			result[username] = UserData{
				ID:    username,
				Login: username,
			}
			continue
		}

		if _, ok := userCacheByID[username]; ok {
			result[username] = userCacheByID[username]
		} else {
			filteredUsernames = append(filteredUsernames, username)
		}
	}

	if len(filteredUsernames) == 1 {
		params := "?login=" + filteredUsernames[0]

		err := c.makeRequest(params)
		if err != nil {
			return nil, err
		}

	} else if len(filteredUsernames) > 1 {
		var params string

		for index, id := range filteredUsernames {
			if index == 0 {
				params += "?login=" + id
			} else {
				params += "&login=" + id
			}
		}

		err := c.makeRequest(params)
		if err != nil {
			return nil, err
		}
	}

	for _, username := range filteredUsernames {
		result[username] = userCacheByUsername[username]
	}

	return result, nil
}

func (c *Client) makeRequest(parameters string) error {
	request, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users"+parameters, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Client-ID", c.clientID)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	log.Infof("%d GET https://api.twitch.tv/helix/users%s", response.StatusCode, parameters)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var userResp userResponse
	err = json.Unmarshal(contents, &userResp)
	if err != nil {
		return err
	}

	for _, user := range userResp.Data {
		userCacheByID[user.ID] = user
		userCacheByUsername[user.Login] = user
	}

	return nil
}
