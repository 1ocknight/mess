package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/1ocknight/mess/shared/auth"
	"github.com/1ocknight/mess/shared/logger"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	cfg         HTTPConfig
	Router      *mux.Router
	AuthService auth.Service
	lg          logger.Logger
	httpServer  *http.Server
}

func NewServer(cfg HTTPConfig, authService auth.Service, handler *Handler, lg logger.Logger) *Server {
	r := mux.NewRouter()

	s := &Server{
		cfg:         cfg,
		Router:      r,
		AuthService: authService,
		lg:          lg,
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000 "}, // фронт
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()), // rs/cors требует int секунд
	})
	r.Use(c.Handler)
	r.Use(SubjectMiddleware(authService, lg))

	// WS endpoint
	r.HandleFunc("/ws", handler.WSHandler)

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
