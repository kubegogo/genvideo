package model

import "time"

// Task represents a video generation task
type Task struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`    // "repurposing", "script_to_video"
	Status    string    `json:"status"`  // "pending", "processing", "completed", "failed"
	Input     string    `json:"input"`
	Output    string    `json:"output,omitempty"`
	Error     string    `json:"error,omitempty"`
	Progress  int       `json:"progress"` // 0-100
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VideoGenerationRequest 视频生成请求
type VideoGenerationRequest struct {
	Input       string `json:"input"`        // 关键词/文档/小说
	InputType   string `json:"input_type"`  // "keywords", "document", "novel"
	Style       string `json:"style"`       // "dramatic", "comedy", "documentary"
	Duration    int    `json:"duration"`     // 目标时长(秒)
	AspectRatio string `json:"aspect_ratio"` // "16:9", "9:16", "1:1"
	Music       string `json:"music"`       // 背景音乐风格
}

// DownloadRequest 下载视频请求
type DownloadRequest struct {
	Platform   string `json:"platform"`  // "douyin", "kuaishou", "bilibili"
	VideoURL   string `json:"video_url"`
	MetricType string `json:"metric_type"` // "likes", "views", "favorites"
}

// RecreateRequest 二次创作请求
type RecreateRequest struct {
	OriginalVideo string `json:"original_video"`
	Style        string `json:"style"`
	KeepAudio    bool   `json:"keep_audio"`
}

// PublishRequest 发布视频请求
type PublishRequest struct {
	VideoPath string   `json:"video_path"`
	Platforms []string `json:"platforms"`
	Caption   string   `json:"caption"`
	Tags      []string `json:"tags"`
}

// AIProvider AI服务商配置
type AIProvider struct {
	Type     string `json:"type"`     // "minimax", "self_hosted"
	APIKey   string `json:"api_key"`
	BaseURL  string `json:"base_url"`
	IsActive bool   `json:"is_active"`
}

// OSSConfig OSS配置
type OSSConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	IsActive  bool   `json:"is_active"`
}

// VideoProvider 视频平台配置
type VideoProvider struct {
	Platform string `json:"platform"`
	Cookie   string `json:"cookie"`
	Status   string `json:"status"`
}
