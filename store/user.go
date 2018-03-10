package store

import (
	"log"

	twitch "github.com/gempir/go-twitch-irc"
)

// SetUser update or add user to db
func (c *Client) SetUser(user twitch.User) {
	stmt, err := c.db.Prepare("REPLACE INTO users(`id`, `display_name`, `username`, `user_type`, `last_typed`, `last_seen`) VALUES(?, ?, ?, ?, NOW(), NOW())")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = stmt.Exec(user.UserID, user.DisplayName, user.Username, user.UserType)
	if err != nil {
		log.Println(err.Error())
	}
}
