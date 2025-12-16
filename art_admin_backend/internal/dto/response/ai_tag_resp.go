package response

import "time"

/*** AI Tag Response DTOs ***/

// AITagResponse represents a single AI tag in responses
// Requirements: 2.3 - Response structure completeness
type AITagResponse struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	KnowledgeBaseID   string    `json:"knowledge_base_id"`
	KnowledgeBaseName string    `json:"knowledge_base_name"`
	SystemPrompt      string    `json:"system_prompt"`
	ChatAPIKey        string    `json:"chat_api_key"`
	Status            int       `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AITagListResponse represents the paginated list response for AI tags
// Requirements: 2.1 - Return all active AI_Tag records with pagination
type AITagListResponse struct {
	List  []*AITagResponse `json:"list"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// AITagTestResponse represents the test result response
// Requirements: 6.2 - Display AI response and knowledge base access status
type AITagTestResponse struct {
	Success           bool   `json:"success"`
	Response          string `json:"response"`
	KnowledgeAccessed bool   `json:"knowledge_accessed"`
	ErrorMessage      string `json:"error_message,omitempty"`
}
