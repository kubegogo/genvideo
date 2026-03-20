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
			Timeout: 120 * time.Second,
		},
	}
}

type GenerateTextRequest struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	MaxTokens  int       `json:"max_tokens,omitempty"`
	Temperature float64  `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GenerateTextResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	TotalTokens int `json:"total_tokens"`
}

func (c *Client) GenerateText(req *GenerateTextRequest) (*GenerateTextResponse, error) {
	if req.Model == "" {
		req.Model = "minimax-01"
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 2048
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/v1/text/chatcompletion_v2", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

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

	var result GenerateTextResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &result, nil
}

// GenerateScript 生成视频剧本
func (c *Client) GenerateScript(input, inputType, style string, duration int) (string, error) {
	var prompt string
	switch inputType {
	case "keywords":
		prompt = fmt.Sprintf(`根据以下关键词生成一个%d秒的视频剧本。
要求：包含开场、发展、结尾，有节奏感，适合视频表现。

关键词：%s
风格：%s

请用中文回复剧本内容。`, duration, input, style)
	case "document":
		prompt = fmt.Sprintf(`根据以下文档内容，提取核心内容生成一个%d秒的视频剧本。
要求：简洁有力，适合视频表现，保留关键信息。

文档：%s
风格：%s

请用中文回复剧本内容。`, duration, input, style)
	case "novel":
		prompt = fmt.Sprintf(`根据以下小说片段，提取精彩情节生成一个%d秒的视频剧本。
要求：戏剧性强，节奏紧凑，有视觉冲击力。

小说片段：%s
风格：%s

请用中文回复剧本内容。`, duration, input, style)
	}

	resp, err := c.GenerateText(&GenerateTextRequest{
		Messages: []Message{{Role: "user", Content: prompt}},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateStoryboard 生成分镜剧本
func (c *Client) GenerateStoryboard(script string, sceneCount int) (string, error) {
	prompt := fmt.Sprintf(`根据以下剧本生成分镜剧本。
要求：分为%d个场景，每个场景描述：镜头类型、画面内容、台词/解说、配乐建议。

剧本：
%s

请用中文回复分镜内容。`, sceneCount, script)

	resp, err := c.GenerateText(&GenerateTextRequest{
		Messages: []Message{{Role: "user", Content: prompt}},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return resp.Choices[0].Message.Content, nil
}
