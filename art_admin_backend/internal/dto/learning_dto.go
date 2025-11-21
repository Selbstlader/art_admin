package dto

import "art_admin_backend/internal/model"

// GenerateMaterialRequest 生成教材请求
type GenerateMaterialRequest struct {
	SubjectID  int64  `json:"subject_id" binding:"required"`
	Grade      string `json:"grade" binding:"required"`
	Topic      string `json:"topic" binding:"required"`
	Difficulty int8   `json:"difficulty" binding:"required,min=1,max=3"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	SubjectID  int64  `json:"subject_id"`
	Grade      string `json:"grade"`
	Topic      string `json:"topic"`
	Difficulty int8   `json:"difficulty"`
	Status     int8   `json:"status"`
	MaterialID int64  `json:"material_id"`
	ErrorMsg   string `json:"error_msg"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// MaterialListRequest 教材列表请求
type MaterialListRequest struct {
	SubjectID int64  `form:"subject_id"`
	Grade     string `form:"grade"`
	Page      int    `form:"page" binding:"required,min=1"`
	Limit     int    `form:"limit" binding:"required,min=1,max=100"`
}

// MaterialResponse 教材响应
type MaterialResponse struct {
	ID            int64                 `json:"id"`
	UserID        int64                 `json:"user_id"`
	SubjectID     int64                 `json:"subject_id"`
	Title         string                `json:"title"`
	Topic         string                `json:"topic"`
	Grade         string                `json:"grade"`
	Difficulty    int8                  `json:"difficulty"`
	Summary       string                `json:"summary"`
	Content       model.MaterialContent `json:"content"`
	AudioURL      string                `json:"audio_url"`
	AudioDuration int                   `json:"audio_duration"`
	TotalTime     int                   `json:"total_time"`
	ViewCount     int                   `json:"view_count"`
	FavoriteCount int                   `json:"favorite_count"`
	CreatedAt     string                `json:"created_at"`
	UpdatedAt     string                `json:"updated_at"`
}

// MaterialListResponse 教材列表响应
type MaterialListResponse struct {
	List  []*MaterialResponse `json:"list"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

// SubjectResponse 学科响应
type SubjectResponse struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Icon        string   `json:"icon"`
	Description string   `json:"description"`
	GradeLevels []string `json:"grade_levels"`
	Sort        int      `json:"sort"`
}

// GradeInfo 年级信息
type GradeInfo struct {
	Grade string `json:"grade"`
	Count int64  `json:"count"`
}

// GradeListResponse 年级列表响应
type GradeListResponse struct {
	SubjectID int64        `json:"subject_id"`
	Grades    []*GradeInfo `json:"grades"`
}
