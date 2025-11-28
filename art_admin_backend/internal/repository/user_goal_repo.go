package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"errors"
	"time"

	"gorm.io/gorm"
)

// UserGoalRepository 用户目标仓库
type UserGoalRepository struct{}

// NewUserGoalRepository 创建用户目标仓库
func NewUserGoalRepository() *UserGoalRepository {
	return &UserGoalRepository{}
}

// CreateGoal 创建目标
func (r *UserGoalRepository) CreateGoal(goal *model.UserGoal) error {
	db := database.GetDB()
	return db.Create(goal).Error
}

// UpdateGoal 更新目标
func (r *UserGoalRepository) UpdateGoal(goal *model.UserGoal) error {
	db := database.GetDB()
	return db.Save(goal).Error
}

// DeleteGoal 删除目标
func (r *UserGoalRepository) DeleteGoal(id int64) error {
	db := database.GetDB()
	return db.Delete(&model.UserGoal{}, id).Error
}

// FindByID 根据ID查找目标
func (r *UserGoalRepository) FindByID(id int64) (*model.UserGoal, error) {
	db := database.GetDB()
	var goal model.UserGoal
	err := db.First(&goal, id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// FindByUserID 获取用户所有目标
func (r *UserGoalRepository) FindByUserID(userID int64) ([]model.UserGoal, error) {
	db := database.GetDB()
	var goals []model.UserGoal
	err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// FindActiveGoals 获取用户激活的目标
func (r *UserGoalRepository) FindActiveGoals(userID int64) ([]model.UserGoal, error) {
	db := database.GetDB()
	var goals []model.UserGoal
	err := db.Where("user_id = ? AND is_active = true", userID).Find(&goals).Error
	return goals, err
}

// UserCheckinRepository 用户打卡仓库
type UserCheckinRepository struct{}

// NewUserCheckinRepository 创建用户打卡仓库
func NewUserCheckinRepository() *UserCheckinRepository {
	return &UserCheckinRepository{}
}

// HasCheckedToday 检查今日是否已打卡特定目标
func (r *UserCheckinRepository) HasCheckedToday(userID, goalID int64, date string) (bool, error) {
	db := database.GetDB()
	var count int64
	err := db.Model(&model.UserCheckin{}).
		Where("user_id = ? AND goal_id = ? AND checkin_date = ?", userID, goalID, date).
		Count(&count).Error
	return count > 0, err
}

// CreateCheckin 创建打卡记录
func (r *UserCheckinRepository) CreateCheckin(checkin *model.UserCheckin) error {
	db := database.GetDB()
	return db.Create(checkin).Error
}

// Checkin 打卡
func (r *UserCheckinRepository) Checkin(userID int64, checkinType, note string) error {
	db := database.GetDB()
	today := time.Now().Format("2006-01-02")

	// 检查今天是否已打卡
	var existing model.UserCheckin
	err := db.Where("user_id = ? AND checkin_date = ? AND checkin_type = ?", userID, today, checkinType).First(&existing).Error
	if err == nil {
		return errors.New("今天已打卡")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建打卡记录
	checkin := &model.UserCheckin{
		UserID:      userID,
		CheckinDate: today,
		CheckinType: checkinType,
		Note:        note,
		CreatedAt:   time.Now(),
	}
	return db.Create(checkin).Error
}

// GetTodayCheckin 获取今日打卡状态
func (r *UserCheckinRepository) GetTodayCheckin(userID int64) ([]model.UserCheckin, error) {
	db := database.GetDB()
	today := time.Now().Format("2006-01-02")
	var checkins []model.UserCheckin
	err := db.Where("user_id = ? AND checkin_date = ?", userID, today).Find(&checkins).Error
	return checkins, err
}

// GetCheckinHistory 获取打卡历史
func (r *UserCheckinRepository) GetCheckinHistory(userID int64, days int) ([]model.UserCheckin, error) {
	db := database.GetDB()
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var checkins []model.UserCheckin
	err := db.Where("user_id = ? AND checkin_date >= ?", userID, startDate).
		Order("checkin_date DESC").Find(&checkins).Error
	return checkins, err
}

// GetCheckinStreak 获取连续打卡天数
func (r *UserCheckinRepository) GetCheckinStreak(userID int64, checkinType string) (int, error) {
	db := database.GetDB()

	// 获取最近30天的打卡记录
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	var checkins []model.UserCheckin

	query := db.Where("user_id = ? AND checkin_date >= ?", userID, startDate)
	if checkinType != "" {
		query = query.Where("checkin_type = ?", checkinType)
	}
	err := query.Order("checkin_date DESC").Find(&checkins).Error
	if err != nil {
		return 0, err
	}

	if len(checkins) == 0 {
		return 0, nil
	}

	// 计算连续天数
	streak := 0
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 创建日期集合
	dateSet := make(map[string]bool)
	for _, c := range checkins {
		dateSet[c.CheckinDate] = true
	}

	// 从今天或昨天开始计算
	currentDate := today
	if !dateSet[today] {
		if !dateSet[yesterday] {
			return 0, nil
		}
		currentDate = yesterday
	}

	for {
		if dateSet[currentDate] {
			streak++
			// 前一天
			t, _ := time.Parse("2006-01-02", currentDate)
			currentDate = t.AddDate(0, 0, -1).Format("2006-01-02")
		} else {
			break
		}
	}

	return streak, nil
}

// GetCheckinCalendar 获取打卡日历数据
func (r *UserCheckinRepository) GetCheckinCalendar(userID int64, year, month int) ([]map[string]interface{}, error) {
	db := database.GetDB()

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	endDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")

	var checkins []model.UserCheckin
	err := db.Where("user_id = ? AND checkin_date >= ? AND checkin_date < ?", userID, startDate, endDate).
		Find(&checkins).Error
	if err != nil {
		return nil, err
	}

	// 按日期分组
	dateMap := make(map[string][]string)
	for _, c := range checkins {
		dateMap[c.CheckinDate] = append(dateMap[c.CheckinDate], c.CheckinType)
	}

	// 转换为结果
	result := make([]map[string]interface{}, 0, len(dateMap))
	for date, types := range dateMap {
		result = append(result, map[string]interface{}{
			"date":          date,
			"checkin_types": types,
			"count":         len(types),
		})
	}

	return result, nil
}

// GoalProgressRepository 目标进度仓库
type GoalProgressRepository struct{}

// NewGoalProgressRepository 创建目标进度仓库
func NewGoalProgressRepository() *GoalProgressRepository {
	return &GoalProgressRepository{}
}

// UpdateProgress 更新目标进度
func (r *GoalProgressRepository) UpdateProgress(userID, goalID int64, value int) error {
	db := database.GetDB()
	today := time.Now().Format("2006-01-02")

	var progress model.GoalProgress
	err := db.Where("user_id = ? AND goal_id = ? AND progress_date = ?", userID, goalID, today).First(&progress).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新记录
		progress = model.GoalProgress{
			UserID:       userID,
			GoalID:       goalID,
			ProgressDate: today,
			CurrentValue: value,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		return db.Create(&progress).Error
	}

	if err != nil {
		return err
	}

	// 更新现有记录
	progress.CurrentValue = value
	progress.UpdatedAt = time.Now()
	return db.Save(&progress).Error
}

// GetTodayProgress 获取今日目标进度
func (r *GoalProgressRepository) GetTodayProgress(userID int64) ([]model.GoalProgress, error) {
	db := database.GetDB()
	today := time.Now().Format("2006-01-02")
	var progress []model.GoalProgress
	err := db.Where("user_id = ? AND progress_date = ?", userID, today).Find(&progress).Error
	return progress, err
}

// GetProgressHistory 获取目标进度历史
func (r *GoalProgressRepository) GetProgressHistory(userID, goalID int64, days int) ([]model.GoalProgress, error) {
	db := database.GetDB()
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var progress []model.GoalProgress
	err := db.Where("user_id = ? AND goal_id = ? AND progress_date >= ?", userID, goalID, startDate).
		Order("progress_date DESC").Find(&progress).Error
	return progress, err
}
