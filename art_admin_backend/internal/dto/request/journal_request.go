package request

import "time"

// CreateJournalEntryRequest 创建日记条目请求
type CreateJournalEntryRequest struct {
	Title        string   `json:"title" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	MoodRecordID *int64   `json:"mood_record_id"`
	IsPrivate    bool     `json:"is_private"`
	Tags         []string `json:"tags"`
	Images       []string `json:"images"`
}

// UpdateJournalEntryRequest 更新日记条目请求
type UpdateJournalEntryRequest struct {
	ID           int64    `json:"id" binding:"required"`
	Title        *string  `json:"title"`
	Content      *string  `json:"content"`
	MoodRecordID *int64   `json:"mood_record_id"`
	IsPrivate    *bool    `json:"is_private"`
	Tags         []string `json:"tags"`
	Images       []string `json:"images"`
}

// JournalEntryListRequest 日记条目列表查询请求
type JournalEntryListRequest struct {
	Page          int        `form:"page" binding:"required,min=1"`
	PageSize      int        `form:"pageSize" binding:"required,min=1,max=100"`
	MoodRecordID  *int64     `form:"mood_record_id"`
	IsPrivate     *bool      `form:"is_private"`
	StartDate     *time.Time `form:"start_date"`
	EndDate       *time.Time `form:"end_date"`
	HasAttachment *bool      `form:"has_attachment"`
}

// UpdateSentimentScoreRequest 更新情感得分请求
type UpdateSentimentScoreRequest struct {
	SentimentScore float64 `json:"sentiment_score" binding:"required,min=-1,max=1"`
}

// UpdateJournalTagsRequest 更新日记标签请求
type UpdateJournalTagsRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

// AddJournalImagesRequest 添加日记图片请求
type AddJournalImagesRequest struct {
	Images []string `json:"images" binding:"required"`
}

// RemoveJournalImageRequest 删除日记图片请求
type RemoveJournalImageRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
}
