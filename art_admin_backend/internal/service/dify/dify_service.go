package dify

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/pkg/dify"
	"fmt"

	"github.com/spf13/viper"
)

type DifyService struct {
	datasetClient *dify.Client
	chatClient    *dify.Client
}

func NewDifyService() *DifyService {
	datasetApiKey := viper.GetString("dify.datasetApiKey")
	chatApiKey := viper.GetString("dify.chatApiKey")
	baseURL := viper.GetString("dify.baseUrl")
	timeout := viper.GetInt("dify.timeout")

	if timeout == 0 {
		timeout = 30
	}

	return &DifyService{
		datasetClient: dify.NewClient(datasetApiKey, baseURL, timeout),
		chatClient:    dify.NewClient(chatApiKey, baseURL, timeout),
	}
}

// GetDatasetList 获取知识库列表
func (s *DifyService) GetDatasetList(req *request.DifyDatasetListRequest) (*response.DifyDatasetListResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}

	limit := req.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.datasetClient.GetDatasetList(page, limit, req.Keyword)
}

// GetDatasetDetail 获取知识库详情
func (s *DifyService) GetDatasetDetail(datasetID string) (*response.DifyDatasetDetailResponse, error) {
	if datasetID == "" {
		return nil, fmt.Errorf("知识库ID不能为空")
	}

	return s.datasetClient.GetDatasetDetail(datasetID)
}

// UploadFile 上传文件到知识库
func (s *DifyService) UploadFile(datasetID, filePath string) (*response.DifyUploadFileResponse, error) {
	if datasetID == "" {
		return nil, fmt.Errorf("知识库ID不能为空")
	}

	if filePath == "" {
		return nil, fmt.Errorf("文件路径不能为空")
	}

	return s.datasetClient.UploadFile(datasetID, filePath)
}

// ChatStreaming 流式对话
func (s *DifyService) ChatStreaming(req *request.DifyChatRequest, callback func(event *response.DifyChatStreamEvent) error) error {
	// 构建请求数据
	data := map[string]interface{}{
		"query":         req.Query,
		"user":          req.User,
		"response_mode": "streaming",
	}

	if req.ConversationID != "" {
		data["conversation_id"] = req.ConversationID
	}

	if req.Inputs != nil {
		data["inputs"] = req.Inputs
	}

	if len(req.Files) > 0 {
		data["files"] = req.Files
	}

	if req.AutoGenerateName {
		data["auto_generate_name"] = true
	}

	return s.chatClient.ChatStreaming(data, callback)
}

// Chat 非流式对话
func (s *DifyService) Chat(req *request.DifyChatRequest) (*response.DifyChatResponse, error) {
	// 构建请求数据
	data := map[string]interface{}{
		"query":         req.Query,
		"user":          req.User,
		"response_mode": "blocking",
	}

	if req.ConversationID != "" {
		data["conversation_id"] = req.ConversationID
	}

	if req.Inputs != nil {
		data["inputs"] = req.Inputs
	}

	if len(req.Files) > 0 {
		data["files"] = req.Files
	}

	if req.AutoGenerateName {
		data["auto_generate_name"] = true
	}

	return s.chatClient.Chat(data)
}

// StopChatMessage 停止消息生成
func (s *DifyService) StopChatMessage(taskID, user string) (*response.DifyStopResponse, error) {
	if taskID == "" {
		return nil, fmt.Errorf("任务ID不能为空")
	}
	if user == "" {
		return nil, fmt.Errorf("用户标识不能为空")
	}

	return s.chatClient.StopChatMessage(taskID, user)
}

// GetSuggestedQuestions 获取建议问题
func (s *DifyService) GetSuggestedQuestions(messageID, user string) (*response.DifySuggestedQuestionsResponse, error) {
	if messageID == "" {
		return nil, fmt.Errorf("消息ID不能为空")
	}
	if user == "" {
		return nil, fmt.Errorf("用户标识不能为空")
	}

	return s.chatClient.GetSuggestedQuestions(messageID, user)
}

// GetConversations 获取会话列表
func (s *DifyService) GetConversations(user, lastID string, limit int) (*response.DifyConversationListResponse, error) {
	if user == "" {
		return nil, fmt.Errorf("用户标识不能为空")
	}

	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.chatClient.GetConversations(user, lastID, limit)
}

// GetConversationMessages 获取会话消息历史
func (s *DifyService) GetConversationMessages(conversationID, user, firstID string, limit int) (*response.DifyMessageListResponse, error) {
	if conversationID == "" {
		return nil, fmt.Errorf("会话ID不能为空")
	}
	if user == "" {
		return nil, fmt.Errorf("用户标识不能为空")
	}

	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.chatClient.GetConversationMessages(conversationID, user, firstID, limit)
}

// DeleteConversation 删除会话
func (s *DifyService) DeleteConversation(conversationID, user string) error {
	if conversationID == "" {
		return fmt.Errorf("会话ID不能为空")
	}
	if user == "" {
		return fmt.Errorf("用户标识不能为空")
	}

	return s.chatClient.DeleteConversation(conversationID, user)
}

// RenameConversation 重命名会话
func (s *DifyService) RenameConversation(conversationID, user, name string) (*response.DifyConversation, error) {
	if conversationID == "" {
		return nil, fmt.Errorf("会话ID不能为空")
	}
	if user == "" {
		return nil, fmt.Errorf("用户标识不能为空")
	}
	if name == "" {
		return nil, fmt.Errorf("会话名称不能为空")
	}

	return s.chatClient.RenameConversation(conversationID, user, name)
}

// DeleteDataset 删除知识库
func (s *DifyService) DeleteDataset(datasetID string) error {
	if datasetID == "" {
		return fmt.Errorf("知识库ID不能为空")
	}

	return s.datasetClient.DeleteDataset(datasetID)
}
