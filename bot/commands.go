package bot

import (
	"fmt"
	"strconv"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/justlog/humanize"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handlePrivateMessage(message twitch.PrivateMessage) {
	if contains(b.cfg.Admins, message.User.Name) {
		if strings.HasPrefix(message.Message, "!justlog status") || strings.HasPrefix(message.Message, "!status") {
			uptime := humanize.TimeSince(b.startTime)
			b.twitchClient.Say(message.Channel, message.User.DisplayName+", uptime: "+uptime)
		}
		if strings.HasPrefix(message.Message, "!justlog join ") {
			b.handleJoin(message)
		}
		if strings.HasPrefix(message.Message, "!justlog messageType ") {
			b.handleMessageType(message)
		}
	}
}

func (b *Bot) handleJoin(message twitch.PrivateMessage) {
	input := strings.TrimPrefix(message.Message, "!justlog join ")

	users, err := b.helixClient.GetUsersByUsernames(strings.Split(input, ","))
	if err != nil {
		log.Error(err)
		b.twitchClient.Say(message.Channel, message.User.DisplayName+", something went wrong requesting the userids")
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
		log.Infof("[bot] joining %s", user.Login)
		b.twitchClient.Join(user.Login)
	}
	b.cfg.AddChannels(ids...)
	b.twitchClient.Say(message.Channel, fmt.Sprintf("%s, added channels: %v", message.User.DisplayName, ids))
}

func (b *Bot) handleMessageType(message twitch.PrivateMessage) {
	input := strings.TrimPrefix(message.Message, "!justlog messageType ")

	parts := strings.Split(input, " ")
	if len(parts) < 2 {
		return
	}

	users, err := b.helixClient.GetUsersByUsernames([]string{parts[0]})
	if err != nil {
		log.Error(err)
		return
	}

	if parts[1] == "reset" {
		b.cfg.ResetMessageTypes(users[parts[0]].ID)
		log.Infof("[bot] setting %s config messageTypes to default", parts[0])
	} else {
		var messageTypes []twitch.MessageType
		for _, msgType := range strings.Split(parts[1], ",") {
			messageType, err := strconv.Atoi(msgType)
			if err != nil {
				log.Error(err)
				return
			}

			messageTypes = append(messageTypes, twitch.MessageType(messageType))
		}

		b.cfg.SetMessageTypes(users[parts[0]].ID, messageTypes)
		b.updateMessageTypesToLog()
		log.Infof("[bot] setting %s config messageTypes to %v", parts[0], messageTypes)
	}
}

func contains(arr []string, str string) bool {
	for _, x := range arr {
		if x == str {
			return true
		}
	}
	return false
}
