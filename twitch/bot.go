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
}

var (
	mainConn   *net.Conn
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
	fmt.Fprintf(*mainConn, "PRIVMSG %s :%s\r\n", channel.Name, text)
}

func (bot *Bot) CreateConnection() error {
	conn, err := net.Dial("tcp", bot.ircAddress)
	mainConn = &conn
	if err != nil {
		bot.logger.Error(err.Error())
		return err
	}
	fmt.Fprintf(*mainConn, "PASS %s\r\n", bot.ircToken)
	fmt.Fprintf(*mainConn, "NICK %s\r\n", bot.ircUser)
	fmt.Fprint(*mainConn, "CAP REQ :twitch.tv/tags\r\n")
	fmt.Fprint(*mainConn, "CAP REQ :twitch.tv/commands\r\n")
	fmt.Fprintf(*mainConn, "JOIN %s\r\n", "#" + bot.ircUser)

	go bot.joinDefault()

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
			if strings.Contains(msg, ".tmi.twitch.tv PRIVMSG ") {
				bot.Messages <- parseMessage(msg)
			}
		}
	}
	defer bot.CreateConnection()
	return nil
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
	fmt.Fprintf(*mainConn, "JOIN %s\r\n", channel.Name)
}