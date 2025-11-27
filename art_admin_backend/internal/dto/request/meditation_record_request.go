package request

// CreateMeditationRecordRequest 创建冥想记录请求
type CreateMeditationRecordRequest struct {
	ContentID         *int64   `json:"content_id" binding:"omitempty"`                       // 关联的冥想内容ID
	Duration          int      `json:"duration" binding:"required,min=1"`                    // 冥想时长(分钟)
	MeditationType    string   `json:"meditation_type" binding:"required"`                   // 冥想类型
	DifficultyLevel   int      `json:"difficulty_level" binding:"omitempty,min=1,max=5"`     // 难度等级 1-5
	PreMoodType       string   `json:"pre_mood_type" binding:"omitempty"`                    // 冥想前情绪类型
	PostMoodType      string   `json:"post_mood_type" binding:"omitempty"`                   // 冥想后情绪类型
	PreMoodIntensity  int      `json:"pre_mood_intensity" binding:"omitempty,min=1,max=10"`  // 冥想前情绪强度 1-10
	PostMoodIntensity int      `json:"post_mood_intensity" binding:"omitempty,min=1,max=10"` // 冥想后情绪强度 1-10
	Experience        string   `json:"experience" binding:"omitempty"`                       // 体验描述
	DistractionLevel  int      `json:"distraction_level" binding:"omitempty,min=1,max=10"`   // 分心程度 1-10
	Techniques        []string `json:"techniques" binding:"omitempty"`                       // 使用的技巧
	Notes             string   `json:"notes" binding:"omitempty"`                            // 备注
}

// UpdateMeditationRecordRequest 更新冥想记录请求
type UpdateMeditationRecordRequest struct {
	ID                int64    `json:"id" binding:"required"`
	ContentID         *int64   `json:"content_id" binding:"omitempty"`
	Duration          *int     `json:"duration" binding:"omitempty,min=1"`
	MeditationType    *string  `json:"meditation_type" binding:"omitempty"`
	DifficultyLevel   *int     `json:"difficulty_level" binding:"omitempty,min=1,max=5"`
	PreMoodType       *string  `json:"pre_mood_type" binding:"omitempty"`
	PostMoodType      *string  `json:"post_mood_type" binding:"omitempty"`
	PreMoodIntensity  *int     `json:"pre_mood_intensity" binding:"omitempty,min=1,max=10"`
	PostMoodIntensity *int     `json:"post_mood_intensity" binding:"omitempty,min=1,max=10"`
	Experience        *string  `json:"experience" binding:"omitempty"`
	DistractionLevel  *int     `json:"distraction_level" binding:"omitempty,min=1,max=10"`
	Techniques        []string `json:"techniques" binding:"omitempty"`
	Notes             *string  `json:"notes" binding:"omitempty"`
}
