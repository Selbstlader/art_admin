package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	achievementResponse "art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserAchievements 获取用户成就列表
// @Summary 获取成就列表
// @Description 获取用户的成就列表，支持筛选和分页
// @Tags 成就系统
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "每页数量" default(10)
// @Param type query string false "成就类型"
// @Param tier query string false "成就等级"
// @Param is_unlocked query bool false "是否已解锁"
// @Param show_hidden query bool false "是否显示隐藏成就" default(false)
// @Success 200 {object} response.Response{data=response.PaginationResponse{items=[]response.AchievementProgressItem}} "获取成功"
// @Router /api/achievements [get]
func GetAchievementList(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	// 解析查询参数
	var req request.AchievementListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取用户成就进度
	achievements, err := achievementService.GetUserAchievements(userID)
	if err != nil {
		response.ServerError(c, "获取成就列表失败")
		return
	}

	// 筛选处理
	var filteredAchievements []achievementResponse.AchievementProgressItem
	for _, achievement := range achievements {
		// 类型筛选（暂时跳过，因为service.AchievementProgress没有Type字段）
		// if req.Type != "" && achievement.Type != req.Type {
		//     continue
		// }
		// 等级筛选
		if req.Tier != "" && achievement.Tier != req.Tier {
			continue
		}
		// 解锁状态筛选
		if req.IsUnlocked != nil && achievement.IsUnlocked != *req.IsUnlocked {
			continue
		}
		// 隐藏成就筛选（暂时跳过，因为service.AchievementProgress没有IsHidden字段）
		// if !req.ShowHidden && achievement.IsHidden && !achievement.IsUnlocked {
		//     continue
		// }

		// 转换为响应类型
		achievementItem := achievementResponse.AchievementProgressItem{
			AchievementItem: achievementResponse.AchievementItem{
				ID:          achievement.AchievementID,
				Name:        achievement.Name,
				Description: achievement.Description,
				Icon:        achievement.Icon,
				Tier:        achievement.Tier,
				TargetValue: achievement.TargetProgress,
				Points:      achievement.Points,
				IsHidden:    false, // 默认值，需要从数据库获取
			},
			CurrentProgress: achievement.CurrentProgress,
			ProgressPercent: achievement.ProgressPercent,
			IsUnlocked:      achievement.IsUnlocked,
			UnlockedAt:      nil, // 需要从数据库获取
		}

		filteredAchievements = append(filteredAchievements, achievementItem)
	}

	// 分页处理
	total := len(filteredAchievements)
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	if end > total {
		end = total
	}

	var pageItems []achievementResponse.AchievementProgressItem
	if start < total {
		pageItems = filteredAchievements[start:end]
	}

	response.SuccessWithPagination(c, pageItems, req.Page, req.PageSize, int64(total))
}

// GetAchievementStats 获取用户成就统计
// @Summary 获取用户成就统计
// @Description 获取用户的成就统计信息，包括积分、等级、成就数量等
// @Tags 成就系统
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.UserStatsItem} "获取成功"
// @Router /api/achievements/stats [get]
func GetAchievementStats(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	// 获取用户统计信息
	stats, err := achievementService.GetUserStats(userID)
	if err != nil {
		response.ServerError(c, "获取统计信息失败")
		return
	}

	// 计算下一等级和进度
	nextLevelPoints := (stats.Level + 1) * 1000
	levelProgress := float64(stats.ExperiencePoints%1000) / 1000 * 100

	// 构建响应
	userStats := &achievementResponse.UserStatsItem{
		TotalPoints:      stats.TotalPoints,
		Level:            stats.Level,
		ExperiencePoints: stats.ExperiencePoints,
		UnlockedCount:    stats.UnlockedCount,
		TierCounts:       stats.TierCounts,
		CategoryPoints:   stats.CategoryPoints,
		NextLevelPoints:  nextLevelPoints,
		LevelProgress:    levelProgress,
	}

	// 转换最近成就
	for _, ua := range stats.RecentAchievements {
		if ua.Achievement != nil {
			progressPercent := float64(ua.CurrentValue) / float64(ua.Achievement.TargetValue) * 100
			if progressPercent > 100 {
				progressPercent = 100
			}

			userStats.RecentAchievements = append(userStats.RecentAchievements, achievementResponse.AchievementProgressItem{
				AchievementItem: achievementResponse.AchievementItem{
					ID:          ua.Achievement.ID,
					Type:        ua.Achievement.Type,
					Name:        ua.Achievement.Name,
					Description: ua.Achievement.Description,
					Icon:        ua.Achievement.Icon,
					Tier:        ua.Achievement.Tier,
					TargetValue: ua.Achievement.TargetValue,
					Points:      ua.Achievement.Points,
					IsHidden:    ua.Achievement.IsHidden,
				},
				CurrentProgress: ua.CurrentValue,
				ProgressPercent: progressPercent,
				IsUnlocked:      ua.IsUnlocked,
				UnlockedAt:      ua.UnlockedAt,
			})
		}
	}

	response.Success(c, userStats)
}

