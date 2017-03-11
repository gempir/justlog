package main

import (
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCanInitLogger(t *testing.T) {
	log := initLogger()

	assert.Equal(t, "logging.Logger", reflect.TypeOf(log).String(), "logger has invalid type")
}

func TestCanReadConfig(t *testing.T) {
	cfg, err := readConfig("config.example.json")
	if err != nil {
		t.Fatal("error reading config", err)
	}

	assert.Equal(t, "irc.chat.twitch.tv:6667", cfg.IrcAddress, "Invalid config data")
	assert.Equal(t, "gempbot", cfg.IrcUser, "Invalid config data")
	assert.Equal(t, "oauth:123123123", cfg.IrcToken, "Invalid config data")
	assert.Equal(t, "gempir", cfg.Admin, "Invalid config data")
	assert.Equal(t, "/var/twitch_logs/", cfg.LogPath, "Invalid config data")
	assert.Equal(t, "8025", cfg.APIPort, "Invalid config data")
	assert.Equal(t, "127.0.0.1:6379", cfg.RedisAddress, "Invalid config data")
	assert.Equal(t, "asdasd", cfg.RedisPassword, "Invalid config data")
	assert.Equal(t, int(0), cfg.RedisDatabase, "Invalid config data")
}
