package cache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/pflag"
)

var _ Cache[string] = &GoCache[string]{}

type GoCache[V any] struct {
	cache cache.Cache
}

func GoCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("go-cache", pflag.ExitOnError)
	return fs
}

func NewGoCache[V any]() GoCache[V] {
	return GoCache[V]{}
}

func (g *GoCache[V]) Get(ctx context.Context, key string) (V, error) {
	var output V
	g.cache.Get(key)
	return output, nil
}

func (g *GoCache[V]) Set(ctx context.Context, key string, value V) error {
	return nil
}
