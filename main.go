package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/gempir/justlog/api"
	"github.com/gempir/justlog/filelog"
	"github.com/gempir/justlog/humanize"

	log "github.com/sirupsen/logrus"
)

type config struct {
	LogsDirectory string   `json:"logsDirectory"`
	Admin         string   `json:"admin"`
	Channels      []string `json:"channels"`
}

var (
	cfg config
)

func main() {
	startTime := time.Now()

	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()
	cfg = loadConfiguration(*configFile)

	apiServer := api.NewServer(cfg.LogsDirectory)
	go apiServer.Init()

	twitchClient := twitch.NewClient("justinfan123123", "oauth:123123123")
	fileLogger := filelog.NewFileLogger(cfg.LogsDirectory)

	for _, channel := range cfg.Channels {
		fmt.Println("Joining " + channel)
		twitchClient.Join(channel)
		apiServer.AddChannel(channel)
	}

	twitchClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {

		go func() {
			err := fileLogger.LogMessageForUser(channel, user, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			err := fileLogger.LogMessageForChannel(channel, user, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		if user.Username == cfg.Admin && strings.HasPrefix(message.Text, "!status") {
			uptime := humanize.TimeSince(startTime)
			twitchClient.Say(channel, cfg.Admin+", uptime: "+uptime)
		}
	})

	twitchClient.OnNewClearchatMessage(func(channel string, user twitch.User, message twitch.Message) {

		go func() {
			err := fileLogger.LogMessageForUser(channel, user, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()

		go func() {
			err := fileLogger.LogMessageForChannel(channel, user, message)
			if err != nil {
				log.Error(err.Error())
			}
		}()
	})

	twitchClient.Connect()
}

func loadConfiguration(file string) config {
	log.Info("Loading config from " + file)
	var cfg config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)

	cfg.LogsDirectory = strings.TrimSuffix(cfg.LogsDirectory, "/")

	return cfg
}
