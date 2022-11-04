package cache

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"time"
)

var _ Cache[string] = &MemCache[string]{}

type MemCache[V any] struct {
	client                *memcache.Client
	expiration            time.Duration
	overrideExistingCache bool
}

const (
	memcacheAddrsFlag      = "memcache-addrs"
	memcacheExpirationFlag = "memcaceh-expiration"
	memcacheOverrideFlag   = "memcache-overide-existing-cache"
)

func MemCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("mem-cache", pflag.ExitOnError)
	fs.StringSlice(memcacheAddrsFlag, []string{}, "")
	fs.Duration(memcacheExpirationFlag, 5*time.Minute, "")
	fs.Bool(memcacheOverrideFlag, false, "")
	return fs
}
func NewMemCacheFromViper[V any]() MemCache[V] {
	return MemCache[V]{
		client:                memcache.New(viper.GetStringSlice(memcacheAddrsFlag)...),
		expiration:            viper.GetDuration(memcacheExpirationFlag),
		overrideExistingCache: viper.GetBool(memcacheOverrideFlag),
	}
}
func NewMemCache[V any](client *memcache.Client, overrideExistingCache bool) MemCache[V] {
	return MemCache[V]{
		client:                client,
		overrideExistingCache: overrideExistingCache,
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
	if !g.overrideExistingCache {
		if _, err := g.Get(ctx, key); err != ErrCacheMiss {
			return ErrCacheForKeyExists
		}
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return g.client.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(g.expiration.Seconds()),
	})
}
