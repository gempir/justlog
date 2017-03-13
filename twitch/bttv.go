package twitch

import (
	"github.com/pajlada/gobttv"
	"strings"
)

var (
	bttv = gobttv.New()
	channel Channel
)

func (bot *Bot) InitBttvEmoteCache() {
	go bot.cacheGlobalBttvEmotes()
	go bot.cacheChannelBttvEmotes()
}

func (bot *Bot) cacheGlobalBttvEmotes() {
	bttv.GetEmotes(bot.onSuccessGlobal, onHTTPError, onInternalError)
}

func (bot *Bot) cacheChannelBttvEmotes() {
	val, _ := bot.rClient.HGetAll("channels").Result()
	for channelStr := range val {
		channelShort := strings.TrimPrefix(channelStr, "#")
		channel = NewChannel(channelShort)
		bttv.GetChannel(channelShort, bot.onSuccessChannel, onHTTPError, onInternalError)
	}
}

func (bot *Bot) addBttvEmotes(msg Message) *Message {

	msgSplit := strings.Split(msg.Text, " ")

	for _, word := range msgSplit {
		if val, ok := bot.channelBttvEmotes[msg.Channel][word]; ok {
			msg.Emotes = append(msg.Emotes, &val)
		}
		if val, ok := bot.globalBttvEmotes[word]; ok {
			msg.Emotes = append(msg.Emotes, &val)
		}
	}

	return &msg
}

func (bot *Bot) onSuccessChannel(emotes gobttv.ChannelResponse) {

	bot.channelBttvEmotes[channel] = make(map[string]Emote)

	for _, bttvEmote := range emotes.Emotes {

		emote := Emote{
			Name: bttvEmote.Code,
			ID: bttvEmote.ID,
			Type: BTTVCHANNELEMOTE,
		}

		bot.channelBttvEmotes[channel][bttvEmote.Code] = emote
	}
	bot.logger.Infof("Loaded BTTV Channel Emotes for %s", channel.Name)
}

func (bot *Bot) onSuccessGlobal(emotes gobttv.EmotesResponse) {
	for _, bttvEmote := range emotes.Emotes {

		emote := Emote{
			Name: bttvEmote.Regex,
			ID: bttvEmote.URL,
			Type: BTTVEMOTE,
		}

		bot.globalBttvEmotes[bttvEmote.Regex]= emote
	}
}

func onHTTPError(statusCode int, statusMessage, errorMessage string) {
	logger.Errorf("statusCode: %d, statusMessage: %s, errorMessage: %s", statusCode, statusMessage, errorMessage)
}

func onInternalError(err error) {
	logger.Errorf("internalError: %s", err.Error())
}