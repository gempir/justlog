package twitch

import "time"

type Message struct {
	Text    string
	User    User
	Channel Channel
	Timestamp time.Time
}

func newMessage(text string, user User, channel Channel) Message {
	return Message{
		Text: text,
		User: user,
		Channel: channel,
		Timestamp: time.Now(),
	}
}