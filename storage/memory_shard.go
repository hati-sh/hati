package storage

import (
	"crypto/sha1"
	"sync"
)

type MemoryShard struct {
	sync.RWMutex
	m map[string][]byte
}

type MemoryShardMap []*MemoryShard

func newShardMap(size int) MemoryShardMap {
	m := make([]*MemoryShard, size)
	for i := 0; i < size; i++ {
		s := MemoryShard{m: make(map[string][]byte)}
		m[i] = &s
	}
	return m
}

func (m MemoryShardMap) getShardKey(key string) int {
	hash := sha1.Sum([]byte(key))
	return int(hash[17]) % len(m)
}

func (m MemoryShardMap) GetShard(key string) *MemoryShard {
	// mShardMutex ?
	shardKey := m.getShardKey(key)
	return m[shardKey]
}

func (m MemoryShardMap) Get(key string) ([]byte, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()
	v, ok := shard.m[key]
	return v, ok
}

func (m MemoryShardMap) Has(key string) bool {
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key] != nil
}

func (m MemoryShardMap) Set(key string, val []byte) {
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.m[key] = val
}

func (m MemoryShardMap) Delete(key string) {
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()

	delete(shard.m, key)
}

func (m MemoryShardMap) CountKeys() int {
	var keysCount = 0
	for _, sm := range m {
		sm.RLock()
		keysCount = keysCount + len(sm.m)
		sm.RUnlock()
	}
	return keysCount
}
