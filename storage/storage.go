package storage

import (
	"context"
	"errors"
)

const DEFAULT_NUMBER_OF_SHARDS = 20

type Type string

const Memory Type = "memory"
const Hdd Type = "hdd"

type Storage struct {
	ctx    context.Context
	Memory *memoryStorage
}

func New(ctx context.Context) Storage {
	return Storage{
		ctx:    ctx,
		Memory: NewMemoryStorage(),
	}
}

func (s *Storage) Set(storageType Type, key []byte, value []byte) error {
	if storageType == Memory && s.Memory.Set(key, value) {
		return nil
	}

	return errors.New("")
}

func (s *Storage) Get(storageType Type, key []byte) ([]byte, error) {
	if storageType == Memory {
		value, err := s.Memory.Get(key)
		if err != nil {
			return nil, err
		}

		return value, nil
	}

	return nil, errors.New("")
}
