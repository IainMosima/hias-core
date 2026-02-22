package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

type schedulerManagerImpl struct {
	cron  *cron.Cron
	tasks []Task
	mu    sync.Mutex
}

func NewSchedulerManager() SchedulerManager {
	return &schedulerManagerImpl{
		cron: cron.New(cron.WithLogger(cron.DefaultLogger)),
	}
}

func (s *schedulerManagerImpl) RegisterTask(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.cron.AddFunc(task.Schedule(), func() {
		ctx := context.Background()
		log.Printf("Executing scheduled task: %s", task.Name())
		if err := task.Execute(ctx); err != nil {
			log.Printf("Scheduled task %s failed: %v", task.Name(), err)
			return
		}
		log.Printf("Scheduled task %s completed", task.Name())
	})
	if err != nil {
		return fmt.Errorf("failed to register task %s: %w", task.Name(), err)
	}

	s.tasks = append(s.tasks, task)
	log.Printf("Registered scheduled task: %s [%s]", task.Name(), task.Schedule())
	return nil
}

func (s *schedulerManagerImpl) Start() error {
	s.cron.Start()
	log.Printf("Scheduler started with %d tasks", len(s.tasks))
	return nil
}

func (s *schedulerManagerImpl) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped")
}

func (s *schedulerManagerImpl) GetRegisteredTasks() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	names := make([]string, len(s.tasks))
	for i, t := range s.tasks {
		names[i] = t.Name()
	}
	return names
}
