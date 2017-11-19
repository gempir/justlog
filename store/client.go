package store

import (
	"database/sql"
	"log"
)

// Client store client for persisting and reading data
type Client struct {
	db *sql.DB
}

// NewClient create new store client
func NewClient(dsn string) (*Client, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

// GetAllChannels get all channels to
func (c *Client) GetAllChannels() []string {
	rows, err := c.db.Query("SELECT name FROM channels")
	if err != nil {
		log.Println(err.Error())
	}

	var channels []string

	for rows.Next() {
		var channel string
		err = rows.Scan(&channel)
		if err != nil {
			log.Println(err)
		}
		channels = append(channels, channel)
	}

	return channels
}
