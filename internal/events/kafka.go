package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

type KafkaEvent interface {
	Event
	Topic() string
}

type KafkaHandler interface {
	Handler
	Group() string
}

func NewKafkaSubscriber(brokerList []string, client sarama.Client, defaultGroup string) *KafkaSubscriber {
	return &KafkaSubscriber{
		BrokerList:   brokerList,
		Client:       client,
		DefaultGroup: defaultGroup,
	}
}

type KafkaSubscriber struct {
	BrokerList   []string
	Client       sarama.Client
	DefaultGroup string
}

// Subscribe creates a consumer group, and returns an channel of events from the specified topic
func (s *KafkaSubscriber) Subscribe(e Event, h Handler) error {
	event := e.(KafkaEvent)
	handler := h.(KafkaHandler)

	group := handler.Group()
	if group == "" {
		group = s.DefaultGroup
	}

	// Create a new consumer, and transform messages into event body
	log.Printf("creating kafka consumer\n")
	consumer := NewKafkaConsumer()
	go func() {
		log.Printf("inside routing listening for consumer message/err\n")
		for {
			select {
			case message := <-consumer.messageCh:
				payload := event.NewPayload()
				err := json.Unmarshal(message.Value, payload)
				if err != nil {
					h.HandleError(err)
				} else {
					h.Handle(payload)
				}
			case err := <-consumer.errCh:
				h.HandleError(err)
			}
		}
	}()

	// Create a new consumer group
	ctx := context.Background()
	client, err := sarama.NewConsumerGroup(s.BrokerList, group, s.Client.Config())
	if err != nil {
		return fmt.Errorf("unable to create consumer group: %s", err)
	}

	go func() {
		for {
			topics := []string{
				event.Topic(),
			}
			consumer.ready = make(chan bool, 0)
			err := client.Consume(ctx, topics, &consumer)
			if err != nil {
				panic(err)
			}
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")
	return nil
}

func NewKafkaConsumer() KafkaConsumer {
	return KafkaConsumer{
		messageCh: make(chan *sarama.ConsumerMessage),
		errCh:     make(chan error),
	}
}

// KafkaConsumer represents a Sarama consumer group consumer
type KafkaConsumer struct {
	ready     chan bool
	messageCh chan *sarama.ConsumerMessage
	errCh     chan error
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
		consumer.messageCh <- message
	}

	return nil
}
