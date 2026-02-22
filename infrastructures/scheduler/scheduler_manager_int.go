package scheduler

import "context"

// Task represents a scheduled task.
type Task interface {
	Name() string
	Schedule() string // cron expression
	Execute(ctx context.Context) error
}

// SchedulerManager manages cron-based scheduled tasks.
type SchedulerManager interface {
	RegisterTask(task Task) error
	Start() error
	Stop()
	GetRegisteredTasks() []string
}
