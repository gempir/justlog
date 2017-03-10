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
			go bot.parseMessage(msg)
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

func (bot *Bot) parseMessage(msg string) {

	if !strings.Contains(msg, ".tmi.twitch.tv PRIVMSG ") {
		return
	}

	fullUser := userReg.FindString(msg)
	userIrc := strings.Split(fullUser, "!")
	username := userIrc[0][1:len(userIrc[0])]
	split2 := strings.Split(msg, ".tmi.twitch.tv PRIVMSG ")
	split3 := channelReg.FindString(split2[1])
	channel := split3[0 : len(split3)-2]
	split4 := strings.Split(split2[1], split3)
	message := split4[1]
	message = actionReg.ReplaceAllLiteralString(message, "")
	message = actionReg2.ReplaceAllLiteralString(message, "")

	split5 := strings.Replace(split2[0], " :" + username + "!" + username + "@" + username, "", -1)
	tags := strings.Split(strings.Replace(split5, "@", "", 1), ";")

	tagMap := make(map[string]string)

	for _,tag := range tags {
		tagSplit := strings.Split(tag, "=")
		tagMap[tagSplit[0]] = tagSplit[1]
	}

	subscriber, turbo, mod := false, false, false

	if tagMap["subscriber"] == "1" {
		subscriber = true
	}
	if tagMap["turbo"] == "1" {
		turbo = true
	}
	if tagMap["mod"] == "1" {
		mod = true
	}

	user := newUser(username, tagMap["User-id"], tagMap["color"], tagMap["display-name"], mod, subscriber, turbo)
	bot.Messages <- newMessage(message, user, channel)
}

func (bot *Bot) join(channel string) {
	fmt.Fprintf(*mainConn, "JOIN %s\r\n", channel)
}