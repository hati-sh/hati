package core

import (
	"errors"
	"github.com/hati-sh/hati/storage"
)

type RpcStorageService struct {
	storageManager *storage.Manager
}

type CountArgs struct {
	Type storage.Type `json:"type"`
}

type SetArgs struct {
	Type  storage.Type `json:"type"`
	Key   string       `json:"key"`
	Value string       `json:"value"`
	Ttl   string       `json:"ttl"`
}

type GetArgs struct {
	Type storage.Type `json:"type"`
	Key  string       `json:"key"`
}

type DeleteArgs struct {
	Type storage.Type `json:"type"`
	Key  string       `json:"key"`
}

type FlushAllArgs struct {
	Type storage.Type `json:"type"`
}

func (s *RpcStorageService) Count(args *CountArgs, reply *int) error {
	if err := s.validateStorageType(args.Type); err != nil {
		return err
	}

	count, _ := s.storageManager.Count(args.Type)
	*reply = count

	return nil
}

func (s *RpcStorageService) Set(args *SetArgs, reply *bool) error {
	if err := s.validateStorageType(args.Type); err != nil {
		return err
	}

	if args.Key == "" {
		return errors.New("invalid key")
	}

	if err := s.storageManager.Set(args.Type, []byte(args.Key), []byte(args.Value), []byte(args.Ttl)); err != nil {
		*reply = false

		return nil
	}

	*reply = true

	return nil
}

func (s *RpcStorageService) Get(args *SetArgs, reply *string) error {
	if err := s.validateStorageType(args.Type); err != nil {
		return err
	}

	if args.Key == "" {
		return errors.New("invalid key")
	}

	value, err := s.storageManager.Get(args.Type, []byte(args.Key))
	if err != nil {
		*reply = err.Error()

		return nil
	}

	*reply = string(value)

	return nil
}

func (s *RpcStorageService) Delete(args *DeleteArgs, reply *bool) error {
	if err := s.validateStorageType(args.Type); err != nil {
		return err
	}

	if args.Key == "" {
		return errors.New("invalid key")
	}

	*reply = s.storageManager.Delete(args.Type, []byte(args.Key))

	return nil
}

func (s *RpcStorageService) FlushAll(args *FlushAllArgs, reply *bool) error {
	if err := s.validateStorageType(args.Type); err != nil {
		return err
	}

	res, err := s.storageManager.FlushAll(args.Type)
	if err != nil {
		return err
	}

	*reply = res

	return nil
}

func (s *RpcStorageService) validateStorageType(storageType storage.Type) error {
	if storageType == storage.Memory || storageType == storage.Hdd {
		return nil
	}

	return storage.ErrInvalidStorageType
}
