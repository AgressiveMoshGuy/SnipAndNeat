package scheduler

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"
)

func Test_Delayed(t *testing.T) {
	var a bool
	mu := sync.RWMutex{}

	t1 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a = true
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	s := New("test", &Config{})

	err := s.AddDelayed(&t1, time.Now().Add(1000*time.Millisecond))
	if err != nil {
		t.Fatal("add task 1 failed")
	}

	s.Start(context.TODO())

	time.Sleep(900 * time.Millisecond)

	mu.Lock()
	if a {
		t.Fatal("too early")
	}
	mu.Unlock()

	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	if !a {
		t.Fatal("too late")
	}
	mu.Unlock()

	stopCtx, stopCancel := context.WithCancel(context.Background())
	stopCancel()

	s.Stop(stopCtx)
	log.Println("Test_Delayed OK")
}

func Test_Periodic(t *testing.T) {
	var a int
	mu := sync.RWMutex{}

	t1 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a++
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	s := New("test", &Config{})

	err := s.AddPeriodic(&t1, 1000*time.Millisecond, false)
	if err != nil {
		t.Fatal("add task 1 failed")
	}

	s.Start(context.TODO())

	time.Sleep(5100 * time.Millisecond)

	mu.Lock()
	switch {
	case a < 5:
		t.Fatal("too little")
	case a > 5:
		t.Fatal("too much")
	}
	mu.Unlock()

	stopCtx, stopCancel := context.WithCancel(context.Background())
	stopCancel()

	s.Stop(stopCtx)
	log.Println("Test_Periodic OK")
}

func Test_Periodic_RunAtOnce(t *testing.T) {
	var a int
	mu := sync.RWMutex{}

	t1 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a++
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	s := New("test", &Config{})

	err := s.AddPeriodic(&t1, 1000*time.Millisecond, true)
	if err != nil {
		t.Fatal("add task 1 failed")
	}

	s.Start(context.TODO())

	time.Sleep(5100 * time.Millisecond)

	mu.Lock()
	switch {
	case a < 6:
		t.Fatal("too little")
	case a > 6:
		t.Fatal("too much")
	}
	mu.Unlock()

	stopCtx, stopCancel := context.WithCancel(context.Background())
	stopCancel()

	s.Stop(stopCtx)
	log.Println("Test_Periodic_RunAtOnce OK")
}

func Test_AddTaskAfterStart(t *testing.T) {
	var a int
	mu := sync.RWMutex{}

	t1 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a++
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	t2 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a++
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	s := New("test", &Config{})

	err := s.AddPeriodic(&t1, 1000*time.Millisecond, false)
	if err != nil {
		t.Fatal("add task 1 failed")
	}

	s.Start(context.TODO())

	err = s.AddPeriodic(&t2, 1000*time.Millisecond, false)
	if err != nil {
		t.Fatal("add task 2 failed")
	}

	time.Sleep(5100 * time.Millisecond)

	mu.Lock()
	switch {
	case a < 10:
		t.Fatal("too little")
	case a > 10:
		t.Fatal("too much")
	}
	mu.Unlock()

	stopCtx, stopCancel := context.WithCancel(context.Background())
	stopCancel()

	s.Stop(stopCtx)
	log.Println("Test_AddTaskAfterStart OK")
}

func Test_DeleteTaskAfterStart(t *testing.T) {
	var a int
	mu := sync.RWMutex{}

	t1 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a++
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	t2 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a++
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	s := New("test", &Config{})

	err := s.AddPeriodic(&t1, 1000*time.Millisecond, false)
	if err != nil {
		t.Fatal("add task 1 failed")
	}

	s.Start(context.TODO())

	err = s.AddPeriodic(&t2, 1000*time.Millisecond, false)
	if err != nil {
		t.Fatal("add task 2 failed")
	}
	time.Sleep(2 * time.Second)

	err = s.RemoveTask(&t2)
	if err != nil {
		t.Fatal("delete task 2 failed")
	}
	time.Sleep(4600 * time.Millisecond)

	mu.Lock()
	switch {
	case a < 7:
		t.Fatal("too little")
	case a > 7:
		t.Fatal("too much")
	}
	mu.Unlock()
	stopCtx, stopCancel := context.WithCancel(context.Background())
	stopCancel()

	s.Stop(stopCtx)
	log.Println("Test_DeleteTaskAfterStart OK")
}

func Test_AddSamePeriodicTaskTwice(t *testing.T) {
	var a int
	mu := sync.RWMutex{}

	t1 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a++
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	s := New("test", &Config{})

	err := s.AddPeriodic(&t1, 1000*time.Millisecond, false)
	if err != nil {
		t.Fatal("add task 1 failed")
	}

	s.Start(context.TODO())

	err = s.AddPeriodic(&t1, 1000*time.Millisecond, false)
	if err != nil {
		t.Log("add same task failed")
	}

	time.Sleep(5100 * time.Millisecond)

	mu.Lock()
	switch {
	case a < 5:
		t.Fatal("too little")
	case a > 5:
		t.Fatal("too much")
	}
	mu.Unlock()
	stopCtx, stopCancel := context.WithCancel(context.Background())
	stopCancel()

	s.Stop(stopCtx)
	log.Println("Test_AddSameTaskTwice OK")
}

func Test_AddSameSingleTaskTwice(t *testing.T) {
	var a int
	mu := sync.RWMutex{}

	t1 := Task(func(ctx context.Context) error {
		errCh := make(chan error)

		go func() {
			mu.Lock()
			a++
			mu.Unlock()
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return errors.New("fail")
		}
	})

	s := New("test", &Config{})

	err := s.AddDelayed(&t1, time.Now().Add(1*time.Second))
	if err != nil {
		t.Fatal("add task 1 failed")
	}

	s.Start(context.TODO())

	err = s.AddDelayed(&t1, time.Now().Add(3*time.Second))
	if err != nil {
		t.Log("add same task failed - OK")
	}

	time.Sleep(5100 * time.Millisecond)

	mu.Lock()
	switch {
	case a < 1:
		t.Fatal("too little")
	case a > 1:
		t.Fatal("too much")
	}
	mu.Unlock()
	stopCtx, stopCancel := context.WithCancel(context.Background())
	stopCancel()

	s.Stop(stopCtx)
	log.Println("Test_AddSameSingleTaskTwice OK")
}
