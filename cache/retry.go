package cache

import (
	"context"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"time"
)

var _ Cache[string, string] = &RetryCache[string, string]{}

type RetryCache[K string, V any] struct {
	timeoutDuration time.Duration
	cache           Cache[K, V]

	maxRetry       uint64
	maxInterval    time.Duration
	maxElapsedTime time.Duration
	logger         *zap.Logger
}

func RetryCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("retry-cache", pflag.ExitOnError)
	return fs
}

func NewRetryCache[K string, V any](cache Cache[K, V], maxRetry uint64,
	maxInterval,
	maxElapsedTime,
	contextTimeout time.Duration, logger *zap.Logger) RetryCache[K, V] {
	return RetryCache[K, V]{
		timeoutDuration: contextTimeout,
		cache:           cache,
		maxRetry:        maxRetry,
		maxInterval:     maxInterval,
		maxElapsedTime:  maxElapsedTime,
		logger:          logger,
	}
}

func (g *RetryCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	var (
		err    error
		output V
	)
	op := backoff.Operation(func() error {
		c, cancel := context.WithTimeout(ctx, g.timeoutDuration)
		defer cancel()
		output, err = g.cache.Get(c, key)
		if err != nil {
			return err
		}

		return nil
	})
	notify := func(err error, backoffDuration time.Duration) {
		g.logger.Info("retrying get", zap.Error(err), zap.Duration("backoff_duration", backoffDuration))
	}

	if err = backoff.RetryNotify(op, g.getBackoff(), notify); err != nil {
		return output, err
	}

	return output, nil
}

func (g *RetryCache[K, V]) Set(ctx context.Context, key K, value V) error {
	op := backoff.Operation(func() error {
		c, cancel := context.WithTimeout(ctx, g.timeoutDuration)
		defer cancel()
		return g.cache.Set(c, key, value)
	})
	notify := func(err error, backoffDuration time.Duration) {
		g.logger.Info("retrying get", zap.Error(err), zap.Duration("backoff_duration", backoffDuration))
	}

	return backoff.RetryNotify(op, g.getBackoff(), notify)
}

func (g *RetryCache[K, V]) getBackoff() backoff.BackOff {
	requestExpBackOff := backoff.NewExponentialBackOff()
	requestExpBackOff.InitialInterval = 2 * time.Millisecond
	requestExpBackOff.RandomizationFactor = 0.5
	requestExpBackOff.Multiplier = 1.5
	requestExpBackOff.MaxInterval = g.maxInterval
	requestExpBackOff.MaxElapsedTime = g.maxElapsedTime
	return backoff.WithMaxRetries(requestExpBackOff, g.maxRetry)
}
