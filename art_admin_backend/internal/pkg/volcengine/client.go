package volcengine

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// RetryConfig 重试配置
// Retry configuration
type RetryConfig struct {
	MaxRetries      int           // 最大重试次数 / Max retry count
	InitialInterval time.Duration // 初始重试间隔 / Initial retry interval
	MaxInterval     time.Duration // 最大重试间隔 / Max retry interval
	Multiplier      float64       // 间隔递增倍数 / Interval multiplier
}

// DefaultRetryConfig 默认重试配置
// Default retry configuration
var DefaultRetryConfig = RetryConfig{
	MaxRetries:      3,
	InitialInterval: 1 * time.Second,
	MaxInterval:     10 * time.Second,
	Multiplier:      2.0,
}

// TokenUsageLog Token使用日志
// Token usage log entry
type TokenUsageLog struct {
	Timestamp        time.Time     `json:"timestamp"`
	RequestID        string        `json:"request_id"`
	Model            string        `json:"model"`
	PromptTokens     int           `json:"prompt_tokens"`
	CompletionTokens int           `json:"completion_tokens"`
	TotalTokens      int           `json:"total_tokens"`
	ResponseTime     time.Duration `json:"response_time_ms"`
	Success          bool          `json:"success"`
	ErrorMessage     string        `json:"error_message,omitempty"`
	RetryCount       int           `json:"retry_count"`
}

// TokenUsageCallback Token使用回调函数类型
// Token usage callback function type
type TokenUsageCallback func(log TokenUsageLog)

// Client 火山引擎AI客户端
// VolcEngine AI client
type Client struct {
	apiKey             string
	baseURL            string
	model              string
	httpClient         *http.Client
	maxTokens          int
	temperature        float64
	retryConfig        RetryConfig
	tokenUsageCallback TokenUsageCallback
	tokenUsageLogs     []TokenUsageLog
	logMutex           sync.RWMutex
	enableLogging      bool
}

