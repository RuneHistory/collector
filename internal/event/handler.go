package event

import (
	"context"
	"github.com/jmwri/go-events"
	"sync"
)

type Handler interface {
	SupportedEventTypes() []string
	Handle(eventType string, payload []byte) error
}

func StartEventHandlers(ctx context.Context, wg *sync.WaitGroup, subscriber go_events.Subscriber, handlers []Handler, errCh chan<- error) {
	wg.Add(len(handlers))
	for _, handler := range handlers {
		go func() {
			err := subscriber.Subscribe(ctx, handler.SupportedEventTypes(), handler.Handle)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}
}
