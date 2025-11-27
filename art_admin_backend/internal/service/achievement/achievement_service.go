package achievement

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/repository"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

// AchievementService 成就服务
type AchievementService struct {
	achievementRepo       *repository.AchievementRepository
	userAchievementRepo   *repository.UserAchievementRepository
	userPointRepo         *repository.UserPointRepository
	achievementEventRepo  *repository.AchievementEventRepository
	leaderboardRepo       *repository.LeaderboardRepository
	moodRecordRepo        *repository.MoodRecordRepository
	meditationContentRepo *repository.MeditationContentRepository
	journalEntryRepo      *repository.JournalEntryRepository
	// 防抖机制：记录用户最后检查成就的时间
	lastCheckTime sync.Map
}

// NewAchievementService 创建成就服务
func NewAchievementService() *AchievementService {
	return &AchievementService{
		achievementRepo:       repository.NewAchievementRepository(),
		userAchievementRepo:   repository.NewUserAchievementRepository(),
		userPointRepo:         repository.NewUserPointRepository(),
		achievementEventRepo:  repository.NewAchievementEventRepository(),
		leaderboardRepo:       repository.NewLeaderboardRepository(),
		moodRecordRepo:        repository.NewMoodRecordRepository(),
		meditationContentRepo: repository.NewMeditationContentRepository(),
		journalEntryRepo:      repository.NewJournalEntryRepository(),
	}
}

