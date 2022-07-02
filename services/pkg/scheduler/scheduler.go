package scheduler

import (
	"context"
	"time"
)

type Scheduler struct {
	callback func(ctx context.Context) error
	err      chan error
	ticker   time.Ticker
	ctx      context.Context
}

func NewScheduler(frequency time.Duration, ctx context.Context, callback func(ctx context.Context) error) *Scheduler {
	s := &Scheduler{
		ticker:   *time.NewTicker(frequency),
		callback: callback,
		err:      make(chan error),
		ctx:      ctx,
	}

	s.start()
	return s
}

func (s *Scheduler) start() {
	go func() {
		for {
			<-s.ticker.C
			err := s.callback(s.ctx)
			if err != nil {
				s.err <- err
			}
		}
	}()
}

func (s *Scheduler) Error() chan error {
	return s.err
}
