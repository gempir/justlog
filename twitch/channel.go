package twitch

import "strings"

type Channel struct {
	Name string
}

func NewChannel(channel string) Channel {
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}
	return Channel{
		Name: strings.ToLower(channel),
	}
}
