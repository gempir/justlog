package bot

import (
	"fmt"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/justlog/humanize"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handlePrivateMessageCommands(message twitch.PrivateMessage) {
	if contains(b.cfg.Admins, message.User.Name) {
		if strings.HasPrefix(message.Message, "!justlog status") || strings.HasPrefix(message.Message, "!status") {
			uptime := humanize.TimeSince(b.startTime)
			b.Say(message.Channel, message.User.DisplayName+", uptime: "+uptime)
		}
		if strings.HasPrefix(strings.ToLower(message.Message), "!justlog join ") {
			b.handleJoin(message)
		}
		if strings.HasPrefix(strings.ToLower(message.Message), "!justlog part ") {
			b.handlePart(message)
		}
		if strings.HasPrefix(strings.ToLower(message.Message), "!justlog optout ") {
			b.handleOptOut(message)
		}
		if strings.HasPrefix(strings.ToLower(message.Message), "!justlog optin ") {
			b.handleOptIn(message)
		}
	}
}

func (b *Bot) handleJoin(message twitch.PrivateMessage) {
	input := strings.TrimPrefix(message.Message, "!justlog join ")

	users, err := b.helixClient.GetUsersByUsernames(strings.Split(input, ","))
	if err != nil {
		log.Error(err)
		b.Say(message.Channel, message.User.DisplayName+", something went wrong requesting the userids")
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
		b.Join(user.Login)
	}
	b.cfg.AddChannels(ids...)
	b.Say(message.Channel, fmt.Sprintf("%s, added channels: %v", message.User.DisplayName, ids))
}

func (b *Bot) handlePart(message twitch.PrivateMessage) {
	input := strings.TrimPrefix(message.Message, "!justlog part ")

	users, err := b.helixClient.GetUsersByUsernames(strings.Split(input, ","))
	if err != nil {
		log.Error(err)
		b.Say(message.Channel, message.User.DisplayName+", something went wrong requesting the userids")
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
		b.Part(user.Login)
	}
	b.cfg.RemoveChannels(ids...)
	b.Say(message.Channel, fmt.Sprintf("%s, removed channels: %v", message.User.DisplayName, ids))
}

func (b *Bot) handleOptOut(message twitch.PrivateMessage) {
	input := strings.TrimPrefix(strings.ToLower(message.Message), "!justlog optout ")

	users, err := b.helixClient.GetUsersByUsernames(strings.Split(input, ","))
	if err != nil {
		log.Error(err)
		b.Say(message.Channel, message.User.DisplayName+", something went wrong requesting the userids")
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
	}
	b.cfg.OptOutUsers(ids...)
	b.Say(message.Channel, fmt.Sprintf("%s, opted out channels: %v", message.User.DisplayName, ids))
}

func (b *Bot) handleOptIn(message twitch.PrivateMessage) {
	input := strings.TrimPrefix(strings.ToLower(message.Message), "!justlog optin ")

	users, err := b.helixClient.GetUsersByUsernames(strings.Split(input, ","))
	if err != nil {
		log.Error(err)
		b.Say(message.Channel, message.User.DisplayName+", something went wrong requesting the userids")
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
	}

	b.cfg.RemoveOptOut(ids...)
	b.Say(message.Channel, fmt.Sprintf("%s, opted in channels: %v", message.User.DisplayName, ids))
}

func contains(arr []string, str string) bool {
	for _, x := range arr {
		if x == str {
			return true
		}
	}
	return false
}
