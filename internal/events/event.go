package events

type Event interface {
	NewPayload() Payload
}

type Payload interface{}

type Subscriber interface {
	Subscribe(e Event, h Handler) error
}

type Handler interface {
	Handle(p Payload)
	HandleError(err error)
}
