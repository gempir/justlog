package config

import (
	"github.com/gempir/gempbotgo/modules"
	"gopkg.in/redis.v5"
)

type UserConfig struct {
	rClient redis.Client
}

func NewUserConfig(rClient redis.Client) UserConfig {
	return UserConfig{
		rClient: rClient,
	}
}

func (uCfg *UserConfig) IsEnabled(channel, key string) bool {
	if key == modules.STATUS {
		return true
	}

	res, err := uCfg.rClient.HGet(channel+":config", key).Result()

	if err != nil {
		return false
	}
	if res == "true" || res == "1" {
		return true
	}
	return false
}
