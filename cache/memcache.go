package cache

import (
	"context"
	"github.com/spf13/pflag"
)

var _ Cache[string, string] = &MemCache[string, string]{}

type MemCache[K string, V any] struct {
}

func MemCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("mem-cache", pflag.ExitOnError)
	return fs
}

func NewMemCache[K string, V any]() MemCache[K, V] {
	return MemCache[K, V]{}
}

func (g *MemCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	var output V
	return output, nil
}

func (g *MemCache[K, V]) Set(ctx context.Context, key K, value V) error {
	return nil
}
