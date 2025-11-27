package v1

import (
	achievementSvc "art_admin_backend/internal/service/achievement"
	moodAnalyticsSvc "art_admin_backend/internal/service/mood_analytics"
)

// 全局服务实例 - 统一声明避免重定义
var (
	moodRecordService         *moodAnalyticsSvc.MoodRecordService
	meditationRecordService   *moodAnalyticsSvc.MeditationRecordService
	journalEntryService       *moodAnalyticsSvc.JournalEntryService
	aiAnalysisService         *moodAnalyticsSvc.AIAnalysisService
	achievementService        *achievementSvc.AchievementService
	meditationFavoriteService *moodAnalyticsSvc.MeditationFavoriteService
)

// SetServices 设置所有服务实例
func SetServices(
	moodService *moodAnalyticsSvc.MoodRecordService,
	meditationService *moodAnalyticsSvc.MeditationRecordService,
	journalService *moodAnalyticsSvc.JournalEntryService,
	aiService *moodAnalyticsSvc.AIAnalysisService,
	achieveService *achievementSvc.AchievementService,
	meditationFavService *moodAnalyticsSvc.MeditationFavoriteService,
) {
	moodRecordService = moodService
	meditationRecordService = meditationService
	journalEntryService = journalService
	aiAnalysisService = aiService
	achievementService = achieveService
	meditationFavoriteService = meditationFavService
}

// GetMoodRecordService 获取情绪记录服务
func GetMoodRecordService() *moodAnalyticsSvc.MoodRecordService {
	return moodRecordService
}

// GetMeditationRecordService 获取冥想记录服务
func GetMeditationRecordService() *moodAnalyticsSvc.MeditationRecordService {
	return meditationRecordService
}

// GetJournalEntryService 获取日记服务
func GetJournalEntryService() *moodAnalyticsSvc.JournalEntryService {
	return journalEntryService
}

// GetAIAnalysisService 获取AI分析服务
func GetAIAnalysisService() *moodAnalyticsSvc.AIAnalysisService {
	return aiAnalysisService
}

// GetAchievementService 获取成就服务
func GetAchievementService() *achievementSvc.AchievementService {
	return achievementSvc.NewAchievementService()
}
