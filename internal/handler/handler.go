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

// Health check
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Video Repurposing

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

// Script-to-Video

func (h *Handler) GenerateScript(c *gin.Context) {
	var req model.ScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.svc.GenerateScript(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, task)
}

func (h *Handler) GenerateStoryboard(c *gin.Context) {
	var req model.StoryboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.svc.GenerateStoryboard(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, task)
}

func (h *Handler) GenerateFrames(c *gin.Context) {
	var req model.FrameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.svc.GenerateFrames(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, task)
}

func (h *Handler) GenerateVideo(c *gin.Context) {
	var req model.VideoGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.svc.GenerateVideo(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, task)
}

// Configuration

func (h *Handler) GetVideoProviders(c *gin.Context) {
	// Implementation
	c.JSON(http.StatusOK, gin.H{"providers": []string{"douyin", "kuaishou", "bilibili", "xiaohongshu"}})
}

func (h *Handler) SetVideoProvider(c *gin.Context) {
	var provider model.VideoProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, provider)
}

func (h *Handler) GetAIProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"providers": []gin.H{
			{"type": "minimax", "name": "Minimax API"},
			{"type": "self_hosted", "name": "Self-hosted (n8n+ComfyUI+Ollama)"},
		},
	})
}

func (h *Handler) SetAIProvider(c *gin.Context) {
	var provider model.AIProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, provider)
}

func (h *Handler) GetOSSConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"configured": true})
}

func (h *Handler) SetOSSConfig(c *gin.Context) {
	var cfg model.OSSConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, cfg)
}

// Task status

func (h *Handler) GetTaskStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	// Get task from service/repository
	c.JSON(http.StatusOK, gin.H{"id": id, "status": "processing", "progress": 50})
}

func (h *Handler) ListTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tasks": []interface{}{}})
}
