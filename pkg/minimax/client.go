package minimax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://api.minimax.chat",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type GenerateTextRequest struct {
	Model     string  `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GenerateTextResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage   `json:"usage"`
}

type Choice struct {
	Index        int       `json:"index"`
	Message      Message   `json:"message"`
	FinishReason string    `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (c *Client) GenerateText(req *GenerateTextRequest) (*GenerateTextResponse, error) {
	if req.Model == "" {
		req.Model = "minimax-01"
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 1024
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/v1/text/chatcompletion_v2", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result GenerateTextResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GenerateScript generates a video script from input
func (c *Client) GenerateScript(input, style string, duration int) (string, error) {
	prompt := fmt.Sprintf(`Generate a video script for the following content.
Style: %s
Target duration: %d seconds

Content: %s

Return ONLY the script in Chinese, nothing else.`, style, duration, input)

	resp, err := c.GenerateText(&GenerateTextRequest{
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateStoryboard generates a storyboard from script
func (c *Client) GenerateStoryboard(script string, sceneCount int) (string, error) {
	prompt := fmt.Sprintf(`Create a storyboard with %d scenes from this script.
For each scene, describe: shot type, camera angle, action, dialogue.
Return in JSON format with scenes array.

Script: %s

Return ONLY the JSON storyboard in Chinese, nothing else.`, sceneCount, script)

	resp, err := c.GenerateText(&GenerateTextRequest{
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return resp.Choices[0].Message.Content, nil
}
