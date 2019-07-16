package handler

import (
	"context"
	"github.com/Shopify/sarama"
)

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
	return StartHandlers(ctx, group, handlers)
}
