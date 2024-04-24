package storage

type Type string

const Memory Type = "memory"
const Hdd Type = "hdd"

type Storage struct {
}

func New() Storage {
	return Storage{}
}
