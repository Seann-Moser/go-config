package cache

import (
	"context"
	"github.com/spf13/pflag"
)

var _ Cache[string, string] = &GoCache[string, string]{}

type GoCache[K string, V any] struct {
}

func GoCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("go-cache", pflag.ExitOnError)
	return fs
}

func NewGoCache[K string, V any]() GoCache[K, V] {
	return GoCache[K, V]{}
}

func (g *GoCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	var output V
	return output, nil
}

func (g *GoCache[K, V]) Set(ctx context.Context, key K, value V) error {
	return nil
}
