package response

import "time"

// DesignSuggestionResponse 设计建议响应
// Design suggestion response
type DesignSuggestionResponse struct {
	ID              uint      `json:"id"`
	ProjectID       uint      `json:"projectId"`
	Content         string    `json:"content"`         // 建议内容 / Suggestion content
	ApplicableScene string    `json:"applicableScene"` // 适用场景 / Applicable scene
	CostImpact      string    `json:"costImpact"`      // 成本影响 / Cost impact
	Category        string    `json:"category"`        // 建议类别 / Category
	Priority        int       `json:"priority"`        // 优先级 / Priority
	Status          string    `json:"status"`          // 状态 / Status
	DetailInfo      any       `json:"detailInfo"`      // 详细信息 / Detail info
	ReferenceImages []string  `json:"referenceImages"` // 参考图片 / Reference images
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// DesignSuggestionListResponse 设计建议列表响应
// Design suggestion list response
type DesignSuggestionListResponse struct {
	Records []DesignSuggestionResponse `json:"records"`
	Current int                        `json:"current"`
	Size    int                        `json:"size"`
	Total   int64                      `json:"total"`
}

// GenerateDesignSuggestionResponse 生成设计建议响应
// Generate design suggestion response
type GenerateDesignSuggestionResponse struct {
	Suggestions []DesignSuggestionResponse `json:"suggestions"` // 生成的建议列表 / Generated suggestions
	ProjectID   uint                       `json:"projectId"`   // 项目ID / Project ID
	Count       int                        `json:"count"`       // 生成数量 / Generated count
}

// SuggestionDetailResponse 建议详情响应
// Suggestion detail response
type SuggestionDetailResponse struct {
	DesignSuggestionResponse
	ProjectName string `json:"projectName"` // 项目名称 / Project name
	ProjectInfo struct {
		Area   float64 `json:"area"`   // 面积 / Area
		Budget float64 `json:"budget"` // 预算 / Budget
		Style  string  `json:"style"`  // 风格 / Style
	} `json:"projectInfo"` // 项目信息 / Project info
}

// SuggestionPreferenceResponse 建议偏好记录响应
// Suggestion preference response
type SuggestionPreferenceResponse struct {
	ID           uint      `json:"id"`
	SuggestionID uint      `json:"suggestionId"`
	UserID       uint      `json:"userId"`
	Action       string    `json:"action"`   // 操作类型 / Action type
	Feedback     string    `json:"feedback"` // 用户反馈 / User feedback
	CreatedAt    time.Time `json:"createdAt"`
}

// SuggestionStatsResponse 建议统计响应
// Suggestion statistics response
type SuggestionStatsResponse struct {
	TotalCount    int64            `json:"totalCount"`    // 总数 / Total count
	AdoptedCount  int64            `json:"adoptedCount"`  // 采纳数 / Adopted count
	IgnoredCount  int64            `json:"ignoredCount"`  // 忽略数 / Ignored count
	PendingCount  int64            `json:"pendingCount"`  // 待处理数 / Pending count
	CategoryStats map[string]int64 `json:"categoryStats"` // 分类统计 / Category statistics
}

// GenerateSuggestionAsyncResponse 异步生成设计建议响应
// Async generate design suggestion response
type GenerateSuggestionAsyncResponse struct {
	TaskID    uint   `json:"taskId"`    // 任务ID / Task ID
	ProjectID uint   `json:"projectId"` // 项目ID / Project ID
	Status    string `json:"status"`    // 任务状态 / Task status
	Message   string `json:"message"`   // 提示信息 / Message
}

// SuggestionTaskStatusResponse 建议生成任务状态响应
// Suggestion task status response
type SuggestionTaskStatusResponse struct {
	TaskID         uint      `json:"taskId"`         // 任务ID / Task ID
	ProjectID      uint      `json:"projectId"`      // 项目ID / Project ID
	Status         string    `json:"status"`         // 状态 / Status
	GeneratedCount int       `json:"generatedCount"` // 已生成数量 / Generated count
	ErrorMessage   string    `json:"errorMessage"`   // 错误信息 / Error message
	CreatedAt      time.Time `json:"createdAt"`      // 创建时间 / Created at
}
