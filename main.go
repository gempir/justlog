package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"strings"

	"github.com/gempir/gempbotgo/api"
	"github.com/gempir/gempbotgo/filelog"
	"github.com/gempir/gempbotgo/humanize"
	"github.com/gempir/gempbotgo/store"
	"github.com/gempir/go-twitch-irc"

	_ "github.com/go-sql-driver/mysql"
)

var (
	admin      string
	fileLogger filelog.Logger
	db         *sql.DB
)

func main() {
	startTime := time.Now()
	admin = os.Getenv("ADMIN")

	apiServer := api.NewServer()
	go apiServer.Init()

	time.Sleep(time.Second * 3)
	store, err := store.NewClient(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	twitchClient := twitch.NewClient(os.Getenv("IRCUSER"), os.Getenv("IRCTOKEN"))

	fileLogger = filelog.NewFileLogger()

	channels := store.GetAllChannels()
	for _, channel := range channels {
		fmt.Println("Joining " + channel)
		twitchClient.Join(channel)
	}

	twitchClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {

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
	})

	twitchClient.Connect()
}
