package baiduocr

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Client 百度OCR客户端
type Client struct {
	appID       string
	apiKey      string
	secretKey   string
	baseURL     string
	timeout     time.Duration
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
	tokenMutex  sync.RWMutex
}

// Config 客户端配置
type Config struct {
	AppID     string // 应用ID
	APIKey    string // API Key
	SecretKey string // Secret Key
	Timeout   int    // 超时时间（秒）
}

// NewClient 创建百度OCR客户端
func NewClient(config Config) *Client {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30
	}

	return &Client{
		appID:      config.AppID,
		apiKey:     config.APIKey,
		secretKey:  config.SecretKey,
		baseURL:    "https://aip.baidubce.com",
		timeout:    time.Duration(timeout) * time.Second,
		httpClient: &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

// TokenResponse 获取token响应
type TokenResponse struct {
	AccessToken   string `json:"access_token"`
	ExpiresIn     int    `json:"expires_in"`
	Error         string `json:"error"`
	ErrorDesc     string `json:"error_description"`
	RefreshToken  string `json:"refresh_token"`
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
}

// OCRResponse OCR识别响应
type OCRResponse struct {
	LogID          int64        `json:"log_id"`
	WordsResultNum int          `json:"words_result_num"`
	WordsResult    []WordResult `json:"words_result"`
	ErrorCode      int          `json:"error_code,omitempty"`
	ErrorMsg       string       `json:"error_msg,omitempty"`
}

// WordResult 单个识别结果
type WordResult struct {
	Words    string   `json:"words"`
	Location Location `json:"location,omitempty"`
}

// Location 文字位置信息
type Location struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// GetAccessToken 获取访问令牌
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	// 检查缓存的token是否有效
	c.tokenMutex.RLock()
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		token := c.accessToken
		c.tokenMutex.RUnlock()
		return token, nil
	}
	c.tokenMutex.RUnlock()

	// 需要刷新token
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// 双重检查
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	tokenURL := fmt.Sprintf("%s/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		c.baseURL, c.apiKey, c.secretKey)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, nil)
	if err != nil {
		return "", fmt.Errorf("create token request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read token response failed: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("parse token response failed: %w", err)
	}

	if tokenResp.Error != "" {
		return "", fmt.Errorf("get token failed: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	// 缓存token，提前5分钟过期
	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-300) * time.Second)

	return c.accessToken, nil
}

// RecognizeImageBase64 识别Base64编码的图片
func (c *Client) RecognizeImageBase64(ctx context.Context, imageBase64 string) (*OCRResponse, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token failed: %w", err)
	}

	// 使用高精度版接口
	requestURL := fmt.Sprintf("%s/rest/2.0/ocr/v1/accurate_basic?access_token=%s", c.baseURL, token)

	// URL编码图片数据
	formData := url.Values{}
	formData.Set("image", imageBase64)

	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create OCR request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OCR request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read OCR response failed: %w", err)
	}

	var ocrResp OCRResponse
	if err := json.Unmarshal(body, &ocrResp); err != nil {
		return nil, fmt.Errorf("parse OCR response failed: %w", err)
	}

	if ocrResp.ErrorCode != 0 {
		return nil, fmt.Errorf("OCR error [%d]: %s", ocrResp.ErrorCode, ocrResp.ErrorMsg)
	}

	return &ocrResp, nil
}

// RecognizeImageBytes 识别图片字节数据
func (c *Client) RecognizeImageBytes(ctx context.Context, imageData []byte) (*OCRResponse, error) {
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return c.RecognizeImageBase64(ctx, imageBase64)
}

// RecognizeImageURL 识别网络图片
func (c *Client) RecognizeImageURL(ctx context.Context, imageURL string) (*OCRResponse, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token failed: %w", err)
	}

	requestURL := fmt.Sprintf("%s/rest/2.0/ocr/v1/accurate_basic?access_token=%s", c.baseURL, token)

	formData := url.Values{}
	formData.Set("url", imageURL)

	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create OCR request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OCR request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read OCR response failed: %w", err)
	}

	var ocrResp OCRResponse
	if err := json.Unmarshal(body, &ocrResp); err != nil {
		return nil, fmt.Errorf("parse OCR response failed: %w", err)
	}

	if ocrResp.ErrorCode != 0 {
		return nil, fmt.Errorf("OCR error [%d]: %s", ocrResp.ErrorCode, ocrResp.ErrorMsg)
	}

	return &ocrResp, nil
}

// ExtractText 提取所有识别文本（合并为一个字符串）
func (c *Client) ExtractText(ctx context.Context, imageBase64 string) (string, error) {
	resp, err := c.RecognizeImageBase64(ctx, imageBase64)
	if err != nil {
		return "", err
	}

	var result string
	for i, word := range resp.WordsResult {
		if i > 0 {
			result += "\n"
		}
		result += word.Words
	}

	return result, nil
}

// ExtractTextFromURL 从URL提取所有识别文本
func (c *Client) ExtractTextFromURL(ctx context.Context, imageURL string) (string, error) {
	resp, err := c.RecognizeImageURL(ctx, imageURL)
	if err != nil {
		return "", err
	}

	var result string
	for i, word := range resp.WordsResult {
		if i > 0 {
			result += "\n"
		}
		result += word.Words
	}

	return result, nil
}

// ExtractTextFromBytes 从图片字节提取文本
func (c *Client) ExtractTextFromBytes(ctx context.Context, imageData []byte) (string, error) {
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return c.ExtractText(ctx, imageBase64)
}
