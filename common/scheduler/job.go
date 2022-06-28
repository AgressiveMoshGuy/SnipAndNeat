package scheduler

import (
	"time"
)

type job struct {
	period    time.Duration
	startAt   time.Time
	task      *Task
	runAtOnce bool
}

func (j job) IsPeriodic() bool { return j.period != 0 }
func (j job) IsSingle() bool   { return !j.IsPeriodic() }

func (s *Scheduler) AddPeriodic(t *Task, d time.Duration, r bool) error {
	err := s.add(&job{
		period:    d,
		task:      t,
		runAtOnce: r,
	})
	if err != nil {
		return err
	}
	return nil
}

// func (s *Scheduler) NewTimedJob(t Task, startAt string) job {
// 	log.Println(startAt)
// 	st, err := time.ParseInLocation("2006-01-02T15:04:05Z", startAt, time.Local)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println(st.String())
// 	return s.AddDelayed(t, st)
// }

func (s *Scheduler) AddDelayed(t *Task, startAt time.Time) error {
	err := s.add(&job{
		startAt: startAt,
		task:    t})
	if err != nil {
		return err
	}
	return nil
}
