package dify

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"art_admin_backend/internal/dto/response"
)

// Client Dify API 客户端
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient 创建 Dify 客户端
func NewClient(apiKey, baseURL string, timeout int) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(method, path string, body io.Reader, contentType string) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	return resp, nil
}

// GetDatasetList 获取知识库列表
func (c *Client) GetDatasetList(page, limit int, keyword string) (*response.DifyDatasetListResponse, error) {
	path := fmt.Sprintf("/datasets?page=%d&limit=%d", page, limit)
	if keyword != "" {
		path += "&keyword=" + keyword
	}

	resp, err := c.doRequest("GET", path, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifyDatasetListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetDatasetDetail 获取知识库详情
func (c *Client) GetDatasetDetail(datasetID string) (*response.DifyDatasetDetailResponse, error) {
	path := fmt.Sprintf("/datasets/%s", datasetID)

	resp, err := c.doRequest("GET", path, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifyDatasetDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// UploadFile 上传文件到知识库
func (c *Client) UploadFile(datasetID, filePath string) (*response.DifyUploadFileResponse, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 关闭writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭writer失败: %w", err)
	}

	// 发送请求
	path := fmt.Sprintf("/datasets/%s/document/create_by_file", datasetID)
	resp, err := c.doRequest("POST", path, body, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifyUploadFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// ChatStreaming 流式对话
func (c *Client) ChatStreaming(data map[string]interface{}, callback func(event *response.DifyChatStreamEvent) error) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	resp, err := c.doRequest("POST", "/chat-messages", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// 读取流式响应
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// 跳过空行
		if line == "" {
			continue
		}

		// SSE格式: data: {...}
		if len(line) > 6 && line[:6] == "data: " {
			jsonStr := line[6:]

			var event response.DifyChatStreamEvent
			if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
				continue // 跳过解析失败的行
			}

			// 调用回调函数
			if err := callback(&event); err != nil {
				return err
			}

			// 如果是结束事件,停止读取
			if event.Event == "message_end" || event.Event == "error" {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流式响应失败: %w", err)
	}

	return nil
}

// Chat 非流式对话
func (c *Client) Chat(data map[string]interface{}) (*response.DifyChatResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	resp, err := c.doRequest("POST", "/chat-messages", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifyChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// handleErrorResponse 处理错误响应
func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取错误响应失败: %w", err)
	}

	var errResp response.DifyErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("Dify API错误 [%s]: %s", errResp.Code, errResp.Message)
}

// StopChatMessage 停止消息生成
func (c *Client) StopChatMessage(taskID, user string) (*response.DifyStopResponse, error) {
	data := map[string]interface{}{
		"user": user,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	path := fmt.Sprintf("/chat-messages/%s/stop", taskID)
	resp, err := c.doRequest("POST", path, bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifyStopResponse
	if err := json.Unmarshal([]byte(`{"result":"success"}`), &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetSuggestedQuestions 获取建议问题
func (c *Client) GetSuggestedQuestions(messageID, user string) (*response.DifySuggestedQuestionsResponse, error) {
	path := fmt.Sprintf("/messages/%s/suggested?user=%s", messageID, user)

	resp, err := c.doRequest("GET", path, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifySuggestedQuestionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetConversations 获取会话列表
func (c *Client) GetConversations(user, lastID string, limit int) (*response.DifyConversationListResponse, error) {
	path := fmt.Sprintf("/conversations?user=%s&limit=%d", user, limit)
	if lastID != "" {
		path += "&last_id=" + lastID
	}

	resp, err := c.doRequest("GET", path, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifyConversationListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetConversationMessages 获取会话消息历史
func (c *Client) GetConversationMessages(conversationID, user, firstID string, limit int) (*response.DifyMessageListResponse, error) {
	path := fmt.Sprintf("/messages?conversation_id=%s&user=%s&limit=%d", conversationID, user, limit)
	if firstID != "" {
		path += "&first_id=" + firstID
	}

	resp, err := c.doRequest("GET", path, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifyMessageListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// DeleteConversation 删除会话
func (c *Client) DeleteConversation(conversationID, user string) error {
	data := map[string]interface{}{
		"user": user,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	path := fmt.Sprintf("/conversations/%s", conversationID)
	resp, err := c.doRequest("DELETE", path, bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// RenameConversation 重命名会话
func (c *Client) RenameConversation(conversationID, user, name string) (*response.DifyConversation, error) {
	data := map[string]interface{}{
		"user": user,
		"name": name,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	path := fmt.Sprintf("/conversations/%s/name", conversationID)
	resp, err := c.doRequest("POST", path, bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result response.DifyConversation
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// DeleteDataset 删除知识库
func (c *Client) DeleteDataset(datasetID string) error {
	path := fmt.Sprintf("/datasets/%s", datasetID)
	resp, err := c.doRequest("DELETE", path, nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 204 No Content 或 200 OK 表示删除成功
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return nil
	}

	// 404 Not Found 表示资源不存在,视为删除成功(幂等性)
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return c.handleErrorResponse(resp)
}
