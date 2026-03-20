package n8n

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
			Timeout: 60 * time.Second,
		},
	}
}

type WorkflowRequest struct {
	WorkflowName string                 `json:"workflow_name"`
	Input        map[string]interface{} `json:"input"`
}

type WorkflowResponse struct {
	ExecutionID string `json:"execution_id"`
	Status      string `json:"status"`
}

// TriggerWorkflow 触发工作流
func (c *Client) TriggerWorkflow(workflowName string, input map[string]interface{}) (*WorkflowResponse, error) {
	req := WorkflowRequest{
		WorkflowName: workflowName,
		Input:        input,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/webhook/trigger", bytes.NewReader(body))
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result WorkflowResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &result, nil
}

// PublishVideo 发布视频到平台
func (c *Client) PublishVideo(platform, videoPath, caption string, tags []string) error {
	input := map[string]interface{}{
		"platform":   platform,
		"video_path": videoPath,
		"caption":    caption,
		"tags":       tags,
	}

	_, err := c.TriggerWorkflow("publish-video", input)
	return err
}

// DownloadVideo 下载视频
func (c *Client) DownloadVideo(platform, url string) error {
	input := map[string]interface{}{
		"platform": platform,
		"url":      url,
	}

	_, err := c.TriggerWorkflow("download-video", input)
	return err
}
