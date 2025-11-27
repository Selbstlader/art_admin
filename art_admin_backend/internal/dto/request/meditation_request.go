package request

// CreateMeditationContentRequest 创建冥想内容请求
type CreateMeditationContentRequest struct {
	Title           string   `json:"title" binding:"required,max=255"`
	Category        string   `json:"category" binding:"required"`
	Duration        int      `json:"duration" binding:"required,min=1"`
	AudioURL        string   `json:"audio_url" binding:"required,url"`
	CoverImage      string   `json:"cover_image" binding:"url"`
	Description     string   `json:"description"`
	DifficultyLevel int      `json:"difficulty_level" binding:"min=1,max=5"`
	Tags            []string `json:"tags"`
	IsPublic        bool     `json:"is_public"`
}

// UpdateMeditationContentRequest 更新冥想内容请求
type UpdateMeditationContentRequest struct {
	ID              int64    `json:"id" binding:"required"`
	Title           *string  `json:"title"`
	Category        *string  `json:"category"`
	Duration        *int     `json:"duration"`
	AudioURL        *string  `json:"audio_url"`
	CoverImage      *string  `json:"cover_image"`
	Description     *string  `json:"description"`
	DifficultyLevel *int     `json:"difficulty_level"`
	Tags            []string `json:"tags"`
	IsPublic        *bool    `json:"is_public"`
}

// MeditationContentListRequest 冥想内容列表查询请求
type MeditationContentListRequest struct {
	Page            int    `form:"page" binding:"required,min=1"`
	PageSize        int    `form:"pageSize" binding:"required,min=1,max=100"`
	Category        string `form:"category"`
	DifficultyLevel *int   `form:"difficulty_level"`
	MinDuration     *int   `form:"min_duration"`
	MaxDuration     *int   `form:"max_duration"`
	IsPublic        *bool  `form:"is_public"`
	Keyword         string `form:"keyword"` // 搜索关键词
}

// LikeMeditationContentRequest 点赞冥想内容请求
type LikeMeditationContentRequest struct {
	ContentID int64 `json:"content_id" binding:"required"`
}
