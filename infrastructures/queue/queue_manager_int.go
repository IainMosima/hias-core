package queue

import "context"

type QueueManager interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Subscribe(ctx context.Context, topic string, handler func(payload []byte) error) error
	Close() error
}
