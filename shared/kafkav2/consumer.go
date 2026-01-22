package kafkav2

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

type ConsumerConfig struct {
	Brokers       []string `yaml:"brokers"`
	Topics        []string `yaml:"topics"`
	MessagesLimit int      `yaml:"messages_limit"`
}

type ConsumerMessage struct {
	Value     []byte
	partition int32
	offset    int64
}

type Consumer struct {
	brokers []string
	topic   string

	client   sarama.Client
	consumer sarama.Consumer

	msgCh chan *ConsumerMessage
	wg    sync.WaitGroup
}

func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &Consumer{
		brokers:  brokers,
		topic:    topic,
		client:   client,
		consumer: consumer,
		msgCh:    make(chan *ConsumerMessage),
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return fmt.Errorf("failed to consume partition %d: %w", partition, err)
		}

		c.wg.Add(1)
		go func(pc sarama.PartitionConsumer, partition int32) {
			defer c.wg.Done()
			for {
				select {
				case msg := <-pc.Messages():
					c.msgCh <- &ConsumerMessage{
						Value:     msg.Value,
						partition: msg.Partition,
						offset:    msg.Offset,
					}
				case <-ctx.Done():
					pc.Close()
					return
				}
			}
		}(pc, partition)
	}

	go func() {
		c.wg.Wait()
		close(c.msgCh)
	}()

	return nil
}

func (c *Consumer) GetMessagesChan() chan *ConsumerMessage {
	return c.msgCh
}

func (c *Consumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		return err
	}
	return c.client.Close()
}
