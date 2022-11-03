package cache

import (
	"context"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var _ Cache[string, string] = &Tiered[string, string]{}

type Tiered[K string, V any] struct {
	cache   Cache[K, V]
	getters []Cache[K, V]
	logger  *zap.Logger
}

func TieredFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("tiered-cache", pflag.ExitOnError)
	return fs
}

func NewTieredCache[K string, V any](logger *zap.Logger, cache Cache[K, V], getters ...Cache[K, V]) Tiered[K, V] {
	return Tiered[K, V]{
		cache:   cache,
		getters: getters,
		logger:  logger,
	}
}

func (g *Tiered[K, V]) Get(ctx context.Context, key K) (V, error) {
	var output V
	var err error
	defer func() {
		if err != nil {
			return
		}
		for _, getter := range g.getters {
			if err = getter.Set(ctx, key, output); err != nil {
				g.logger.Warn("failed setting cache", zap.Any("key", key), zap.Error(err))
			}
		}
	}()
	for _, getter := range g.getters {
		if output, err = getter.Get(ctx, key); err == nil {
			return output, nil
		}
	}
	output, err = g.cache.Get(ctx, key)
	if err != nil {
		return output, err
	}

	return output, err
}

func (g *Tiered[K, V]) Set(ctx context.Context, key K, value V) error {
	return nil
}
