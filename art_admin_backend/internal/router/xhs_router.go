package router

import (
	"art_admin_backend/internal/api"
	"art_admin_backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupXHSRoutes 设置小红书总结路由
func SetupXHSRoutes(r *gin.Engine, xhsAPI *api.XHSAPI) {
	// /api/xhs 路由组
	xhsGroup := r.Group("/api/xhs")
	xhsGroup.Use(middleware.JWTAuth())
	{
		// 核心功能
		xhsGroup.POST("/summarize", xhsAPI.Summarize) // 生成总结
		xhsGroup.POST("/ocr", xhsAPI.OCR)             // OCR识别

		// 历史记录
		xhsGroup.GET("/history", xhsAPI.GetHistory)           // 获取总结历史
		xhsGroup.GET("/summary/:id", xhsAPI.GetSummaryDetail) // 获取总结详情
		xhsGroup.DELETE("/summary/:id", xhsAPI.DeleteSummary) // 删除总结

		// 收藏
		xhsGroup.PUT("/summary/:id/favorite", xhsAPI.UpdateFavorite) // 更新收藏状态
		xhsGroup.GET("/favorites", xhsAPI.GetFavorites)              // 获取收藏列表

		// 搜索
		xhsGroup.GET("/search", xhsAPI.SearchSummaries) // 搜索总结

		// 收藏夹管理
		xhsGroup.POST("/folder", xhsAPI.CreateFolder)       // 创建收藏夹
		xhsGroup.GET("/folders", xhsAPI.GetFolders)         // 获取收藏夹列表
		xhsGroup.PUT("/folder/:id", xhsAPI.UpdateFolder)    // 更新收藏夹
		xhsGroup.DELETE("/folder/:id", xhsAPI.DeleteFolder) // 删除收藏夹
	}
}
