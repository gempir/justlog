package store

import "log"

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

// AddChannel persist channel
func (c *Client) AddChannel(channel string) {
	stmt, err := c.db.Prepare("INSERT channels SET name=?")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = stmt.Exec(channel)
	if err != nil {
		log.Println(err.Error())
	}
}
