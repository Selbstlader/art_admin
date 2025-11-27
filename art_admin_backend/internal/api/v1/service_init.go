package v1

import (
	"art_admin_backend/internal/service"
	achievementSvc "art_admin_backend/internal/service/achievement"
)

// 全局服务实例 - 统一声明避免重定义
var (
	moodRecordService         *service.MoodRecordService
	meditationRecordService   *service.MeditationRecordService
	journalEntryService       *service.JournalEntryService
	aiAnalysisService         *service.AIAnalysisService
	achievementService        *achievementSvc.AchievementService
	meditationFavoriteService *service.MeditationFavoriteService
)

// SetServices 设置所有服务实例
func SetServices(
	moodService *service.MoodRecordService,
	meditationService *service.MeditationRecordService,
	journalService *service.JournalEntryService,
	aiService *service.AIAnalysisService,
	achieveService *achievementSvc.AchievementService,
	meditationFavService *service.MeditationFavoriteService,
) {
	moodRecordService = moodService
	meditationRecordService = meditationService
	journalEntryService = journalService
	aiAnalysisService = aiService
	achievementService = achieveService
	meditationFavoriteService = meditationFavService
}

// GetMoodRecordService 获取情绪记录服务
func GetMoodRecordService() *service.MoodRecordService {
	return moodRecordService
}

// GetMeditationRecordService 获取冥想记录服务
func GetMeditationRecordService() *service.MeditationRecordService {
	return meditationRecordService
}

// GetJournalEntryService 获取日记服务
func GetJournalEntryService() *service.JournalEntryService {
	return journalEntryService
}

// GetAIAnalysisService 获取AI分析服务
func GetAIAnalysisService() *service.AIAnalysisService {
	return aiAnalysisService
}

// GetAchievementService 获取成就服务
func GetAchievementService() *achievementSvc.AchievementService {
	return achievementSvc.NewAchievementService()
}
