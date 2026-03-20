package model

import "time"

// Task represents a video generation task
type Task struct {
	ID          int64     `json:"id"`
	Type        string    `json:"type"` // "repurposing" or "script_to_video"
	Status      string    `json:"status"` // "pending", "processing", "completed", "failed"
	Input       string    `json:"input"`
	Output      string    `json:"output,omitempty"`
	Error       string    `json:"error,omitempty"`
	Progress    int       `json:"progress"` // 0-100
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

// ScriptRequest is the request to generate a script
type ScriptRequest struct {
	Input     string `json:"input"` // keywords, document path, or novel content
	InputType string `json:"input_type"` // "keywords", "document", "novel"
	Style     string `json:"style"` // "dramatic", "comedy", "documentary", etc.
	Duration  int    `json:"duration"` // target duration in seconds
}

// StoryboardRequest is the request to generate storyboard from script
type StoryboardRequest struct {
	Script    string `json:"script"`
	SceneCount int   `json:"scene_count"`
}

// FrameRequest is the request to generate first/last frame images
type FrameRequest struct {
	Storyboard string `json:"storyboard"`
	Style      string `json:"style"`
}

// VideoGenerationRequest is the request to generate video
type VideoGenerationRequest struct {
	Storyboard string   `json:"storyboard"`
	Frames     []string `json:"frames"` // URLs or paths to frame images
	Duration   int      `json:"duration"`
}

// DownloadRequest is the request to download a video
type DownloadRequest struct {
	Platform   string `json:"platform"` // "douyin", "kuaishou", etc.
	VideoURL   string `json:"video_url"`
	MetricType string `json:"metric_type"` // "likes", "views", "favorites"
}

// RecreateRequest is the request to recreate a video
type RecreateRequest struct {
	OriginalVideo string `json:"original_video"`
	Style         string `json:"style"`
	KeepAudio     bool   `json:"keep_audio"`
}

// PublishRequest is the request to publish a video
type PublishRequest struct {
	VideoPath   string   `json:"video_path"`
	Platforms   []string `json:"platforms"` // target platforms
	Caption     string   `json:"caption"`
	Tags        []string `json:"tags"`
}