// GetLeaderboard 获取排行榜
// @Summary 获取排行榜
// @Description 获取指定类型和周期的排行榜
// @Tags 成就系统
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string true "排行榜类型" Enums(points, streak, meditation, journal)
// @Param period query string true "周期" Enums(daily, weekly, monthly, all_time)
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(50)
// @Success 200 {object} response.Response{data=response.LeaderboardResponse} "获取成功"
// @Router /api/achievements/leaderboard [get]
func GetLeaderboard(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	// 解析查询参数
	var req request.LeaderboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 50
	}

	// 获取排行榜数据
	rankingData, err := achievementService.GetLeaderboard(req.Type, req.Period)
	if err != nil {
		response.ServerError(c, "获取排行榜失败")
		return
	}

	// 转换为响应格式
	var topUsers []achievementResponse.LeaderboardItem
	var currentUser *achievementResponse.LeaderboardItem

	for i, item := range rankingData {
		rank := item["rank"].(int)
		userIDItem := item["user_id"].(int64)
		value := item["total_points"].(int)
		level := item["level"].(int)

		leaderboardItem := achievementResponse.LeaderboardItem{
			Rank:          rank,
			UserID:        userIDItem,
			Username:      "", // 需要从用户表获取
			Avatar:        "",
			Value:         value,
			Level:         level,
			Tier:          calculateTier(level),
			IsCurrentUser: userIDItem == userID,
		}

		if userIDItem == userID {
			currentUser = &leaderboardItem
		}

		// 只返回前N名
		if i < req.Limit {
			topUsers = append(topUsers, leaderboardItem)
		}
	}

	// 构建响应
	leaderboardResponse := &achievementResponse.LeaderboardResponse{
		Type:        req.Type,
		Period:      req.Period,
		TotalCount:  len(rankingData),
		CurrentUser: currentUser,
		TopUsers:    topUsers,
		UpdatedAt:   time.Now(),
	}

	response.Success(c, leaderboardResponse)
}

// GetRecentAchievements 获取最近解锁的成就
// @Summary 获取最近解锁的成就
// @Description 获取用户最近解锁的成就列表，用于通知和展示
// @Tags 成就系统
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "数量限制" default(5)
// @Success 200 {object} response.Response{data=[]response.AchievementProgressItem} "获取成功"
// @Router /api/achievements/recent [get]
func GetRecentAchievements(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	// 解析查询参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	// 获取用户统计信息（包含最近成就）
	stats, err := achievementService.GetUserStats(userID)
	if err != nil {
		response.ServerError(c, "获取最近成就失败")
		return
	}

	// 转换并限制数量
	var recentAchievements []achievementResponse.AchievementProgressItem
	for i, ua := range stats.RecentAchievements {
		if i >= limit {
			break
		}

		if ua.Achievement != nil {
			progressPercent := float64(ua.CurrentValue) / float64(ua.Achievement.TargetValue) * 100
			if progressPercent > 100 {
				progressPercent = 100
			}

			recentAchievements = append(recentAchievements, achievementResponse.AchievementProgressItem{
				AchievementItem: achievementResponse.AchievementItem{
					ID:          ua.Achievement.ID,
					Type:        ua.Achievement.Type,
					Name:        ua.Achievement.Name,
					Description: ua.Achievement.Description,
					Icon:        ua.Achievement.Icon,
					Tier:        ua.Achievement.Tier,
					TargetValue: ua.Achievement.TargetValue,
					Points:      ua.Achievement.Points,
					IsHidden:    ua.Achievement.IsHidden,
				},
				CurrentProgress: ua.CurrentValue,
				ProgressPercent: progressPercent,
				IsUnlocked:      ua.IsUnlocked,
				UnlockedAt:      ua.UnlockedAt,
			})
		}
	}

	response.Success(c, recentAchievements)
}

// ClaimAchievement 领取成就奖励
// @Summary 领取成就奖励
// @Description 领取特定成就的奖励（积分等）
// @Tags 成就系统
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.ClaimAchievementRequest true "领取请求"
// @Success 200 {object} response.Response{data=response.AchievementNotification} "领取成功"
// @Router /api/achievements/claim [post]
func ClaimAchievement(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	// 解析请求参数
	var req request.ClaimAchievementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// TODO: 实现领取奖励逻辑
	// 这里可以添加额外的验证和奖励发放逻辑

	response.SuccessWithMsg(c, "奖励领取成功", nil)
}

// 辅助函数
func calculateTier(level int) string {
	switch {
	case level >= 100:
		return "platinum"
	case level >= 50:
		return "gold"
	case level >= 20:
		return "silver"
	default:
		return "bronze"
	}
}
