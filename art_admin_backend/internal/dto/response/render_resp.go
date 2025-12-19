package response

import "time"

// RenderResult 渲染结果
// Render result for CAD to image generation
type RenderResult struct {
	CadFileID      uint      `json:"cadFileId"`
	ImageURL       string    `json:"imageUrl"`       // 效果图URL / Render image URL
	DesignProposal string    `json:"designProposal"` // AI生成的设计方案 / AI-generated design proposal
	Style          string    `json:"style"`
	RoomType       string    `json:"roomType"`
	Prompt         string    `json:"prompt"`
	GeneratedAt    time.Time `json:"generatedAt"`
}

// RenderHistoryResponse 渲染历史响应
type RenderHistoryResponse struct {
	CadFileID uint           `json:"cadFileId"`
	Records   []RenderResult `json:"records"`
	Total     int            `json:"total"`
}
