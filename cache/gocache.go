package cache

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Cache[string] = &GoCache[string]{}

type GoCache[V any] struct {
	cache             *cache.Cache
	defaultExpiration time.Duration
}

const (
	goCacheDefaultExpirationFlag = "gocache-default-expiration"
	goCacheCleanUpIntervalFlag   = "gocache-cleanup-interval"
)

func GoCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("go-cache", pflag.ExitOnError)
	fs.Duration(goCacheDefaultExpirationFlag, 1*time.Minute, "")
	fs.Duration(goCacheCleanUpIntervalFlag, 1*time.Minute, "")
	return fs
}

func NewGoCacheFromViper[V any]() GoCache[V] {
	return NewGoCache[V](cache.New(viper.GetDuration(goCacheDefaultExpirationFlag), viper.GetDuration(goCacheCleanUpIntervalFlag)), viper.GetDuration(goCacheDefaultExpirationFlag))
}

func NewGoCache[V any](gocache *cache.Cache, defaultExpiration time.Duration) GoCache[V] {
	return GoCache[V]{
		cache:             gocache,
		defaultExpiration: defaultExpiration,
	}
}
func (g *GoCache[V]) Ping() error {
	return nil
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
