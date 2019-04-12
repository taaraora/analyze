package scheduler

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Interface interface {
	ScheduleJob(jobID string, interval time.Duration, job func() error) error
	RemoveJob(jobID string) error
	Stop()
}

type workItem struct {
	ID     string
	ticker *time.Ticker
	done   chan struct{}
}

type scheduler struct {
	logger    logrus.FieldLogger
	close     chan struct{}
	m         sync.Mutex
	workItems map[string]*workItem
	isClosed  bool
}

func NewScheduler(logger logrus.FieldLogger) Interface {
	s := &scheduler{
		logger:    logger,
		close:     make(chan struct{}),
		workItems: make(map[string]*workItem),
	}

	go func() {
		<-s.close
		s.m.Lock()
		defer s.m.Unlock()
		// TODO: shall I make it synchronous?
		for _, wi := range s.workItems {
			go func(w *workItem) { w.done <- struct{}{} }(wi)
		}
		s.isClosed = true
	}()

	return s
}

func (s *scheduler) Stop() {
	s.close <- struct{}{}
}

func (s *scheduler) ScheduleJob(jobID string, interval time.Duration, job func() error) error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.isClosed {
		return errors.New("scheduler is closed")
	}

	if interval <= 0 {
		return errors.New("interval need to be more than 0")
	}

	wi := &workItem{
		ID:     jobID,
		ticker: time.NewTicker(interval),
		done:   make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-wi.ticker.C:
				err := job()
				if err != nil {
					s.logger.Errorf("Job: %v, failed, error: %v", wi.ID, err)
				}
			case <-wi.done:
				wi.ticker.Stop()
				return
			}
		}
	}()

	s.workItems[jobID] = wi

	return nil
}

func (s *scheduler) RemoveJob(jobID string) error {
	s.m.Lock()
	defer s.m.Unlock()

	v, exists := s.workItems[jobID]
	if !exists {
		return errors.Errorf("there is no scheduled job with ID: %v", jobID)
	}

	v.done <- struct{}{}

	return nil
}
