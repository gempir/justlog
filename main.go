package main

import (
	"encoding/json"
	"io/ioutil"
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
	cfg sysConfig
)

type sysConfig struct {
	IrcUser  string   `json:"irc_user"`
	IrcToken string   `json:"irc_token"`
	Admin    string   `json:"admin"`
	Channels []string `json:"channels"`
}

var (
	fileLogger filelog.Logger
)

func main() {
	startTime := time.Now()
	var err error
	cfg, err = readConfig("configs/config.json")
	if err != nil {
		log.Fatal("failed to read config.json")
	}

	apiServer := api.NewServer()
	go apiServer.Init()

	twitchClient := twitch.NewClient(cfg.IrcUser, cfg.IrcToken)
	twitchClient.SetIrcAddress(getEnv("IRCHOST", "irc.chat.twitch.tv:6667"))

	fileLogger = filelog.NewFileLogger()

	for _, channel := range cfg.Channels {
		log.Println("Joining " + channel)
		go twitchClient.Join(strings.TrimPrefix(channel, "#"))
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

			if user.Username == cfg.Admin && strings.HasPrefix(message.Text, "!status") {
				uptime := humanize.TimeSince(startTime)
				twitchClient.Say(channel, cfg.Admin+", uptime: "+uptime)
			}
		}
	})

	twitchClient.Connect()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func readConfig(path string) (sysConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	return unmarshalConfig(file)
}

func unmarshalConfig(file []byte) (sysConfig, error) {
	err := json.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
