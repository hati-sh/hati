package storage

import (
	"errors"
	"sync"
)

var ErrKeyNotExist = errors.New("KEY_NOT_EXIST")

type memoryStorage struct {
	store  map[string][]byte
	rwLock sync.RWMutex
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		store: make(map[string][]byte),
	}
}

func (s *memoryStorage) Has(key string) bool {
	return s.store[key] != nil
}

// Set
func (s *memoryStorage) Set(key string, value []byte) bool {

	s.rwLock.Lock()
	s.store[key] = value
	s.rwLock.Unlock()

	return true
}

func (s *memoryStorage) Get(key string) ([]byte, error) {
	if s.store[key] != nil {
		return nil, ErrKeyNotExist
	}

	return s.store[key], nil
}

func (s *memoryStorage) Delete(key string) {
	if s.store[key] != nil {
		s.rwLock.Lock()
		s.store[key] = nil
		s.rwLock.Unlock()
	}
}
