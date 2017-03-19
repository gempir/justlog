package command

import (
	"github.com/CleverbotIO/go-cleverbot.io"
	"time"
	"math/rand"
)

func NewCleverBot(apiUser, apiKey string) *cleverbot.Session  {
	bot, err := cleverbot.New(apiUser, apiKey, RandomString(30))
	if err != nil {
		logger.Error(err.Error())
	}
	return bot
}


func RandomString(strLen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strLen)
	for i := 0; i < strLen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}