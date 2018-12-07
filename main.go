package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"strings"

	"github.com/gempir/justlog/api"
	"github.com/gempir/justlog/archiver"
	"github.com/gempir/justlog/bot"
	"github.com/gempir/justlog/filelog"
	"github.com/gempir/justlog/helix"

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

func main() {
	startTime := time.Now()

	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()
	cfg := loadConfiguration(*configFile)

	setupLogger(cfg)
	fileLogger := filelog.NewFileLogger(cfg.LogsDirectory)
	helixClient := helix.NewClient(cfg.ClientID)
	archiver := archiver.NewArchiver(cfg.LogsDirectory)
	go archiver.Boot()

	apiServer := api.NewServer(cfg.LogsDirectory, cfg.ListenAddress, &fileLogger, &helixClient)
	go apiServer.Init()

	bot := bot.NewBot(cfg.Admin, cfg.Username, cfg.OAuth, &startTime, &helixClient, &fileLogger)
	bot.Connect(cfg.Channels)
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
