package combo

import (
	"github.com/gempir/gempbotgo/twitch"
	"fmt"
	"github.com/gempir/gempbotgo/modules"
	"sync"
)

type Handler struct {
	bot       *twitch.Bot
	lastEmote map[twitch.Channel]twitch.Emote
	comboCount map[twitch.Channel]int
	mutex sync.Mutex
}

func NewHandler(bot *twitch.Bot) Handler {
	return Handler{
		bot: bot,
		lastEmote: make(map[twitch.Channel]twitch.Emote),
		comboCount:make(map[twitch.Channel]int),
		mutex: sync.Mutex{},
	}
}

func (h *Handler) HandleMessage(msg twitch.Message) {

	h.mutex.Lock()
	if h.comboCount[msg.Channel] > 3 && (len(msg.Emotes) != 1 || h.lastEmote[msg.Channel] != *msg.Emotes[0]) {
		h.bot.Say(msg.Channel, fmt.Sprintf("/me %dx %s COMBO", h.comboCount[msg.Channel], h.lastEmote[msg.Channel].Name), modules.COMBO)
		h.comboCount[msg.Channel] = 1
	}


	// filter out messages without emotes or with different emotes
	if len(msg.Emotes) != 1 {
		return
	}

	if h.lastEmote[msg.Channel].ID == "" {
		h.lastEmote[msg.Channel] = *msg.Emotes[0]
		h.comboCount[msg.Channel] = 1
		return
	}

	if h.lastEmote[msg.Channel] == *msg.Emotes[0] {
		h.comboCount[msg.Channel]++
	}

	if h.lastEmote[msg.Channel] != *msg.Emotes[0] {
		h.lastEmote[msg.Channel] = *msg.Emotes[0]
		h.comboCount[msg.Channel] = 1
	}
	h.mutex.Unlock()
}
