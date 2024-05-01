package storage

import (
	"crypto/sha1"
	"errors"
	"github.com/hati-sh/hati/common"
	"github.com/hati-sh/hati/common/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"strconv"
	"sync"
)

type HddShardMap []*HddShard

type HddShard struct {
	db      *leveldb.DB
	dataDir string
	counter int
	sync.RWMutex
}

var hddWriteOptions = &opt.Options{
	WriteBuffer: common.STORAGE_HDD_WRITE_BUFFER,
}

func newHddShardMap(size int, dataDir string) HddShardMap {
	var err error
	m := make([]*HddShard, size)

	for i := 0; i < size; i++ {
		m[i] = &HddShard{db: nil, dataDir: dataDir}

		m[i].db, err = common.OpenDatabase(dataDir, "kv_shard_"+strconv.Itoa(i), hddWriteOptions)
		if err != nil {
			panic(err)
		}

		iter := m[i].db.NewIterator(nil, nil)
		for iter.Next() {
			m[i].counter++
		}
		iter.Release()
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

func (m HddShardMap) Set(key string, value []byte) bool {
	shard := m.GetShard(key)

	_, err := shard.db.Get([]byte(key), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			shard.Lock()
			shard.counter++
			shard.Unlock()
		}
	}

	if err := shard.db.Put([]byte(key), value, nil); err != nil {
		logger.Error(err.Error())
		return false
	}
	return true
}

func (m HddShardMap) Delete(key string) bool {
	shard := m.GetShard(key)
	if err := shard.db.Delete([]byte(key), nil); err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return false
		}

		logger.Error(err.Error())
		return false
	}

	if shard.counter > 0 {
		shard.Lock()
		shard.counter--
		shard.Unlock()
	}

	return true
}

func (m HddShardMap) FlushAll() (bool, error) {
	// go shard by shard and delete data
	for i, shard := range m {
		m[i].Lock()
		_ = shard.db.Close()

		dbName := "kv_shard_" + strconv.Itoa(i)
		if err := common.DeleteDatabase(shard.dataDir, dbName); err != nil {
			logger.Error(err.Error())
			return false, err
		}

		var err error
		m[i].counter = 0

		m[i].db, err = common.OpenDatabase(shard.dataDir, dbName, hddWriteOptions)
		if err != nil {
			panic(err)
		}

		m[i].Unlock()
	}

	return true, nil
}

func (m HddShardMap) CountKeys() int {
	var keysCount = 0
	for _, shard := range m {
		shard.RLock()
		keysCount = keysCount + shard.counter
		shard.RUnlock()
	}

	return keysCount
}
