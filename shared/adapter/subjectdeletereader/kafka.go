package subjectdeletereader

import (
	"context"
	"fmt"

	mqdto "github.com/1ocknight/mess/shared/dto/mq"
	"github.com/IBM/sarama"
)

type MessageIMPL struct {
	ds      *mqdto.DeleteSubject
	message *sarama.ConsumerMessage
	session sarama.ConsumerGroupSession
}

func (m *MessageIMPL) GetSubjectID() string {
	return m.ds.GetSubjectID()
}

type handler struct {
	msgCh chan<- *MessageIMPL
}

func (h *handler) Setup(session sarama.ConsumerGroupSession) error { return nil }

func (h *handler) Cleanup(session sarama.ConsumerGroupSession) error { return nil }

func (h handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ds, err := mqdto.UnmarshallDeleteSubject(message.Value)
		if err != nil {
			return fmt.Errorf("unmarshal delete subject message: %w", err)
		}

		h.msgCh <- &MessageIMPL{
			ds:      ds,
			message: message,
			session: session,
		}
	}

	return nil
}

type Config struct {
	Brokers []string `yaml:"brokers"`
	Topics  []string `yaml:"topics"`
	GroupID string   `yaml:"group_id"`
}

type Consumer struct {
	cfg    Config
	client sarama.ConsumerGroup

	handler   *handler
	messageCh chan *MessageIMPL
	errCh     chan error
	cancel    context.CancelFunc
	done      chan struct{}
}

func New(cfg Config) (Service, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("new consumer group: %w", err)
	}

	msgCh := make(chan *MessageIMPL)
	errCh := make(chan error, 1)
	handler := &handler{
		msgCh: msgCh,
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &Consumer{
		cfg:       cfg,
		client:    client,
		handler:   handler,
		messageCh: msgCh,
		errCh:     errCh,
		cancel:    cancel,
		done:      make(chan struct{}),
	}

	go func() {
		defer close(c.done)
		for {
			if err := client.Consume(ctx, cfg.Topics, handler); err != nil {
				select {
				case errCh <- fmt.Errorf("consume: %w", err):
				default:
				}
				return
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return c, nil
}

func (c *Consumer) FetchMessage(ctx context.Context) (Message, error) {
	select {
	case msg := <-c.messageCh:
		return msg, nil
	case err := <-c.errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Consumer) Commit(msg Message) error {
	if impl, ok := msg.(*MessageIMPL); ok {
		impl.session.MarkMessage(impl.message, "")
		return nil
	}

	return fmt.Errorf("invalid message type")
}

func (c *Consumer) Close() error {
	c.cancel()
	<-c.done
	return c.client.Close()
}
