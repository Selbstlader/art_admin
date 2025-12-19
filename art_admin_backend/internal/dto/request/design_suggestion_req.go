package request

// GenerateDesignSuggestionRequest 生成设计建议请求
// Generate design suggestion request
type GenerateDesignSuggestionRequest struct {
	ProjectID uint   `json:"projectId" binding:"required"` // 项目ID / Project ID
	Category  string `json:"category"`                     // 建议类别筛选 / Category filter (layout/material/style/function)
	Count     int    `json:"count"`                        // 期望生成数量 / Expected count (default: 3)
}

// GetDesignSuggestionListRequest 获取设计建议列表请求
// Get design suggestion list request
type GetDesignSuggestionListRequest struct {
	Current   int    `form:"current" binding:"required,min=1" example:"1"`
	Size      int    `form:"size" binding:"required,min=1,max=100" example:"10"`
	ProjectID uint   `form:"projectId" binding:"required"` // 项目ID / Project ID
	Category  string `form:"category"`                     // 类别筛选 / Category filter
	Status    string `form:"status"`                       // 状态筛选 / Status filter (pending/adopted/ignored)
}

// UpdateSuggestionStatusRequest 更新建议状态请求
// Update suggestion status request
type UpdateSuggestionStatusRequest struct {
	ID       uint   `json:"id" binding:"required"`                                   // 建议ID / Suggestion ID
	Status   string `json:"status" binding:"required,oneof=adopted ignored pending"` // 状态 / Status
	Feedback string `json:"feedback"`                                                // 用户反馈 / User feedback
}

// BatchUpdateSuggestionStatusRequest 批量更新建议状态请求
// Batch update suggestion status request
type BatchUpdateSuggestionStatusRequest struct {
	IDs      []uint `json:"ids" binding:"required,min=1"`                            // 建议ID列表 / Suggestion ID list
	Status   string `json:"status" binding:"required,oneof=adopted ignored pending"` // 状态 / Status
	Feedback string `json:"feedback"`                                                // 用户反馈 / User feedback
}

// DeleteDesignSuggestionRequest 删除设计建议请求
// Delete design suggestion request
type DeleteDesignSuggestionRequest struct {
	ID uint `json:"id" binding:"required"` // 建议ID / Suggestion ID
}

// BatchDeleteDesignSuggestionRequest 批量删除设计建议请求
// Batch delete design suggestions request
type BatchDeleteDesignSuggestionRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 建议ID列表 / Suggestion ID list
}
