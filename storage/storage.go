package storage

type Type string

const Memory Type = "memory"
const Hdd Type = "hdd"

type Storage struct {
	Memory *memoryStorage
}

func New() Storage {
	return Storage{
		Memory: NewMemoryStorage(),
	}
}
