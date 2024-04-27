package core

import (
	"github.com/hati-sh/hati/storage"
)

type RpcStorageService struct {
	storage *storage.Storage
}

type Args struct{}

func (s *RpcStorageService) Count(args *Args, reply *int) error {
	// Fill reply pointer to send the data back
	count, _ := s.storage.Count(storage.Memory)
	*reply = count

	return nil
}
