package rest

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port int
}

type Server struct {
	router *gin.Engine
}

func NewServer(cfg Config) *Server {
	return &Server{}
}

// http://localhost:8090/exchange?session_state=377efe15-1762-4d06-b795-04d0739c4ae0&code=6d3aa28a-9852-456a-9897-a8b2c0b3d262.377efe15-1762-4d06-b795-04d0739c4ae0.29e1c1cd-84e8-4a79-8e79-251529228195