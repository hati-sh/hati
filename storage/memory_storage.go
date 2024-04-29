package storage

import (
	"context"

	"github.com/hati-sh/hati/common"
)

type memoryStorage struct {
	ctx   context.Context
	store MemoryShardMap
}

func NewMemoryStorage(ctx context.Context) *memoryStorage {
	return &memoryStorage{
		ctx:   ctx,
		store: newShardMap(common.STORAGE_DEFAULT_NUMBER_OF_SHARDS),
	}
}

func (s *memoryStorage) CountKeys() int {
	return s.store.CountKeys()
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
