package workers

import "context"

// MessageHandler processes a single message from a queue.
type MessageHandler interface {
	HandleMessage(ctx context.Context, payload []byte) error
	GetName() string
}

// QueueHandler binds a handler to a specific queue topic.
type QueueHandler struct {
	Topic   string
	Handler MessageHandler
}

// ConsumerManager manages message consumption from SQS queues.
type ConsumerManager interface {
	RegisterHandler(qh QueueHandler) error
	Start(ctx context.Context) error
	Stop() error
	IsHealthy() bool
}
