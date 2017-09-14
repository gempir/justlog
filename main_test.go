package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanReadConfig(t *testing.T) {
	cfg, err := readConfig("configs/config.example.json")
	if err != nil {
		t.Fatal("error reading sysConfig", err)
	}

	assert.Equal(t, "gempbot", cfg.IrcUser, "Invalid sysConfig data")
	assert.Equal(t, "oauth:123123123", cfg.IrcToken, "Invalid sysConfig data")
	assert.Equal(t, "gempir", cfg.Admin, "Invalid sysConfig data")
}
