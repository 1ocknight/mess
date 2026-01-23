package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/TATAROmangol/mess/profile/internal/domain"
	"github.com/TATAROmangol/mess/shared/auth"
	"github.com/TATAROmangol/mess/shared/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	DebugMode bool   `yaml:"debug_mode"`
}

type HTTPServer struct {
	cfg    *Config
	srv    *gin.Engine
	httpSv *http.Server
}

func NewServer(cfg Config, lg logger.Logger, domain domain.Service, auth auth.Service) *HTTPServer {
	h := NewHandler(domain)

	if !cfg.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // фронт
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(InitLoggerMiddleware(lg))
	r.Use(SetRequestMetadataMiddleware())
	r.Use(LogResponseMiddleware())
	r.Use(InitSubjectMiddleware(auth))

	r.GET("/profile", h.GetProfile)
	r.GET("/profile/:id", h.GetProfile)

	r.GET("/profiles", h.GetProfiles)

	r.POST("/profile", h.AddProfile)

	r.PUT("/profile", h.UpdateProfileMetadata)
	r.PUT("/avatar", h.UploadAvatar)

	r.DELETE("/avatar", h.DeleteAvatar)
	r.DELETE("/profile", h.DeleteProfile)

	return &HTTPServer{
		cfg: &cfg,
		srv: r,
	}
}

func (s *HTTPServer) Run() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)

	s.httpSv = &http.Server{
		Addr:    addr,
		Handler: s.srv,
	}

	return s.httpSv.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.httpSv.Shutdown(ctx)
}
