package twitch

import (
	"bufio"
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/redis.v5"
	"net"
	"net/textproto"
	"strings"
	"github.com/gempir/gempbotgo/config"
	"io/ioutil"
	"net/http"
)

var (
	logger logging.Logger
)

type Bot struct {
	Messages   chan Message
	ircAddress string
	ircUser    string
	ircToken   string
	userConfig config.UserConfig
	rClient    redis.Client
	logger     logging.Logger
	connection net.Conn
	globalBttvEmotes map[string]Emote
	channelBttvEmotes map[Channel]map[string]Emote
}

func NewBot(ircAddress string, ircUser string, ircToken string, uCfg config.UserConfig,rClient redis.Client, loggerMain logging.Logger) Bot {
	channels := make(map[Channel]bool)
	channels[NewChannel(ircUser)] = true

	logger = logger

	return Bot{
		Messages:   make(chan Message),
		ircAddress: ircAddress,
		ircUser:    strings.ToLower(ircUser),
		ircToken:   ircToken,
		userConfig: uCfg,
		rClient:    rClient,
		logger:     loggerMain,
		globalBttvEmotes: make(map[string]Emote),
		channelBttvEmotes: make(map[Channel]map[string]Emote),
	}
}

func (bot *Bot) Say(channel Channel, text string, responseType string) {
	if bot.userConfig.IsEnabled(channel.Name, responseType) {
		bot.send(fmt.Sprintf("PRIVMSG %s :%s", channel.Name, text))
	}
}

func (bot *Bot) CreatePersistentConnection() error {
	for {
		conn, err := net.Dial("tcp", bot.ircAddress)
		bot.connection = conn
		if err != nil {
			bot.logger.Error(err.Error())
			return err
		}

		bot.setupConnection()
		bot.joinDefault()

		err = bot.readConnection(conn)
		if err != nil {
			bot.logger.Error("connection read error, redialing")
			continue
		}
	}
	return nil
}

func (bot *Bot) readConnection(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			bot.logger.Error(err.Error())
			return err
		}
		messages := strings.Split(line, "\r\n")
		if len(messages) == 0 {
			continue
		}
		for _, msg := range messages {
			bot.handleLine(msg)
		}
	}
}

func (bot *Bot) setupConnection() {
	bot.send(fmt.Sprintf("PASS %s", bot.ircToken))
	bot.send(fmt.Sprintf("NICK %s", bot.ircUser))
	bot.send("CAP REQ :twitch.tv/tags")
	bot.send("CAP REQ :twitch.tv/commands")
	bot.send(fmt.Sprintf("JOIN %s", "#"+bot.ircUser))
}

func (bot *Bot) send(line string) {
	fmt.Fprint(bot.connection, line+"\r\n")
}

func (bot *Bot) handleLine(line string) {
	if strings.HasPrefix(line, "PING") {
		bot.send(fmt.Sprintf(strings.Replace(line, "PING", "PONG", 1)))
	}
	if strings.HasPrefix(line, "@") {
		bot.Messages <- *bot.parseMessage(line)
	}
}

func (bot *Bot) joinDefault() {
	val, _ := bot.rClient.HGetAll("channels").Result()
	for channelStr := range val {
		channel := NewChannel(channelStr)
		go bot.join(channel)
	}
}

func (bot *Bot) join(channel Channel) {
	bot.send(fmt.Sprintf("JOIN %s", channel.Name))
}

func (bot *Bot) httpRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		bot.logger.Error(err.Error())
		return nil, err
	}
	return contents, nil
}