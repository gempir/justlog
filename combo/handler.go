package combo

import (
	"github.com/gempir/gempbotgo/twitch"
	"fmt"
	"github.com/gempir/gempbotgo/modules"
)

type Handler struct {
	bot       *twitch.Bot
	lastEmote twitch.Emote
	comboCount int
}

func NewHandler(bot *twitch.Bot) Handler {
	return Handler{
		bot: bot,
	}
}

func (h *Handler) HandleMessage(msg twitch.Message) {

	if h.comboCount > 3 && (len(msg.Emotes) != 1 || h.lastEmote != *msg.Emotes[0]) {
		h.bot.Say(msg.Channel, fmt.Sprintf("/me %dx %s COMBO", h.comboCount, h.lastEmote.Name), modules.COMBO)
		h.comboCount = 1
	}


	// filter out messages without emotes or with different emotes
	if len(msg.Emotes) != 1 {
		return
	}

	if h.lastEmote.ID == "" {
		h.lastEmote = *msg.Emotes[0]
		h.comboCount = 1
		return
	}

	if h.lastEmote == *msg.Emotes[0] {
		h.comboCount++
	}

	if h.lastEmote != *msg.Emotes[0] {
		h.lastEmote = *msg.Emotes[0]
		h.comboCount = 1
	}
}
