package storage

import "context"

type Interface interface {
	GetAll(ctx context.Context, prefix string) ([][]byte, error)
	Get(ctx context.Context, prefix string, key string) ([]byte, error)
	Put(ctx context.Context, prefix string, key string, value []byte) error
	Delete(ctx context.Context, prefix string, key string) error
	Close() error
}