// NewClient 创建火山引擎AI客户端
// Create VolcEngine AI client
func NewClient(apiKey, baseURL, model string, timeout, maxTokens int, temperature float64) *Client {
	return &Client{
		apiKey:         apiKey,
		baseURL:        baseURL,
		model:          model,
		maxTokens:      maxTokens,
		temperature:    temperature,
		retryConfig:    DefaultRetryConfig,
		tokenUsageLogs: make([]TokenUsageLog, 0),
		enableLogging:  true,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// NewClientWithConfig 创建带自定义配置的火山引擎AI客户端
// Create VolcEngine AI client with custom configuration
func NewClientWithConfig(apiKey, baseURL, model string, timeout, maxTokens int, temperature float64, retryConfig RetryConfig) *Client {
	client := NewClient(apiKey, baseURL, model, timeout, maxTokens, temperature)
	client.retryConfig = retryConfig
	return client
}

// SetRetryConfig 设置重试配置
// Set retry configuration
func (c *Client) SetRetryConfig(config RetryConfig) {
	c.retryConfig = config
}

// SetTokenUsageCallback 设置Token使用回调
// Set token usage callback
func (c *Client) SetTokenUsageCallback(callback TokenUsageCallback) {
	c.tokenUsageCallback = callback
}

// SetEnableLogging 设置是否启用日志记录
// Set whether to enable logging
func (c *Client) SetEnableLogging(enable bool) {
	c.enableLogging = enable
}

// GetTokenUsageLogs 获取Token使用日志
// Get token usage logs
func (c *Client) GetTokenUsageLogs() []TokenUsageLog {
	c.logMutex.RLock()
	defer c.logMutex.RUnlock()
	logs := make([]TokenUsageLog, len(c.tokenUsageLogs))
	copy(logs, c.tokenUsageLogs)
	return logs
}

// ClearTokenUsageLogs 清除Token使用日志
// Clear token usage logs
func (c *Client) ClearTokenUsageLogs() {
	c.logMutex.Lock()
	defer c.logMutex.Unlock()
	c.tokenUsageLogs = make([]TokenUsageLog, 0)
}

// GetTotalTokenUsage 获取总Token使用量
// Get total token usage
func (c *Client) GetTotalTokenUsage() (promptTokens, completionTokens, totalTokens int) {
	c.logMutex.RLock()
	defer c.logMutex.RUnlock()
	for _, log := range c.tokenUsageLogs {
		if log.Success {
			promptTokens += log.PromptTokens
			completionTokens += log.CompletionTokens
			totalTokens += log.TotalTokens
		}
	}
	return
}

// logTokenUsage 记录Token使用
// Log token usage
func (c *Client) logTokenUsage(usageLog TokenUsageLog) {
	if !c.enableLogging {
		return
	}

	c.logMutex.Lock()
	c.tokenUsageLogs = append(c.tokenUsageLogs, usageLog)
	c.logMutex.Unlock()

	// 调用回调函数 / Call callback function
	if c.tokenUsageCallback != nil {
		c.tokenUsageCallback(usageLog)
	}

	// 输出日志 / Output log
	if usageLog.Success {
		log.Printf("[VolcEngine] 请求成功 - RequestID: %s, Model: %s, PromptTokens: %d, CompletionTokens: %d, TotalTokens: %d, ResponseTime: %dms, Retries: %d",
			usageLog.RequestID, usageLog.Model, usageLog.PromptTokens, usageLog.CompletionTokens, usageLog.TotalTokens, usageLog.ResponseTime.Milliseconds(), usageLog.RetryCount)
	} else {
		log.Printf("[VolcEngine] 请求失败 - RequestID: %s, Model: %s, Error: %s, ResponseTime: %dms, Retries: %d",
			usageLog.RequestID, usageLog.Model, usageLog.ErrorMessage, usageLog.ResponseTime.Milliseconds(), usageLog.RetryCount)
	}
}

// calculateRetryInterval 计算重试间隔（指数退避）
// Calculate retry interval (exponential backoff)
func (c *Client) calculateRetryInterval(retryCount int) time.Duration {
	interval := c.retryConfig.InitialInterval
	for i := 0; i < retryCount; i++ {
		interval = time.Duration(float64(interval) * c.retryConfig.Multiplier)
		if interval > c.retryConfig.MaxInterval {
			interval = c.retryConfig.MaxInterval
			break
		}
	}
	return interval
}

// isRetryableError 判断是否为可重试的错误
// Check if error is retryable
func isRetryableError(statusCode int, err error) bool {
	// 网络错误可重试 / Network errors are retryable
	if err != nil {
		return true
	}
	// 5xx服务端错误可重试 / 5xx server errors are retryable
	if statusCode >= 500 {
		return true
	}
	// 429 Too Many Requests 可重试 / 429 Too Many Requests is retryable
	if statusCode == 429 {
		return true
	}
	return false
}

// generateRequestID 生成请求ID
// Generate request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// ChatMessage 聊天消息
// Chat message
type ChatMessage struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content,omitempty"`
}

// ContentPart 内容部分（支持文本和图片）
// Content part (supports text and image)
type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL 图片URL
// Image URL
type ImageURL struct {
	URL string `json:"url"`
}

// ChatRequest 聊天请求
// Chat request
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream"`
}

// ChatResponse 聊天响应
// Chat response
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
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// ChatResult 聊天结果（包含响应和元数据）
// Chat result (includes response and metadata)
type ChatResult struct {
	Response     *ChatResponse
	RequestID    string
	RetryCount   int
	ResponseTime time.Duration
}

// Chat 发送聊天请求
// Send chat request
func (c *Client) Chat(messages []ChatMessage) (*ChatResponse, error) {
	result, err := c.ChatWithResult(messages)
	if err != nil {
		return nil, err
	}
	return result.Response, nil
}

// ChatWithResult 发送聊天请求并返回完整结果（包含重试信息和响应时间）
// Send chat request and return full result (includes retry info and response time)
func (c *Client) ChatWithResult(messages []ChatMessage) (*ChatResult, error) {
	requestID := generateRequestID()
	startTime := time.Now()
	retryCount := 0

	var lastErr error
	var chatResp *ChatResponse

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			retryCount = attempt
			interval := c.calculateRetryInterval(attempt - 1)
			log.Printf("[VolcEngine] 重试请求 - RequestID: %s, 第%d次重试, 等待%v后重试",
				requestID, attempt, interval)
			time.Sleep(interval)
		}

		chatResp, lastErr = c.doRequest(messages)
		if lastErr == nil {
			// 请求成功 / Request successful
			responseTime := time.Since(startTime)
			usageLog := TokenUsageLog{
				Timestamp:        time.Now(),
				RequestID:        requestID,
				Model:            c.model,
				PromptTokens:     chatResp.Usage.PromptTokens,
				CompletionTokens: chatResp.Usage.CompletionTokens,
				TotalTokens:      chatResp.Usage.TotalTokens,
				ResponseTime:     responseTime,
				Success:          true,
				RetryCount:       retryCount,
			}
			c.logTokenUsage(usageLog)

			return &ChatResult{
				Response:     chatResp,
				RequestID:    requestID,
				RetryCount:   retryCount,
				ResponseTime: responseTime,
			}, nil
		}

		// 检查是否应该重试 / Check if should retry
		if !c.shouldRetry(lastErr, attempt) {
			break
		}
	}

	// 所有重试都失败 / All retries failed
	responseTime := time.Since(startTime)
	usageLog := TokenUsageLog{
		Timestamp:    time.Now(),
		RequestID:    requestID,
		Model:        c.model,
		ResponseTime: responseTime,
		Success:      false,
		ErrorMessage: lastErr.Error(),
		RetryCount:   retryCount,
	}
	c.logTokenUsage(usageLog)

	return nil, fmt.Errorf("请求失败（已重试%d次）: %w", retryCount, lastErr)
}

