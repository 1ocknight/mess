package transport

import (
	"github.com/1ocknight/mess/shared/logger"
	"github.com/1ocknight/mess/websocket/internal/loglables"
	"github.com/1ocknight/mess/websocket/internal/model"
)

type Hub struct {
	lg logger.Logger

	clients    map[string]map[*Client]struct{}
	register   chan *Client
	unregister chan *Client

	messageChan chan *model.Message
}

func NewHub(messageChan chan *model.Message, lg logger.Logger) *Hub {
	return &Hub{
		lg: lg,

		clients:    make(map[string]map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		messageChan: messageChan,
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			if _, ok := h.clients[client.SubjectID]; !ok {
				h.clients[client.SubjectID] = make(map[*Client]struct{})
			}
			h.clients[client.SubjectID][client] = struct{}{}
			h.lg.With(loglables.Subject, client.SubjectID).Info("register")

		case client := <-h.unregister:
			if clients, ok := h.clients[client.SubjectID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					h.lg.With(loglables.Subject, client.SubjectID).Info("unregister")
				}
				if len(clients) == 0 {
					delete(h.clients, client.SubjectID)
				}
			}

		case message := <-h.messageChan:
			clients, ok := h.clients[message.SubjectID]
			if !ok {
				continue
			}
			for c := range clients {
				select {
				case c.Send <- message.WSMessage:
					h.lg.With(loglables.Subject, c.SubjectID).Info("send message")
				default:
					delete(clients, c)
					close(c.Send)
				}
			}
			if len(clients) == 0 {
				delete(h.clients, message.SubjectID)
			}
		}
	}
}
