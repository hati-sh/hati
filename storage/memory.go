package storage

import (
	"errors"
	"sync"
)

var ErrKeyNotExist = errors.New("KEY_NOT_EXIST\n")

type memoryStorage struct {
	store  ShardMap
	rwLock sync.RWMutex
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		store: newShardMap(DEFAULT_NUMBER_OF_SHARDS),
	}
}

func (s *memoryStorage) Has(key []byte) bool {

	return s.store.Has(string(key))
}

func (s *memoryStorage) Set(key []byte, value []byte) bool {
	s.store.Set(string(key), value)

	return true
}

func (s *memoryStorage) Get(key []byte) ([]byte, error) {
	value, ok := s.store.Get(string(key))
	if !ok {
		return nil, ErrKeyNotExist
	}

	return value, nil
}

func (s *memoryStorage) Delete(key []byte) {
	s.store.Delete(string(key))
}