// AchievementProgress 成就进度
type AchievementProgress struct {
	AchievementID   int64   `json:"achievement_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Icon            string  `json:"icon"`
	Tier            string  `json:"tier"`
	CurrentProgress int     `json:"current_progress"`
	TargetProgress  int     `json:"target_progress"`
	ProgressPercent float64 `json:"progress_percent"`
	IsUnlocked      bool    `json:"is_unlocked"`
	Points          int     `json:"points"`
}

// UserStats 用户统计信息
type UserStats struct {
	TotalPoints        int                     `json:"total_points"`
	Level              int                     `json:"level"`
	ExperiencePoints   int                     `json:"experience_points"`
	UnlockedCount      int64                   `json:"unlocked_count"`
	TierCounts         map[string]int64        `json:"tier_counts"`
	CategoryPoints     map[string]int          `json:"category_points"`
	RecentAchievements []model.UserAchievement `json:"recent_achievements"`
}

// ProcessUserEvent 处理用户事件，触发成就检查
func (s *AchievementService) ProcessUserEvent(userID int64, eventType string, eventData map[string]interface{}) error {
	// 防抖机制：检查用户是否在1分钟内已经检查过成就
	if lastTime, exists := s.lastCheckTime.Load(userID); exists {
		if time.Since(lastTime.(time.Time)) < time.Minute {
			// 跳过此次检查，避免频繁触发
			return nil
		}
	}

	// 更新最后检查时间
	s.lastCheckTime.Store(userID, time.Now())

	// 创建成就事件
	event := &model.AchievementEvent{
		UserID:    userID,
		EventType: eventType,
		EventData: s.marshalEventData(eventData),
		CreatedAt: time.Now(),
	}

	if err := s.achievementEventRepo.Create(event); err != nil {
		return fmt.Errorf("创建成就事件失败: %w", err)
	}

	// 异步处理成就检查，添加错误恢复
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("成就检查发生panic: %v, 用户ID: %d, 事件类型: %s\n", r, userID, eventType)
			}
		}()

		if err := s.checkAndUnlockAchievements(userID, eventType, eventData); err != nil {
			fmt.Printf("成就检查失败: %v, 用户ID: %d, 事件类型: %s\n", err, userID, eventType)
		}
	}()

	return nil
}

// checkAndUnlockAchievements 检查并解锁成就
func (s *AchievementService) checkAndUnlockAchievements(userID int64, eventType string, eventData map[string]interface{}) error {
	// 获取所有激活的成就
	achievements, err := s.achievementRepo.FindActive()
	if err != nil {
		return fmt.Errorf("获取成就列表失败: %w", err)
	}

	// 获取用户当前成就状态
	userAchievements, err := s.userAchievementRepo.FindByUserID(userID)
	if err != nil {
		return fmt.Errorf("获取用户成就失败: %w", err)
	}

	// 构建用户成就映射
	achievementMap := make(map[int64]*model.UserAchievement)
	for i := range userAchievements {
		achievementMap[userAchievements[i].AchievementID] = &userAchievements[i]
	}

	// 检查每个成就
	for _, achievement := range achievements {
		if s.shouldCheckAchievement(&achievement, eventType) {
			s.checkSingleAchievement(userID, &achievement, achievementMap, eventData)
		}
	}

	// 更新用户积分和等级
	s.updateUserPointsAndLevel(userID)

	// 更新排行榜
	s.updateLeaderboards(userID)

	return nil
}

// shouldCheckAchievement 判断是否应该检查该成就
func (s *AchievementService) shouldCheckAchievement(achievement *model.Achievement, eventType string) bool {
	// 根据成就类型判断是否需要检查
	switch achievement.Type {
	case model.AchievementTypeMoodRecord:
		return eventType == model.EventTypeMoodRecord
	case model.AchievementTypeMeditation:
		return eventType == model.EventTypeMeditation
	case model.AchievementTypeJournal:
		return eventType == model.EventTypeJournal
	case model.AchievementTypeAnalysis:
		return eventType == model.EventTypeAnalysis
	case model.AchievementTypeStreak:
		return eventType == model.EventTypeStreakUpdate
	default:
		return false
	}
}

// checkSingleAchievement 检查单个成就
func (s *AchievementService) checkSingleAchievement(userID int64, achievement *model.Achievement,
	achievementMap map[int64]*model.UserAchievement, eventData map[string]interface{}) {

	// 获取或创建用户成就记录
	userAchievement, exists := achievementMap[achievement.ID]
	if !exists {
		userAchievement = &model.UserAchievement{
			UserID:        userID,
			AchievementID: achievement.ID,
			CurrentValue:  0,
			IsUnlocked:    false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	}

	// 如果已解锁，跳过
	if userAchievement.IsUnlocked {
		return
	}

	// 计算当前进度
	newProgress := s.calculateAchievementProgress(userID, achievement, eventData)

	// 如果进度有提升且达到目标，需要事务处理
	if newProgress > userAchievement.CurrentValue && newProgress >= achievement.TargetValue {
		// 使用事务确保成就解锁和积分奖励的原子性
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			// 更新用户成就状态
			userAchievement.CurrentValue = newProgress
			userAchievement.UpdatedAt = time.Now()
			userAchievement.IsUnlocked = true
			now := time.Now()
			userAchievement.UnlockedAt = &now

			// 保存成就记录
			if exists {
				if err := tx.Save(userAchievement).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(userAchievement).Error; err != nil {
					return err
				}
			}

			// 奖励积分 - 直接在事务中操作
			var userPoint model.UserPoint
			err := tx.Model(&model.UserPoint{}).Where("user_id = ?", userID).First(&userPoint).Error
			if err != nil {
				// 创建新的积分记录
				userPoint := &model.UserPoint{
					UserID:           userID,
					TotalPoints:      achievement.Points,
					Level:            1,
					ExperiencePoints: achievement.Points,
					LastUpdated:      time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				}

				// 根据分类设置积分
				switch achievement.Type {
				case model.AchievementTypeMoodRecord:
					userPoint.MoodPoints = achievement.Points
				case model.AchievementTypeMeditation:
					userPoint.MeditationPoints = achievement.Points
				case model.AchievementTypeJournal:
					userPoint.JournalPoints = achievement.Points
				case model.AchievementTypeAnalysis:
					userPoint.AnalysisPoints = achievement.Points
				}

				if err := tx.Create(userPoint).Error; err != nil {
					return err
				}
			} else {
				// 更新现有积分记录
				updates := map[string]interface{}{
					"total_points":      gorm.Expr("total_points + ?", achievement.Points),
					"experience_points": gorm.Expr("experience_points + ?", achievement.Points),
					"last_updated":      time.Now(),
					"updated_at":        time.Now(),
				}

				// 根据分类更新积分
				switch achievement.Type {
				case model.AchievementTypeMoodRecord:
					updates["mood_points"] = gorm.Expr("mood_points + ?", achievement.Points)
				case model.AchievementTypeMeditation:
					updates["meditation_points"] = gorm.Expr("meditation_points + ?", achievement.Points)
				case model.AchievementTypeJournal:
					updates["journal_points"] = gorm.Expr("journal_points + ?", achievement.Points)
				case model.AchievementTypeAnalysis:
					updates["analysis_points"] = gorm.Expr("analysis_points + ?", achievement.Points)
				}

				if err := tx.Model(&model.UserPoint{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			fmt.Printf("成就解锁事务失败: %v, 用户ID: %d, 成就ID: %d\n", err, userID, achievement.ID)
			return
		}

		// 记录解锁事件（事务外，不影响核心功能）
		s.onAchievementUnlocked(userID, achievement)
	} else if newProgress > userAchievement.CurrentValue {
		// 仅更新进度，不需要事务
		userAchievement.CurrentValue = newProgress
		userAchievement.UpdatedAt = time.Now()

		if exists {
			s.userAchievementRepo.Update(userAchievement)
		} else {
			s.userAchievementRepo.Create(userAchievement)
		}
	}
}

// calculateAchievementProgress 计算成就进度
func (s *AchievementService) calculateAchievementProgress(userID int64, achievement *model.Achievement, eventData map[string]interface{}) int {
	switch achievement.Type {
	case model.AchievementTypeMoodRecord:
		return s.calculateMoodRecordProgress(userID, achievement)
	case model.AchievementTypeMeditation:
		return s.calculateMeditationProgress(userID, achievement)
	case model.AchievementTypeJournal:
		return s.calculateJournalProgress(userID, achievement)
	case model.AchievementTypeAnalysis:
		return s.calculateAnalysisProgress(userID, achievement)
	case model.AchievementTypeStreak:
		return s.calculateStreakProgress(userID, achievement)
	default:
		return 0
	}
}

// calculateMoodRecordProgress 计算情绪记录成就进度
func (s *AchievementService) calculateMoodRecordProgress(userID int64, achievement *model.Achievement) int {
	// 根据成就的具体条件计算进度
	// 这里简化为获取总记录数
	moodRecords, err := s.moodRecordRepo.FindByUserID(userID)
	if err != nil {
		return 0
	}
	return len(moodRecords)
}

// calculateMeditationProgress 计算冥想成就进度
func (s *AchievementService) calculateMeditationProgress(userID int64, achievement *model.Achievement) int {
	// 这里需要实现冥想会话的统计逻辑
	// 暂时返回模拟数据
	return 0
}

// calculateJournalProgress 计算日记成就进度
func (s *AchievementService) calculateJournalProgress(userID int64, achievement *model.Achievement) int {
	// 获取用户日记总数
	journalEntries, _, err := s.journalEntryRepo.FindByDateRange(userID, time.Time{}, time.Now(), 1, 1000)
	if err != nil {
		return 0
	}
	return len(journalEntries)
}

// calculateAnalysisProgress 计算AI分析成就进度
func (s *AchievementService) calculateAnalysisProgress(userID int64, achievement *model.Achievement) int {
	// 这里需要实现AI分析次数的统计逻辑
	// 暂时返回模拟数据
	return 0
}

// calculateStreakProgress 计算连续记录成就进度
func (s *AchievementService) calculateStreakProgress(userID int64, achievement *model.Achievement) int {
	// 获取情绪记录连续天数
	moodStreak, err := s.moodRecordRepo.GetMoodStreak(userID)
	if err != nil {
		return 0
	}
	return moodStreak
}

// awardPoints 奖励积分
func (s *AchievementService) awardPoints(userID int64, points int, category string) error {
	// 确保用户积分记录存在
	userPoint, err := s.userPointRepo.FindByUserID(userID)
	if err != nil {
		// 创建新的积分记录
		userPoint = &model.UserPoint{
			UserID:           userID,
			TotalPoints:      points,
			Level:            1,
			ExperiencePoints: points,
			LastUpdated:      time.Now(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		// 根据分类设置积分
		switch category {
		case model.AchievementTypeMoodRecord:
			userPoint.MoodPoints = points
		case model.AchievementTypeMeditation:
			userPoint.MeditationPoints = points
		case model.AchievementTypeJournal:
			userPoint.JournalPoints = points
		case model.AchievementTypeAnalysis:
			userPoint.AnalysisPoints = points
		}

		return s.userPointRepo.Create(userPoint)
	}

	// 更新现有积分记录
	return s.userPointRepo.AddPoints(userID, points, category)
}

// updateUserPointsAndLevel 更新用户积分和等级
func (s *AchievementService) updateUserPointsAndLevel(userID int64) {
	userPoint, err := s.userPointRepo.FindByUserID(userID)
	if err != nil {
		// 如果找不到用户积分记录，创建初始记录
		userPoint = &model.UserPoint{
			UserID:           userID,
			TotalPoints:      0,
			Level:            1,
			ExperiencePoints: 0,
			LastUpdated:      time.Now(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := s.userPointRepo.Create(userPoint); err != nil {
			return
		}
	}

	// 计算新等级（每1000经验值升一级）
	newLevel := int(userPoint.ExperiencePoints / 1000)
	if newLevel > userPoint.Level {
		userPoint.Level = newLevel

		// 触发升级事件
		s.ProcessUserEvent(userID, model.EventTypeLevelUp, map[string]interface{}{
			"old_level": userPoint.Level,
			"new_level": newLevel,
		})
	}

	s.userPointRepo.Update(userPoint)
}

// updateLeaderboards 更新排行榜
func (s *AchievementService) updateLeaderboards(userID int64) {
	// 更新积分排行榜
	s.updatePointsLeaderboard()
	// 可以添加其他类型的排行榜更新
}

// updatePointsLeaderboard 更新积分排行榜
func (s *AchievementService) updatePointsLeaderboard() {
	// 获取前100名用户
	topUsers, err := s.userPointRepo.GetTopUsers(100)
	if err != nil {
		return
	}

	// 构建排行榜数据
	rankingData := make([]map[string]interface{}, len(topUsers))
	for i, user := range topUsers {
		rankingData[i] = map[string]interface{}{
			"rank":         i + 1,
			"user_id":      user.UserID,
			"username":     "", // 需要从User关联获取
			"total_points": user.TotalPoints,
			"level":        user.Level,
		}
	}

	// 更新全时段排行榜
	leaderboard := &model.Leaderboard{
		Type:        model.LeaderboardTypePoints,
		Period:      model.PeriodAllTime,
		UserRanking: s.marshalRankingData(rankingData),
		UpdatedAt:   time.Now(),
	}

	s.leaderboardRepo.Upsert(leaderboard)
}

// onAchievementUnlocked 成就解锁回调
func (s *AchievementService) onAchievementUnlocked(userID int64, achievement *model.Achievement) {
	// 这里可以发送通知、记录日志等
	fmt.Printf("用户 %d 解锁了成就: %s\n", userID, achievement.Name)
}

// GetUserAchievements 获取用户成就列表
func (s *AchievementService) GetUserAchievements(userID int64) ([]AchievementProgress, error) {
	// 使用分页查询避免内存问题，默认获取前100条记录
	userAchievements, _, err := s.userAchievementRepo.FindByUserIDWithPagination(userID, 1, 100)
	if err != nil {
		return nil, err
	}

	var progress []AchievementProgress
	for _, ua := range userAchievements {
		progressPercent := float64(ua.CurrentValue) / float64(ua.Achievement.TargetValue) * 100
		if progressPercent > 100 {
			progressPercent = 100
		}

		progress = append(progress, AchievementProgress{
			AchievementID:   ua.AchievementID,
			Name:            ua.Achievement.Name,
			Description:     ua.Achievement.Description,
			Icon:            ua.Achievement.Icon,
			Tier:            ua.Achievement.Tier,
			CurrentProgress: ua.CurrentValue,
			TargetProgress:  ua.Achievement.TargetValue,
			ProgressPercent: progressPercent,
			IsUnlocked:      ua.IsUnlocked,
			Points:          ua.Achievement.Points,
		})
	}

	return progress, nil
}

// GetUserStats 获取用户统计信息
func (s *AchievementService) GetUserStats(userID int64) (*UserStats, error) {
	// 获取用户积分信息
	userPoint, err := s.userPointRepo.FindByUserID(userID)
	if err != nil {
		// 如果没有积分记录，返回默认值
		userPoint = &model.UserPoint{
			UserID:           userID,
			TotalPoints:      0,
			Level:            1,
			ExperiencePoints: 0,
		}
	}

	// 获取成就统计
	unlockedCount, _ := s.userAchievementRepo.GetUnlockedCount(userID)
	tierCounts, _ := s.userAchievementRepo.GetUnlockedByTier(userID)

	// 获取最近解锁的成就
	userAchievements, _ := s.userAchievementRepo.FindByUserID(userID)
	var recentAchievements []model.UserAchievement
	for _, ua := range userAchievements {
		if ua.IsUnlocked && ua.UnlockedAt != nil {
			recentAchievements = append(recentAchievements, ua)
		}
	}

	// 构建分类积分
	categoryPoints := map[string]int{
		"mood":       userPoint.MoodPoints,
		"meditation": userPoint.MeditationPoints,
		"journal":    userPoint.JournalPoints,
		"analysis":   userPoint.AnalysisPoints,
		"bonus":      userPoint.BonusPoints,
	}

	return &UserStats{
		TotalPoints:        userPoint.TotalPoints,
		Level:              userPoint.Level,
		ExperiencePoints:   userPoint.ExperiencePoints,
		UnlockedCount:      unlockedCount,
		TierCounts:         tierCounts,
		CategoryPoints:     categoryPoints,
		RecentAchievements: recentAchievements,
	}, nil
}

// GetLeaderboard 获取排行榜 - 实时查询user_points表
func (s *AchievementService) GetLeaderboard(leaderboardType, period string) ([]map[string]interface{}, error) {
	// 实时查询user_points表获取排行榜数据，避免依赖leaderboards表的缓存
	topUsers, err := s.userPointRepo.GetTopUsers(100)
	if err != nil {
		return nil, fmt.Errorf("获取排行榜数据失败: %w", err)
	}

	// 构建排行榜数据
	rankingData := make([]map[string]interface{}, len(topUsers))
	for i, user := range topUsers {
		rankingData[i] = map[string]interface{}{
			"rank":         i + 1,
			"user_id":      user.UserID,
			"username":     "", // 需要从User关联获取，暂时留空
			"total_points": user.TotalPoints,
			"level":        user.Level,
		}
	}

	return rankingData, nil
}

// 辅助方法
func (s *AchievementService) marshalEventData(data map[string]interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}

func (s *AchievementService) marshalRankingData(data []map[string]interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "[]"
	}
	return string(jsonData)
}

func (s *AchievementService) unmarshalRankingData(data string, result *[]map[string]interface{}) error {
	return json.Unmarshal([]byte(data), result)
}
