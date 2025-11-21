package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// KnowledgeClient Dify 知识库客户端
type KnowledgeClient struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

// KnowledgeConfig 知识库客户端配置
type KnowledgeConfig struct {
	APIKey  string
	BaseURL string
	Timeout int
}

// NewKnowledgeClient 创建知识库客户端
func NewKnowledgeClient(config KnowledgeConfig) *KnowledgeClient {
	return &KnowledgeClient{
		apiKey:  config.APIKey,
		baseURL: config.BaseURL,
		timeout: time.Duration(config.Timeout) * time.Second,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// RetrievalRequest 检索请求
type RetrievalRequest struct {
	Query QueryContent `json:"query"`
}

// QueryContent 查询内容
type QueryContent struct {
	Content string `json:"content"`
}

// RetrievalModel 检索模型配置
type RetrievalModel struct {
	SearchMethod          string          `json:"search_method"`             // semantic_search, full_text_search, hybrid_search
	RerankingEnable       bool            `json:"reranking_enable"`          // 是否启用重排序
	RerankingModel        *RerankingModel `json:"reranking_model,omitempty"` // 重排序模型
	TopK                  int             `json:"top_k"`                     // 返回结果数量
	ScoreThresholdEnabled bool            `json:"score_threshold_enabled"`   // 是否启用分数阈值
	ScoreThreshold        float64         `json:"score_threshold"`           // 分数阈值
}

// RerankingModel 重排序模型
type RerankingModel struct {
	RerankingProviderName string `json:"reranking_provider_name"` // cohere, jina
	RerankingModelName    string `json:"reranking_model_name"`    // 模型名称
}

// RetrievalResponse 检索响应
type RetrievalResponse struct {
	Query   QueryContent      `json:"query"`
	Records []RetrievalRecord `json:"records"`
}

// RetrievalRecord 检索记录
type RetrievalRecord struct {
	Segment Segment `json:"segment"`
	Score   float64 `json:"score"`
}

// Segment 文档片段
type Segment struct {
	ID         string                 `json:"id"`
	Position   int                    `json:"position"`
	DocumentID string                 `json:"document_id"`
	Content    string                 `json:"content"`
	Answer     string                 `json:"answer"`
	WordCount  int                    `json:"word_count"`
	Tokens     int                    `json:"tokens"`
	Keywords   []string               `json:"keywords"`
	Document   map[string]interface{} `json:"document"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Retrieve 检索知识库
func (c *KnowledgeClient) Retrieve(ctx context.Context, datasetID string, query string, topK int) (*RetrievalResponse, error) {
	// 构建检索请求
	reqBody := RetrievalRequest{
		Query: QueryContent{
			Content: query,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 构建请求 URL
	url := fmt.Sprintf("%s/datasets/%s/retrieve", c.baseURL, datasetID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error [%s]: %s", errResp.Code, errResp.Message)
	}

	// 解析响应
	var retrievalResp RetrievalResponse
	if err := json.Unmarshal(body, &retrievalResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &retrievalResp, nil
}

// SimpleRetrieve 简单检索（使用默认配置）
func (c *KnowledgeClient) SimpleRetrieve(ctx context.Context, datasetID string, query string) ([]string, error) {
	resp, err := c.Retrieve(ctx, datasetID, query, 5)
	if err != nil {
		return nil, err
	}

	// 提取文档内容
	var contents []string
	for _, record := range resp.Records {
		contents = append(contents, record.Segment.Content)
	}

	return contents, nil
}

// RetrieveWithScore 检索并返回带分数的结果
func (c *KnowledgeClient) RetrieveWithScore(ctx context.Context, datasetID string, query string, topK int) ([]RetrievalRecord, error) {
	resp, err := c.Retrieve(ctx, datasetID, query, topK)
	if err != nil {
		return nil, err
	}

	return resp.Records, nil
}