// doRequest 执行单次请求
// Execute single request
func (c *Client) doRequest(messages []ChatMessage) (*ChatResponse, error) {
	req := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
		Stream:      false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, &NonRetryableError{Err: fmt.Errorf("序列化请求失败: %w", err)}
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &NonRetryableError{Err: fmt.Errorf("创建请求失败: %w", err)}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, &RetryableError{Err: fmt.Errorf("网络请求失败: %w", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RetryableError{Err: fmt.Errorf("读取响应失败: %w", err)}
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, &NonRetryableError{Err: fmt.Errorf("解析响应失败: %w, body: %s", err, string(body))}
	}

	// 检查API错误 / Check API error
	if chatResp.Error != nil {
		errMsg := fmt.Sprintf("API错误 [%s]: %s", chatResp.Error.Code, chatResp.Error.Message)
		if isRetryableError(resp.StatusCode, nil) {
			return nil, &RetryableError{Err: errors.New(errMsg), StatusCode: resp.StatusCode}
		}
		return nil, &NonRetryableError{Err: errors.New(errMsg)}
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("HTTP错误: %d, body: %s", resp.StatusCode, string(body))
		if isRetryableError(resp.StatusCode, nil) {
			return nil, &RetryableError{Err: errors.New(errMsg), StatusCode: resp.StatusCode}
		}
		return nil, &NonRetryableError{Err: errors.New(errMsg)}
	}

	return &chatResp, nil
}

// shouldRetry 判断是否应该重试
// Check if should retry
func (c *Client) shouldRetry(err error, attempt int) bool {
	if attempt >= c.retryConfig.MaxRetries {
		return false
	}

	// 检查是否为可重试错误 / Check if retryable error
	if _, ok := err.(*RetryableError); ok {
		return true
	}
	if _, ok := err.(*NonRetryableError); ok {
		return false
	}

	// 默认重试 / Default retry
	return true
}

// RetryableError 可重试错误
// Retryable error
type RetryableError struct {
	Err        error
	StatusCode int
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// NonRetryableError 不可重试错误
// Non-retryable error
type NonRetryableError struct {
	Err error
}

func (e *NonRetryableError) Error() string {
	return e.Err.Error()
}

func (e *NonRetryableError) Unwrap() error {
	return e.Err
}

// GetRetryConfig 获取当前重试配置
// Get current retry configuration
func (c *Client) GetRetryConfig() RetryConfig {
	return c.retryConfig
}

// GetMaxRetries 获取最大重试次数
// Get max retries
func (c *Client) GetMaxRetries() int {
	return c.retryConfig.MaxRetries
}

// ChatWithText 发送纯文本聊天请求
// Send text-only chat request
func (c *Client) ChatWithText(systemPrompt, userPrompt string) (string, error) {
	messages := []ChatMessage{
		{
			Role: "system",
			Content: []ContentPart{
				{Type: "text", Text: systemPrompt},
			},
		},
		{
			Role: "user",
			Content: []ContentPart{
				{Type: "text", Text: userPrompt},
			},
		},
	}

	resp, err := c.Chat(messages)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("没有返回结果")
	}

	return resp.Choices[0].Message.Content, nil
}

// ChatWithImage 发送带图片的聊天请求
// Send chat request with image
func (c *Client) ChatWithImage(systemPrompt, userPrompt, imageBase64 string) (string, error) {
	messages := []ChatMessage{
		{
			Role: "system",
			Content: []ContentPart{
				{Type: "text", Text: systemPrompt},
			},
		},
		{
			Role: "user",
			Content: []ContentPart{
				{Type: "text", Text: userPrompt},
				{
					Type: "image_url",
					ImageURL: &ImageURL{
						URL: "data:image/jpeg;base64," + imageBase64,
					},
				},
			},
		},
	}

	resp, err := c.Chat(messages)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("没有返回结果")
	}

	return resp.Choices[0].Message.Content, nil
}

