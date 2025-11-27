package router

import (
	"art_admin_backend/internal/api"
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/api/project"
	v1 "art_admin_backend/internal/api/v1"
	achievementSvc "art_admin_backend/internal/service/achievement"
	moodAnalyticsSvc "art_admin_backend/internal/service/mood_analytics"

	"github.com/gin-gonic/gin"
)

var learningAPI *api.LearningAPI
var chatAPI *api.ChatAPI
var aiAnalysisService *moodAnalyticsSvc.AIAnalysisService
var achievementService *achievementSvc.AchievementService

// SetLearningAPI 设置学习系统 API（由 main 函数调用）
func SetLearningAPI(api *api.LearningAPI) {
	learningAPI = api
}

// SetChatAPI 设置聊天室 API（由 main 函数调用）
func SetChatAPI(api *api.ChatAPI) {
	chatAPI = api
}

// SetAIAnalysisService 设置AI分析服务（由 main 函数调用）
func SetAIAnalysisService(service *moodAnalyticsSvc.AIAnalysisService) {
	aiAnalysisService = service
}

// SetAchievementService 设置成就服务（由 main 函数调用）
func SetAchievementService(service *achievementSvc.AchievementService) {
	achievementService = service
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 统一使用 /api 前缀
	apiGroup := r.Group("/api")
	{
		// 认证相关路由（无需JWT认证）
		apiGroup.POST("/auth/login", v1.Login)

		// 需要JWT认证的路由
		authApiGroup := apiGroup.Group("")
		authApiGroup.Use(middleware.JWTAuth())
		authApiGroup.Use(middleware.OperationLog()) // 操作日志中间件
		{
			// 用户信息
			authApiGroup.GET("/user/info", v1.GetUserInfo)
			authApiGroup.PUT("/user/info", v1.UpdateUserInfo)
			authApiGroup.POST("/user/change-password", v1.ChangePassword)

			// 用户管理
			authApiGroup.GET("/user/list", v1.GetUserList)
			authApiGroup.POST("/user/create", v1.CreateUser)
			authApiGroup.PUT("/user/update", v1.UpdateUser)
			authApiGroup.DELETE("/user/delete/:id", v1.DeleteUser)
			authApiGroup.POST("/user/reset-password", v1.ResetPassword)

			// APP用户管理
			authApiGroup.GET("/app-user/list", v1.GetAppUserList)
			authApiGroup.POST("/app-user/create", v1.CreateAppUser)
			authApiGroup.PUT("/app-user/update", v1.UpdateAppUser)
			authApiGroup.DELETE("/app-user/delete/:id", v1.DeleteAppUser)
			authApiGroup.POST("/app-user/reset-password", v1.ResetAppUserPassword)

			// 角色管理
			authApiGroup.GET("/role/list", v1.GetRoleList)
			authApiGroup.POST("/role/create", v1.CreateRole)
			authApiGroup.PUT("/role/update", v1.UpdateRole)
			authApiGroup.DELETE("/role/delete/:id", v1.DeleteRole)
			authApiGroup.GET("/role/permissions/:id", v1.GetRolePermissions)
			authApiGroup.PUT("/role/permissions/:id", v1.UpdateRolePermissions)

			// 部门管理
			authApiGroup.GET("/department/list", v1.GetDepartmentList)
			authApiGroup.GET("/department/:id", v1.GetDepartmentByID)
			authApiGroup.POST("/department", v1.CreateDepartment)
			authApiGroup.PUT("/department", v1.UpdateDepartment)
			authApiGroup.DELETE("/department/:id", v1.DeleteDepartment)

			// 菜单管理
			authApiGroup.GET("/system/menus", v1.GetMenuList)
			authApiGroup.GET("/menu/all", v1.GetAllMenus)
			authApiGroup.POST("/menu/create", v1.CreateMenu)
			authApiGroup.PUT("/menu/update", v1.UpdateMenu)
			authApiGroup.DELETE("/menu/delete/:id", v1.DeleteMenu)

			// 字典类型管理
			authApiGroup.GET("/dictionary/type/list", v1.GetDictionaryTypeList)
			authApiGroup.GET("/dictionary/type/:id", v1.GetDictionaryTypeByID)
			authApiGroup.POST("/dictionary/type", v1.CreateDictionaryType)
			authApiGroup.PUT("/dictionary/type", v1.UpdateDictionaryType)
			authApiGroup.DELETE("/dictionary/type", v1.DeleteDictionaryType)

			// 字典数据管理
			authApiGroup.GET("/dictionary/list", v1.GetDictionaryList)
			authApiGroup.GET("/dictionary/:id", v1.GetDictionaryByID)
			authApiGroup.GET("/dictionary/by-type", v1.GetDictionaryByTypeCode)
			authApiGroup.POST("/dictionary", v1.CreateDictionary)
			authApiGroup.PUT("/dictionary", v1.UpdateDictionary)
			authApiGroup.DELETE("/dictionary", v1.DeleteDictionary)

			// 操作日志管理
			authApiGroup.GET("/operation-log/list", v1.GetOperationLogList)
			authApiGroup.GET("/operation-log/:id", v1.GetOperationLogDetail)
			authApiGroup.DELETE("/operation-log/:id", v1.DeleteOperationLog)
			authApiGroup.DELETE("/operation-log/batch-delete", v1.BatchDeleteOperationLog)
			authApiGroup.POST("/operation-log/clean", v1.CleanOperationLog)

			// 项目管理模块 - 独立分区
			projectGroup := authApiGroup.Group("/project")
			{
				// 项目模板管理
				projectGroup.GET("/template/list", project.GetTemplateList)
				projectGroup.GET("/template/:id", project.GetTemplateDetail)
				projectGroup.POST("/template", project.CreateTemplate)
				projectGroup.PUT("/template", project.UpdateTemplate)
				projectGroup.DELETE("/template/:id", project.DeleteTemplate)
				projectGroup.POST("/template/dependency", project.CreateTemplateDependency)
				projectGroup.DELETE("/template/dependency/:id", project.DeleteTemplateDependency)

				// 项目管理
				projectGroup.GET("/list", project.GetProjectList)
				projectGroup.GET("/:id", project.GetProjectDetail)
				projectGroup.POST("", project.CreateProject)
				projectGroup.POST("/from-template", project.CreateProjectFromTemplate)
				projectGroup.PUT("", project.UpdateProject)
				projectGroup.DELETE("/:id", project.DeleteProject)
				projectGroup.GET("/statistics", project.GetProjectStatistics)

				// 任务管理
				projectGroup.GET("/task/list", project.GetTaskList)
				projectGroup.GET("/task/:id", project.GetTaskDetail)
				projectGroup.POST("/task", project.CreateTask)
				projectGroup.PUT("/task", project.UpdateTask)
				projectGroup.PUT("/task/batch", project.BatchUpdateTask)
				projectGroup.PUT("/task/quick", project.QuickUpdateTask)
				projectGroup.DELETE("/task/:id", project.DeleteTask)
				projectGroup.GET("/task/gantt", project.GetTaskGanttData)

				// 任务依赖管理
				projectGroup.POST("/task/dependency", project.CreateTaskDependency)
				projectGroup.DELETE("/task/dependency/:id", project.DeleteTaskDependency)

				// 任务评论
				projectGroup.POST("/task/comment", project.CreateTaskComment)
			}

			// 情绪管理模块
			moodGroup := authApiGroup.Group("/mood-record")
			{
				moodGroup.POST("/create", v1.CreateMoodRecord)
				moodGroup.PUT("/update", v1.UpdateMoodRecord)
				moodGroup.DELETE("/delete/:id", v1.DeleteMoodRecord)
				moodGroup.GET("/list", v1.GetMoodRecordList)
				moodGroup.GET("/statistics", v1.GetMoodStatistics)
				moodGroup.POST("/analytics", v1.GetMoodAnalytics)
				moodGroup.GET("/:id", v1.GetMoodRecordDetail)         // 新增：获取单条记录详情
				moodGroup.GET("/date/:date", v1.GetMoodRecordsByDate) // 新增：按日期查询
				moodGroup.GET("/calendar", v1.GetMoodCalendar)        // 新增：日历视图数据

				// AI分析相关接口
				moodGroup.POST("/analyze-patterns", v1.AnalyzeMoodPatterns)
				moodGroup.GET("/trends", v1.GetMoodTrends)
			}

			// 冥想记录管理模块
			meditationRecordGroup := authApiGroup.Group("/meditation-record")
			{
				meditationRecordGroup.POST("/create", v1.CreateMeditationRecord)
				meditationRecordGroup.PUT("/update", v1.UpdateMeditationRecord)
				meditationRecordGroup.DELETE("/delete/:id", v1.DeleteMeditationRecord)
				meditationRecordGroup.GET("/list", v1.GetMeditationRecordList)
				meditationRecordGroup.GET("/:id", v1.GetMeditationRecordDetail)
				meditationRecordGroup.GET("/statistics", v1.GetMeditationStatistics)

				// AI分析相关接口
				meditationRecordGroup.POST("/:id/analyze", v1.AnalyzeMeditationSession)
				meditationRecordGroup.GET("/trends", v1.GetMeditationTrends)
				meditationRecordGroup.GET("/recommendations", v1.GetMeditationRecommendations)
			}

			// 冥想内容管理模块
			meditationGroup := authApiGroup.Group("/meditation-content")
			{
				meditationGroup.POST("/create", v1.CreateMeditationContent)
				meditationGroup.PUT("/update", v1.UpdateMeditationContent)
				meditationGroup.DELETE("/delete/:id", v1.DeleteMeditationContent)
				meditationGroup.GET("/list", v1.GetMeditationContentList)
				meditationGroup.GET("/:id", v1.GetMeditationContentDetail)
				meditationGroup.GET("/category/:category", v1.GetMeditationContentByCategory)
				meditationGroup.GET("/difficulty/:difficulty", v1.GetMeditationContentByDifficulty)
				meditationGroup.GET("/popular", v1.GetPopularMeditationContent)
				meditationGroup.POST("/like/:id", v1.LikeMeditationContent)
				meditationGroup.POST("/unlike/:id", v1.UnlikeMeditationContent)
				meditationGroup.GET("/search", v1.SearchMeditationContent)
				meditationGroup.GET("/categories", v1.GetMeditationCategories)
				meditationGroup.GET("/difficulty-levels", v1.GetMeditationDifficultyLevels)

				// 收藏功能
				meditationGroup.POST("/:id/favorite", v1.AddMeditationFavorite)        // 新增：添加收藏
				meditationGroup.DELETE("/:id/favorite", v1.RemoveMeditationFavorite)   // 新增：取消收藏
				meditationGroup.GET("/:id/favorite/check", v1.CheckMeditationFavorite) // 新增：检查是否收藏
				meditationGroup.GET("/favorites", v1.GetUserFavorites)                 // 新增：获取收藏列表

				// 播放记录功能
				meditationGroup.POST("/play-record", v1.UpdatePlayRecord)       // 新增：更新播放记录
				meditationGroup.GET("/:id/play-record", v1.GetPlayRecord)       // 新增：获取播放记录
				meditationGroup.GET("/play-history", v1.GetUserPlayHistory)     // 新增：获取播放历史
				meditationGroup.GET("/recently-played", v1.GetRecentlyPlayed)   // 新增：最近播放
				meditationGroup.GET("/continue-playing", v1.GetContinuePlaying) // 新增：继续播放
				meditationGroup.GET("/play-stats", v1.GetPlayStats)             // 新增：播放统计
			}

			// 日记管理模块
			journalGroup := authApiGroup.Group("/journal-entry")
			{
				journalGroup.POST("/create", v1.CreateJournalEntry)
				journalGroup.PUT("/update", v1.UpdateJournalEntry)
				journalGroup.DELETE("/delete/:id", v1.DeleteJournalEntry)
				journalGroup.GET("/list", v1.GetJournalEntryList)
				journalGroup.GET("/:id", v1.GetJournalEntryDetail)
				journalGroup.GET("/date-range", v1.GetJournalEntriesByDateRange)
				journalGroup.GET("/recent", v1.GetRecentJournalEntries)
				journalGroup.GET("/search", v1.SearchJournalEntries)
				journalGroup.GET("/statistics", v1.GetJournalStatistics)
				journalGroup.PUT("/sentiment/:id", v1.UpdateSentimentScore)

				// 标签和图片管理接口
				journalGroup.GET("/tags", v1.GetAllUserTags)              // 新增：获取所有标签
				journalGroup.GET("/by-tag", v1.GetJournalsByTag)          // 新增：按标签查询
				journalGroup.PUT("/:id/tags", v1.UpdateJournalTags)       // 新增：更新标签
				journalGroup.POST("/:id/images", v1.AddJournalImages)     // 新增：添加图片
				journalGroup.DELETE("/:id/images", v1.RemoveJournalImage) // 新增：删除图片

				// AI分析相关接口
				journalGroup.POST("/:id/analyze", v1.AnalyzeJournalEntry)
				journalGroup.GET("/insights", v1.GetJournalInsights)
				journalGroup.GET("/trends", v1.GetJournalTrends)
			}

			// AI分析相关路由
			authApiGroup.POST("/ai-analysis/emotions", v1.AnalyzeUserEmotions)
			authApiGroup.GET("/ai-analysis/history", v1.GetEmotionAnalysisHistory)
			authApiGroup.GET("/ai-analysis/:id", v1.GetEmotionAnalysisDetail)
			authApiGroup.DELETE("/ai-analysis/:id", v1.DeleteEmotionAnalysis)
			authApiGroup.GET("/ai-analysis/:id/export", v1.ExportEmotionAnalysis)
			authApiGroup.POST("/ai-analysis/batch-analyze", v1.BatchAnalyze)
			authApiGroup.GET("/ai-analysis/status", v1.GetAnalysisStatus)

			// 成就系统相关路由
			authApiGroup.GET("/achievements", v1.GetAchievementList)
			authApiGroup.GET("/achievements/stats", v1.GetAchievementStats)
			authApiGroup.GET("/achievements/leaderboard", v1.GetLeaderboard)
			authApiGroup.GET("/achievements/recent", v1.GetRecentAchievements)
			authApiGroup.POST("/achievements/claim", v1.ClaimAchievement)

			// 目标管理模块
			goalGroup := authApiGroup.Group("/goals")
			{
				goalGroup.POST("", v1.CreateGoal)                     // 创建目标
				goalGroup.GET("", v1.GetUserGoals)                    // 获取目标列表
				goalGroup.PUT("/:id", v1.UpdateGoal)                  // 更新目标
				goalGroup.DELETE("/:id", v1.DeleteGoal)               // 删除目标
				goalGroup.PUT("/:id/progress", v1.UpdateGoalProgress) // 更新目标进度
				goalGroup.GET("/:id/progress", v1.GetGoalProgress)    // 获取目标进度历史
			}

			// 打卡管理模块
			checkinGroup := authApiGroup.Group("/checkin")
			{
				checkinGroup.POST("", v1.Checkin)                    // 打卡
				checkinGroup.GET("/today", v1.GetTodayCheckin)       // 获取今日打卡状态
				checkinGroup.GET("/history", v1.GetCheckinHistory)   // 获取打卡历史
				checkinGroup.GET("/streak", v1.GetCheckinStreak)     // 获取连续打卡天数
				checkinGroup.GET("/calendar", v1.GetCheckinCalendar) // 获取打卡日历
				checkinGroup.GET("/stats", v1.GetCheckinStats)       // 获取打卡统计
			}

			// Dify AI 模块
			difyGroup := authApiGroup.Group("/dify")
			{
				// 知识库管理
				difyGroup.GET("/dataset/list", v1.GetDatasetList)
				difyGroup.GET("/dataset/:id", v1.GetDatasetDetail)
				difyGroup.POST("/dataset/upload", v1.UploadFileToDataset)
				difyGroup.DELETE("/dataset/:id", v1.DeleteDataset)

				// AI对话
				difyGroup.POST("/chat", v1.ChatWithAI)
				difyGroup.POST("/chat/stream", v1.ChatWithAIStreaming)
				difyGroup.POST("/chat/stop/:task_id", v1.StopChatMessage)

				// 会话管理
				difyGroup.GET("/conversations", v1.GetConversations)
				difyGroup.GET("/messages", v1.GetConversationMessages)
				difyGroup.DELETE("/conversations/:conversation_id", v1.DeleteConversation)
				difyGroup.POST("/conversations/:conversation_id/name", v1.RenameConversation)

				// 建议问题
				difyGroup.GET("/messages/:message_id/suggested", v1.GetSuggestedQuestions)
			}
		}
	}

	// 学习系统模块（独立模块，不在 v1 中）
	if learningAPI != nil {
		SetupLearningRoutes(r, learningAPI)
	}

	// 聊天室模块
	if chatAPI != nil {
		SetupChatRoutes(r, chatAPI)
	}
}
