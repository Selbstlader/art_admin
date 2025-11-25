package router

import (
	"art_admin_backend/internal/api"
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/api/project"
	v1 "art_admin_backend/internal/api/v1"

	"github.com/gin-gonic/gin"
)

var learningAPI *api.LearningAPI
var chatAPI *api.ChatAPI

// SetLearningAPI 设置学习系统 API（由 main 函数调用）
func SetLearningAPI(api *api.LearningAPI) {
	learningAPI = api
}

// SetChatAPI 设置聊天室 API（由 main 函数调用）
func SetChatAPI(api *api.ChatAPI) {
	chatAPI = api
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
