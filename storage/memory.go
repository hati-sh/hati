package storage

import (
	"errors"
	"sync"
)

var ErrKeyNotExist = errors.New("KEY_NOT_EXIST\n")

type memoryStorage struct {
	store  map[string][]byte
	rwLock sync.RWMutex
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		store: make(map[string][]byte),
	}
}

func (s *memoryStorage) Has(key []byte) bool {
	return s.store[string(key)] != nil
}

func (s *memoryStorage) Set(key []byte, value []byte) bool {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	s.store[string(key)] = value

	return true
}

func (s *memoryStorage) Get(key []byte) ([]byte, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	if s.store[string(key)] != nil {
		return s.store[string(key)], nil
	}

	return nil, ErrKeyNotExist
}

func (s *memoryStorage) Delete(key []byte) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	keyString := string(key)
	if s.store[keyString] != nil {
		delete(s.store, keyString)
	}
}
