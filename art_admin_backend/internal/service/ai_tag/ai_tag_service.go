package ai_tag

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/dify"
	"art_admin_backend/internal/repository"
	"errors"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// CreateAITagRequest represents the request to create an AI tag
type CreateAITagRequest struct {
	Name              string `json:"name" binding:"required,max=100"`
	Description       string `json:"description" binding:"max=500"`
	KnowledgeBaseID   string `json:"knowledge_base_id"`
	KnowledgeBaseName string `json:"knowledge_base_name"`
	SystemPrompt      string `json:"system_prompt" binding:"required"`
	ChatAPIKey        string `json:"chat_api_key"`
}

// UpdateAITagRequest represents the request to update an AI tag
type UpdateAITagRequest struct {
	Name              string `json:"name" binding:"max=100"`
	Description       string `json:"description" binding:"max=500"`
	KnowledgeBaseID   string `json:"knowledge_base_id"`
	KnowledgeBaseName string `json:"knowledge_base_name"`
	SystemPrompt      string `json:"system_prompt"`
	ChatAPIKey        string `json:"chat_api_key"`
	Status            *int   `json:"status"`
}

// AITagListResponse represents the paginated list response
type AITagListResponse struct {
	List  []*model.AITag `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// AITagTestResponse represents the test result response
type AITagTestResponse struct {
	Success           bool   `json:"success"`
	Response          string `json:"response"`
	KnowledgeAccessed bool   `json:"knowledge_accessed"`
	ErrorMessage      string `json:"error_message,omitempty"`
}

// AITagConfig represents the configuration for chat integration
type AITagConfig struct {
	KnowledgeBaseID string `json:"knowledge_base_id"`
	SystemPrompt    string `json:"system_prompt"`
	ChatAPIKey      string `json:"chat_api_key"`
}

// AITagService handles business logic for AI tags
type AITagService struct {
	repo *repository.AITagRepository
}

// NewAITagService creates a new AITagService instance
func NewAITagService() *AITagService {
	return &AITagService{
		repo: repository.NewAITagRepository(),
	}
}

// NewAITagServiceWithRepo creates a new AITagService with a custom repository (for testing)
func NewAITagServiceWithRepo(repo *repository.AITagRepository) *AITagService {
	return &AITagService{
		repo: repo,
	}
}

// CreateTag creates a new AI tag with validation
// Requirements: 1.1 - Create new AI_Tag record and return created tag details
// Requirements: 1.2 - Reject creation with duplicate name
// Requirements: 1.3 - Reject creation with empty name or system prompt
func (s *AITagService) CreateTag(req *CreateAITagRequest) (*model.AITag, error) {
	// Validate empty/whitespace name
	// Requirements: 1.3 - Reject empty name
	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, errors.New("tag name cannot be empty or whitespace only")
	}

	// Validate empty/whitespace system prompt
	// Requirements: 1.3 - Reject empty system prompt
	trimmedPrompt := strings.TrimSpace(req.SystemPrompt)
	if trimmedPrompt == "" {
		return nil, errors.New("system prompt cannot be empty or whitespace only")
	}

	// Check for duplicate name
	// Requirements: 1.2 - Reject duplicate name
	existingTag, err := s.repo.FindByName(trimmedName)
	if err == nil && existingTag != nil {
		return nil, errors.New("tag name already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create the tag
	tag := &model.AITag{
		Name:              trimmedName,
		Description:       strings.TrimSpace(req.Description),
		KnowledgeBaseID:   req.KnowledgeBaseID,
		KnowledgeBaseName: req.KnowledgeBaseName,
		SystemPrompt:      trimmedPrompt,
		ChatAPIKey:        req.ChatAPIKey,
		Status:            1, // Active by default
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// UpdateTag updates an existing AI tag with conflict check
// Requirements: 3.1 - Update AI_Tag record and return updated tag details
// Requirements: 3.2 - Reject update with name conflict
// Requirements: 3.3 - Record update timestamp
func (s *AITagService) UpdateTag(id int64, req *UpdateAITagRequest) (*model.AITag, error) {
	// Find existing tag
	tag, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	// Validate and update name if provided
	if req.Name != "" {
		trimmedName := strings.TrimSpace(req.Name)
		if trimmedName == "" {
			return nil, errors.New("tag name cannot be empty or whitespace only")
		}

		// Check for name conflict with other tags
		// Requirements: 3.2 - Reject name conflict
		existingTag, err := s.repo.FindByNameExcludingID(trimmedName, id)
		if err == nil && existingTag != nil {
			return nil, errors.New("tag name already exists")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		tag.Name = trimmedName
	}

	// Validate and update system prompt if provided
	if req.SystemPrompt != "" {
		trimmedPrompt := strings.TrimSpace(req.SystemPrompt)
		if trimmedPrompt == "" {
			return nil, errors.New("system prompt cannot be empty or whitespace only")
		}
		tag.SystemPrompt = trimmedPrompt
	}

	// Update other fields
	if req.Description != "" {
		tag.Description = strings.TrimSpace(req.Description)
	}
	if req.KnowledgeBaseID != "" {
		tag.KnowledgeBaseID = req.KnowledgeBaseID
	}
	if req.KnowledgeBaseName != "" {
		tag.KnowledgeBaseName = req.KnowledgeBaseName
	}
	if req.ChatAPIKey != "" {
		tag.ChatAPIKey = req.ChatAPIKey
	}
	if req.Status != nil {
		tag.Status = *req.Status
	}

	// Update timestamp is handled by GORM autoUpdateTime
	// Requirements: 3.3 - Record update timestamp
	tag.UpdatedAt = time.Now()

	if err := s.repo.Update(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag performs a soft delete on an AI tag
// Requirements: 4.1 - Soft delete by setting status to inactive
func (s *AITagService) DeleteTag(id int64) error {
	// Check if tag exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tag not found")
		}
		return err
	}

	// Perform soft delete (sets status to 0)
	return s.repo.Delete(id)
}

// GetTag retrieves an AI tag by ID
// Requirements: 1.1 - Return tag details
func (s *AITagService) GetTag(id int64) (*model.AITag, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return tag, nil
}

// ListTags returns paginated AI tags with optional keyword search
// Requirements: 2.1 - Return all active AI_Tag records with pagination
// Requirements: 2.2 - Search tags by keyword
func (s *AITagService) ListTags(page, pageSize int, keyword string) (*AITagListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	tags, total, err := s.repo.List(page, pageSize, keyword)
	if err != nil {
		return nil, err
	}

	return &AITagListResponse{
		List:  tags,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// GetTagConfig retrieves the configuration for chat integration
// Requirements: 5.1 - Get tag config for chat integration
func (s *AITagService) GetTagConfig(id int64) (*AITagConfig, error) {
	tag, err := s.repo.FindActiveByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("active tag not found")
		}
		return nil, err
	}

	return &AITagConfig{
		KnowledgeBaseID: tag.KnowledgeBaseID,
		SystemPrompt:    tag.SystemPrompt,
		ChatAPIKey:      tag.ChatAPIKey,
	}, nil
}

// TestTag tests an AI tag configuration by calling Dify API
// Requirements: 6.1 - Initiate test chat session with predefined test query
// Requirements: 6.2 - Display AI response and knowledge base access status
// Requirements: 6.3 - Display clear error message on failure
func (s *AITagService) TestTag(id int64) (*AITagTestResponse, error) {
	// Get the tag
	tag, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &AITagTestResponse{
				Success:      false,
				ErrorMessage: "tag not found",
			}, nil
		}
		return &AITagTestResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Determine which API key to use
	apiKey := tag.ChatAPIKey
	if apiKey == "" {
		apiKey = viper.GetString("dify.chatApiKey")
	}

	if apiKey == "" {
		return &AITagTestResponse{
			Success:      false,
			ErrorMessage: "no API key configured for this tag or in system defaults",
		}, nil
	}

	// Create Dify client with the tag's API key
	baseURL := viper.GetString("dify.baseUrl")
	timeout := viper.GetInt("dify.timeout")
	if timeout == 0 {
		timeout = 30
	}

	client := dify.NewClient(apiKey, baseURL, timeout)

	// Build test request data
	testQuery := "Hello, this is a test message to verify the configuration."
	if tag.KnowledgeBaseID != "" {
		testQuery = "Please provide a brief summary of the knowledge base content."
	}

	data := map[string]interface{}{
		"query":         testQuery,
		"user":          "test_user",
		"response_mode": "blocking",
	}

	// Add system prompt as input if available
	if tag.SystemPrompt != "" {
		data["inputs"] = map[string]interface{}{
			"system_prompt": tag.SystemPrompt,
		}
	}

	// Call Dify API
	resp, err := client.Chat(data)
	if err != nil {
		return &AITagTestResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Check if knowledge base was accessed
	knowledgeAccessed := false
	if resp.Metadata != nil {
		if retrieverResources, ok := resp.Metadata["retriever_resources"]; ok {
			if resources, ok := retrieverResources.([]interface{}); ok && len(resources) > 0 {
				knowledgeAccessed = true
			}
		}
	}

	return &AITagTestResponse{
		Success:           true,
		Response:          resp.Answer,
		KnowledgeAccessed: knowledgeAccessed,
	}, nil
}

// ValidateName checks if a name is valid (not empty/whitespace)
func ValidateName(name string) bool {
	return strings.TrimSpace(name) != ""
}

// ValidateSystemPrompt checks if a system prompt is valid (not empty/whitespace)
func ValidateSystemPrompt(prompt string) bool {
	return strings.TrimSpace(prompt) != ""
}
