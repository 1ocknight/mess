package kafka

import (
	"context"
	"time"

	"github.com/TATAROmangol/mess/shared/messagequeue"
	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg ProducerConfig) messagequeue.Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, key []byte, val []byte) error {
	kMsg := kafka.Message{
		Key:   key,
		Value: val,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, kMsg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
