package event

import (
	"context"
	"github.com/Shopify/sarama"
	saramaEvents "github.com/jmwri/go-events/sarama"
	"sync"
)

type Handler interface {
	SupportedEventTypes() []string
	Handle(eventType string, payload []byte) error
	GroupName() string
}

func StartEventHandler(ctx context.Context, client sarama.Client, handler Handler) error {
	group, err := sarama.NewConsumerGroupFromClient(handler.GroupName(), client)
	if err != nil {
		return err
	}
	defer func() {
		err := group.Close()
		if err != nil {
			panic(err)
		}
	}()
	subscriber := saramaEvents.NewSubscriber(group)

	err = subscriber.Subscribe(ctx, handler.SupportedEventTypes(), handler.Handle)
	if err != nil {
		return err
	}
	return nil
}

func StartEventHandlers(ctx context.Context, client sarama.Client, handlers []Handler, wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	for _, handler := range handlers {
		go func(h Handler) {
			err := StartEventHandler(ctx, client, h)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}(handler)
	}
}
