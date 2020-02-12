package storage

import "sync"

type Storage interface {
	Get(key string) (data Data, exists bool)
	Set(key string, value Data) (new bool)
	Delete(key string)
}

type Data struct {
	Type    string
	Payload []byte
}

type storage struct {
	sync.RWMutex
	inner map[string]Data
}

func (s *storage) Get(key string) (data Data, exists bool) {
	s.RLock()
	defer s.RUnlock()
	data, exists = s.inner[key]
	return
}

func (s *storage) Set(key string, value Data) (new bool) {
	s.Lock()
	defer s.Unlock()
	_, exists := s.inner[key]
	s.inner[key] = value
	new = !exists
	return
}

func (s *storage) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.inner, key)
}

func New() Storage {
	return &storage{
		RWMutex: sync.RWMutex{},
		inner:   make(map[string]Data),
	}
}
