package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/gempir/justlog/api"
	"github.com/gempir/justlog/archiver"
	"github.com/gempir/justlog/filelog"
	"github.com/gempir/justlog/helix"
	"github.com/gempir/justlog/humanize"

	log "github.com/sirupsen/logrus"
)

type config struct {
	LogsDirectory string   `json:"logsDirectory"`
	Username      string   `json:"username"`
	OAuth         string   `json:"oauth"`
	ListenAddress string   `json:"listenAddress"`
	Admin         string   `json:"admin"`
	Channels      []string `json:"channels"`
	ClientID      string   `json:"clientID"`
	LogLevel      string   `json:"logLevel"`
}

var (
	cfg config
)

func main() {
	startTime := time.Now()

	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()
	cfg = loadConfiguration(*configFile)

	setupLogger(cfg)
	twitchClient := twitch.NewClient(cfg.Username, "oauth:"+cfg.OAuth)
	fileLogger := filelog.NewFileLogger(cfg.LogsDirectory)
	helixClient := helix.NewClient(cfg.ClientID)
	archiver := archiver.NewArchiver(cfg.LogsDirectory)
	archiver.Boot()

	if strings.HasPrefix(cfg.Username, "justinfan") {
		log.Info("Bot joining anonymous")
	} else {
		log.Info("Bot joining as user " + cfg.Username)
	}

	apiServer := api.NewServer(cfg.LogsDirectory, cfg.ListenAddress, &fileLogger, &helixClient)
	go apiServer.Init()

	for _, channel := range cfg.Channels {
		log.Info("Joining " + channel)
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

	log.Fatal(twitchClient.Connect())
}

func setupLogger(cfg config) {
	switch cfg.LogLevel {
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}
}

func loadConfiguration(file string) config {
	log.Info("Loading config from " + file)

	// setup defaults
	cfg := config{
		LogsDirectory: "./logs",
		ListenAddress: "127.0.0.1:8025",
		Username:      "justinfan777777",
		OAuth:         "oauth:777777777",
		Channels:      []string{},
		Admin:         "gempir",
		LogLevel:      "info",
	}

	configFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// normalize
	cfg.LogsDirectory = strings.TrimSuffix(cfg.LogsDirectory, "/")
	cfg.OAuth = strings.TrimPrefix(cfg.OAuth, "oauth:")
	cfg.LogLevel = strings.ToLower(cfg.LogLevel)

	// ensure required
	if cfg.ClientID == "" {
		log.Fatal("No clientID specified")
	}

	return cfg
}
