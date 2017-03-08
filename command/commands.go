package command

import (
	"github.com/gempir/gempbotgo/twitch"
	"github.com/op/go-logging"
)

type handler struct {
	log logging.Logger
}

func NewHandler(logger logging.Logger) handler {
	return handler{
		log: logger,
	}
}

func (h *handler) HandleMessage(msg twitch.Message) {
	h.log.Debug(msg.Text)
}