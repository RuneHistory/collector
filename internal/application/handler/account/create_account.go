package account

import (
	"fmt"
	"github.com/RuneHistory/collector/internal/application/service"
	"github.com/RuneHistory/events"
	"log"
)

func NewCreateAccountHandler(accountService service.Account, bucketService service.Bucket) *CreateAccountHandler {
	return &CreateAccountHandler{
		AccountService: accountService,
		BucketService:  bucketService,
	}
}

type CreateAccountHandler struct {
	AccountService service.Account
	BucketService  service.Bucket
}

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
	account, err := h.AccountService.GetById(event.ID)
	if err != nil {
		return fmt.Errorf("failed when getting account: %s", err)
	}
	if account != nil {
		log.Printf("Asked to create account %s - already exists", event.ID)
		return nil
	}
	bucket, err := h.BucketService.GetPriorityBucket()
	if err != nil {
		return fmt.Errorf("failed when getting priority bucket: %s", err)
	}
	account, err = h.AccountService.Create(event.ID, bucket.ID, event.Nickname)
	if err != nil {
		return fmt.Errorf("failed when creating account: %s", err)
	}
	err = h.BucketService.IncrementAmount(bucket)
	if err != nil {
		return fmt.Errorf("failed when incrementing bucket count: %s", err)
	}
	log.Printf("created account: %s\n", event.ID)
	return nil
}
