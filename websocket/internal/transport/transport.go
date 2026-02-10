package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/1ocknight/mess/shared/logger"
	"github.com/1ocknight/mess/shared/verify"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type HTTPConfig struct {
	CorsUrl []string `yaml:"cors_url"`
	Host    string   `yaml:"host"`
	Port    string   `yaml:"port"`
}

type Server struct {
	cfg        HTTPConfig
	Router     *mux.Router
	ver        verify.Service
	lg         logger.Logger
	httpServer *http.Server
}

func NewServer(cfg HTTPConfig, ver verify.Service, handler *Handler, lg logger.Logger) *Server {
	r := mux.NewRouter()

	s := &Server{
		cfg:    cfg,
		Router: r,
		ver:    ver,
		lg:     lg,
	}

	if len(cfg.CorsUrl) != 0 {
		c := cors.New(cors.Options{
			AllowedOrigins:   cfg.CorsUrl,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
			ExposedHeaders:   []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           int((12 * time.Hour).Seconds()),
		})
		r.Use(c.Handler)
	}

	r.Use(SubjectMiddleware(ver, lg))

	r.HandleFunc("/ws", handler.General)
	// r.HandleFunc("/ws/chat/{chat_id}", handler.Chat)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		Handler: r,
	}

	return s
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
