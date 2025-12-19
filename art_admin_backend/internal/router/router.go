package router

import (
	"art_admin_backend/internal/api"
	"art_admin_backend/internal/api/middleware"
	v1 "art_admin_backend/internal/api/v1"

	"github.com/gin-gonic/gin"
)

var chatAPI *api.ChatAPI

// SetChatAPI 设置聊天室 API（由 main 函数调用）
func SetChatAPI(api *api.ChatAPI) {
	chatAPI = api
}

// RegisterRoutes 注册所有路由
// 基础功能 + 系统功能 + Dify AI + 聊天室
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
			// 用户效果图配额管理（管理员）/ User render quota management (admin)
			authApiGroup.GET("/system/user/render-quota", v1.GetUserRenderQuotaByAdmin)
			authApiGroup.PUT("/system/user/render-quota", v1.UpdateUserRenderQuota)
			authApiGroup.POST("/system/user/render-quota/reset", v1.ResetUserRenderQuota)

			// 通知管理 / Notification management
			authApiGroup.GET("/notification/list", v1.GetNotificationList)
			authApiGroup.GET("/notification/unread-count", v1.GetUnreadCount)
			authApiGroup.POST("/notification/read/:id", v1.MarkNotificationRead)
			authApiGroup.POST("/notification/read-all", v1.MarkAllNotificationsRead)
			authApiGroup.DELETE("/notification/:id", v1.DeleteNotification)

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

			// 文件管理模块
			fileGroup := authApiGroup.Group("/file")
			{
				fileGroup.POST("/upload", v1.UploadFile)
				fileGroup.POST("/upload-multiple", v1.UploadMultipleFiles)
				fileGroup.GET("/list", v1.GetFileList)
				fileGroup.GET("/:id", v1.GetFileDetail)
				fileGroup.DELETE("/:id", v1.DeleteFile)
				fileGroup.POST("/batch-delete", v1.BatchDeleteFiles)
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

			// 设计师AI辅助系统模块 / Designer AI Assistant Module
			designerGroup := authApiGroup.Group("/designer")
			{
				// 项目管理 / Project Management
				designerGroup.POST("/projects", v1.CreateDesignerProject)
				designerGroup.GET("/projects", v1.GetDesignerProjectList)
				designerGroup.GET("/projects/search", v1.SearchDesignerProjects)
				designerGroup.GET("/projects/stats", v1.GetDesignerProjectStats)
				designerGroup.GET("/projects/:id", v1.GetDesignerProjectDetail)
				designerGroup.PUT("/projects/:id", v1.UpdateDesignerProject)
				designerGroup.DELETE("/projects/:id", v1.DeleteDesignerProject)
				designerGroup.POST("/projects/batch-delete", v1.BatchDeleteDesignerProjects)

				// 项目材料清单管理 / Project Material Management
				projectMaterialAPI := v1.NewProjectMaterialAPI()
				designerGroup.GET("/projects/:id/materials", projectMaterialAPI.GetProjectMaterials)
				designerGroup.POST("/projects/:id/materials", projectMaterialAPI.SaveProjectMaterials)
				designerGroup.POST("/projects/:id/materials/add", projectMaterialAPI.AddProjectMaterial)
				designerGroup.PUT("/projects/:id/materials/:itemId", projectMaterialAPI.UpdateProjectMaterial)
				designerGroup.DELETE("/projects/:id/materials/:itemId", projectMaterialAPI.DeleteProjectMaterial)
				designerGroup.POST("/projects/:id/materials/import", projectMaterialAPI.ImportMaterialsToProject)
				designerGroup.POST("/projects/:id/materials/export", projectMaterialAPI.ExportProjectMaterials)
				designerGroup.GET("/projects/:id/cost", projectMaterialAPI.GetProjectCostSummary)
				designerGroup.POST("/projects/:id/cost", projectMaterialAPI.SaveProjectCost)

				// 文档分析 / Document Analysis
				designerGroup.POST("/documents/upload", v1.UploadDocument)
				designerGroup.POST("/documents/analyze", v1.AnalyzeDocument)
				designerGroup.GET("/documents", v1.GetDocumentList)
				designerGroup.GET("/documents/:id", v1.GetDocumentDetail)
				designerGroup.GET("/documents/:id/keywords", v1.GetDocumentKeywords)
				designerGroup.GET("/documents/:id/summary", v1.GetDocumentSummary)
				designerGroup.DELETE("/documents/:id", v1.DeleteDocument)
				designerGroup.POST("/documents/batch-delete", v1.BatchDeleteDocuments)

				// 设计比对 / Design Compare
				designerGroup.POST("/compare/upload", v1.UploadDesignImages)
				designerGroup.POST("/compare/analyze", v1.AnalyzeDesignCompare)
				designerGroup.GET("/compare", v1.GetDesignCompareList)
				designerGroup.GET("/compare/:id", v1.GetDesignCompareDetail)
				designerGroup.GET("/compare/:id/result", v1.GetDesignCompareResult)
				designerGroup.DELETE("/compare/:id", v1.DeleteDesignCompare)
				designerGroup.POST("/compare/batch-delete", v1.BatchDeleteDesignCompare)

				// CAD图纸预览 / CAD Drawing Preview
				designerGroup.POST("/cad/upload", v1.UploadCadFile)
				designerGroup.POST("/cad/parse", v1.ParseCadFile)
				designerGroup.GET("/cad", v1.GetCadFileList)
				designerGroup.GET("/cad/:id", v1.GetCadFileDetail)
				designerGroup.GET("/cad/:id/parse", v1.GetCadParseResult)
				designerGroup.GET("/cad/:id/layers", v1.GetCadLayers)
				designerGroup.DELETE("/cad/:id", v1.DeleteCadFile)
				designerGroup.POST("/cad/batch-delete", v1.BatchDeleteCadFiles)

				// CAD效果图渲染（异步）/ CAD Rendering (async)
				designerGroup.POST("/cad/render", v1.GenerateCadRender)
				designerGroup.GET("/cad/render/history", v1.GetRenderHistory)
				designerGroup.GET("/cad/render/styles", v1.GetRenderStyles)
				// 效果图配额与记录 / Render quota and records
				designerGroup.GET("/cad/render/quota", v1.GetUserRenderQuota)
				designerGroup.GET("/cad/render/records", v1.ListRenderRecords)
				designerGroup.DELETE("/cad/render/records/:id", v1.DeleteRenderRecord)

				// 材料库管理 / Material Library Management
				designerGroup.POST("/materials", v1.CreateMaterial)
				designerGroup.PUT("/materials", v1.UpdateMaterial)
				designerGroup.GET("/materials", v1.GetMaterialList)
				designerGroup.GET("/materials/categories", v1.GetMaterialCategories)
				designerGroup.GET("/materials/brands", v1.GetMaterialBrands)
				designerGroup.GET("/materials/stats", v1.GetMaterialStats)
				designerGroup.GET("/materials/:id", v1.GetMaterialDetail)
				designerGroup.DELETE("/materials/:id", v1.DeleteMaterial)
				designerGroup.POST("/materials/batch-delete", v1.BatchDeleteMaterials)
				designerGroup.POST("/materials/batch-import", v1.BatchImportMaterials)
				designerGroup.POST("/materials/calculate", v1.CalculateMaterialUsage)
				designerGroup.POST("/materials/batch-calculate", v1.BatchCalculateMaterialUsage)
				designerGroup.POST("/materials/recommend", v1.RecommendMaterials)

				// 成本估算管理 / Cost Estimate Management
				costAPI := v1.NewCostAPI()
				designerGroup.POST("/cost", costAPI.CreateCostEstimate)
				designerGroup.PUT("/cost", costAPI.UpdateCostEstimate)
				designerGroup.GET("/cost/project", costAPI.GetCostEstimateByProjectID)
				designerGroup.GET("/cost/summary", costAPI.GetCostSummary)
				designerGroup.POST("/cost/calculate", costAPI.CalculateCost)
				designerGroup.POST("/cost/export", costAPI.ExportCostReport)
				designerGroup.GET("/cost/:id", costAPI.GetCostEstimateByID)
				designerGroup.DELETE("/cost/:id", costAPI.DeleteCostEstimate)

				// AI对话模块 / AI Chat Module
				// Requirements: 5.1, 5.2, 5.4
				designerGroup.POST("/chat", v1.DesignerChat)
				designerGroup.GET("/chat/history", v1.GetDesignerChatHistory)
				designerGroup.GET("/chat/sessions", v1.GetDesignerChatSessions)
				designerGroup.DELETE("/chat/sessions/:sessionId", v1.DeleteDesignerChatSession)
				designerGroup.GET("/chat/export/:sessionId", v1.ExportDesignerChatHistory)

				// 智能设计建议模块 / Design Suggestion Module
				// Requirements: 3.1, 3.2, 3.3, 3.4
				designerGroup.POST("/suggestions/generate", v1.GenerateDesignSuggestions)
				designerGroup.GET("/suggestions/task/:taskId", v1.GetSuggestionTaskStatus)
				designerGroup.GET("/suggestions", v1.GetDesignSuggestionList)
				designerGroup.GET("/suggestions/stats", v1.GetDesignSuggestionStats)
				designerGroup.GET("/suggestions/:id", v1.GetDesignSuggestionDetail)
				designerGroup.PUT("/suggestions/status", v1.UpdateDesignSuggestionStatus)
				designerGroup.PUT("/suggestions/batch-status", v1.BatchUpdateDesignSuggestionStatus)
				designerGroup.DELETE("/suggestions/:id", v1.DeleteDesignSuggestion)
				designerGroup.POST("/suggestions/batch-delete", v1.BatchDeleteDesignSuggestions)

				// 设计规范合规检查模块 / Design Compliance Check Module
				// Requirements: 11.1, 11.2, 11.3, 11.4, 11.5
				complianceAPI := v1.NewComplianceAPI()
				// 设计规范管理 / Design Standard Management
				designerGroup.POST("/compliance/standards", complianceAPI.CreateDesignStandard)
				designerGroup.PUT("/compliance/standards", complianceAPI.UpdateDesignStandard)
				designerGroup.GET("/compliance/standards", complianceAPI.GetDesignStandardList)
				designerGroup.GET("/compliance/standards/stats", complianceAPI.GetDesignStandardStats)
				designerGroup.GET("/compliance/standards/:id", complianceAPI.GetDesignStandardDetail)
				designerGroup.DELETE("/compliance/standards/:id", complianceAPI.DeleteDesignStandard)
				designerGroup.POST("/compliance/standards/batch-delete", complianceAPI.BatchDeleteDesignStandards)
				// 合规检查 / Compliance Check
				designerGroup.POST("/compliance/check", complianceAPI.PerformComplianceCheck)
				designerGroup.GET("/compliance/results", complianceAPI.GetComplianceCheckList)
				designerGroup.GET("/compliance/results/:id", complianceAPI.GetComplianceCheckDetail)
				designerGroup.GET("/compliance/latest", complianceAPI.GetLatestComplianceCheck)
				designerGroup.GET("/compliance/summary", complianceAPI.GetComplianceCheckSummary)
				designerGroup.DELETE("/compliance/results/:id", complianceAPI.DeleteComplianceCheck)
				designerGroup.POST("/compliance/results/batch-delete", complianceAPI.BatchDeleteComplianceChecks)

				// 施工图标注模块 / Construction Annotation Module
				// Requirements: 9.1, 9.2, 9.3, 9.4, 9.5
				// Note: Controller needs to be initialized with service dependencies
				// annotationAPI := v1.NewConstructionAnnotationController(annotationService)
				// designerGroup.POST("/construction-annotation/analyze", annotationAPI.AnalyzeConstructionDrawing)
				// designerGroup.GET("/construction-annotation/list", annotationAPI.GetAnnotationsList)
				// designerGroup.GET("/construction-annotation/:id", annotationAPI.GetAnnotation)
				// designerGroup.PUT("/construction-annotation/:id/annotations", annotationAPI.UpdateAnnotations)
				// designerGroup.DELETE("/construction-annotation/:id", annotationAPI.DeleteAnnotation)
				// designerGroup.POST("/construction-annotation/:id/items", annotationAPI.AddAnnotationItem)
				// designerGroup.GET("/construction-annotation/:id/items/:itemId", annotationAPI.GetAnnotationItem)
				// designerGroup.PUT("/construction-annotation/:id/items/:itemId", annotationAPI.UpdateAnnotationItem)
				// designerGroup.DELETE("/construction-annotation/:id/items/:itemId", annotationAPI.DeleteAnnotationItem)
				// designerGroup.POST("/construction-annotation/:id/batch", annotationAPI.BatchUpdateAnnotations)
				// designerGroup.POST("/construction-annotation/export", annotationAPI.ExportAnnotation)
				// designerGroup.GET("/construction-annotation/export-formats", annotationAPI.GetExportFormats)
				// designerGroup.GET("/construction-annotation/:id/export-history", annotationAPI.GetExportHistory)

				// AI生成CAD模块 / AI CAD Generation Module
				// Requirements: AI based CAD generation from project documents
				cadGenAPI := v1.NewCadGenerationAPI()
				designerGroup.POST("/cad-generations", cadGenAPI.Create)
				designerGroup.GET("/cad-generations", cadGenAPI.List)
				designerGroup.GET("/cad-generations/type-options", cadGenAPI.GetTypeOptions)
				designerGroup.GET("/cad-generations/:id", cadGenAPI.GetByID)
				designerGroup.DELETE("/cad-generations/:id", cadGenAPI.Delete)
				designerGroup.POST("/cad-generations/:id/retry", cadGenAPI.Retry)
				designerGroup.POST("/cad-generations/confirm", cadGenAPI.Confirm)

				// 设计版本管理模块 / Design Version Management Module
				// Requirements: 7.1, 7.2, 7.3, 7.4
				designerGroup.POST("/versions", v1.CreateDesignVersion)
				designerGroup.GET("/versions", v1.GetDesignVersionList)
				designerGroup.GET("/versions/:id", v1.GetDesignVersionDetail)
				designerGroup.PUT("/versions/:id", v1.UpdateDesignVersion)
				designerGroup.DELETE("/versions/:id", v1.DeleteDesignVersion)
				designerGroup.POST("/versions/compare", v1.CompareDesignVersions)
				designerGroup.GET("/versions/compares", v1.GetVersionCompareList)
				designerGroup.GET("/versions/compares/:id", v1.GetVersionCompareDetail)
				designerGroup.GET("/versions/diff", v1.GetVersionDiff)
				designerGroup.DELETE("/versions/compares/:id", v1.DeleteVersionCompare)
			}
		}
	}

	// 聊天室模块
	if chatAPI != nil {
		SetupChatRoutes(r, chatAPI)
	}
}
