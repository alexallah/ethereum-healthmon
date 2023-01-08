package common

import (
	"log"
	"sync"
)

type State struct {
	sync.RWMutex
	healthy bool
	errors  int
}

func (s *State) IsHealthy() bool {
	s.RLock()
	defer s.RUnlock()
	return s.healthy
}

func (s *State) Error(err error) {
	log.Printf("update error: %q", err)
	s.Lock()
	s.errors += 1
	s.Unlock()
	if s.errors == 3 && s.IsHealthy() {
		log.Println("marking unhealthy")
		s.setUnhealthy()
	}
}

func (s *State) setUnhealthy() {
	s.Lock()
	s.healthy = false
	s.Unlock()
}

func (s *State) SetHealthy() {
	s.Lock()
	if !s.healthy {
		log.Println("marking healthy")
	}
	s.healthy = true
	s.errors = 0
	s.Unlock()
}
