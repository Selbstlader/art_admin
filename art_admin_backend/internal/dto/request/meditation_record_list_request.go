package request

import "time"

// MeditationRecordListRequest 冥想记录列表请求
type MeditationRecordListRequest struct {
	Page            int        `json:"page" binding:"omitempty,min=1"`
	PageSize        int        `json:"page_size" binding:"omitempty,min=1,max=100"`
	ContentID       *int64     `json:"content_id" binding:"omitempty"`
	MeditationType  string     `json:"meditation_type" binding:"omitempty"`
	DifficultyLevel *int       `json:"difficulty_level" binding:"omitempty,min=1,max=5"`
	AnalysisStatus  string     `json:"analysis_status" binding:"omitempty"`
	MinDuration     *int       `json:"min_duration" binding:"omitempty"`
	MaxDuration     *int       `json:"max_duration" binding:"omitempty"`
	StartDate       *time.Time `json:"start_date" binding:"omitempty"`
	EndDate         *time.Time `json:"end_date" binding:"omitempty"`
}
