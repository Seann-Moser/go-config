package cache

import (
	"context"
	"errors"
)

type Cache[V any] interface {
	Getter[V]
	Set(ctx context.Context, key string, value V) error
}

type Getter[V any] interface {
	Get(ctx context.Context, key string) (V, error)
}

var (
	ErrCacheMiss = errors.New("cache miss")
)
