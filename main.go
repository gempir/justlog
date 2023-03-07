package main

import (
	"embed"
	"flag"

	"github.com/gempir/justlog/api"
	"github.com/gempir/justlog/archiver"
	"github.com/gempir/justlog/bot"
	"github.com/gempir/justlog/config"
	"github.com/gempir/justlog/filelog"
	"github.com/gempir/justlog/helix"
)

// content holds our static web server content.
//
//go:embed web/dist/*
var assets embed.FS

func main() {

	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)

	fileLogger := filelog.NewFileLogger(cfg.LogsDirectory)
	helixClient := helix.NewClient(cfg.ClientID, cfg.ClientSecret)
	go helixClient.StartRefreshTokenRoutine()

	if cfg.Archive {
		archiver := archiver.NewArchiver(cfg.LogsDirectory)
		go archiver.Boot()
	}

	bot := bot.NewBot(cfg, &helixClient, &fileLogger)

	apiServer := api.NewServer(cfg, bot, &fileLogger, &helixClient, assets)
	go apiServer.Init()

	bot.Connect()
}
