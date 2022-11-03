package cache

import (
	"context"
	"github.com/spf13/pflag"
)

var _ Cache[string] = &MemCache[string]{}

type MemCache[V any] struct {
}

func MemCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("mem-cache", pflag.ExitOnError)
	return fs
}

func NewMemCache[V any]() MemCache[V] {
	return MemCache[V]{}
}

func (g *MemCache[V]) Get(ctx context.Context, key string) (V, error) {
	var output V
	return output, nil
}

func (g *MemCache[V]) Set(ctx context.Context, key string, value V) error {
	return nil
}
