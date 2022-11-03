package cache

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/spf13/pflag"
)

var _ Cache[string] = &MemCache[string]{}

type MemCache[V any] struct {
	client *memcache.Client
}

func MemCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("mem-cache", pflag.ExitOnError)
	return fs
}

func NewMemCache[V any](client *memcache.Client) MemCache[V] {
	return MemCache[V]{
		client: client,
	}
}

func (g *MemCache[V]) Get(ctx context.Context, key string) (V, error) {
	var output V
	data, err := g.client.Get(key)
	if err != nil && err != ErrCacheMiss {
		return output, err
	}
	if err == ErrCacheMiss {
		return output, ErrCacheMiss
	}
	err = json.Unmarshal(data.Value, &output)
	return output, err
}

func (g *MemCache[V]) Set(ctx context.Context, key string, value V) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return g.client.Set(&memcache.Item{
		Key:   key,
		Value: data,
	})
}
