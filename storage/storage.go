package storage

import "context"

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
