package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/gempir/justlog/config"
	"github.com/gempir/justlog/filelog"

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/justlog/helix"
	"github.com/gempir/justlog/humanize"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	startTime         time.Time
	cfg               *config.Config
	helixClient       *helix.Client
	twitchClient      *twitch.Client
	fileLogger        *filelog.Logger
	messageTypesToLog map[string][]twitch.MessageType
}

// NewBot create new bot instance
func NewBot(cfg *config.Config, helixClient *helix.Client, fileLogger *filelog.Logger, messageTypesToLog map[string][]twitch.MessageType) *Bot {
	return &Bot{
		cfg:               cfg,
		helixClient:       helixClient,
		fileLogger:        fileLogger,
		messageTypesToLog: messageTypesToLog,
	}
}

// Connect startup the logger and bot
func (b *Bot) Connect(channelIds []string) {
	b.startTime = time.Now()
	b.twitchClient = twitch.NewClient(b.cfg.Username, "oauth:"+b.cfg.OAuth)

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("Bot joining anonymous")
	} else {
		log.Info("Bot joining as user " + b.cfg.Username)
	}

	channels, err := b.helixClient.GetUsersByUserIds(channelIds)
	if err != nil {
		log.Fatalf("Failed to load configured channels %s", err.Error())
	}

	messageTypesToLog := make(map[string][]twitch.MessageType)

	for _, channel := range channels {
		log.Info("Joining " + channel.Login)
		b.twitchClient.Join(channel.Login)

		if _, ok := b.messageTypesToLog[channel.ID]; ok {
			messageTypesToLog[channel.Login] = b.messageTypesToLog[channel.ID]
		} else {
			messageTypesToLog[channel.Login] = []twitch.MessageType{twitch.PRIVMSG, twitch.CLEARCHAT, twitch.USERNOTICE}
		}
	}

	b.twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {

		go func() {
			if !shouldLog(messageTypesToLog, message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogPrivateMessageForUser(message.User, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			if !shouldLog(messageTypesToLog, message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogPrivateMessageForChannel(message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		b.handlePrivateMessage(message)
	})

	b.twitchClient.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		log.Debug(message.Raw)

		go func() {
			if !shouldLog(messageTypesToLog, message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogUserNoticeMessageForUser(message.User.ID, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		if _, ok := message.Tags["msg-param-recipient-id"]; ok {
			go func() {
				if !shouldLog(messageTypesToLog, message.Channel, message.GetType()) {
					return
				}

				err := b.fileLogger.LogUserNoticeMessageForUser(message.Tags["msg-param-recipient-id"], message)
				if err != nil {
					log.Error(err.Error())
				}
			}()
		}

		go func() {
			if !shouldLog(messageTypesToLog, message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogUserNoticeMessageForChannel(message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

	})

	b.twitchClient.OnClearChatMessage(func(message twitch.ClearChatMessage) {

		go func() {
			if !shouldLog(messageTypesToLog, message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogClearchatMessageForUser(message.TargetUserID, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			if !shouldLog(messageTypesToLog, message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogClearchatMessageForChannel(message)
			if err != nil {
				log.Error(err.Error())
			}
		}()
	})

	log.Fatal(b.twitchClient.Connect())
}

func shouldLog(messageTypesToLog map[string][]twitch.MessageType, channelName string, receivedMsgType twitch.MessageType) bool {
	for _, msgType := range messageTypesToLog[channelName] {
		if msgType == receivedMsgType {
			return true
		}
	}

	return false
}

func (b *Bot) handlePrivateMessage(message twitch.PrivateMessage) {
	if message.User.Name == b.cfg.Admin {
		if strings.HasPrefix(message.Message, "!status") {
			uptime := humanize.TimeSince(b.startTime)
			b.twitchClient.Say(message.Channel, message.User.DisplayName+", uptime: "+uptime)
		}
		if strings.HasPrefix(message.Message, "!justlog join ") {
			input := strings.TrimPrefix(message.Message, "!justlog join ")

			users, err := b.helixClient.GetUsersByUsernames(strings.Split(input, ","))
			if err != nil {
				log.Error(err)
				b.twitchClient.Say(message.Channel, message.User.DisplayName+", something went wrong requesting the userids")
			}

			ids := []string{}
			for _, user := range users {
				ids = append(ids, user.ID)
			}
			b.cfg.AddChannels(ids...)
			b.twitchClient.Say(message.Channel, fmt.Sprintf("%s, added channels: %v", message.User.DisplayName, ids))
		}
	}
}
