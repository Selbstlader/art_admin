package request

/*** AI Generate Request DTOs ***/

// GeneratePromptRequest represents the request to generate/polish system prompt
type GeneratePromptRequest struct {
	Keywords      string `json:"keywords" binding:"required,max=500"`            // User input keywords
	CurrentPrompt string `json:"current_prompt" binding:"max=5000"`              // Existing prompt to polish (optional)
	Mode          string `json:"mode" binding:"omitempty,oneof=generate polish"` // generate: create new, polish: improve existing
}
