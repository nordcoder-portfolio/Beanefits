package kafka

import (
	"context"
	"errors"
	"strconv"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	w            *kafkago.Writer
	writeTimeout time.Duration
}

type ProducerConfig struct {
	Brokers      []string
	Topic        string
	WriteTimeout time.Duration
	BatchSize    int
	BatchTimeout time.Duration
}

func NewProducer(cfg ProducerConfig) (*Producer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka: brokers is empty")
	}
	if cfg.Topic == "" {
		return nil, errors.New("kafka: topic is empty")
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 1
	}

	w := &kafkago.Writer{
		Addr:         kafkago.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafkago.LeastBytes{},
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
		Async:        false,
	}

	return &Producer{
		w:            w,
		writeTimeout: cfg.WriteTimeout,
	}, nil
}

func (p *Producer) Close() error {
	if p == nil || p.w == nil {
		return nil
	}
	return p.w.Close()
}

func (p *Producer) SendInt(ctx context.Context, n int64) error {
	if p == nil || p.w == nil {
		return errors.New("kafka: producer not initialized")
	}

	ctx, cancel := withOptionalTimeout(ctx, p.writeTimeout)
	defer cancel()

	msg := kafkago.Message{
		Value: []byte(strconv.FormatInt(n, 10)),
	}

	return p.w.WriteMessages(ctx, msg)
}

func withOptionalTimeout(ctx context.Context, d time.Duration) (context.Context, func()) {
	if d <= 0 {
		return ctx, func() {}
	}
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, d)
}
