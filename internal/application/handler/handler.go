package handler

import (
	"context"
	"github.com/Shopify/sarama"
	saramaEvents "github.com/jmwri/go-events/sarama"
)

type Handler interface {
	SupportedEventTypes() []string
	Handle(eventType string, payload []byte) error
}

func StartHandlers(ctx context.Context, group sarama.ConsumerGroup, handlers []Handler) error {
	subscriber := saramaEvents.NewSubscriber(group)
	for _, handler := range handlers {
		err := subscriber.AddHandler(handler.SupportedEventTypes(), handler.Handle)
		if err != nil {
			return err
		}
	}

	return subscriber.Subscribe(ctx)
}
