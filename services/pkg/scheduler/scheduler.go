package scheduler

import (
	"context"
	"time"
)

type Scheduler struct {
	callback func(ctx context.Context) error
	err      chan error
	ticker   time.Ticker
}

func NewScheduler(frequency time.Duration, callback func(ctx context.Context) error) *Scheduler {
	s := &Scheduler{
		ticker:   *time.NewTicker(frequency),
		callback: callback,
		err:      make(chan error),
	}

	return s
}

func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		for {
			<-s.ticker.C
			err := s.callback(ctx)
			if err != nil {
				s.err <- err
			}
		}
	}()
}

func (s *Scheduler) Error() chan error {
	return s.err
}
