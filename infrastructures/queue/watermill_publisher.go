package queue

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type WatermillPublisher struct {
	publisher message.Publisher
	config    QueueConfig
	logger    watermill.LoggerAdapter
}

func NewWatermillPublisher(publisher message.Publisher, config QueueConfig) *WatermillPublisher {
	return &WatermillPublisher{
		publisher: publisher,
		config:    config,
		logger:    watermill.NewStdLogger(false, false),
	}
}

func (p *WatermillPublisher) Publish(_ context.Context, topic string, payload []byte) error {
	queueURL := p.config.GetQueueURL(topic)
	if queueURL == "" {
		return fmt.Errorf("unknown queue topic: %s", topic)
	}

	msg := message.NewMessage(uuid.New().String(), payload)

	if err := p.publisher.Publish(queueURL, msg); err != nil {
		return fmt.Errorf("failed to publish message to %s: %w", topic, err)
	}

	p.logger.Info("Published message", watermill.LogFields{
		"topic":      topic,
		"message_id": msg.UUID,
	})

	return nil
}

func (p *WatermillPublisher) Close() error {
	return p.publisher.Close()
}
