package cache

import (
	"context"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var _ Cache[string] = &RetryCache[string]{}

type RetryCache[V any] struct {
	timeoutDuration time.Duration
	cache           Cache[V]

	maxRetry       uint64
	maxInterval    time.Duration
	maxElapsedTime time.Duration
	logger         *zap.Logger
}

const (
	retryCacheTimeoutDurationFlag = "retry-cache-timeout-duration"
	retryCacheMaxRetriesFlag      = "retry-cache-max-retries"
	retryCacheMaxIntervalFlag     = "retry-cache-max-interval"
	retryCacheMaxElapsedTimeFlag  = "retry-cache-max-elapsed-time"
)

func RetryCacheFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("retry-cache", pflag.ExitOnError)
	fs.Duration(retryCacheTimeoutDurationFlag, 5*time.Second, "")
	fs.Duration(retryCacheMaxIntervalFlag, 20*time.Millisecond, "")
	fs.Duration(retryCacheMaxElapsedTimeFlag, 1*time.Minute, "")
	fs.Uint64(retryCacheMaxRetriesFlag, 3, "")
	return fs
}
func NewRetryCacheFromViper[V any](cache Cache[V], logger *zap.Logger) RetryCache[V] {
	return RetryCache[V]{
		timeoutDuration: viper.GetDuration(retryCacheTimeoutDurationFlag),
		cache:           cache,
		maxRetry:        viper.GetUint64(retryCacheMaxRetriesFlag),
		maxInterval:     viper.GetDuration(retryCacheMaxIntervalFlag),
		maxElapsedTime:  viper.GetDuration(retryCacheMaxElapsedTimeFlag),
		logger:          logger,
	}
}
func NewRetryCache[V any](cache Cache[V], maxRetry uint64,
	maxInterval,
	maxElapsedTime,
	contextTimeout time.Duration, logger *zap.Logger) RetryCache[V] {
	return RetryCache[V]{
		timeoutDuration: contextTimeout,
		cache:           cache,
		maxRetry:        maxRetry,
		maxInterval:     maxInterval,
		maxElapsedTime:  maxElapsedTime,
		logger:          logger,
	}
}

func (g *RetryCache[V]) Get(ctx context.Context, key string) (V, error) {
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

func (g *RetryCache[V]) Set(ctx context.Context, key string, value V) error {
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

func (g *RetryCache[V]) getBackoff() backoff.BackOff {
	requestExpBackOff := backoff.NewExponentialBackOff()
	requestExpBackOff.InitialInterval = 2 * time.Millisecond
	requestExpBackOff.RandomizationFactor = 0.5
	requestExpBackOff.Multiplier = 1.5
	requestExpBackOff.MaxInterval = g.maxInterval
	requestExpBackOff.MaxElapsedTime = g.maxElapsedTime
	return backoff.WithMaxRetries(requestExpBackOff, g.maxRetry)
}
