package request

// AchievementListRequest 成就列表查询请求
type AchievementListRequest struct {
	Page       int    `form:"page" binding:"required,min=1"`
	PageSize   int    `form:"pageSize" binding:"required,min=1,max=100"`
	Type       string `form:"type"`        // 成就类型筛选
	Tier       string `form:"tier"`        // 等级筛选
	IsUnlocked *bool  `form:"is_unlocked"` // 是否已解锁筛选
	ShowHidden bool   `form:"show_hidden"` // 是否显示隐藏成就
}

// ClaimAchievementRequest 领取成就奖励请求
type ClaimAchievementRequest struct {
	AchievementID int64 `json:"achievement_id" binding:"required"`
}

// LeaderboardRequest 排行榜查询请求
type LeaderboardRequest struct {
	Type   string `form:"type" binding:"required,oneof=points streak meditation journal"` // 排行榜类型
	Period string `form:"period" binding:"required,oneof=daily weekly monthly all_time"`  // 周期
	Page   int    `form:"page" binding:"min=1"`
	Limit  int    `form:"limit" binding:"min=1,max=100"`
}
