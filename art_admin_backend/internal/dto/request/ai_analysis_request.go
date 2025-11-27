package request

import "time"

// EmotionAnalysisRequest 情绪分析请求
type EmotionAnalysisRequest struct {
	UserID       int64     `json:"user_id" binding:"required"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	AnalysisType string    `json:"analysis_type" binding:"required,oneof=pattern trend insight"`
}

// BatchEmotionAnalysisRequest 批量情绪分析请求
type BatchEmotionAnalysisRequest struct {
	UserIDs      []int64   `json:"user_ids" binding:"required"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	AnalysisType string    `json:"analysis_type" binding:"required,oneof=pattern trend insight"`
}

// ExportAnalysisRequest 导出分析请求
type ExportAnalysisRequest struct {
	Format string `json:"format" binding:"required,oneof=json pdf"`
}
