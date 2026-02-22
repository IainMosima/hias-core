package workers

import (
	"github.com/bitbiz/hias-core/infrastructures/queue"
)

func NewConsumerManager(queueMgr queue.QueueManager, config ConsumerConfig) ConsumerManager {
	return NewWatermillConsumerManager(queueMgr, config)
}
