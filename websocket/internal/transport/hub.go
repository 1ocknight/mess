package transport

import "github.com/TATAROmangol/mess/websocket/internal/model"

type Hub struct {
	clients     map[string]*Client
	register    chan *Client
	unregister  chan *Client
	messageChan chan *model.Message
}

func NewHub(messageChan chan *model.Message) *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		messageChan: messageChan,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.SubjectID] = client

		case client := <-h.unregister:
			if client, ok := h.clients[client.SubjectID]; ok {
				delete(h.clients, client.SubjectID)
				close(client.Send)
			}

		case message := <-h.messageChan:
			client, ok := h.clients[message.SubjectID]
			if !ok {
				continue
			}

			select {
			case client.Send <- message.Data:
			default:
				delete(h.clients, client.SubjectID)
				close(client.Send)
			}
		}
	}
}
