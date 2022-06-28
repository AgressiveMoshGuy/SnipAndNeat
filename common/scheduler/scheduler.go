package scheduler

import (
	"context"
	"errors"

	"sync"
	"time"

	"github.com/rs/zerolog"
)

func New(name string, cfg *Config) *Scheduler {
	if name == "" {
		name = "scheduler"
	}

	s := Scheduler{
		name: name,
		mu:   sync.Mutex{},

		startCh:  make(chan struct{}),
		addCh:    make(chan struct{}),
		delCh:    make(chan *Task),
		stopCh:   make(chan struct{}),
		stopedCh: make(chan struct{}),

		jobs:  make(map[*job]contol),
		errCh: make(chan error),

		cfg: cfg,
	}

	go s.worker()

	return &s
}

type contol struct {
	isSingleStarted   bool
	isPeriodicStarted bool
	started           bool
	stopedCh          chan struct{}
	ctx               context.Context
	cancel            context.CancelFunc
}

type Scheduler struct {
	name string
	mu   sync.Mutex

	startCh  chan struct{}
	addCh    chan struct{}
	delCh    chan *Task
	stopCh   chan struct{}
	stopedCh chan struct{}

	jobs  map[*job]contol
	errCh chan error

	cfg *Config
	log zerolog.Logger
}

func (s *Scheduler) worker() {
	<-s.startCh
	s.log.Debug().Msg("worker started")
	for {
		select {
		case <-s.addCh:
			s.mu.Lock()
			jobs := s.jobs
			s.mu.Unlock()

			for job := range jobs {
				if !s.getJob(job).started {
					s.mu.Lock()
					ctx, cancel := context.WithCancel(context.Background())
					s.jobs[job] = contol{
						ctx:      ctx,
						cancel:   cancel,
						started:  true,
						stopedCh: make(chan struct{}),
					}
					s.mu.Unlock()
					go s.starter(job, s.getJob(job))
				}
			}
		case task := <-s.delCh:
			s.delTask(task)
		case <-s.stopCh:
			close(s.startCh)
			return
		}
	}
}

func (s *Scheduler) add(job *job) error {
	if s.existTask(job.task) {
		return errors.New("task already exist")
	}
	s.mu.Lock()
	s.jobs[job] = contol{}
	s.mu.Unlock()
	go func() { s.addCh <- struct{}{} }()
	return nil
}

func (s *Scheduler) RemoveTask(t *Task) error {
	if !s.existTask(t) {
		return errors.New("task is not exist")
	}
	go func() { s.delCh <- t }()
	return nil
}

func (s *Scheduler) Start(context.Context) error {
	go func() { s.startCh <- struct{}{} }()
	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	close(s.stopCh)

	select {
	case <-ctx.Done():
		return nil
	case <-s.stopedCh:
		return nil
	}
}

func (s *Scheduler) starter(job *job, ctrl contol) {
	j := *job

	switch {
	case job.IsSingle():
	loopS:
		for {
			select {
			case <-time.After(time.Until(job.startAt)):
				if !s.existJob(job) {
					s.log.Error().Msg("task is not exist")
					break loopS
				}
				if ctrl.started && !ctrl.isSingleStarted {
					go s.single(j, ctrl.ctx, ctrl)
					ctrl.isSingleStarted = true
				}
			case <-s.stopCh:
				ctrl.cancel()
				<-ctrl.stopedCh
				s.stopedCh <- struct{}{}
				return
			}
		}
	case job.IsPeriodic():
		if ctrl.started && job.runAtOnce {
			if !s.existJob(job) {
				s.log.Error().Msg("task is not exist")
			}
			go s.periodic(j, ctrl.ctx, ctrl)
			ctrl.isPeriodicStarted = true
		}
	loopP:
		for {
			select {
			case <-time.After(job.period):
				if !s.existJob(job) {
					s.log.Error().Msg("task is not exist")
					break loopP
				}
				if ctrl.started {
					go s.periodic(j, ctrl.ctx, ctrl)
					ctrl.isPeriodicStarted = true
				}
			case <-s.stopCh:
				ctrl.cancel()
				<-ctrl.stopedCh
				s.stopedCh <- struct{}{}
				return
			}
		}
	}
}

func (s *Scheduler) existJob(job *job) bool {
	s.mu.Lock()
	_, ok := s.jobs[job]
	s.mu.Unlock()
	return ok
}

func (s *Scheduler) existTask(t *Task) bool {
	for j := range s.jobs {
		if j.task == t {
			return true
		}
	}
	return false
}

func (s *Scheduler) delTask(t *Task) {
	for j := range s.jobs {
		if j.task == t {
			s.mu.Lock()
			delete(s.jobs, j)
			s.mu.Unlock()
		}
	}
}

func (s *Scheduler) getJob(job *job) contol {
	var ctrl contol
	if s.existJob(job) {
		s.mu.Lock()
		ctrl = s.jobs[job]
		s.mu.Unlock()
	}
	return ctrl
}
