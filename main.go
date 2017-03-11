package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/op/go-logging"
	"gopkg.in/redis.v5"

	"github.com/gempir/gempbotgo/twitch"
	"github.com/gempir/gempbotgo/command"
	"github.com/gempir/gempbotgo/filelog"
	"github.com/gempir/gempbotgo/api"
	"github.com/gempir/gempbotgo/combo"
)

var (
	cfg    config
	logger logging.Logger
)

type config struct {
	IrcAddress       string `json:"irc_address"`
	IrcUser          string `json:"irc_user"`
	IrcToken         string `json:"irc_token"`
	Admin			 string `json:"admin"`
	LogPath			 string `json:"log_path"`
	APIPort          string `json:"api_port"`
	RedisAddress     string `json:"redis_address"`
	RedisPassword    string `json:"redis_password"`
	RedisDatabase    int    `json:"redis_database"`
}

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

	bot := twitch.NewBot(cfg.IrcAddress, cfg.IrcUser, cfg.IrcToken, *rClient, logger)
	go bot.CreateConnection()

	fileLogger := filelog.NewFileLogger(cfg.LogPath)
	cmdHandler := command.NewHandler(cfg.Admin, bot, startTime, logger)
	comboHandler := combo.NewHandler()

	for msg := range bot.Messages {

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

func initLogger() logging.Logger {
	var logger *logging.Logger
	logger = logging.MustGetLogger("gempbotgo")
	backend := logging.NewLogBackend(os.Stdout, "", 0)

	format := logging.MustStringFormatter(
		`%{color}%{shortfile:-15s} %{level:.4s}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
	backendLeveled := logging.AddModuleLevel(backend)
	logging.SetBackend(backendLeveled)
	return *logger
}

func readConfig(path string) (config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	return unmarshalConfig(file)
}

func unmarshalConfig(file []byte) (config, error) {
	err := json.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}