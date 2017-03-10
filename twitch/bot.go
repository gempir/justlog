package twitch

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"gopkg.in/redis.v3"
)

type Bot struct {
	Messages chan Message
	ircAddress string
	ircUser string
	ircToken string
	rClient redis.Client
}

var (
	mainConn   *net.Conn
	userReg        = regexp.MustCompile(`:\w+!\w+@\w+\.tmi\.twitch\.tv`)
	channelReg  = regexp.MustCompile(`#\w+\s:`)
	actionReg    = regexp.MustCompile(`^\x{0001}ACTION\s`)
	actionReg2  = regexp.MustCompile(`([\x{0001}]+)`)
)

func NewBot(ircAddress string, ircUser string, ircToken string, rClient redis.Client) Bot {
	return Bot{
		Messages: make(chan Message),
		ircAddress: ircAddress,
		ircUser: strings.ToLower(ircUser),
		ircToken: ircToken,
		rClient: rClient,
	}
}

func (bot *Bot) Say(text string, channel string) {
	fmt.Fprintf(*mainConn, "PRIVMSG %s :%s\r\n", channel, text)
}

func (bot *Bot) CreateConnection() error {
	conn, err := net.Dial("tcp", bot.ircAddress)
	mainConn = &conn
	if err != nil {
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
	for _, element := range val {
		if element == "1" || element == "0" {
			continue
		}
		go bot.join(element)
	}
}

func (bot *Bot) join(channel string) {
	fmt.Fprintf(*mainConn, "JOIN %s\r\n", channel)
}