package router

import (
	"art_admin_backend/internal/api/middleware"
	v1 "art_admin_backend/internal/api/v1"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
// 基础功能 + 系统功能
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
		}
	}
}
