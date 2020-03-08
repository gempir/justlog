package bot

import (
	"strings"
	"time"

	"github.com/gempir/justlog/config"
	"github.com/gempir/justlog/filelog"

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/justlog/helix"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	startTime         time.Time
	cfg               *config.Config
	helixClient       *helix.Client
	twitchClient      *twitch.Client
	fileLogger        *filelog.Logger
	channels          map[string]helix.UserData
	messageTypesToLog map[string][]twitch.MessageType
}

// NewBot create new bot instance
func NewBot(cfg *config.Config, helixClient *helix.Client, fileLogger *filelog.Logger) *Bot {
	channels, err := helixClient.GetUsersByUserIds(cfg.Channels)
	if err != nil {
		log.Fatalf("[bot] failed to load configured channels %s", err.Error())
	}

	return &Bot{
		cfg:         cfg,
		helixClient: helixClient,
		fileLogger:  fileLogger,
		channels:    channels,
	}
}

// Connect startup the logger and bot
func (b *Bot) Connect() {
	b.startTime = time.Now()
	b.twitchClient = twitch.NewClient(b.cfg.Username, "oauth:"+b.cfg.OAuth)
	b.updateMessageTypesToLog()
	b.initialJoins()

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("[bot] joining anonymous")
	} else {
		log.Info("[bot] joining as user " + b.cfg.Username)
	}

	b.twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {

		go func() {
			if !b.shouldLog(message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogPrivateMessageForUser(message.User, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			if !b.shouldLog(message.Channel, message.GetType()) {
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
			if !b.shouldLog(message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogUserNoticeMessageForUser(message.User.ID, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		if _, ok := message.Tags["msg-param-recipient-id"]; ok {
			go func() {
				if !b.shouldLog(message.Channel, message.GetType()) {
					return
				}

				err := b.fileLogger.LogUserNoticeMessageForUser(message.Tags["msg-param-recipient-id"], message)
				if err != nil {
					log.Error(err.Error())
				}
			}()
		}

		go func() {
			if !b.shouldLog(message.Channel, message.GetType()) {
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
			if !b.shouldLog(message.Channel, message.GetType()) {
				return
			}

			err := b.fileLogger.LogClearchatMessageForUser(message.TargetUserID, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			if !b.shouldLog(message.Channel, message.GetType()) {
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

func (b *Bot) shouldLog(channelName string, receivedMsgType twitch.MessageType) bool {
	for _, msgType := range b.messageTypesToLog[channelName] {
		if msgType == receivedMsgType {
			return true
		}
	}

	return false
}

func (b *Bot) updateMessageTypesToLog() {
	messageTypesToLog := make(map[string][]twitch.MessageType)

	for _, channel := range b.channels {
		if _, ok := b.cfg.ChannelConfigs[channel.ID]; ok && b.cfg.ChannelConfigs[channel.ID].MessageTypes != nil {
			messageTypesToLog[channel.Login] = b.cfg.ChannelConfigs[channel.ID].MessageTypes
		} else {
			messageTypesToLog[channel.Login] = []twitch.MessageType{twitch.PRIVMSG, twitch.CLEARCHAT, twitch.USERNOTICE}
		}
	}

	b.messageTypesToLog = messageTypesToLog
}

func (b *Bot) initialJoins() {
	for _, channel := range b.channels {
		log.Info("[bot] joining " + channel.Login)
		b.twitchClient.Join(channel.Login)
	}
}
