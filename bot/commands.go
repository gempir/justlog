package bot

import (
	"fmt"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/gempir/justlog/humanize"
	log "github.com/sirupsen/logrus"
)

const (
	commandPrefix        = "!justlog"
	errNoUsernames       = ", at least 1 username has to be provided. multiple usernames have to be separated with a space"
	errRequestingUserIDs = ", something went wrong requesting the userids"
)

func (b *Bot) handlePrivateMessageCommands(message twitch.PrivateMessage) {
	if !strings.HasPrefix(strings.ToLower(message.Message), commandPrefix) {
		return
	}

	args := strings.Fields(message.Message[len(commandPrefix):])
	if len(args) < 1 {
		return
	}
	commandName := args[0]
	args = args[1:]

	switch commandName {
	case "status":
		if !contains(b.cfg.Admins, message.User.Name) {
			return
		}
		uptime := humanize.TimeSince(b.startTime)
		b.Say(message.Channel, fmt.Sprintf("%s, uptime: %s", message.User.DisplayName, uptime))

	case "join":
		if !contains(b.cfg.Admins, message.User.Name) {
			return
		}
		b.handleJoin(message, args)

	case "part":
		if !contains(b.cfg.Admins, message.User.Name) {
			return
		}
		b.handlePart(message, args)

	case "optout":
		b.handleOptOut(message, args)
	case "optin":
		if !contains(b.cfg.Admins, message.User.Name) {
			return
		}
		b.handleOptIn(message, args)
	}
}

// Commands

func (b *Bot) handleJoin(message twitch.PrivateMessage, args []string) {
	if len(args) < 1 {
		b.Say(message.Channel, message.User.DisplayName+errNoUsernames)
		return
	}

	users, err := b.helixClient.GetUsersByUsernames(args)
	if err != nil {
		log.Error(err)
		b.Say(message.Channel, message.User.DisplayName+errRequestingUserIDs)
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
		b.Join(user.Login)
	}
	b.cfg.AddChannels(ids...)
	b.Say(message.Channel, fmt.Sprintf("%s, added channels: %v", message.User.DisplayName, ids))
}

func (b *Bot) handlePart(message twitch.PrivateMessage, args []string) {
	if len(args) < 1 {
		b.Say(message.Channel, message.User.DisplayName+errNoUsernames)
		return
	}

	users, err := b.helixClient.GetUsersByUsernames(args)
	if err != nil {
		log.Error(err)
		b.Say(message.Channel, message.User.DisplayName+errRequestingUserIDs)
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
		b.Part(user.Login)
	}
	b.cfg.RemoveChannels(ids...)
	b.Say(message.Channel, fmt.Sprintf("%s, removed channels: %v", message.User.DisplayName, ids))
}

func (b *Bot) handleOptOut(message twitch.PrivateMessage, args []string) {
	if len(args) < 1 {
		b.Say(message.Channel, message.User.DisplayName+errNoUsernames)
		return
	}

	if _, ok := b.OptoutCodes.LoadAndDelete(args[0]); ok {
		b.cfg.OptOutUsers(message.User.ID)
		b.Say(message.Channel, fmt.Sprintf("%s, opted you out", message.User.DisplayName))
		return
	}

	if !contains(b.cfg.Admins, message.User.Name) {
		return
	}

	users, err := b.helixClient.GetUsersByUsernames(args)
	if err != nil {
		log.Error(err)
		b.Say(message.Channel, message.User.DisplayName+errRequestingUserIDs)
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
	}
	b.cfg.OptOutUsers(ids...)
	b.Say(message.Channel, fmt.Sprintf("%s, opted out channels: %v", message.User.DisplayName, ids))
}

func (b *Bot) handleOptIn(message twitch.PrivateMessage, args []string) {
	if len(args) < 1 {
		b.Say(message.Channel, message.User.DisplayName+errNoUsernames)
		return
	}

	users, err := b.helixClient.GetUsersByUsernames(args)
	if err != nil {
		log.Error(err)
		b.Say(message.Channel, message.User.DisplayName+errRequestingUserIDs)
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
	}

	b.cfg.RemoveOptOut(ids...)
	b.Say(message.Channel, fmt.Sprintf("%s, opted in channels: %v", message.User.DisplayName, ids))
}

// Utilities

func contains(arr []string, str string) bool {
	for _, x := range arr {
		if x == str {
			return true
		}
	}
	return false
}
