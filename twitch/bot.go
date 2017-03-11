package twitch

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"regexp"
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

var (
	userReg        = regexp.MustCompile(`:\w+!\w+@\w+\.tmi\.twitch\.tv`)
	channelReg  = regexp.MustCompile(`#\w+\s:`)
	actionReg    = regexp.MustCompile(`^\x{0001}ACTION\s`)
	actionReg2  = regexp.MustCompile(`([\x{0001}]+)`)
)

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
	bot.send(fmt.Sprintf("PRIVMSG %s :%s\r\n", channel.Name, text))
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
	bot.send(fmt.Sprintf("PASS %s\r\n", bot.ircToken))
	bot.send(fmt.Sprintf("NICK %s\r\n", bot.ircUser))
	bot.send("CAP REQ :twitch.tv/tags\r\n")
	bot.send("CAP REQ :twitch.tv/commands\r\n")
	bot.send(fmt.Sprintf("JOIN %s\r\n", "#" + bot.ircUser))
}

func (bot *Bot) send(line string) {
	bot.logger.Debug("SEND | " + line)
	fmt.Fprint(bot.connection, line)
}

func (bot *Bot) handleLine(line string) {
	if strings.HasPrefix(line, "PING") {
		bot.send(fmt.Sprintf(strings.Replace(line, "PING", "PONG", 1)))
	}
	if strings.Contains(line, ".tmi.twitch.tv PRIVMSG ") {
		bot.Messages <- parseMessage(line)
	} else if !strings.Contains(line, "tmi.twitch.tv USERSTATE ") {
		bot.logger.Debug(line)
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
	bot.send(fmt.Sprintf("JOIN %s\r\n", channel.Name))
}