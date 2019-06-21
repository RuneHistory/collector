package events

import "time"

var NewAccountEvent = &newAccountEvent{}

type newAccountEvent struct{}

func (e *newAccountEvent) NewPayload() Payload {
	return &NewAccountEventPayload{}
}

func (e *newAccountEvent) Topic() string {
	return "queue.account.new"
}

type NewAccountEventPayload struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}
