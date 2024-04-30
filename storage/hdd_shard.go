package storage

import (
	"crypto/sha1"
	"errors"
	"github.com/hati-sh/hati/common/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"path"
	"strconv"
)

type HddShardMap []*HddShard

type HddShard struct {
	db *leveldb.DB
}

func newHddShardMap(size int, dataDir string) HddShardMap {
	var err error
	m := make([]*HddShard, size)

	options := &opt.Options{
		WriteBuffer: 1024 * 1024 * 16,
	}

	for i := 0; i < size; i++ {
		m[i] = &HddShard{db: nil}

		dbPath := path.Join(dataDir, "db", "kv_shard_"+strconv.Itoa(i))
		m[i].db, err = leveldb.OpenFile(dbPath, options)
		if err != nil {
			panic(err)
		}
	}

	return m
}

func (m HddShardMap) getShardKey(key string) int {
	hash := sha1.Sum([]byte(key))
	return int(hash[17]) % len(m)
}

func (m HddShardMap) GetShard(key string) *HddShard {
	shardKey := m.getShardKey(key)

	return m[shardKey]
}

func (m HddShardMap) Get(key string) ([]byte, bool) {
	shard := m.GetShard(key)

	data, err := shard.db.Get([]byte(key), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, false
		}
		logger.Error(err.Error())

		return nil, false
	}

	return data, true
}

func (m HddShardMap) Has(key string) bool {
	shard := m.GetShard(key)

	_, err := shard.db.Get([]byte(key), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return false
		}
		logger.Error(err.Error())

		return false
	}

	return true
}

func (m HddShardMap) Set(key string, value []byte) {
	shard := m.GetShard(key)

	if err := shard.db.Put([]byte(key), value, nil); err != nil {
		logger.Error(err.Error())
	}
}

func (m HddShardMap) Delete(key string) {
	shard := m.GetShard(key)
	if err := shard.db.Delete([]byte(key), nil); err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return
		}

		logger.Error(err.Error())
	}
}

func (m HddShardMap) FlushAll() bool {
	// go shard by shard and delete data
	for _, shard := range m {
		iter := shard.db.NewIterator(nil, nil)
		for iter.Next() {
			key := iter.Key()

			if err := shard.db.Delete(key, nil); err != nil {
				if errors.Is(err, leveldb.ErrNotFound) {
					continue
				}

				logger.Error(err.Error())
				break
			}
		}
		iter.Release()

		if err := iter.Error(); err != nil {
			logger.Error(err.Error())
			return false
		}
	}

	return true
}

func (m HddShardMap) CountKeys() int {
	var keysCount = 0
	for _, shard := range m {
		iter := shard.db.NewIterator(nil, nil)
		for iter.Next() {
			keysCount++
		}
	}

	return keysCount
}
