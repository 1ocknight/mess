package transport

import (
	"net/http"

	"github.com/1ocknight/mess/websocket/internal/ctxkey"
	"github.com/1ocknight/mess/websocket/internal/hub/general"
	"github.com/gorilla/websocket"
)

type HandlerConfig struct {
	ReadBufferSizeBytes  int  `yaml:"read_buffer_size_bytes"`
	WriteBufferSizeBytes int  `yaml:"write_buffer_size_bytes"`
	CheckOrigin          bool `yaml:"check_origin"`
}

type Handler struct {
	cfg        HandlerConfig
	generalHub *general.Hub
	upgrader   *websocket.Upgrader
}

func NewHandler(cfg HandlerConfig, generalHub *general.Hub) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  cfg.ReadBufferSizeBytes,
		WriteBufferSize: cfg.WriteBufferSizeBytes,
	}

	if !cfg.CheckOrigin {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	return &Handler{
		cfg:        cfg,
		generalHub: generalHub,
		upgrader:   &upgrader,
	}
}

func (h *Handler) General(w http.ResponseWriter, r *http.Request) {
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
	client := general.NewClient(subj.GetSubjectId(), conn, h.generalHub)
	h.generalHub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
