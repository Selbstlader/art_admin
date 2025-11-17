package router

import (
	"art_admin_backend/internal/api/middleware"
	v1 "art_admin_backend/internal/api/v1"

	"github.com/gin-gonic/gin"
)

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
		{
			// 用户信息
			authApiGroup.GET("/user/info", v1.GetUserInfo)

			// 用户管理
			authApiGroup.GET("/user/list", v1.GetUserList)
			authApiGroup.POST("/user/create", v1.CreateUser)
			authApiGroup.PUT("/user/update", v1.UpdateUser)
			authApiGroup.DELETE("/user/delete/:id", v1.DeleteUser)
			authApiGroup.POST("/user/reset-password", v1.ResetPassword)

			// 角色管理
			authApiGroup.GET("/role/list", v1.GetRoleList)
			authApiGroup.POST("/role/create", v1.CreateRole)
			authApiGroup.PUT("/role/update", v1.UpdateRole)
			authApiGroup.DELETE("/role/delete/:id", v1.DeleteRole)
			authApiGroup.GET("/role/permissions/:id", v1.GetRolePermissions)
			authApiGroup.PUT("/role/permissions/:id", v1.UpdateRolePermissions)

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
		}
	}
}
