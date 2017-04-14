package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/op/go-logging"
	"gopkg.in/redis.v5"

	"github.com/gempir/gempbotgo/api"
	"github.com/gempir/gempbotgo/filelog"
	"github.com/gempir/go-twitch-irc"
	"fmt"
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
				uptime := formatDiff(diff(startTime, time.Now()))
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

func formatDiff(years, months, days, hours, mins, secs int) string {
	since := ""
	if years > 0 {
		switch years {
		case 1:
			since += fmt.Sprintf("%d year ", years)
			break
		default:
			since += fmt.Sprintf("%d years ", years)
			break
		}
	}
	if months > 0 {
		switch months {
		case 1:
			since += fmt.Sprintf("%d month ", months)
			break
		default:
			since += fmt.Sprintf("%d months ", months)
			break
		}
	}
	if days > 0 {
		switch days {
		case 1:
			since += fmt.Sprintf("%d day ", days)
			break
		default:
			since += fmt.Sprintf("%d days ", days)
			break
		}
	}
	if hours > 0 {
		switch hours {
		case 1:
			since += fmt.Sprintf("%d hour ", hours)
			break
		default:
			since += fmt.Sprintf("%d hours ", hours)
			break
		}
	}
	if mins > 0 && days == 0 && months == 0 && years == 0 {
		switch mins {
		case 1:
			since += fmt.Sprintf("%d min ", mins)
			break
		default:
			since += fmt.Sprintf("%d mins ", mins)
			break
		}
	}
	if secs > 0 && days == 0 && months == 0 && years == 0 && hours == 0 {
		switch secs {
		case 1:
			since += fmt.Sprintf("%d sec ", secs)
			break
		default:
			since += fmt.Sprintf("%d secs ", secs)
			break
		}
	}
	return strings.TrimSpace(since)
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}
