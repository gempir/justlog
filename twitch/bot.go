package twitch

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"gopkg.in/redis.v5"
	"github.com/op/go-logging"
)

type Bot struct {
	Messages chan Message
	ircAddress string
	ircUser string
	ircToken string
	rClient redis.Client
	logger logging.Logger
	channels map[Channel]bool
	connection net.Conn
}

func NewBot(ircAddress string, ircUser string, ircToken string, rClient redis.Client, logger logging.Logger) Bot {
	channels := make(map[Channel]bool)
	channels[NewChannel(ircUser)] = true

	return Bot{
		Messages: make(chan Message),
		ircAddress: ircAddress,
		ircUser: strings.ToLower(ircUser),
		ircToken: ircToken,
		rClient: rClient,
		logger: logger,
		channels: channels,
	}
}

func (bot *Bot) Say(channel Channel, text string, adminCommand bool) {
	if !bot.channels[channel] && !adminCommand {
		return
	}
	bot.send(fmt.Sprintf("Message %s :%s", channel.Name, text))
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
	bot.send(fmt.Sprintf("JOIN %s", "#" + bot.ircUser))
}

func (bot *Bot) send(line string) {
	bot.logger.Debug("SEND | " + line)
	fmt.Fprint(bot.connection, line + "\r\n")
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
	val,_ := bot.rClient.HGetAll("channels").Result()
	for channelStr, activeNum := range val {
		channel := NewChannel(channelStr)
		active := false
		if activeNum == "1" {
			active = true
		}
		bot.channels[channel] = active
		go bot.join(channel)
	}
}

func (bot *Bot) join(channel Channel) {
	bot.send(fmt.Sprintf("JOIN %s", channel.Name))
}