package request

// CreateGoalRequest 创建目标请求
type CreateGoalRequest struct {
	GoalType    string `json:"goal_type" binding:"required"` // mood_record, meditation, journal
	TargetValue int    `json:"target_value" binding:"required,min=1"`
	Period      string `json:"period" binding:"required"` // daily, weekly, monthly
}

// UpdateGoalRequest 更新目标请求
type UpdateGoalRequest struct {
	TargetValue int  `json:"target_value" binding:"min=0"`
	IsActive    bool `json:"is_active"`
}

// CheckinRequest 打卡请求
type CheckinRequest struct {
	CheckinType string `json:"checkin_type" binding:"required"` // mood, meditation, journal, general
	Note        string `json:"note"`
}

// UpdateGoalProgressRequest 更新目标进度请求
type UpdateGoalProgressRequest struct {
	Value int `json:"value" binding:"required,min=0"`
}
