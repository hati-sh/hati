package storage

import (
	"context"
	"github.com/hati-sh/hati/common"
)

type memoryStorage struct {
	ctx   context.Context
	store MemoryShardMap
	ttlGc *memoryStorageTtlGc
}

func NewMemoryStorage(ctx context.Context) *memoryStorage {
	ms := &memoryStorage{
		ctx:   ctx,
		store: newMemoryShardMap(common.STORAGE_DEFAULT_NUMBER_OF_SHARDS),
	}
	ms.ttlGc = NewMemoryStorageTtlGc(ms)

	return ms
}

func (s *memoryStorage) Start() {
	s.ttlGc.Start()
}

func (s *memoryStorage) Stop() {
	s.ttlGc.Stop()
}

func (s *memoryStorage) CountKeys() int {
	return s.store.CountKeys()
}

func (s *memoryStorage) Has(key []byte) bool {
	return s.store.Has(string(key))
}

func (s *memoryStorage) Set(key []byte, value []byte, ttl int64) bool {
	s.store.Set(string(key), value)

	if ttl > 0 {
		s.ttlGc.SetTtl(ttl, key)
	}

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

func (s *memoryStorage) FlushAll() (bool, error) {
	return s.store.FlushAll()
}
