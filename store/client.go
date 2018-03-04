package store

import (
	"bufio"
	"log"
	"os"
)

// Client store client for persisting and reading data
type Client struct {
	channelsFile string
}

// NewClient create new store client
func NewClient(channelsFile string) (*Client, error) {
	return &Client{
		channelsFile: channelsFile,
	}, nil
}

// GetAllChannels get all channels to
func (c *Client) GetAllChannels() []string {

	channels := []string{}

	f, err := os.Open(c.channelsFile)
	if err != nil {
		log.Println(err.Error())
		return []string{}
	}

	scanner := bufio.NewScanner(f)
	if err != nil {
		log.Println(err.Error())
		return []string{}
	}

	for scanner.Scan() {
		channels = append(channels, scanner.Text())
	}

	return channels
}

// AddChannel persist channel
func (c *Client) AddChannel(channel string) error {
	f, err := os.OpenFile(c.channelsFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(channel + "\n")
	if err != nil {
		return err
	}
	return nil
}
