package storage

import "context"

type hddStorage struct {
	ctx context.Context
}

func NewHddStorage(ctx context.Context) *hddStorage {
	return &hddStorage{
		ctx: ctx,
	}
}
