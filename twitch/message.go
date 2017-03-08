package twitch

import "time"

type Message struct {
	Text    string
	User    User
	Channel string
	Timestamp time.Time
}

func newMessage(text string, user User, channel string) Message {
	return Message{
		Text: text,
		User: user,
		Channel: channel,
		Timestamp: time.Now(),
	}
}