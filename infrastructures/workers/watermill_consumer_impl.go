package workers

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/bitbiz/hias-core/infrastructures/queue"
)

type watermillConsumerManager struct {
	queueMgr queue.QueueManager
	config   ConsumerConfig
	handlers []QueueHandler
	mu       sync.Mutex
	running  bool
	cancel   context.CancelFunc
}

func NewWatermillConsumerManager(queueMgr queue.QueueManager, config ConsumerConfig) ConsumerManager {
	return &watermillConsumerManager{
		queueMgr: queueMgr,
		config:   config,
	}
}

func (m *watermillConsumerManager) RegisterHandler(qh QueueHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("cannot register handler while consumer is running")
	}

	m.handlers = append(m.handlers, qh)
	log.Printf("Registered handler %s for topic %s", qh.Handler.GetName(), qh.Topic)
	return nil
}

func (m *watermillConsumerManager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("consumer manager already running")
	}
	m.running = true
	m.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	for _, qh := range m.handlers {
		handler := qh
		err := m.queueMgr.Subscribe(ctx, handler.Topic, func(payload []byte) error {
			procCtx, procCancel := context.WithTimeout(ctx, m.config.ProcessingTimeout)
			defer procCancel()

			if err := handler.Handler.HandleMessage(procCtx, payload); err != nil {
				log.Printf("Error processing message on %s: %v", handler.Topic, err)
				return err
			}
			return nil
		})
		if err != nil {
			cancel()
			return fmt.Errorf("failed to subscribe to %s: %w", handler.Topic, err)
		}
		log.Printf("Consumer subscribed to topic: %s", handler.Topic)
	}

	log.Println("Consumer manager started")
	<-ctx.Done()
	return nil
}

func (m *watermillConsumerManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cancel != nil {
		m.cancel()
	}
	m.running = false
	log.Println("Consumer manager stopped")
	return nil
}

func (m *watermillConsumerManager) IsHealthy() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}
