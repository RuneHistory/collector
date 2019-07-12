package account

import (
	"fmt"
	"github.com/RuneHistory/collector/internal/application/service"
	"github.com/RuneHistory/events"
	"log"
)

func NewRenameAccountHandler(accountService service.Account) *RenameAccountHandler {
	return &RenameAccountHandler{
		AccountService: accountService,
	}
}

type RenameAccountHandler struct {
	AccountService service.Account
}

func (h *RenameAccountHandler) GroupName() string {
	return "rh-collector.RenameAccountHandler"
}
func (h *RenameAccountHandler) SupportedEventTypes() []string {
	return []string{
		(&events.RenameAccountEvent{}).Type(),
	}
}
func (h *RenameAccountHandler) Handle(eventType string, payload []byte) error {
	event := &events.RenameAccountEvent{}
	if eventType != event.Type() {
		return fmt.Errorf("unexpected event type: %s", eventType)
	}
	err := event.WithPayload(payload)
	if err != nil {
		return fmt.Errorf("unexpected event payload: %s", err)
	}

	log.Printf("renaming account: %s", event.ID)
	account, err := h.AccountService.GetById(event.ID)
	if err != nil {
		return fmt.Errorf("failed when getting account: %s", err)
	}
	if account == nil {
		log.Printf("can't rename account that doesn't exist: %s", event.ID)
		return nil
	}
	account.Nickname = event.Nickname
	account, err = h.AccountService.Update(account)
	if err != nil {
		return fmt.Errorf("failed when update account: %s", err)
	}
	log.Printf("updated account: %s", account.ID)
	return nil
}
