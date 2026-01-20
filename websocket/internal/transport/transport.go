package transport

import (
	"log"
	"net/http"

	"github.com/TATAROmangol/mess/shared/auth"
	"github.com/TATAROmangol/mess/shared/logger"
	"github.com/gorilla/mux"
)

type Server struct {
	Router      *mux.Router
	AuthService auth.Service
	Logger      logger.Logger
}

func NewServer(authService auth.Service, lg logger.Logger, handler *Handler) *Server {
	s := &Server{
		Router:      mux.NewRouter(),
		AuthService: authService,
		Logger:      lg,
	}

	s.Router.Use(LoggerMiddleware(lg))
	s.Router.Use(RequestMetadataMiddleware())
	s.Router.Use(LogResponseMiddleware())
	s.Router.Use(SubjectMiddleware(authService))

	return s
}

func (s *Server) Run(addr string) {
	log.Printf("server running at %s", addr)
	if err := http.ListenAndServe(addr, s.Router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
