package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/op/go-logging"
	"github.com/gempir/gempbotgo/twitch"
	"gopkg.in/redis.v3"
	"github.com/gempir/gempbotgo/command"
	"strings"
	"github.com/gempir/gempbotgo/filelog"
	"time"
)

var (
	cfg config
	Log logging.Logger
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
	RedisDatabase    int64  `json:"redis_database"`
}

func main() {
	startTime := time.Now()
	Log = initLogger()
	var err error
	cfg, err = readConfig("config.json")
	if err != nil {
		Log.Fatal(err)
	}

	rClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})

	bot := twitch.NewBot(cfg.IrcAddress, cfg.IrcUser, cfg.IrcToken, Log, *rClient)
	go bot.CreateConnection()

	fileLogger := filelog.NewFileLogger(cfg.LogPath, Log)
	cmdHandler := command.NewHandler(cfg.Admin, startTime, Log)


	for msg := range bot.Messages {

		fileLogger.LogMessage(msg)

		if strings.HasPrefix(msg.Text, "!") {
			cmdHandler.HandleCommand(msg)
		}

	}
}

func initLogger() logging.Logger {
	var logger *logging.Logger
	logger = logging.MustGetLogger("gempbotgo")
	backend := logging.NewLogBackend(os.Stdout, "", 0)

	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{level:.4s}%{color:reset} %{message}`,
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