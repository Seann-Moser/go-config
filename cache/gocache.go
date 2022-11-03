package cache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/pflag"
	"time"
)

var _ Cache[string] = &GoCache[string]{}

type GoCache[V any] struct {
	cache             *cache.Cache
	defaultExpiration time.Duration
}

func GoCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("go-cache", pflag.ExitOnError)
	return fs
}

func NewGoCache[V any](gocache *cache.Cache, defaultExpiration time.Duration) GoCache[V] {
	return GoCache[V]{
		cache:             gocache,
		defaultExpiration: defaultExpiration,
	}
}

func (g *GoCache[V]) Get(ctx context.Context, key string) (V, error) {
	var output V
	if data, found := g.cache.Get(key); found {
		return data.(V), nil
	}
	return output, ErrCacheMiss
}

func (g *GoCache[V]) Set(ctx context.Context, key string, value V) error {
	g.cache.Set(key, value, g.defaultExpiration)
	return nil
}
