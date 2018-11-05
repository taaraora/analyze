package scheduler

import (
	"time"
)

type Scheduler struct {
	workItem func()
	ticker   *time.Ticker
	done     chan struct{}
}

func NewScheduler(interval time.Duration, workItem func()) *Scheduler {
	scheduler := &Scheduler{
		ticker:   time.NewTicker(interval),
		workItem: workItem,
		done:     make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-scheduler.ticker.C:
				workItem()
			case <-scheduler.done:
				return
			}
		}
	}()

	return scheduler
}

func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.done <- struct{}{}
}
