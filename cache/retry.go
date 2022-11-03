package cache

import (
	"context"
	"github.com/spf13/pflag"
)

var _ Cache[string, string] = &RetryCache[string, string]{}

type RetryCache[K string, V any] struct {
}

func RetryCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("retry-cache", pflag.ExitOnError)
	return fs
}

func NewRetryCache[K string, V any]() RetryCache[K, V] {
	return RetryCache[K, V]{}
}

func (g *RetryCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	var output V
	return output, nil
}

func (g *RetryCache[K, V]) Set(ctx context.Context, key K, value V) error {
	return nil
}
