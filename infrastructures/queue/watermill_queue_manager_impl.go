package queue

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type watermillQueueManager struct {
	publisher  *WatermillPublisher
	subscriber message.Subscriber
	config     QueueConfig
	logger     watermill.LoggerAdapter
	router     *message.Router
}

func NewWatermillQueueManager(publisher *WatermillPublisher, subscriber message.Subscriber, config QueueConfig) (QueueManager, error) {
	logger := watermill.NewStdLogger(false, false)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create message router: %w", err)
	}

	return &watermillQueueManager{
		publisher:  publisher,
		subscriber: subscriber,
		config:     config,
		logger:     logger,
		router:     router,
	}, nil
}

func (m *watermillQueueManager) Publish(ctx context.Context, topic string, payload []byte) error {
	return m.publisher.Publish(ctx, topic, payload)
}

func (m *watermillQueueManager) Subscribe(_ context.Context, topic string, handler func(payload []byte) error) error {
	queueURL := m.config.GetQueueURL(topic)
	if queueURL == "" {
		return fmt.Errorf("unknown queue topic: %s", topic)
	}

	m.router.AddNoPublisherHandler(
		fmt.Sprintf("handler-%s", topic),
		queueURL,
		m.subscriber,
		func(msg *message.Message) error {
			if err := handler(msg.Payload); err != nil {
				msg.Nack()
				return err
			}
			msg.Ack()
			return nil
		},
	)

	return nil
}

func (m *watermillQueueManager) Close() error {
	if err := m.publisher.Close(); err != nil {
		return err
	}
	return m.router.Close()
}

func (m *watermillQueueManager) RunRouter(ctx context.Context) error {
	return m.router.Run(ctx)
}
