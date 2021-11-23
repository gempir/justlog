package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Config application configuration
type Config struct {
	configFile            string
	configFilePermissions os.FileMode
	BotVerified           bool            `json:"botVerified"`
	LogsDirectory         string          `json:"logsDirectory"`
	Archive               bool            `json:"archive"`
	AdminAPIKey           string          `json:"adminAPIKey"`
	Username              string          `json:"username"`
	OAuth                 string          `json:"oauth"`
	ListenAddress         string          `json:"listenAddress"`
	Admins                []string        `json:"admins"`
	Channels              []string        `json:"channels"`
	ClientID              string          `json:"clientID"`
	ClientSecret          string          `json:"clientSecret"`
	LogLevel              string          `json:"logLevel"`
	OptOut                map[string]bool `json:"optOut"`
}

// NewConfig create configuration from file
func NewConfig(filePath string) *Config {
	cfg := loadConfiguration(filePath)

	log.Info("Loaded config from " + filePath)

	return cfg
}

// AddChannels adds channels to the config
func (cfg *Config) AddChannels(channelIDs ...string) {
	for _, id := range channelIDs {
		cfg.Channels = appendIfMissing(cfg.Channels, id)
	}

	cfg.persistConfig()
}

// OptOutUsers will opt out a user
func (cfg *Config) OptOutUsers(userIDs ...string) {
	for _, id := range userIDs {
		cfg.OptOut[id] = true
	}

	cfg.persistConfig()
}

// IsOptedOut check if a user is opted out
func (cfg *Config) IsOptedOut(userID string) bool {
	_, ok := cfg.OptOut[userID]

	return ok
}

// AddChannels remove user from opt out
func (cfg *Config) RemoveOptOut(userIDs ...string) {
	for _, id := range userIDs {
		delete(cfg.OptOut, id)
	}

	cfg.persistConfig()
}

// RemoveChannels removes channels from the config
func (cfg *Config) RemoveChannels(channelIDs ...string) {
	channels := cfg.Channels

	for i, channel := range channels {
		for _, removeChannel := range channelIDs {
			if channel == removeChannel {
				channels[i] = channels[len(channels)-1]
				channels[len(channels)-1] = ""
				channels = channels[:len(channels)-1]
			}
		}
	}

	cfg.Channels = channels
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
		configFile:    filePath,
		LogsDirectory: "./logs",
		ListenAddress: ":8025",
		Username:      "justinfan777777",
		OAuth:         "oauth:777777777",
		Channels:      []string{},
		Admins:        []string{"gempir"},
		LogLevel:      "info",
		Archive:       true,
		OptOut:        map[string]bool{},
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
