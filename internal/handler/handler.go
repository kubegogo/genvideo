package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kubegogo/genvideo/internal/model"
	"github.com/kubegogo/genvideo/internal/service"
	"github.com/kubegogo/genvideo/pkg/response"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Health 健康检查
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GenerateVideo 生成视频（脚本转视频）
// POST /api/v1/video/generate
func (h *Handler) GenerateVideo(c *gin.Context) {
	var req model.VideoGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 设置默认值
	if req.InputType == "" {
		req.InputType = "keywords"
	}
	if req.Style == "" {
		req.Style = "documentary"
	}
	if req.Duration == 0 {
		req.Duration = 60
	}
	if req.AspectRatio == "" {
		req.AspectRatio = "16:9"
	}

	task, err := h.svc.GenerateVideo(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, task)
}

// DownloadVideo 下载视频
// POST /api/v1/video/download
func (h *Handler) DownloadVideo(c *gin.Context) {
	var req model.DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.svc.DownloadVideo(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, task)
}

// RecreateVideo 二次创作
// POST /api/v1/video/recreate
func (h *Handler) RecreateVideo(c *gin.Context) {
	var req model.RecreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.svc.RecreateVideo(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, task)
}

// PublishVideo 发布视频
// POST /api/v1/video/publish
func (h *Handler) PublishVideo(c *gin.Context) {
	var req model.PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.svc.PublishVideo(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, task)
}

// GetTask 获取任务状态
// GET /api/v1/task/:id
func (h *Handler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	// TODO: 从数据库获取任务
	c.JSON(http.StatusOK, gin.H{"id": id, "status": "processing", "progress": 50})
}

// ListTasks 列出任务
// GET /api/v1/tasks
func (h *Handler) ListTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tasks": []interface{}{}})
}

// 配置相关 API
func (h *Handler) GetAIProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"providers": []gin.H{
			{"type": "minimax", "name": "Minimax API"},
			{"type": "self_hosted", "name": "Self-hosted (n8n+ComfyUI+Ollama)"},
		},
	})
}

func (h *Handler) GetVideoProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"providers": []string{"douyin", "kuaishou", "bilibili", "xiaohongshu"},
	})
}
