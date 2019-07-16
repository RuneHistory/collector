package event

import (
	"context"
	"github.com/Shopify/sarama"
	saramaEvents "github.com/jmwri/go-events/sarama"
)

type Handler interface {
	SupportedEventTypes() []string
	Handle(eventType string, payload []byte) error
	GroupName() string
}

func StartAccountManagementHandlers(ctx context.Context, client sarama.Client, handlers []Handler) error {
	group, err := sarama.NewConsumerGroupFromClient("rh-collector.accounts", client)
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
	for _, handler := range handlers {
		err := subscriber.AddHandler(handler.SupportedEventTypes(), handler.Handle)
		if err != nil {
			return err
		}
	}

	err = subscriber.Subscribe(ctx)
	if err != nil {
		return err
	}
	return nil
}
