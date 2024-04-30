package storage

import (
	"context"
	"github.com/hati-sh/hati/common"
)

type hddStorage struct {
	ctx   context.Context
	store HddShardMap
}

func NewHddStorage(ctx context.Context, dataDir string) *hddStorage {
	return &hddStorage{
		ctx:   ctx,
		store: newHddShardMap(common.STORAGE_DEFAULT_NUMBER_OF_SHARDS, dataDir),
	}
}

func (s *hddStorage) CountKeys() int {
	return s.store.CountKeys()
}

func (s *hddStorage) Has(key []byte) bool {
	return s.store.Has(string(key))
}

func (s *hddStorage) Set(key []byte, value []byte) bool {
	s.store.Set(string(key), value)

	return true
}

func (s *hddStorage) Get(key []byte) ([]byte, error) {
	value, ok := s.store.Get(string(key))
	if !ok {
		return nil, ErrKeyNotExist
	}

	return value, nil
}

func (s *hddStorage) Delete(key []byte) {
	s.store.Delete(string(key))
}

func (s *hddStorage) FlushAll() bool {
	return s.store.FlushAll()
}
