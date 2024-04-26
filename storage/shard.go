package storage

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	m           map[string][]byte
	mShardMutex sync.RWMutex
}

type ShardMap []*Shard

func newShardMap(size int) ShardMap {
	m := make([]*Shard, size)
	for i := 0; i < size; i++ {
		s := Shard{m: make(map[string][]byte)}
		m[i] = &s
	}
	return m
}

func (m ShardMap) getShardKey(key string) int {
	hash := sha1.Sum([]byte(key))
	return int(hash[17]) % len(m)
}

func (m ShardMap) GetShard(key string) *Shard {
	// mShardMutex
	shardKey := m.getShardKey(key)
	return m[shardKey]
}

func (m ShardMap) Get(key string) ([]byte, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()
	v, ok := shard.m[key]
	return v, ok
}

func (m ShardMap) Has(key string) bool {
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key] != nil
}

func (m ShardMap) Set(key string, val []byte) {
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.m[key] = val
}

func (m ShardMap) Delete(key string) {
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	if _, ok := shard.m[key]; ok {
		delete(shard.m, key)
	}
}
