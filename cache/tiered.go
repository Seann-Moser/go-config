package cache

import (
	"context"
	"github.com/spf13/pflag"
)

var _ Cache[string, string] = &Tiered[string, string]{}

type Tiered[K string, V any] struct {
}

func TieredFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("tiered-cache", pflag.ExitOnError)
	return fs
}

func NewTieredCache[K string, V any]() Tiered[K, V] {
	return Tiered[K, V]{}
}

func (g *Tiered[K, V]) Get(ctx context.Context, key K) (V, error) {
	var output V
	return output, nil
}

func (g *Tiered[K, V]) Set(ctx context.Context, key K, value V) error {
	return nil
}
