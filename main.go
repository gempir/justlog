package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/op/go-logging"

	"fmt"
	"strings"

	"github.com/gempir/gempbotgo/api"
	"github.com/gempir/gempbotgo/filelog"
	"github.com/gempir/gempbotgo/humanize"
	"github.com/gempir/go-twitch-irc"
)

var (
	cfg    sysConfig
	logger logging.Logger
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
	logger = initLogger()
	var err error
	cfg, err = readConfig("configs/config.json")
	if err != nil {
		logger.Fatal(err)
	}

	apiServer := api.NewServer()
	go apiServer.Init()

	twitchClient := twitch.NewClient(cfg.IrcUser, cfg.IrcToken)
	twitchClient.SetIrcAddress(getEnv("IRCHOST", "irc.chat.twitch.tv:6667"))

	fileLogger = filelog.NewFileLogger()

	for _, channel := range cfg.Channels {
		fmt.Println("Joining " + channel)
		go twitchClient.Join(strings.TrimPrefix(channel, "#"))
	}

	twitchClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {

		if message.Type == twitch.PRIVMSG || message.Type == twitch.CLEARCHAT {
			go func() {
				err := fileLogger.LogMessageForUser(channel, user, message)
				if err != nil {
					logger.Error(err.Error())
				}
			}()

			go func() {
				err := fileLogger.LogMessageForChannel(channel, user, message)
				if err != nil {
					logger.Error(err.Error())
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

func initLogger() logging.Logger {
	var logger *logging.Logger
	logger = logging.MustGetLogger("gempbotgo")
	backend := logging.NewLogBackend(os.Stdout, "", 0)

	format := logging.MustStringFormatter(
		`%{color}%{level} %{shortfile}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
	backendLeveled := logging.AddModuleLevel(backend)
	logging.SetBackend(backendLeveled)
	return *logger
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
