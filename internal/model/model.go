package model

import "time"

// Task represents a video generation task
type Task struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`    // "video_generation", "repurposing", "publish"
	Status    string    `json:"status"`  // "pending", "processing", "completed", "failed"
	Input     string    `json:"input"`   // input data or reference
	Output    string    `json:"output,omitempty"` // output result or path
	Error     string    `json:"error,omitempty"`
	Progress  int       `json:"progress"` // 0-100
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VideoProvider represents a video platform configuration
type VideoProvider struct {
	ID       int64  `json:"id"`
	Platform string `json:"platform"` // "douyin", "kuaishou", etc.
	Cookie   string `json:"cookie,omitempty"`
	Status   string `json:"status"` // "active", "inactive"
}

// AIProvider represents AI service configuration
type AIProvider struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"` // "minimax", "self_hosted"
	APIKey   string `json:"api_key,omitempty"`
	BaseURL  string `json:"base_url,omitempty"`
	IsActive bool   `json:"is_active"`
}

// OSSConfig represents OSS configuration
type OSSConfig struct {
	ID        int64  `json:"id"`
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	Bucket    string `json:"bucket"`
	IsActive  bool   `json:"is_active"`
}

// ============================================
// API Request/Response Models
// ============================================

// VideoGenerationRequest 视频生成请求
type VideoGenerationRequest struct {
	Input       string `json:"input"`        // 文案/关键词/脚本
	InputType   string `json:"input_type"`  // "keywords", "script", "article"
	Style       string `json:"style"`       // "dramatic", "comedy", "documentary"
	Duration    int    `json:"duration"`     // 目标时长(秒)
	Music       string `json:"music"`        // 背景音乐风格
	AspectRatio string `json:"aspect_ratio"` // "16:9", "9:16", "1:1"
}

// DownloadRequest 下载视频请求
type DownloadRequest struct {
	Platform   string `json:"platform"`  // "douyin", "kuaishou", etc.
	VideoURL   string `json:"video_url"`
	MetricType string `json:"metric_type"` // "likes", "views", "favorites"
}

// RecreateRequest 二次创作请求
type RecreateRequest struct {
	OriginalVideo string `json:"original_video"`
	Style         string `json:"style"`
	KeepAudio     bool   `json:"keep_audio"`
}

// PublishRequest 发布视频请求
type PublishRequest struct {
	VideoPath string   `json:"video_path"`
	Platforms []string `json:"platforms"`
	Caption   string   `json:"caption"`
	Tags      []string `json:"tags"`
}
