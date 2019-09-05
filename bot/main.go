package bot

import (
	"strings"
	"time"

	"github.com/gempir/justlog/filelog"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/justlog/helix"
	"github.com/gempir/justlog/humanize"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	admin       string
	username    string
	oauth       string
	startTime   *time.Time
	helixClient *helix.Client
	fileLogger  *filelog.Logger
}

func NewBot(admin, username, oauth string, startTime *time.Time, helixClient *helix.Client, fileLogger *filelog.Logger) *Bot {
	return &Bot{
		admin:       admin,
		username:    username,
		oauth:       oauth,
		startTime:   startTime,
		helixClient: helixClient,
		fileLogger:  fileLogger,
	}
}

func (b *Bot) Connect(channelIds []string) {
	twitchClient := twitch.NewClient(b.username, "oauth:"+b.oauth)

	if strings.HasPrefix(b.username, "justinfan") {
		log.Info("Bot joining anonymous")
	} else {
		log.Info("Bot joining as user " + b.username)
	}

	channels, err := b.helixClient.GetUsersByUserIds(channelIds)
	if err != nil {
		log.Fatalf("Failed to load configured channels %s", err.Error())
	}

	for _, channel := range channels {
		log.Info("Joining " + channel.Login)
		twitchClient.Join(channel.Login)
	}

	twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {

		go func() {
			err := b.fileLogger.LogPrivateMessageForUser(message.User, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			err := b.fileLogger.LogPrivateMessageForChannel(message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		if message.User.Name == b.admin && strings.HasPrefix(message.Message, "!status") {
			uptime := humanize.TimeSince(*b.startTime)
			twitchClient.Say(message.Channel, message.User.DisplayName+", uptime: "+uptime)
		}
	})

	twitchClient.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		log.Debug(message.Raw)

		go func() {
			err := b.fileLogger.LogUserNoticeMessageForUser(message.User.ID, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		if _, ok := message.Tags["msg-param-recipient-id"]; ok {
			go func() {
				err := b.fileLogger.LogUserNoticeMessageForUser(message.Tags["msg-param-recipient-id"], message)
				if err != nil {
					log.Error(err.Error())
				}
			}()
		}

		go func() {
			err := b.fileLogger.LogUserNoticeMessageForChannel(message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

	})

	twitchClient.OnClearChatMessage(func(message twitch.ClearChatMessage) {

		go func() {
			err := b.fileLogger.LogClearchatMessageForUser(message.TargetUserID, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			err := b.fileLogger.LogClearchatMessageForChannel(message)
			if err != nil {
				log.Error(err.Error())
			}
		}()
	})

	log.Fatal(twitchClient.Connect())
}
