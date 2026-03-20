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
	cfg := config.Load()

	db, err := repository.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	rdb := repository.NewRedis(cfg)
	defer rdb.Close()

	repo := repository.NewRepository(db, rdb)
	svc := service.NewService(repo, cfg)
	h := handler.NewHandler(svc)

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// API 路由
	api := r.Group("/api/v1")
	{
		api.POST("/video/generate", h.GenerateVideo)
		api.POST("/video/download", h.DownloadVideo)
		api.POST("/video/recreate", h.RecreateVideo)
		api.POST("/video/publish", h.PublishVideo)
		api.GET("/task/:id", h.GetTask)
		api.GET("/tasks", h.ListTasks)
		api.GET("/config/ai-providers", h.GetAIProviders)
		api.GET("/config/video-providers", h.GetVideoProviders)
	}

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
