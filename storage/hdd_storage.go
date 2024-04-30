package storage

import (
	"context"
	"github.com/hati-sh/hati/common"
)

type hddStorage struct {
	ctx   context.Context
	store HddShardMap
	ttlGc *hddStorageTtlGc
}

func NewHddStorage(ctx context.Context, dataDir string) *hddStorage {
	hs := &hddStorage{
		ctx:   ctx,
		store: newHddShardMap(common.STORAGE_DEFAULT_NUMBER_OF_SHARDS, dataDir),
	}

	hs.ttlGc = NewHddStorageTtlGc(hs)
	return hs
}

func (s *hddStorage) Start() {
	s.ttlGc.Start()
}

func (s *hddStorage) Stop() error {
	s.ttlGc.Stop()

	for _, shard := range s.store {
		if err := shard.db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (s *hddStorage) CountKeys() int {
	return s.store.CountKeys()
}

func (s *hddStorage) Has(key []byte) bool {
	return s.store.Has(string(key))
}

func (s *hddStorage) Set(key []byte, value []byte, ttl int64) bool {
	ok := s.store.Set(string(key), value)

	if ok {
		if ttl > 0 {
			s.ttlGc.SetTtl(ttl, key)
		}
		return true
	}

	return false
}

func (s *hddStorage) Get(key []byte) ([]byte, error) {
	value, ok := s.store.Get(string(key))
	if !ok {
		return nil, ErrKeyNotExist
	}

	return value, nil
}

func (s *hddStorage) Delete(key []byte) bool {
	return s.store.Delete(string(key))
}

func (s *hddStorage) FlushAll() (bool, error) {
	return s.store.FlushAll()
}