// EncodeImageToBase64 将图片文件编码为Base64
// Encode image file to Base64
func EncodeImageToBase64(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取图片文件失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// EncodeImageToBase64FromBytes 将图片字节数据编码为Base64
// Encode image bytes to Base64
func EncodeImageToBase64FromBytes(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// GetTokenUsage 获取最后一次请求的token使用量
// Get token usage from last request
func (c *Client) GetTokenUsage(resp *ChatResponse) (promptTokens, completionTokens, totalTokens int) {
	if resp == nil {
		return 0, 0, 0
	}
	return resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens
}

// =====================================================
// 大文档分段处理功能 / Large document chunking functionality
// =====================================================

// ChunkConfig 分段配置
// Chunk configuration
type ChunkConfig struct {
	MaxTokensPerChunk int    // 每段最大token数 / Max tokens per chunk
	OverlapTokens     int    // 段间重叠token数 / Overlap tokens between chunks
	Separator         string // 分隔符 / Separator
}

// DefaultChunkConfig 默认分段配置
// Default chunk configuration
var DefaultChunkConfig = ChunkConfig{
	MaxTokensPerChunk: 3000,
	OverlapTokens:     200,
	Separator:         "\n\n",
}

// ChunkResult 分段处理结果
// Chunk processing result
type ChunkResult struct {
	ChunkIndex   int           `json:"chunk_index"`
	Content      string        `json:"content"`
	TokensUsed   int           `json:"tokens_used"`
	ResponseTime time.Duration `json:"response_time"`
}

// ChunkedResponse 分段响应
// Chunked response
type ChunkedResponse struct {
	Results      []ChunkResult `json:"results"`
	MergedResult string        `json:"merged_result"`
	TotalChunks  int           `json:"total_chunks"`
	TotalTokens  int           `json:"total_tokens"`
	TotalTime    time.Duration `json:"total_time"`
}

// EstimateTokenCount 估算文本的token数量（简单估算：中文约1.5字符/token，英文约4字符/token）
// Estimate token count for text (simple estimation: ~1.5 chars/token for Chinese, ~4 chars/token for English)
func EstimateTokenCount(text string) int {
	chineseCount := 0
	englishCount := 0
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			chineseCount++
		} else {
			englishCount++
		}
	}
	return int(float64(chineseCount)/1.5) + int(float64(englishCount)/4)
}

// SplitTextIntoChunks 将文本分割成多个块
// Split text into multiple chunks
func SplitTextIntoChunks(text string, config ChunkConfig) []string {
	estimatedTokens := EstimateTokenCount(text)
	if estimatedTokens <= config.MaxTokensPerChunk {
		return []string{text}
	}

	// 按分隔符分割 / Split by separator
	paragraphs := strings.Split(text, config.Separator)
	var chunks []string
	var currentChunk strings.Builder
	currentTokens := 0

	for _, para := range paragraphs {
		paraTokens := EstimateTokenCount(para)

		// 如果单个段落超过限制，需要进一步分割 / If single paragraph exceeds limit, need further split
		if paraTokens > config.MaxTokensPerChunk {
			// 先保存当前块 / Save current chunk first
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
				currentTokens = 0
			}
			// 按句子分割大段落 / Split large paragraph by sentences
			subChunks := splitLargeParagraph(para, config.MaxTokensPerChunk)
			chunks = append(chunks, subChunks...)
			continue
		}

		// 检查是否需要开始新块 / Check if need to start new chunk
		if currentTokens+paraTokens > config.MaxTokensPerChunk {
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				// 添加重叠部分 / Add overlap
				if config.OverlapTokens > 0 {
					overlap := getOverlapText(currentChunk.String(), config.OverlapTokens)
					currentChunk.Reset()
					currentChunk.WriteString(overlap)
					currentTokens = EstimateTokenCount(overlap)
				} else {
					currentChunk.Reset()
					currentTokens = 0
				}
			}
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString(config.Separator)
		}
		currentChunk.WriteString(para)
		currentTokens += paraTokens
	}

	// 添加最后一块 / Add last chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// splitLargeParagraph 分割大段落
