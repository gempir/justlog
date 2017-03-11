package combo

import (
	"github.com/gempir/gempbotgo/twitch"
)

type Handler struct {
	lastEmote twitch.Emote
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) HandleMessage(msg twitch.Message) {

	// filter out messages without emotes or with different emotes
	if len(msg.Emotes) != 1 {
		return
	}

}
