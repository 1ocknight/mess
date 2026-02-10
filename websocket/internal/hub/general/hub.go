package general

import (
	"context"
	"fmt"

	"github.com/1ocknight/mess/shared/logger"
	"github.com/1ocknight/mess/shared/redisclient"
	"github.com/1ocknight/mess/websocket/internal/loglables"
	"github.com/1ocknight/mess/websocket/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis  redisclient.Config `yaml:"redis"`
	Client ClientConfig       `yaml:"client"`
}

type Hub struct {
	lg logger.Logger

	clients    map[string]map[*Client]struct{}
	register   chan *Client
	unregister chan *Client

	redis *redis.PubSub
	cfg   Config
}

func NewHub(ctx context.Context, lg logger.Logger, cfg Config) *Hub {
	redis := redisclient.NewClient(cfg.Redis)
	ps := redis.PSubscribe(ctx)

	return &Hub{
		lg:  lg,
		cfg: cfg,

		clients:    make(map[string]map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		redis: ps,
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) registerClient(ctx context.Context, client *Client) {
	if _, ok := h.clients[client.SubjectID]; !ok {
		h.clients[client.SubjectID] = make(map[*Client]struct{})
	}
	h.clients[client.SubjectID][client] = struct{}{}
	h.lg.With(loglables.Subject, client.SubjectID).Info("register")

	if err := h.redis.PSubscribe(ctx, redisclient.BuildListenChannel(&client.SubjectID, nil, nil)); err != nil {
		h.lg.Error(fmt.Errorf("redis subscribe: %w", err))
	}
}

func (h *Hub) unregisterClient(ctx context.Context, client *Client) {
	clients, ok := h.clients[client.SubjectID]
	if !ok {
		return
	}
	if _, exists := clients[client]; exists {
		delete(clients, client)
		close(client.send)
		h.lg.With(loglables.Subject, client.SubjectID).Info("unregister")
	}
	if len(clients) == 0 {
		delete(h.clients, client.SubjectID)
		if err := h.redis.PUnsubscribe(ctx, redisclient.BuildListenChannel(&client.SubjectID, nil, nil)); err != nil {
			h.lg.Error(fmt.Errorf("redis unsubscribe: %w", err))
		}
	}
}

func (h *Hub) sendMessage(rMsg *redis.Message) {
	subjID, _, msg, err := utils.RedisMessageToWSMessage(rMsg)
	if err != nil {
		h.lg.Error(fmt.Errorf("redis message to ws message: %w", err))
		return
	}

	clients, ok := h.clients[subjID]
	if !ok {
		return
	}

	for c := range clients {
		select {
		case c.send <- msg:
			h.lg.With(loglables.Subject, c.SubjectID).Info("send message")
		default:
			delete(clients, c)
			close(c.send)
		}
	}
	if len(clients) == 0 {
		delete(h.clients, subjID)
	}
}

func (h *Hub) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				h.lg.Info("hub context done, shutting down")
				return

			case client := <-h.register:
				h.registerClient(ctx, client)

			case client := <-h.unregister:
				h.unregisterClient(ctx, client)

			case rMsg := <-h.redis.Channel():
				h.sendMessage(rMsg)
			}
		}
	}()
}
