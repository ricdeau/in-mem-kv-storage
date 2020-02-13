package storage

import "sync"

// Storage - key/value storage.
type Storage interface {
	// Get - returns Data with given key and flag representing data existence.
	Get(key string) (data *Data, exists bool)

	// Set - creates or updates data in storage. Returns true if data is new.
	Set(key string, value *Data) (new bool)

	// Delete - deletes date with given key from storage. Do nothing if no such element.
	Delete(key string)
}

// Data - structure of the stored data.
type Data struct {
	// Type - client defined name of the data type.
	Type string

	// Payload - data representation as a bytes slice.
	Payload []byte
}

type storage struct {
	sync.RWMutex
	inner map[string]*Data
}

// New - creates new concurrent safe storage.
func New() Storage {
	return &storage{
		RWMutex: sync.RWMutex{},
		inner:   make(map[string]*Data),
	}
}

func (s *storage) Get(key string) (data *Data, exists bool) {
	s.RLock()
	defer s.RUnlock()
	data, exists = s.inner[key]
	return
}

func (s *storage) Set(key string, value *Data) (new bool) {
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
