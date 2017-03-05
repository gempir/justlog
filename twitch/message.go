package twitch

type message struct {
	Text string
	User user
	Channel string
}

func newMessage(text string, user user, channel string) message {
	return message{
		Text: text,
		User: user,
		Channel: channel,
	}
}