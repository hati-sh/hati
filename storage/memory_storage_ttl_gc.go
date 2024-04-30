package storage

import (
	"fmt"
	"github.com/hati-sh/hati/common/logger"
	"sync"
	"time"
)

// memoryStorageTtlGc is a garbage collector for messages with TTL in memory storage
type memoryStorageTtlGc struct {
	storage *memoryStorage
	ttlKeys map[int64][][]byte
	sync.WaitGroup
	sync.RWMutex
	stopChan chan bool
}

func NewMemoryStorageTtlGc(storage *memoryStorage) *memoryStorageTtlGc {
	return &memoryStorageTtlGc{
		storage:  storage,
		ttlKeys:  make(map[int64][][]byte),
		stopChan: make(chan bool),
	}
}

func (s *memoryStorageTtlGc) SetTtl(expiresAt int64, key []byte) {
	s.Lock()
	defer s.Unlock()

	expiresAt = time.Now().UnixMilli() + expiresAt

	if s.ttlKeys[expiresAt] == nil {
		s.ttlKeys[expiresAt] = make([][]byte, 0)
	}
	s.ttlKeys[expiresAt] = append(s.ttlKeys[expiresAt], key)
}

func (s *memoryStorageTtlGc) DeleteTtl(expiresAt int64, key []byte) {
	s.Lock()
	defer s.Unlock()

	delete(s.ttlKeys, expiresAt)
}

func (s *memoryStorageTtlGc) daemon() {
	ticker := time.NewTicker(time.Millisecond)

OuterLoop:
	for {
		select {
		case <-ticker.C:
			s.deleteExpiredKeys()
			continue
		case <-s.stopChan:
			break OuterLoop
		case <-s.storage.ctx.Done():
			break OuterLoop
		}
	}
	logger.Debug("stopping memory storage ttl gc")
	s.Done()
}

func (s *memoryStorageTtlGc) deleteExpiredKeys() {
	s.Lock()
	defer s.Unlock()

	if len(s.ttlKeys) < 1 {
		return
	}

	now := time.Now().UnixMilli()
	for expiresAt, keys := range s.ttlKeys {
		if expiresAt <= now {
			for _, key := range keys {
				s.storage.Delete(key)
			}
			delete(s.ttlKeys, expiresAt)
		}
	}
}

func (s *memoryStorageTtlGc) Start() {
	s.Add(1)

	go s.daemon()
}

func (s *memoryStorageTtlGc) Stop() {
	fmt.Println("memoryStorageTtlGc Stop")
	if s.stopChan != nil {
		close(s.stopChan)
	}
	s.Wait()
}
