package cache

import (
	"context"
	"fmt"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var _ Cache[string] = &Tiered[string]{}

type Tiered[V any] struct {
	caches []Cache[V]
	logger *zap.Logger
}

func NewTieredCache[V any](logger *zap.Logger, caches ...Cache[V]) Tiered[V] {
	return Tiered[V]{
		caches: caches,
		logger: logger,
	}
}

func (g *Tiered[V]) Ping() error {
	for _, c := range g.caches {
		if err := c.Ping(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Tiered[V]) Get(ctx context.Context, key string) (V, error) {
	var output V
	var err error
	defer func() {
		if err != nil {
			return
		}
		for _, c := range g.caches {
			if err = c.Set(ctx, key, output); err != nil {
				g.logger.Warn("failed setting cache", zap.Any("key", key), zap.Error(err))
			}
		}
	}()
	for _, getter := range g.caches {
		output, err = getter.Get(ctx, key)
		if err == ErrCacheMiss {
			continue
		}
		if err != ErrCacheMiss && err != nil {
			return output, err
		}
		return output, nil
	}
	return output, ErrCacheMiss
}

func (g *Tiered[V]) Set(ctx context.Context, key string, value V) error {
	var combinedErr error
	for _, c := range g.caches {
		if err := c.Set(ctx, key, value); err != nil {
			combinedErr = multierr.Append(combinedErr, fmt.Errorf("failed setting cache: key: %s : %w", key, err))
		}
	}
	return combinedErr
}
