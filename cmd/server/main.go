package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kubegogo/genvideo/internal/config"
	"github.com/kubegogo/genvideo/internal/handler"
	"github.com/kubegogo/genvideo/internal/middleware"
	"github.com/kubegogo/genvideo/internal/repository"
	"github.com/kubegogo/genvideo/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := repository.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	rdb := repository.NewRedis(cfg)
	defer rdb.Close()

	// Initialize repositories
	repo := repository.NewRepository(db, rdb)

	// Initialize services
	svc := service.NewService(repo, cfg)

	// Initialize handlers
	h := handler.NewHandler(svc)

	// Setup Gin router
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Routes
	api := r.Group("/api/v1")
	{
		// 视频生成（核心功能）
		// 输入文案/关键词 → AI生成素材 → 自动剪辑 → 成片
		api.POST("/video/generate", h.GenerateVideo)

		// 视频搬运
		api.POST("/video/download", h.DownloadVideo)
		api.POST("/video/recreate", h.RecreateVideo)
		api.POST("/video/publish", h.PublishVideo)

		// 配置
		api.GET("/config/video-providers", h.GetVideoProviders)
		api.POST("/config/video-providers", h.SetVideoProvider)
		api.GET("/config/ai-providers", h.GetAIProviders)
		api.POST("/config/ai-providers", h.SetAIProvider)
		api.GET("/config/oss", h.GetOSSConfig)
		api.POST("/config/oss", h.SetOSSConfig)

		// 任务状态
		api.GET("/task/:id", h.GetTaskStatus)
		api.GET("/tasks", h.ListTasks)
	}

	// Health check
	r.GET("/health", h.Health)

	port := cfg.ServerPort
	if port == "" {
		port = "3004"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
