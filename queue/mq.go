package queue

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/Seann-Moser/go-config/flags"
)

var _ Queue[any] = &MQ[any]{}

const (
	mqAddressFlag      = "mq-address"
	mqQueueFlag        = "mq-queue"
	mqPortFlag         = "mq-port"
	mqBatchSizeFlag    = "mq-batch-size"
	mqBatchTimeoutFlag = "mq-batch-timeout"
)

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("mq", pflag.ExitOnError)
	fs.String(mqAddressFlag, "", "")
	fs.String(mqQueueFlag, "default", "")
	fs.Int(mqPortFlag, 50000, "")
	fs.Int(mqBatchSizeFlag, 10, "")
	fs.Duration(mqBatchTimeoutFlag, 10*time.Second, "")
	return fs
}

type MQ[V any] struct {
	logger  *zap.Logger
	qc      *kubemq.QueuesClient
	channel string

	batchSize    int
	batchTimeout time.Duration
	msgBatch     []*kubemq.QueueMessage

	batchMutex *sync.Mutex
}

func NewMQ[V any](ctx context.Context, logger *zap.Logger) (*MQ[V], error) {
	address, err := flags.RequiredString(mqAddressFlag)
	if err != nil {
		return nil, err
	}
	port, err := flags.RequiredInt(mqPortFlag)
	if err != nil {
		return nil, err
	}

	batchSize, err := flags.RequiredInt(mqBatchSizeFlag)
	if err != nil {
		return nil, err
	}
	batchTimeout, err := flags.RequiredDuration(mqBatchTimeoutFlag)
	if err != nil {
		return nil, err
	}

	qc, err := kubemq.NewQueuesStreamClient(ctx,
		kubemq.WithAddress(address, port),
		kubemq.WithClientId("stream-queue-sender"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		return nil, err
	}

	return &MQ[V]{
		qc:           qc,
		batchSize:    batchSize,
		logger:       logger.With(zap.String("queue_service", "MQ")),
		channel:      viper.GetString(mqQueueFlag),
		batchTimeout: batchTimeout,
		msgBatch:     []*kubemq.QueueMessage{},
		batchMutex:   &sync.Mutex{},
	}, nil
}
func (M *MQ[V]) batchTicker(ctx context.Context) {
	go func() {
		tick := time.NewTicker(M.batchTimeout)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				err := M.SendBatch(ctx)
				if err != nil {
					M.logger.Error("failed to send batch", zap.Error(err))
				}
			}
		}
	}()
}
func (M *MQ[V]) Produce(ctx context.Context, msg V) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	queueMsg := kubemq.NewQueueMessage()
	queueMsg.SetBody(b)
	queueMsg.SetChannel(M.channel)
	_, err = M.qc.Send(ctx, queueMsg)
	return err
}

func (M *MQ[V]) ProduceBatch(ctx context.Context, msgs ...V) error {
	for _, msg := range msgs {
		queueMsg := kubemq.NewQueueMessage()
		b, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		queueMsg.SetBody(b)
		queueMsg.SetChannel(M.channel)
		M.batchMutex.Lock()
		M.msgBatch = append(M.msgBatch, queueMsg)
		M.batchMutex.Unlock()
		if len(M.msgBatch) >= M.batchSize {
			_ = M.SendBatch(ctx)
		}
	}

	return nil
}

func (M *MQ[V]) SendBatch(ctx context.Context) error {
	if len(M.msgBatch) == 0 {
		return nil
	}
	defer func() {
		M.msgBatch = []*kubemq.QueueMessage{}
		M.batchMutex.Unlock()
	}()
	M.batchMutex.Lock()
	_, err := M.qc.Batch(ctx, M.msgBatch)
	return err
}
func (M *MQ[V]) Consume(ctx context.Context, topic string) (chan V, error) {
	if topic == "" {
		topic = M.channel
	}
	data := make(chan V, 10)
	_, err := M.qc.Subscribe(ctx, kubemq.NewReceiveQueueMessagesRequest().SetChannel(topic).SetWaitTimeSeconds(15), func(response *kubemq.ReceiveQueueMessagesResponse, err error) {
		if err != nil {
			M.logger.Fatal("failed pulling from queue", zap.Error(err))
		}
		for _, msg := range response.Messages {
			var output V
			err := json.Unmarshal(msg.GetBody(), &output)
			if err != nil {
				M.logger.Error("failed unmarshalling data", zap.Error(err))
				if err := msg.Ack(); err != nil {
					M.logger.Error("failed acking msg", zap.Error(err))
				}
				continue
			}
			data <- output
			if err := msg.Ack(); err != nil {
				M.logger.Error("failed acking msg", zap.Error(err))
			}
		}
	})

	return data, err
}

func (M *MQ[V]) Close() error {
	return M.Close()
}
