package volcengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 火山引擎AI客户端 (豆包大模型)
type Client struct {
	apiKey      string
	baseURL     string
	model       string
	timeout     time.Duration
	maxTokens   int
	temperature float64
	httpClient  *http.Client
}

// Config 客户端配置
type Config struct {
	APIKey      string  // API Key
	BaseURL     string  // API基础URL，默认 https://ark.cn-beijing.volces.com/api/v3
	Model       string  // 模型ID/接入点ID，如 doubao-pro-32k 或 ep-xxxxxxxx
	Timeout     int     // 超时时间（秒）
	MaxTokens   int     // 最大token数
	Temperature float64 // 温度参数
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		BaseURL:     "https://ark.cn-beijing.volces.com/api/v3",
		Model:       "doubao-pro-32k",
		Timeout:     60,
		MaxTokens:   4096,
		Temperature: 0.7,
	}
}

// NewClient 创建火山引擎AI客户端
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	}
	if config.Timeout == 0 {
		config.Timeout = 60
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4096
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}

	return &Client{
		apiKey:      config.APIKey,
		baseURL:     config.BaseURL,
		model:       config.Model,
		timeout:     time.Duration(config.Timeout) * time.Second,
		maxTokens:   config.MaxTokens,
		temperature: config.Temperature,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"` // 消息内容
}

// ChatRequest 对话请求
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse 对话响应
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Chat 发送对话请求
func (c *Client) Chat(ctx context.Context, messages []Message) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		Temperature: c.temperature,
		MaxTokens:   c.maxTokens,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error [%s]: %s", errResp.Error.Code, errResp.Error.Message)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &chatResp, nil
}

// SimpleChat 简单对话（单轮）
func (c *Client) SimpleChat(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userMessage,
		},
	}

	resp, err := c.Chat(ctx, messages)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return resp.Choices[0].Message.Content, nil
}

// ChatWithHistory 带历史记录的对话
func (c *Client) ChatWithHistory(ctx context.Context, messages []Message) (string, error) {
	resp, err := c.Chat(ctx, messages)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return resp.Choices[0].Message.Content, nil
}

// SummarizeContent 总结内容（专用于小红书笔记总结）
func (c *Client) SummarizeContent(ctx context.Context, content string, style string, maxLength int) (string, error) {
	stylePrompts := map[string]string{
		"concise":  "请用简洁凝练的语言总结，仅保留核心信息，适合快速浏览。",
		"detailed": "请详细全面地总结，包含完整逻辑链及细节信息，适合深度了解。",
		"casual":   "请用口语化、生活化的语言转述，让新手用户也能轻松理解。",
	}

	styleHint := stylePrompts["concise"]
	if hint, ok := stylePrompts[style]; ok {
		styleHint = hint
	}

	systemPrompt := fmt.Sprintf(`你是一位专业的内容分析师，擅长分析和总结小红书笔记内容。
%s
总结字数控制在%d字以内。
请以JSON格式返回结果，包含以下字段：
{
  "title": "总结标题",
  "summary": "核心观点概述",
  "key_points": ["关键信息1", "关键信息2", ...],
  "tags": ["标签1", "标签2", ...],
  "sentiment": "positive/neutral/negative"
}`, styleHint, maxLength)

	return c.SimpleChat(ctx, systemPrompt, content)
}
