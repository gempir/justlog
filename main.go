package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/op/go-logging"
	"gopkg.in/redis.v5"

	"fmt"
	"github.com/gempir/gempbotgo/api"
	"github.com/gempir/gempbotgo/filelog"
	"github.com/gempir/gempbotgo/humanize"
	"github.com/gempir/go-twitch-irc"
	"strings"
)

var (
	cfg    sysConfig
	logger logging.Logger
)

type sysConfig struct {
	IrcAddress    string `json:"irc_address"`
	IrcUser       string `json:"irc_user"`
	IrcToken      string `json:"irc_token"`
	Admin         string `json:"admin"`
	LogPath       string `json:"log_path"`
	APIPort       string `json:"api_port"`
	RedisAddress  string `json:"redis_address"`
	RedisPassword string `json:"redis_password"`
	RedisDatabase int    `json:"redis_database"`
	CleverBotUser string `json:"cleverbot_user"`
	CleverBotKey  string `json:"cleverbot_key"`
}

var (
	fileLogger filelog.Logger
)

func main() {
	startTime := time.Now()
	logger = initLogger()
	var err error
	cfg, err = readConfig("config.json")
	if err != nil {
		logger.Fatal(err)
	}

	apiServer := api.NewServer(cfg.APIPort, cfg.LogPath)
	go apiServer.Init()

	rClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})

	twitchClient := twitch.NewClient(cfg.IrcUser, cfg.IrcToken)
	twitchClient.SetIrcAddress(cfg.IrcAddress)

	fileLogger = filelog.NewFileLogger(cfg.LogPath)

	val, _ := rClient.HGetAll("channels").Result()
	for channelStr := range val {
		fmt.Println("Joining " + channelStr)
		go twitchClient.Join(strings.TrimPrefix(channelStr, "#"))
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

			if user.Username == cfg.Admin && strings.HasPrefix(message.Text, "!status") {
				uptime := humanize.TimeSince(startTime)
				twitchClient.Say(channel, cfg.Admin+", uptime: "+uptime)
			}
		}
	})

	fmt.Println(twitchClient.Connect())
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
