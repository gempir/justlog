package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/op/go-logging"
	"gopkg.in/redis.v5"

	"github.com/gempir/gempbotgo/api"
	"github.com/gempir/gempbotgo/combo"
	"github.com/gempir/gempbotgo/command"
	"github.com/gempir/gempbotgo/filelog"
	"github.com/gempir/gempbotgo/twitch"
	"github.com/gempir/gempbotgo/config"
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
}

var (
	fileLogger   filelog.Logger
	cmdHandler   command.Handler
	comboHandler combo.Handler
)

func main() {
	startTime := time.Now()
	logger = initLogger()
	var err error
	cfg, err = readConfig("config.json")
	if err != nil {
		logger.Fatal(err)
	}

	rClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})

	apiServer := api.NewServer(cfg.APIPort, cfg.LogPath)
	go apiServer.Init()

	userConfig := config.NewUserConfig(*rClient)

	bot := twitch.NewBot(cfg.IrcAddress, cfg.IrcUser, cfg.IrcToken, userConfig, *rClient, logger)
	go bot.InitBttvEmoteCache()
	go func() {
		err := bot.CreatePersistentConnection()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	fileLogger = filelog.NewFileLogger(cfg.LogPath)
	cmdHandler = command.NewHandler(cfg.Admin, &bot, startTime, logger)
	comboHandler = combo.NewHandler(&bot)

	for msg := range bot.Messages {

		if msg.Type == twitch.PRIVMSG || msg.Type == twitch.CLEARCHAT {
			go func() {
				err := fileLogger.LogMessageForUser(msg)
				if err != nil {
					logger.Error(err.Error())
				}
			}()

			go func() {
				err := fileLogger.LogMessageForChannel(msg)
				if err != nil {
					logger.Error(err.Error())
				}
			}()

			go comboHandler.HandleMessage(msg)

			if strings.HasPrefix(msg.Text, "!") {
				go func() {
					err := cmdHandler.HandleCommand(msg)
					if err != nil {
						logger.Error(err.Error())
					}
				}()
			}
		}
	}
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
