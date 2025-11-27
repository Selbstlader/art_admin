package response

import "time"

// EmotionPattern 情绪模式
type EmotionPattern struct {
	Type        string  `json:"type"`        // 情绪类型
	Frequency   int     `json:"frequency"`   // 频率
	Percentage  float64 `json:"percentage"`  // 百分比
	Description string  `json:"description"` // 描述
}

// EmotionTrend 情绪趋势
type EmotionTrend struct {
	Period     string  `json:"period"`      // 周期
	Direction  string  `json:"direction"`   // 方向：上升/下降/稳定
	ChangeRate float64 `json:"change_rate"` // 变化率
	Confidence float64 `json:"confidence"`  // 置信度
}

// MoodPatternAnalysis 情绪模式分析
type MoodPatternAnalysis struct {
	UserID          int64            `json:"user_id"`
	AnalysisType    string           `json:"analysis_type"`
	StartDate       time.Time        `json:"start_date"`
	EndDate         time.Time        `json:"end_date"`
	Summary         string           `json:"summary"`
	Patterns        []EmotionPattern `json:"patterns"`
	Trends          []EmotionTrend   `json:"trends"`
	Insights        []string         `json:"insights"`
	Recommendations []string         `json:"recommendations"`
	AnalyzedAt      time.Time        `json:"analyzed_at"`
	DataPoints      int              `json:"data_points"`
}

// MoodTrendData 情绪趋势数据
type MoodTrendData struct {
	Date      string `json:"date"`
	MoodType  string `json:"mood_type"`
	Intensity int    `json:"intensity"`
}

// JournalEntryAnalysis 日记条目分析
type JournalEntryAnalysis struct {
	EntryID         int64     `json:"entry_id"`
	Title           string    `json:"title"`
	AnalysisType    string    `json:"analysis_type"`
	Summary         string    `json:"summary"`
	Insights        []string  `json:"insights"`
	Recommendations []string  `json:"recommendations"`
	AnalyzedAt      time.Time `json:"analyzed_at"`
}

// JournalInsights 日记洞察
type JournalInsights struct {
	UserID           int64            `json:"user_id"`
	AnalysisType     string           `json:"analysis_type"`
	StartDate        time.Time        `json:"start_date"`
	EndDate          time.Time        `json:"end_date"`
	TotalEntries     int              `json:"total_entries"`
	AverageWordCount float64          `json:"average_word_count"`
	MostCommonMood   string           `json:"most_common_mood"`
	WritingStreak    int              `json:"writing_streak"`
	Summary          string           `json:"summary"`
	Patterns         []EmotionPattern `json:"patterns"`
	Trends           []EmotionTrend   `json:"trends"`
	Insights         []string         `json:"insights"`
	Recommendations  []string         `json:"recommendations"`
	AnalyzedAt       time.Time        `json:"analyzed_at"`
	DataPoints       int              `json:"data_points"`
}

// JournalTrendData 日记趋势数据
type JournalTrendData struct {
	Date             string  `json:"date"`
	EntryCount       int     `json:"entry_count"`
	TotalWordCount   int     `json:"total_word_count"`
	AverageWordCount int     `json:"average_word_count"`
	AverageSentiment float64 `json:"average_sentiment"`
}

// MeditationRecordItem 冥想记录项
type MeditationRecordItem struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	ContentID           *int64    `json:"content_id"`
	Duration            int       `json:"duration"`
	MeditationType      string    `json:"meditation_type"`
	MeditationTypeLabel string    `json:"meditation_type_label"`
	DifficultyLevel     int       `json:"difficulty_level"`
	PreMoodType         string    `json:"pre_mood_type"`
	PostMoodType        string    `json:"post_mood_type"`
	PreMoodIntensity    int       `json:"pre_mood_intensity"`
	PostMoodIntensity   int       `json:"post_mood_intensity"`
	Experience          string    `json:"experience"`
	DistractionLevel    int       `json:"distraction_level"`
	Techniques          []string  `json:"techniques"`
	Notes               string    `json:"notes"`
	AnalysisStatus      string    `json:"analysis_status"`
	MoodImprovement     float64   `json:"mood_improvement"`
	CreatedAt           time.Time `json:"created_at"`
}

// MeditationRecordDetail 冥想记录详情
type MeditationRecordDetail struct {
	MeditationRecordItem
	UpdatedAt time.Time `json:"updated_at"`
}

// MeditationSessionAnalysis 冥想会话分析
type MeditationSessionAnalysis struct {
	RecordID           int64     `json:"record_id"`
	Effectiveness      string    `json:"effectiveness"`
	MoodImprovement    string    `json:"mood_improvement"`
	TechniquesFeedback string    `json:"techniques_feedback"`
	Recommendations    []string  `json:"recommendations"`
	NextPractice       string    `json:"next_practice"`
	AnalyzedAt         time.Time `json:"analyzed_at"`
}

// MeditationTrendData 冥想趋势数据
type MeditationTrendData struct {
	Date               string  `json:"date"`
	TotalDuration      int     `json:"total_duration"`
	SessionCount       int     `json:"session_count"`
	AverageImprovement float64 `json:"average_improvement"`
	MostUsedType       string  `json:"most_used_type"`
}

// MeditationRecommendation 冥想推荐
type MeditationRecommendation struct {
	Type           string `json:"type"`
	Reason         string `json:"reason"`
	ExpectedEffect string `json:"expected_effect"`
	Duration       string `json:"duration"`
}

// AnalysisStatus 分析状态
type AnalysisStatus struct {
	Status       string     `json:"status"`        // pending, analyzing, completed, failed
	Progress     float64    `json:"progress"`      // 进度百分比
	StartedAt    *time.Time `json:"started_at"`    // 开始时间
	CompletedAt  *time.Time `json:"completed_at"`  // 完成时间
	ErrorMessage string     `json:"error_message"` // 错误信息
}

// BatchAnalysisRequest 批量分析请求
type BatchAnalysisRequest struct {
	UserIDs      []int64   `json:"user_ids"`      // 用户ID列表
	StartDate    time.Time `json:"start_date"`    // 开始日期
	EndDate      time.Time `json:"end_date"`      // 结束日期
	AnalysisType string    `json:"analysis_type"` // 分析类型
}

// BatchAnalysisResponse 批量分析响应
type BatchAnalysisResponse struct {
	TotalUsers     int                   `json:"total_users"`
	ProcessedUsers int                   `json:"processed_users"`
	FailedUsers    int                   `json:"failed_users"`
	Results        map[int64]interface{} `json:"results"` // 用户ID -> 分析结果
	Errors         map[int64]string      `json:"errors"`  // 用户ID -> 错误信息
	ProcessedAt    time.Time             `json:"processed_at"`
}
