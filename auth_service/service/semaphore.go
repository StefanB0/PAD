package service

import (
	"sync"
)

type Semaphore struct {
	maximum int
	current int
	mutex   sync.Mutex
	cond    *sync.Cond
}

func NewSemaphore(maximum int) *Semaphore {
	sem := &Semaphore{
		maximum: maximum,
		current: 0,
	}
	sem.cond = sync.NewCond(&sem.mutex)
	return sem
}

func (s *Semaphore) Acquire() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for s.current >= s.maximum {
		s.cond.Wait()
	}

	s.current++
}

func (s *Semaphore) Release() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.current--
	s.cond.Signal()
}
