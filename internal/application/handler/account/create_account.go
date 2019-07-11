package account

import (
	"fmt"
	"github.com/RuneHistory/events"
	"log"
)

type CreateAccountHandler struct{}

func (h *CreateAccountHandler) SupportedEventTypes() []string {
	return []string{
		(&events.NewAccountEvent{}).Type(),
	}
}
func (h *CreateAccountHandler) Handle(eventType string, payload []byte) error {
	event := &events.NewAccountEvent{}
	if eventType != event.Type() {
		return fmt.Errorf("unexpected event type: %s", eventType)
	}
	err := event.WithPayload(payload)
	if err != nil {
		return fmt.Errorf("unexpected event payload: %s", err)
	}

	log.Printf("creating account: %s\n", event.ID)
	// TODO: Insert to DB
	log.Printf("created account: %s\n", event)
	return nil
}
