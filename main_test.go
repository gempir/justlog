package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCanInitLogger(t *testing.T) {
	log := initLogger()

	assert.Equal(t, "logging.Logger", reflect.TypeOf(log).String(), "logger has invalid type")
}

func TestCanReadConfig(t *testing.T) {
	cfg, err := readConfig("sysConfig.example.json")
	if err != nil {
		t.Fatal("error reading sysConfig", err)
	}

	assert.Equal(t, "irc.chat.twitch.tv:6667", cfg.IrcAddress, "Invalid sysConfig data")
	assert.Equal(t, "gempbot", cfg.IrcUser, "Invalid sysConfig data")
	assert.Equal(t, "oauth:123123123", cfg.IrcToken, "Invalid sysConfig data")
	assert.Equal(t, "gempir", cfg.Admin, "Invalid sysConfig data")
	assert.Equal(t, "/var/twitch_logs/", cfg.LogPath, "Invalid sysConfig data")
	assert.Equal(t, "8025", cfg.APIPort, "Invalid sysConfig data")
	assert.Equal(t, "127.0.0.1:6379", cfg.RedisAddress, "Invalid sysConfig data")
	assert.Equal(t, "asdasd", cfg.RedisPassword, "Invalid sysConfig data")
	assert.Equal(t, int(0), cfg.RedisDatabase, "Invalid sysConfig data")
}
