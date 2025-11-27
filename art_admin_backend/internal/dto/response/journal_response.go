package response

import "time"

// JournalEntryItem 日记条目项
type JournalEntryItem struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	MoodRecordID   *int64    `json:"mood_record_id"`
	SentimentScore float64   `json:"sentiment_score"`
	IsPublic       bool      `json:"is_public"`
	AnalysisStatus string    `json:"analysis_status"`
	CreatedAt      time.Time `json:"created_at"`
}

// JournalEntryDetail 日记条目详情
type JournalEntryDetail struct {
	JournalEntryItem
	AttachmentURL  string    `json:"attachment_url"`
	PrivacySetting string    `json:"privacy_setting"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// JournalStatistics 日记统计信息
type JournalStatistics struct {
	TotalEntries       int       `json:"total_entries"`
	TotalCount         int64     `json:"total_count"`
	AverageWordCount   float64   `json:"average_word_count"`
	MostCommonMoodType string    `json:"most_common_mood_type"`
	WritingStreak      int       `json:"writing_streak"`
	AnalysisDate       time.Time `json:"analysis_date"`
}
