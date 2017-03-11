package combo

import (
	"github.com/gempir/gempbotgo/twitch"
)

type Handler struct {
	lastEmote twitch.Emote
}

func NewHandler() Handler {
	return Handler{

	}
}


func (h *Handler) HandleMessage(msg twitch.Message) {

	// filter out messages without emotes
	if len(msg.User.Emotes) == 0 {
		return
	}

	// filter out messages with multiple emotes
	for _, emote := range msg.User.Emotes {
		if msg.User.Emotes[0].Id != emote.Id {
			return
		}
	}

}