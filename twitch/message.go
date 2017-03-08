package twitch

type Message struct {
	Text    string
	User    User
	Channel string
}

func newMessage(text string, user User, channel string) Message {
	return Message{
		Text: text,
		User: user,
		Channel: channel,
	}
}