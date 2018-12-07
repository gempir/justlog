package bot

import (
	"strings"
	"time"

	"github.com/gempir/justlog/filelog"

	"github.com/gempir/go-twitch-irc"
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

	twitchClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {

		go func() {
			err := b.fileLogger.LogMessageForUser(channel, user, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			err := b.fileLogger.LogMessageForChannel(channel, user, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		if user.Username == b.admin && strings.HasPrefix(message.Text, "!status") {
			uptime := humanize.TimeSince(*b.startTime)
			twitchClient.Say(channel, user.DisplayName+", uptime: "+uptime)
		}
	})

	twitchClient.OnNewClearchatMessage(func(channel string, user twitch.User, message twitch.Message) {

		go func() {
			err := b.fileLogger.LogMessageForUser(channel, user, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			err := b.fileLogger.LogMessageForChannel(channel, user, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()
	})

	log.Fatal(twitchClient.Connect())
}
