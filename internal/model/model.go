package model

import "time"

type Task struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Input     string    `json:"input"`
	Output    string    `json:"output,omitempty"`
	Error     string    `json:"error,omitempty"`
	Progress  int       `json:"progress"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VideoGenerationRequest struct {
	Input       string `json:"input"`
	InputType   string `json:"input_type"`
	Style       string `json:"style"`
	Duration    int    `json:"duration"`
	AspectRatio string `json:"aspect_ratio"`
	Music       string `json:"music"`
}

type DownloadRequest struct {
	Platform   string `json:"platform"`
	VideoURL   string `json:"video_url"`
	MetricType string `json:"metric_type"`
}

type RecreateRequest struct {
	OriginalVideo string `json:"original_video"`
	Style        string `json:"style"`
	KeepAudio    bool   `json:"keep_audio"`
}

type PublishRequest struct {
	VideoPath string   `json:"video_path"`
	Platforms []string `json:"platforms"`
	Caption   string   `json:"caption"`
	Tags      []string `json:"tags"`
}
