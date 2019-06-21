package handler

import (
	"collector/internal/events"
	"log"
)

type KafkaHandler interface {
	events.Handler
	Group() string
}

type NewAccountEventHandler struct{}

func (h *NewAccountEventHandler) Handle(p events.Payload) {
	payload := p.(*events.NewAccountEventPayload)
	log.Printf("got new account event for: %s %s\n", payload.Nickname, payload.ID)
}

func (h *NewAccountEventHandler) HandleError(err error) {
	log.Printf("got error: %s\n", err)
}

func (h *NewAccountEventHandler) Group() string {
	return "rh-collector.assign-bucket"
}
