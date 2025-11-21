package server

import (
	"fmt"
	"log"

	"qiniu-uploader/internal/config"
	"qiniu-uploader/internal/routes"

	"github.com/gin-gonic/gin"
)

// Server 服务器结构体
type Server struct {
	config *config.Config
	router *gin.Engine
}

// NewServer 创建新的服务器实例
func NewServer(cfg *config.Config) *Server {
	// 设置Gin模式
	gin.SetMode(cfg.GinMode)

	router := gin.Default()

	return &Server{
		config: cfg,
		router: router,
	}
}

// Setup 设置服务器
func (s *Server) Setup() {
	// 设置路由
	routes.SetupRoutes(s.router, s.config)
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("服务器启动在 http://localhost%s", addr)
	return s.router.Run(addr)
}