package router

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/api/travel"

	"github.com/gin-gonic/gin"
)

// SetupTravelRoutes 设置旅游模块路由
func SetupTravelRoutes(r *gin.Engine) {
	// 旅游模块路由组
	travelGroup := r.Group("/api/v1/travel")
	travelGroup.Use(middleware.JWTAuth())
	travelGroup.Use(middleware.OperationLog())
	{
		// 路书管理
		travelGroup.POST("/roadbooks", travel.CreateRoadbook)
		travelGroup.GET("/roadbooks", travel.GetRoadbookList)
		travelGroup.GET("/roadbooks/:id", travel.GetRoadbookDetail)
		travelGroup.PUT("/roadbooks/:id", travel.UpdateRoadbook)
		travelGroup.DELETE("/roadbooks/:id", travel.DeleteRoadbook)

		// 途经点管理
		travelGroup.GET("/roadbooks/:id/waypoints", travel.GetWaypoints)
		travelGroup.GET("/roadbooks/:id/waypoints/grouped", travel.GetWaypointsGroupedByDay)
		travelGroup.POST("/roadbooks/:id/waypoints", travel.AddWaypoint)
		travelGroup.PUT("/roadbooks/:id/waypoints", travel.BatchUpdateWaypoints)
		travelGroup.PUT("/roadbooks/:id/waypoints/sort", travel.UpdateWaypointSortOrder)
		travelGroup.PUT("/roadbooks/:id/waypoints/:waypointId", travel.UpdateWaypoint)
		travelGroup.DELETE("/roadbooks/:id/waypoints/:waypointId", travel.DeleteWaypoint)

		// 评论管理
		travelGroup.POST("/roadbooks/:id/comments", travel.CreateComment)
		travelGroup.GET("/roadbooks/:id/comments", travel.GetCommentList)
		travelGroup.DELETE("/comments/:id", travel.DeleteComment)
		travelGroup.POST("/comments/:id/like", travel.LikeComment)

		// 收藏管理
		travelGroup.POST("/roadbooks/:id/favorite", travel.FavoriteRoadbook)
		travelGroup.DELETE("/roadbooks/:id/favorite", travel.UnfavoriteRoadbook)
		travelGroup.GET("/favorites", travel.GetFavoriteList)

		// 分享管理
		travelGroup.POST("/roadbooks/:id/share", travel.CreateShare)
		travelGroup.GET("/share/:code", travel.GetShareDetail)

		// 导航导出
		travelGroup.GET("/roadbooks/:id/nav-link", travel.GetNavLink)

		// 路书广场
		travelGroup.GET("/explore", travel.GetExploreList)
		travelGroup.GET("/search", travel.SearchRoadbooks)

		// 标签管理
		travelGroup.GET("/tags", travel.GetTagList)
		travelGroup.POST("/tags", travel.CreateTag)

		// AI规划
		travelGroup.POST("/ai/plan", travel.AIPlan)
		travelGroup.POST("/ai/chat", travel.AIChat)

		// 模板管理
		travelGroup.GET("/templates", travel.GetTemplateList)
		travelGroup.POST("/templates/:id/use", travel.UseTemplate)

		// 统计
		travelGroup.GET("/roadbooks/:id/stats", travel.GetRoadbookStats)
		travelGroup.GET("/my/stats", travel.GetMyStats)
	}
}
