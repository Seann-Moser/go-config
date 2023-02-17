package queue

import "context"

type Queue[V any] interface {
	Produce(ctx context.Context, msg V) error
	ProduceBatch(ctx context.Context, msg ...V) error
	Consume(ctx context.Context, topic string) (chan V, error)
	Close() error
}
