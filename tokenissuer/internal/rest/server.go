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
