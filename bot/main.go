package bot

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gempir/justlog/config"
	"github.com/gempir/justlog/filelog"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/gempir/justlog/helix"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	startTime   time.Time
	cfg         *config.Config
	helixClient helix.TwitchApiClient
	logger      filelog.Logger
	worker      []*worker
	channels    map[string]helix.UserData
	clearchats  sync.Map
	OptoutCodes sync.Map
}

type worker struct {
	client         *twitch.Client
	joinedChannels map[string]bool
}

func newWorker(client *twitch.Client) *worker {
	return &worker{
		client:         client,
		joinedChannels: map[string]bool{},
	}
}

// NewBot create new bot instance
func NewBot(cfg *config.Config, helixClient helix.TwitchApiClient, fileLogger filelog.Logger) *Bot {
	channels, err := helixClient.GetUsersByUserIds(cfg.Channels)
	if err != nil {
		log.Fatalf("[bot] failed to load configured channels %s", err.Error())
	}

	return &Bot{
		cfg:         cfg,
		helixClient: helixClient,
		logger:      fileLogger,
		channels:    channels,
		worker:      []*worker{},
		OptoutCodes: sync.Map{},
	}
}

func (b *Bot) Say(channel, text string) {
	randomIndex := rand.Intn(len(b.worker))
	b.worker[randomIndex].client.Say(channel, text)
}

// Connect startup the logger and bot
func (b *Bot) Connect() {
	b.startTime = time.Now()
	client := b.newClient()
	go b.startJoinLoop()

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("[bot] joining as anonymous user")
	} else {
		log.Info("[bot] joining as user " + b.cfg.Username)
	}

	log.Fatal(client.Connect())
}

// constantly join channels to rejoin some channels that got unbanned over time
func (b *Bot) startJoinLoop() {
	for {
		for _, channel := range b.channels {
			b.Join(channel.Login)
		}

		time.Sleep(time.Hour * 1)
		log.Info("[bot] running hourly join loop")
	}
}

func (b *Bot) Part(channelNames ...string) {
	for _, channelName := range channelNames {
		log.Info("[bot] leaving " + channelName)

		for _, worker := range b.worker {
			worker.client.Depart(channelName)
		}
	}
}

func (b *Bot) Join(channelNames ...string) {
	for _, channel := range channelNames {
		channel = strings.ToLower(channel)

		joined := false

		for _, worker := range b.worker {
			if _, ok := worker.joinedChannels[channel]; ok {
				// already joined but join again in case it was a temporary ban
				worker.client.Join(channel)
				joined = true
			} else if len(worker.joinedChannels) < 50 {
				log.Info("[bot] joining " + channel)
				worker.client.Join(channel)
				worker.joinedChannels[channel] = true
				joined = true
				break
			}
		}
		if !joined {
			client := b.newClient()
			go client.Connect()
			b.Join(channel)
		}
	}
}

func (b *Bot) newClient() *twitch.Client {
	client := twitch.NewClient(b.cfg.Username, "oauth:"+b.cfg.OAuth)
	if b.cfg.BotVerified {
		client.SetJoinRateLimiter(twitch.CreateVerifiedRateLimiter())
	}

	b.worker = append(b.worker, newWorker(client))
	log.Infof("[bot] creating new twitch connection, new total: %d", len(b.worker))

	client.OnPrivateMessage(b.handlePrivateMessage)
	client.OnUserNoticeMessage(b.handleUserNotice)
	client.OnClearChatMessage(b.handleClearChat)

	return client
}

func (b *Bot) handlePrivateMessage(message twitch.PrivateMessage) {
	b.handlePrivateMessageCommands(message)

	if b.cfg.IsOptedOut(message.User.ID) || b.cfg.IsOptedOut(message.RoomID) {
		return
	}

	go func() {
		err := b.logger.LogPrivateMessageForUser(message.User, message)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	go func() {
		err := b.logger.LogPrivateMessageForChannel(message)
		if err != nil {
			log.Error(err.Error())
		}
	}()
}

func (b *Bot) handleUserNotice(message twitch.UserNoticeMessage) {
	if b.cfg.IsOptedOut(message.User.ID) || b.cfg.IsOptedOut(message.RoomID) {
		return
	}

	go func() {
		err := b.logger.LogUserNoticeMessageForUser(message.User.ID, message)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	if _, ok := message.Tags["msg-param-recipient-id"]; ok {
		go func() {
			err := b.logger.LogUserNoticeMessageForUser(message.Tags["msg-param-recipient-id"], message)
			if err != nil {
				log.Error(err.Error())
			}
		}()
	}

	go func() {
		err := b.logger.LogUserNoticeMessageForChannel(message)
		if err != nil {
			log.Error(err.Error())
		}
	}()
}

func (b *Bot) handleClearChat(message twitch.ClearChatMessage) {
	if b.cfg.IsOptedOut(message.TargetUserID) || b.cfg.IsOptedOut(message.RoomID) {
		return
	}

	if message.BanDuration == 0 {
		count, ok := b.clearchats.Load(message.RoomID)
		if !ok {
			count = 0
		}
		newCount := count.(int) + 1
		b.clearchats.Store(message.RoomID, newCount)

		go func() {
			time.Sleep(time.Second * 1)
			count, ok := b.clearchats.Load(message.RoomID)
			if ok {
				b.clearchats.Store(message.RoomID, count.(int)-1)
			}
		}()

		if newCount > 50 {
			if newCount == 51 {
				log.Infof("Stopped recording CLEARCHAT permabans in: %s", message.Channel)
			}
			return
		}
	}

	go func() {
		err := b.logger.LogClearchatMessageForUser(message.TargetUserID, message)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	go func() {
		err := b.logger.LogClearchatMessageForChannel(message)
		if err != nil {
			log.Error(err.Error())
		}
	}()
}
