package cache

import (
	"context"
	"errors"
)

type Cache[K string, V any] interface {
	Getter[K, V]
	Set(ctx context.Context, key K, value V) error
}

type Getter[K string, V any] interface {
	Get(ctx context.Context, key K) (V, error)
}

var (
	ErrCacheMiss = errors.New("cache miss")
)
