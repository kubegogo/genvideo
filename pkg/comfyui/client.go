package comfyui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	Client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 300 * time.Second, // 5分钟超时
		},
	}
}

type GenerateImageRequest struct {
	Prompt       string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	Steps        int    `json:"steps,omitempty"`
	Seed         int64  `json:"seed,omitempty"`
}

type GenerateImageResponse struct {
	Images []string `json:"images"`
	Seed   int64    `json:"seed"`
}

type GenerateVideoRequest struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Frames      int    `json:"frames,omitempty"`
	FPS         int    `json:"fps,omitempty"`
	Seed        int64  `json:"seed,omitempty"`
}

type GenerateVideoResponse struct {
	VideoPath string `json:"video_path"`
	Status    string `json:"status"`
}

// GenerateImage 生成图片
func (c *Client) GenerateImage(req *GenerateImageRequest) (*GenerateImageResponse, error) {
	if req.Width == 0 {
		req.Width = 512
	}
	if req.Height == 0 {
		req.Height = 512
	}
	if req.Steps == 0 {
		req.Steps = 20
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/api/generate/image", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result GenerateImageResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &result, nil
}

// GenerateVideo 生成视频
func (c *Client) GenerateVideo(req *GenerateVideoRequest) (*GenerateVideoResponse, error) {
	if req.Frames == 0 {
		req.Frames = 24
	}
	if req.FPS == 0 {
		req.FPS = 8
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/api/generate/video", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result GenerateVideoResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &result, nil
}

// GenerateFrame 生成首尾帧
func (c *Client) GenerateFrame(prompt, style string, isFirst bool) (string, error) {
	frameType := "开头"
	if !isFirst {
		frameType = "结尾"
	}

	fullPrompt := fmt.Sprintf("%s风格，%s帧画面：%s", style, frameType, prompt)

	resp, err := c.GenerateImage(&GenerateImageRequest{
		Prompt:       fullPrompt,
		NegativePrompt: "low quality, blurry, distorted",
		Width:        512,
		Height:       512,
		Steps:        25,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Images) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	return resp.Images[0], nil
}
