package core

import (
	"errors"

	"github.com/hati-sh/hati/storage"
)

type RpcStorageService struct {
	storageManager *storage.StorageManager
}

type CountArgs struct{}

type SetArgs struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetArgs struct {
	Key string `json:"key"`
}

func (s *RpcStorageService) Count(args *CountArgs, reply *int) error {
	count, _ := s.storageManager.Count(storage.Memory)
	*reply = count

	return nil
}

func (s *RpcStorageService) Set(args *SetArgs, reply *bool) error {
	if args.Key == "" {
		return errors.New("invalid key")
	}

	// Fill reply pointer to send the data back
	if err := s.storageManager.Set(storage.Memory, []byte(args.Key), []byte(args.Value)); err != nil {
		*reply = false

		return nil
	}

	*reply = true

	return nil
}

func (s *RpcStorageService) Get(args *SetArgs, reply *string) error {
	if args.Key == "" {
		return errors.New("invalid key")
	}

	value, err := s.storageManager.Get(storage.Memory, []byte(args.Key))
	if err != nil {
		*reply = err.Error()

		return nil
	}

	*reply = string(value)

	return nil
}
