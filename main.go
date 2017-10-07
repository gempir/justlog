package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"strings"

	"github.com/gempir/gempbotgo/api"
	"github.com/gempir/gempbotgo/filelog"
	"github.com/gempir/gempbotgo/humanize"
	"github.com/gempir/go-twitch-irc"
)

var (
	admin      string
	fileLogger filelog.Logger
)

func main() {
	startTime := time.Now()
	admin = getEnv("ADMIN")

	apiServer := api.NewServer()
	go apiServer.Init()

	twitchClient := twitch.NewClient(getEnv("IRCUSER"), getEnv("IRCTOKEN"))

	fileLogger = filelog.NewFileLogger()

	channels := strings.Split(getEnv("CHANNELS"), ",")
	for _, channel := range channels {
		fmt.Println("Joining " + channel)
		twitchClient.Join(channel)
	}

	twitchClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {

		if message.Type == twitch.PRIVMSG || message.Type == twitch.CLEARCHAT {
			go func() {
				err := fileLogger.LogMessageForUser(channel, user, message)
				if err != nil {
					log.Println(err.Error())
				}
			}()

			go func() {
				err := fileLogger.LogMessageForChannel(channel, user, message)
				if err != nil {
					log.Println(err.Error())
				}
			}()

			if strings.HasPrefix(message.Text, "!pingall") {
				uptime := humanize.TimeSince(startTime)
				twitchClient.Say(channel, "uptime: "+uptime)
			}

			if user.Username == admin && strings.HasPrefix(message.Text, "!status") {
				uptime := humanize.TimeSince(startTime)
				twitchClient.Say(channel, admin+", uptime: "+uptime)
			}
		}
	})

	twitchClient.Connect()
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("failed to read env: " + key)
}
