package helix

import (
	"net/http"
	"strings"
	"sync"
	"time"

	helixClient "github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

// Client wrapper for helix
type Client struct {
	clientID       string
	clientSecret   string
	appAccessToken string
	client         *helixClient.Client
	httpClient     *http.Client
}

var (
	userCacheByID       sync.Map
	userCacheByUsername sync.Map
)

type TwitchApiClient interface {
	GetUsersByUserIds([]string) (map[string]UserData, error)
	GetUsersByUsernames([]string) (map[string]UserData, error)
}

// NewClient Create helix client
func NewClient(clientID string, clientSecret string) Client {
	client, err := helixClient.NewClient(&helixClient.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.RequestAppAccessToken([]string{})
	if err != nil {
		panic(err)
	}
	log.Infof("Requested access token, response: %d, expires in: %d", resp.StatusCode, resp.Data.ExpiresIn)
	client.SetAppAccessToken(resp.Data.AccessToken)

	return Client{
		clientID:       clientID,
		clientSecret:   clientSecret,
		appAccessToken: resp.Data.AccessToken,
		client:         client,
		httpClient:     &http.Client{},
	}
}

// UserData exported data from twitch
type UserData struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

// StartRefreshTokenRoutine refresh our token
func (c *Client) StartRefreshTokenRoutine() {
	ticker := time.NewTicker(24 * time.Hour)

	for range ticker.C {
		resp, err := c.client.RequestAppAccessToken([]string{})
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Requested access token from routine, response: %d, expires in: %d", resp.StatusCode, resp.Data.ExpiresIn)

		c.client.SetAppAccessToken(resp.Data.AccessToken)
	}
}

func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}

// GetUsersByUserIds receive userData for given ids
func (c *Client) GetUsersByUserIds(userIDs []string) (map[string]UserData, error) {
	var filteredUserIDs []string

	for _, id := range userIDs {
		if _, ok := userCacheByID.Load(id); !ok {
			filteredUserIDs = append(filteredUserIDs, id)
		}
	}

	if len(filteredUserIDs) > 0 {
		chunks := chunkBy(filteredUserIDs, 100)

		for _, chunk := range chunks {
			resp, err := c.client.GetUsers(&helixClient.UsersParams{
				IDs: chunk,
			})
			if err != nil {
				return map[string]UserData{}, err
			}

			log.Infof("%d GetUsersByUserIds %v", resp.StatusCode, chunk)

			for _, user := range resp.Data.Users {
				data := &UserData{
					ID:              user.ID,
					Login:           user.Login,
					DisplayName:     user.Login,
					Type:            user.Type,
					BroadcasterType: user.BroadcasterType,
					Description:     user.Description,
					ProfileImageURL: user.ProfileImageURL,
					OfflineImageURL: user.OfflineImageURL,
					ViewCount:       user.ViewCount,
					Email:           user.Email,
				}
				userCacheByID.Store(user.ID, data)
				userCacheByUsername.Store(user.Login, data)
			}
		}
	}

	result := make(map[string]UserData)

	for _, id := range userIDs {
		value, ok := userCacheByID.Load(id)
		if !ok {
			log.Debugf("Could not find userId, channel might be banned: %s", id)
			continue
		}
		result[id] = *(value.(*UserData))
	}

	return result, nil
}

// GetUsersByUsernames fetches userdata from helix
func (c *Client) GetUsersByUsernames(usernames []string) (map[string]UserData, error) {
	var filteredUsernames []string

	for _, username := range usernames {
		username = strings.ToLower(username)
		if _, ok := userCacheByUsername.Load(username); !ok {
			filteredUsernames = append(filteredUsernames, username)
		}
	}

	if len(filteredUsernames) > 0 {
		chunks := chunkBy(filteredUsernames, 100)

		for _, chunk := range chunks {
			resp, err := c.client.GetUsers(&helixClient.UsersParams{
				Logins: chunk,
			})
			if err != nil {
				return map[string]UserData{}, err
			}

			log.Infof("%d GetUsersByUsernames %v", resp.StatusCode, chunk)

			for _, user := range resp.Data.Users {
				data := &UserData{
					ID:              user.ID,
					Login:           user.Login,
					DisplayName:     user.Login,
					Type:            user.Type,
					BroadcasterType: user.BroadcasterType,
					Description:     user.Description,
					ProfileImageURL: user.ProfileImageURL,
					OfflineImageURL: user.OfflineImageURL,
					ViewCount:       user.ViewCount,
					Email:           user.Email,
				}
				userCacheByID.Store(user.ID, data)
				userCacheByUsername.Store(user.Login, data)
			}
		}
	}

	result := make(map[string]UserData)

	for _, username := range usernames {
		username = strings.ToLower(username)
		value, ok := userCacheByUsername.Load(username)
		if !ok {
			log.Debugf("Could not find username, channel might be banned: %s", username)
			continue
		}
		result[username] = *(value.(*UserData))
	}

	return result, nil
}
