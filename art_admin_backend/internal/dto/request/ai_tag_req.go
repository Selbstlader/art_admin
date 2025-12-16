package request

/*** AI Tag Request DTOs ***/

// AITagListRequest represents the list query request for AI tags
// Requirements: 2.1 - Return all active AI_Tag records with pagination
// Requirements: 2.2 - Search tags by keyword
type AITagListRequest struct {
	Current int    `form:"current" binding:"required,min=1" example:"1"`
	Size    int    `form:"size" binding:"required,min=1,max=100" example:"10"`
	Keyword string `form:"keyword" example:"customer service"`
}

// CreateAITagRequest represents the request to create an AI tag
// Requirements: 1.1 - Create new AI_Tag record
// Requirements: 1.3 - Reject creation with empty name or system prompt
type CreateAITagRequest struct {
	Name              string `json:"name" binding:"required,max=100"`
	Description       string `json:"description" binding:"max=500"`
	KnowledgeBaseID   string `json:"knowledge_base_id"`
	KnowledgeBaseName string `json:"knowledge_base_name"`
	SystemPrompt      string `json:"system_prompt" binding:"required"`
	ChatAPIKey        string `json:"chat_api_key"`
}

// UpdateAITagRequest represents the request to update an AI tag
// Requirements: 3.1 - Update AI_Tag record
// Requirements: 3.2 - Reject update with name conflict
type UpdateAITagRequest struct {
	Name              string `json:"name" binding:"max=100"`
	Description       string `json:"description" binding:"max=500"`
	KnowledgeBaseID   string `json:"knowledge_base_id"`
	KnowledgeBaseName string `json:"knowledge_base_name"`
	SystemPrompt      string `json:"system_prompt"`
	ChatAPIKey        string `json:"chat_api_key"`
	Status            *int   `json:"status"`
}
