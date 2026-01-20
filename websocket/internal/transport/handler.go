package transport

import (
	"net/http"

	"github.com/TATAROmangol/mess/websocket/internal/ctxkey"
	"github.com/TATAROmangol/mess/websocket/internal/model"
	"github.com/gorilla/websocket"
)

type Handler struct {
	cfg      HandlerConfig
	hub      *Hub
	upgrader *websocket.Upgrader
}

func NewHandler(cfg HandlerConfig, messageChan chan *model.Message) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  cfg.ReadBufferSizeBytes,
		WriteBufferSize: cfg.WriteBufferSizeBytes,
	}

	hub := NewHub(messageChan)

	return &Handler{
		cfg:      cfg,
		hub:      hub,
		upgrader: &upgrader,
	}

}

func (h *Handler) WSHandler(w http.ResponseWriter, r *http.Request) {
	subj, err := ctxkey.ExtractSubject(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := NewClient(subj.GetSubjectId(), conn, h.cfg.ClientConfig, h.hub)
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