// Split large paragraph
func splitLargeParagraph(text string, maxTokens int) []string {
	// 按句子分割 / Split by sentences
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '。' || r == '！' || r == '？' || r == '.' || r == '!' || r == '?'
	})

	var chunks []string
	var currentChunk strings.Builder
	currentTokens := 0

	for _, sentence := range sentences {
		sentenceTokens := EstimateTokenCount(sentence)
		if currentTokens+sentenceTokens > maxTokens {
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
				currentTokens = 0
			}
		}
		currentChunk.WriteString(sentence)
		currentTokens += sentenceTokens
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// getOverlapText 获取重叠文本
// Get overlap text
func getOverlapText(text string, overlapTokens int) string {
	runes := []rune(text)
	// 估算需要的字符数 / Estimate needed characters
	estimatedChars := overlapTokens * 2 // 保守估计
	if estimatedChars > len(runes) {
		return text
	}
	return string(runes[len(runes)-estimatedChars:])
}

// ChatWithChunking 分段处理大文档
// Process large document with chunking
func (c *Client) ChatWithChunking(systemPrompt, userPrompt string, config ChunkConfig) (*ChunkedResponse, error) {
	startTime := time.Now()

	// 检查是否需要分段 / Check if chunking is needed
	estimatedTokens := EstimateTokenCount(userPrompt)
	if estimatedTokens <= config.MaxTokensPerChunk {
		// 不需要分段，直接处理 / No chunking needed, process directly
		result, err := c.ChatWithResult([]ChatMessage{
			{Role: "system", Content: []ContentPart{{Type: "text", Text: systemPrompt}}},
			{Role: "user", Content: []ContentPart{{Type: "text", Text: userPrompt}}},
		})
		if err != nil {
			return nil, err
		}

		content := ""
		if len(result.Response.Choices) > 0 {
			content = result.Response.Choices[0].Message.Content
		}

		return &ChunkedResponse{
			Results: []ChunkResult{{
				ChunkIndex:   0,
				Content:      content,
				TokensUsed:   result.Response.Usage.TotalTokens,
				ResponseTime: result.ResponseTime,
			}},
			MergedResult: content,
			TotalChunks:  1,
			TotalTokens:  result.Response.Usage.TotalTokens,
			TotalTime:    time.Since(startTime),
		}, nil
	}

	// 分段处理 / Process with chunking
	chunks := SplitTextIntoChunks(userPrompt, config)
	results := make([]ChunkResult, 0, len(chunks))
	totalTokens := 0

	// 构建分段处理的系统提示 / Build system prompt for chunked processing
	chunkSystemPrompt := systemPrompt + "\n\n注意：这是一个大文档的分段处理，请处理当前段落内容。"

	for i, chunk := range chunks {
		chunkPrompt := fmt.Sprintf("【第%d段，共%d段】\n%s", i+1, len(chunks), chunk)

		result, err := c.ChatWithResult([]ChatMessage{
			{Role: "system", Content: []ContentPart{{Type: "text", Text: chunkSystemPrompt}}},
			{Role: "user", Content: []ContentPart{{Type: "text", Text: chunkPrompt}}},
		})
		if err != nil {
			return nil, fmt.Errorf("处理第%d段失败: %w", i+1, err)
		}

		content := ""
		if len(result.Response.Choices) > 0 {
			content = result.Response.Choices[0].Message.Content
		}

		results = append(results, ChunkResult{
			ChunkIndex:   i,
			Content:      content,
			TokensUsed:   result.Response.Usage.TotalTokens,
			ResponseTime: result.ResponseTime,
		})
		totalTokens += result.Response.Usage.TotalTokens
	}

	// 合并结果 / Merge results
	mergedResult := mergeChunkResults(results)

	return &ChunkedResponse{
		Results:      results,
		MergedResult: mergedResult,
		TotalChunks:  len(chunks),
		TotalTokens:  totalTokens,
		TotalTime:    time.Since(startTime),
	}, nil
}

// mergeChunkResults 合并分段结果
// Merge chunk results
func mergeChunkResults(results []ChunkResult) string {
	var merged strings.Builder
	for i, result := range results {
		if i > 0 {
			merged.WriteString("\n\n")
		}
		merged.WriteString(result.Content)
	}
	return merged.String()
}

// ChatWithAutoChunking 自动分段处理（使用默认配置）
// Auto chunking with default configuration
func (c *Client) ChatWithAutoChunking(systemPrompt, userPrompt string) (*ChunkedResponse, error) {
	return c.ChatWithChunking(systemPrompt, userPrompt, DefaultChunkConfig)
}
