package cache

import (
	"context"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var _ Cache[string] = &Tiered[string]{}

type Tiered[V any] struct {
	cache   Cache[V]
	getters []Cache[V]
	logger  *zap.Logger
}

func TieredFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("tiered-cache", pflag.ExitOnError)
	return fs
}

func NewTieredCache[V any](logger *zap.Logger, cache Cache[V], getters ...Cache[V]) Tiered[V] {
	return Tiered[V]{
		cache:   cache,
		getters: getters,
		logger:  logger,
	}
}

func (g *Tiered[V]) Get(ctx context.Context, key string) (V, error) {
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

func (g *Tiered[V]) Set(ctx context.Context, key string, value V) error {
	return nil
}
