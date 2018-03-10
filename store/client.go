package store

import (
	"database/sql"
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
