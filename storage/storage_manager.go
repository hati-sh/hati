package storage

import (
	"context"
	"errors"
	"github.com/hati-sh/hati/common/logger"
	"sync"
)

// StorageManager is responsible for managing storages
// it offers two storage types: MEMORY and HDD
type Manager struct {
	ctx    context.Context
	memory *memoryStorage
	hdd    *hddStorage
	sync.WaitGroup
	stopChan chan bool
}

func NewStorageManager(ctx context.Context, dataDir string) *Manager {
	sm := &Manager{
		ctx:      ctx,
		stopChan: make(chan bool),
	}

	sm.memory = NewMemoryStorage(sm.ctx)
	sm.hdd = NewHddStorage(sm.ctx, dataDir)

	return sm
}

func (s *Manager) Start() {
	s.Add(1)
	go func(sm *Manager) {
		select {
		case <-sm.ctx.Done():
			s.hdd.Stop()
			break
		case <-sm.stopChan:
			s.hdd.Stop()
			break
		}
		sm.Done()
	}(s)
}

func (s *Manager) Stop() error {
	s.stopChan <- true
	s.Wait()

	logger.Debug("storage manager stopped")

	return nil
}

func (s *Manager) WaitForStop() {
	s.Wait()
	logger.Debug("storage manager stopped")
}

func (s *Manager) Count(storageType Type) (int, error) {
	if storageType == Memory {
		return s.memory.CountKeys(), nil
	}
	if storageType == Hdd {
		return s.hdd.CountKeys(), nil
	}

	return 0, nil
}

func (s *Manager) Set(storageType Type, key []byte, value []byte) error {
	if storageType == Memory && s.memory.Set(key, value) {
		return nil
	} else if storageType == Hdd && s.hdd.Set(key, value) {
		return nil
	}

	return errors.New("INVALID_STORAGE_TYPE")
}

func (s *Manager) Get(storageType Type, key []byte) ([]byte, error) {
	if storageType == Memory {
		value, err := s.memory.Get(key)
		if err != nil {
			return nil, err
		}

		return value, nil
	} else if storageType == Hdd {
		value, err := s.hdd.Get(key)
		if err != nil {
			return nil, err
		}

		return value, nil
	}

	return nil, errors.New("")
}

func (s *Manager) Has(storageType Type, key []byte) bool {
	if storageType == Memory {
		return s.memory.Has(key)
	} else if storageType == Hdd {
		return s.hdd.Has(key)
	}

	return false
}

func (s *Manager) Delete(storageType Type, key []byte) bool {
	switch storageType {
	case Memory:
		s.memory.Delete(key)
		return true
	case Hdd:
		s.hdd.Delete(key)
		return true
	default:
		return false
	}
}

func (s *Manager) FlushAll(storageType Type) (bool, error) {
	switch storageType {
	case Memory:
		return s.memory.FlushAll()
	case Hdd:
		return s.hdd.FlushAll()
	default:
		return false, errors.New("INVALID_STORAGE_TYPE")
	}
}
