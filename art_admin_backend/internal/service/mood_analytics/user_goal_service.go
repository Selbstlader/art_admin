package mood_analytics

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"errors"
	"time"
)

// UserGoalService 用户目标服务
type UserGoalService struct {
	goalRepo     *repository.UserGoalRepository
	checkinRepo  *repository.UserCheckinRepository
	progressRepo *repository.GoalProgressRepository
}

// NewUserGoalService 创建用户目标服务
func NewUserGoalService() *UserGoalService {
	return &UserGoalService{
		goalRepo:     repository.NewUserGoalRepository(),
		checkinRepo:  repository.NewUserCheckinRepository(),
		progressRepo: repository.NewGoalProgressRepository(),
	}
}

// CreateGoal 创建目标
func (s *UserGoalService) CreateGoal(userID int64, goalType string, targetValue int, period string) error {
	// 验证目标类型
	if goalType != model.GoalTypeMoodRecord && goalType != model.GoalTypeMeditation && goalType != model.GoalTypeJournal {
		return errors.New("无效的目标类型")
	}

	// 验证周期
	if period != model.GoalPeriodDaily && period != model.GoalPeriodWeekly && period != model.GoalPeriodMonthly {
		return errors.New("无效的目标周期")
	}

	// 验证目标值
	if targetValue <= 0 {
		return errors.New("目标值必须大于0")
	}

	goal := &model.UserGoal{
		UserID:      userID,
		GoalType:    goalType,
		TargetValue: targetValue,
		Period:      period,
		IsActive:    true,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.goalRepo.CreateGoal(goal)
}

// UpdateGoal 更新目标
func (s *UserGoalService) UpdateGoal(userID, goalID int64, targetValue int, isActive bool) error {
	goal, err := s.goalRepo.FindByID(goalID)
	if err != nil {
		return errors.New("目标不存在")
	}

	if goal.UserID != userID {
		return errors.New("无权限操作此目标")
	}

	if targetValue > 0 {
		goal.TargetValue = targetValue
	}
	goal.IsActive = isActive
	goal.UpdatedAt = time.Now()

	return s.goalRepo.UpdateGoal(goal)
}

// DeleteGoal 删除目标
func (s *UserGoalService) DeleteGoal(userID, goalID int64) error {
	goal, err := s.goalRepo.FindByID(goalID)
	if err != nil {
		return errors.New("目标不存在")
	}

	if goal.UserID != userID {
		return errors.New("无权限操作此目标")
	}

	return s.goalRepo.DeleteGoal(goalID)
}

// GetUserGoals 获取用户目标列表
func (s *UserGoalService) GetUserGoals(userID int64) ([]map[string]interface{}, error) {
	goals, err := s.goalRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 获取今日进度
	todayProgress, _ := s.progressRepo.GetTodayProgress(userID)
	progressMap := make(map[int64]int)
	for _, p := range todayProgress {
		progressMap[p.GoalID] = p.CurrentValue
	}

	result := make([]map[string]interface{}, len(goals))
	for i, goal := range goals {
		currentValue := progressMap[goal.ID]
		result[i] = map[string]interface{}{
			"id":            goal.ID,
			"goal_type":     goal.GoalType,
			"target_value":  goal.TargetValue,
			"current_value": currentValue,
			"period":        goal.Period,
			"is_active":     goal.IsActive,
			"progress":      float64(currentValue) / float64(goal.TargetValue) * 100,
			"is_completed":  currentValue >= goal.TargetValue,
			"start_date":    goal.StartDate,
			"created_at":    goal.CreatedAt,
		}
	}

	return result, nil
}

// CheckinGoal 为特定目标打卡
func (s *UserGoalService) CheckinGoal(userID, goalID int64, value int, note string) error {
	// 验证目标是否存在且属于当前用户
	goal, err := s.goalRepo.FindByID(goalID)
	if err != nil {
		return errors.New("目标不存在")
	}

	if goal.UserID != userID {
		return errors.New("无权限操作此目标")
	}

	if !goal.IsActive {
		return errors.New("目标已暂停，无法打卡")
	}

	// 检查今日是否已打卡过此目标
	today := time.Now().Format("2006-01-02")
	hasChecked, err := s.checkinRepo.HasCheckedToday(userID, goalID, today)
	if err != nil {
		return errors.New("检查打卡状态失败")
	}
	if hasChecked {
		return errors.New("今日已打卡此目标")
	}

	// 创建打卡记录
	checkin := &model.UserCheckin{
		UserID:      userID,
		GoalID:      goalID,
		CheckinType: goal.GoalType,
		Value:       value,
		Note:        note,
		CheckinDate: today,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.checkinRepo.CreateCheckin(checkin); err != nil {
		return errors.New("创建打卡记录失败")
	}

	// 更新目标进度
	if err := s.progressRepo.UpdateProgress(userID, goalID, value); err != nil {
		return errors.New("更新目标进度失败")
	}

	return nil
}

// Checkin 打卡
func (s *UserGoalService) Checkin(userID int64, checkinType, note string) error {
	// 验证打卡类型
	validTypes := []string{model.CheckinTypeMood, model.CheckinTypeMeditation, model.CheckinTypeJournal, model.CheckinTypeGeneral}
	isValid := false
	for _, t := range validTypes {
		if checkinType == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("无效的打卡类型")
	}

	return s.checkinRepo.Checkin(userID, checkinType, note)
}

// GetTodayCheckin 获取今日打卡状态
func (s *UserGoalService) GetTodayCheckin(userID int64) (map[string]interface{}, error) {
	checkins, err := s.checkinRepo.GetTodayCheckin(userID)
	if err != nil {
		return nil, err
	}

	// 构建打卡状态
	checkinStatus := map[string]bool{
		model.CheckinTypeMood:       false,
		model.CheckinTypeMeditation: false,
		model.CheckinTypeJournal:    false,
		model.CheckinTypeGeneral:    false,
	}

	for _, c := range checkins {
		checkinStatus[c.CheckinType] = true
	}

	return map[string]interface{}{
		"date":           time.Now().Format("2006-01-02"),
		"checkin_status": checkinStatus,
		"checkin_count":  len(checkins),
	}, nil
}

// GetCheckinHistory 获取打卡历史
func (s *UserGoalService) GetCheckinHistory(userID int64, days int) ([]model.UserCheckin, error) {
	return s.checkinRepo.GetCheckinHistory(userID, days)
}

// GetCheckinStreak 获取连续打卡天数
func (s *UserGoalService) GetCheckinStreak(userID int64, checkinType string) (int, error) {
	return s.checkinRepo.GetCheckinStreak(userID, checkinType)
}

// GetCheckinCalendar 获取打卡日历
func (s *UserGoalService) GetCheckinCalendar(userID int64, year, month int) ([]map[string]interface{}, error) {
	return s.checkinRepo.GetCheckinCalendar(userID, year, month)
}

// GetCheckinStats 获取打卡统计
func (s *UserGoalService) GetCheckinStats(userID int64) (map[string]interface{}, error) {
	// 获取各类型连续打卡天数
	moodStreak, _ := s.checkinRepo.GetCheckinStreak(userID, model.CheckinTypeMood)
	meditationStreak, _ := s.checkinRepo.GetCheckinStreak(userID, model.CheckinTypeMeditation)
	journalStreak, _ := s.checkinRepo.GetCheckinStreak(userID, model.CheckinTypeJournal)
	totalStreak, _ := s.checkinRepo.GetCheckinStreak(userID, "")

	// 获取最近30天打卡记录
	history, _ := s.checkinRepo.GetCheckinHistory(userID, 30)

	// 统计各类型打卡次数
	typeCount := make(map[string]int)
	for _, c := range history {
		typeCount[c.CheckinType]++
	}

	return map[string]interface{}{
		"mood_streak":       moodStreak,
		"meditation_streak": meditationStreak,
		"journal_streak":    journalStreak,
		"total_streak":      totalStreak,
		"last_30_days":      len(history),
		"type_count":        typeCount,
	}, nil
}

// UpdateGoalProgress 更新目标进度
func (s *UserGoalService) UpdateGoalProgress(userID, goalID int64, value int) error {
	goal, err := s.goalRepo.FindByID(goalID)
	if err != nil {
		return errors.New("目标不存在")
	}

	if goal.UserID != userID {
		return errors.New("无权限操作此目标")
	}

	return s.progressRepo.UpdateProgress(userID, goalID, value)
}

// GetGoalProgress 获取目标进度历史
func (s *UserGoalService) GetGoalProgress(userID, goalID int64, days int) ([]model.GoalProgress, error) {
	goal, err := s.goalRepo.FindByID(goalID)
	if err != nil {
		return nil, errors.New("目标不存在")
	}

	if goal.UserID != userID {
		return nil, errors.New("无权限查看此目标")
	}

	return s.progressRepo.GetProgressHistory(userID, goalID, days)
}
