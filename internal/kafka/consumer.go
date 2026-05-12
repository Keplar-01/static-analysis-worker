package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/diploma/worker-static-analyzer/internal/model"
	kafkago "github.com/segmentio/kafka-go"
)

type MessageHandler interface {
	HandleStartEvent(ctx context.Context, event model.StartEvent)
}

type Consumer struct {
	reader  *kafkago.Reader
	handler MessageHandler
}

func NewConsumer(brokers string, handler MessageHandler) *Consumer {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  []string{brokers},
		Topic:    TopicStartStatic,
		GroupID:  "worker-static-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &Consumer{reader: reader, handler: handler}
}

func (c *Consumer) Listen(ctx context.Context) {
	log.Printf("[kafka-consumer] listening for %s...", TopicStartStatic)
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("[kafka-consumer] read error: %v", err)
			continue
		}
		log.Printf("[kafka-consumer] received: %s", string(msg.Value))

		var event model.StartEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("[kafka-consumer] unmarshal error: %v", err)
			continue
		}
		c.handler.HandleStartEvent(ctx, event)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
