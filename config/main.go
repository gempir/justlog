package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v2"
	log "github.com/sirupsen/logrus"
)

// Config application configuration
type Config struct {
	configFile            string
	configFilePermissions os.FileMode
	LogsDirectory         string                   `json:"logsDirectory"`
	Archive               bool                     `json:"archive"`
	AdminAPIKey           string                   `json:"adminAPIKey"`
	Username              string                   `json:"username"`
	OAuth                 string                   `json:"oauth"`
	ListenAddress         string                   `json:"listenAddress"`
	Admins                []string                 `json:"admins"`
	Channels              []string                 `json:"channels"`
	ClientID              string                   `json:"clientID"`
	ClientSecret          string                   `json:"clientSecret"`
	LogLevel              string                   `json:"logLevel"`
	ChannelConfigs        map[string]ChannelConfig `json:"channelConfigs"`
}

// ChannelConfig config for individual channels
type ChannelConfig struct {
	MessageTypes []twitch.MessageType `json:"messageTypes,omitempty"`
}

// NewConfig create configuration from file
func NewConfig(filePath string) *Config {
	cfg := loadConfiguration(filePath)

	log.Info("Loaded config from " + filePath)

	return cfg
}

// AddChannels adds channels to the config
func (cfg *Config) AddChannels(channelIDs ...string) {
	cfg.Channels = append(cfg.Channels, channelIDs...)
	for _, id := range channelIDs {
		cfg.Channels = appendIfMissing(cfg.Channels, id)
	}

	cfg.persistConfig()
}

// SetMessageTypes sets recorded message types for a channel
func (cfg *Config) SetMessageTypes(channelID string, messageTypes []twitch.MessageType) {
	if _, ok := cfg.ChannelConfigs[channelID]; ok {
		channelCfg := cfg.ChannelConfigs[channelID]
		channelCfg.MessageTypes = messageTypes

		cfg.ChannelConfigs[channelID] = channelCfg
	} else {
		cfg.ChannelConfigs[channelID] = ChannelConfig{
			MessageTypes: messageTypes,
		}
	}

	cfg.persistConfig()
}

// ResetMessageTypes removed message type option and therefore resets it
func (cfg *Config) ResetMessageTypes(channelID string) {
	if _, ok := cfg.ChannelConfigs[channelID]; ok {
		channelCfg := cfg.ChannelConfigs[channelID]
		channelCfg.MessageTypes = nil

		cfg.ChannelConfigs[channelID] = channelCfg
	}

	cfg.persistConfig()
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func (cfg *Config) persistConfig() {
	fileContents, err := json.MarshalIndent(*cfg, "", "    ")
	if err != nil {
		log.Error(err)
		return
	}

	err = ioutil.WriteFile(cfg.configFile, fileContents, cfg.configFilePermissions)
	if err != nil {
		log.Error(err)
	}
}

func loadConfiguration(filePath string) *Config {
	// setup defaults
	cfg := Config{
		configFile:     filePath,
		LogsDirectory:  "./logs",
		ListenAddress:  "127.0.0.1:8025",
		Username:       "justinfan777777",
		OAuth:          "oauth:777777777",
		Channels:       []string{},
		ChannelConfigs: make(map[string]ChannelConfig),
		Admins:         []string{"gempir"},
		LogLevel:       "info",
		Archive:        true,
	}

	info, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	cfg.configFilePermissions = info.Mode()

	configFile, err := os.Open(filePath)
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
	cfg.setupLogger()

	// ensure required
	if cfg.ClientID == "" {
		log.Fatal("No clientID specified")
	}

	return &cfg
}

func (cfg *Config) setupLogger() {
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
