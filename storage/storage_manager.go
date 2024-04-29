package storage

import (
	"context"
	"errors"
)

// StorageManager is responsible for managing storages
// it offers two storage types: MEMORY and HDD
type StorageManager struct {
	ctx    context.Context
	memory *memoryStorage
	hdd    *hddStorage
}

func NewStorageManager(ctx context.Context) *StorageManager {
	sm := &StorageManager{
		ctx: ctx,
	}

	sm.memory = NewMemoryStorage(sm.ctx)
	sm.hdd = NewHddStorage(sm.ctx)

	return sm
}

func (s *StorageManager) Count(storageType Type) (int, error) {
	if storageType == Memory {
		return s.memory.CountKeys(), nil
	}

	return 0, nil
}

func (s *StorageManager) Set(storageType Type, key []byte, value []byte) error {
	if storageType == Memory && s.memory.Set(key, value) {
		return nil
	}

	return errors.New("")
}

func (s *StorageManager) Get(storageType Type, key []byte) ([]byte, error) {
	if storageType == Memory {
		value, err := s.memory.Get(key)
		if err != nil {
			return nil, err
		}

		return value, nil
	}

	return nil, errors.New("")
}

func (s *StorageManager) Has(storageType Type, key []byte) bool {
	if storageType == Memory && s.memory.Has(key) {
		return true
	}

	return false
}

func (s *StorageManager) Delete(storageType Type, key []byte) bool {
	switch storageType {
	case Memory:
		s.memory.Delete(key)
		return true
	case Hdd:
		return false
	default:
		return false
	}
}

func (s *StorageManager) FlushAll(storageType Type) bool {
	switch storageType {
	case Memory:
		return s.memory.FlushAll()
	case Hdd:
		return false
	default:
		return false
	}
}
