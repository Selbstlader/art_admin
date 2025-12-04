package travel

import (
	"art_admin_backend/internal/pkg/response"
	travelSvc "art_admin_backend/internal/service/travel"
	"strconv"

	"github.com/gin-gonic/gin"
)

var waypointService = travelSvc.NewWaypointService()

// AddWaypoint 添加途经点
// POST /api/v1/travel/roadbooks/:id/waypoints
func AddWaypoint(c *gin.Context) {
	// 获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	var req travelSvc.AddWaypointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	req.RoadbookID = roadbookID

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")
	req.UserID = userID.(int64)

	waypoint, err := waypointService.AddWaypoint(&req)
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权操作该路书")
			return
		}
		if err.Error() == "roadbook not found" {
			response.NotFound(c, "路书不存在")
			return
		}
		response.ServerError(c, "添加途经点失败: "+err.Error())
		return
	}

	response.Success(c, waypoint)
}

// GetWaypoints 获取路书的途经点列表
// GET /api/v1/travel/roadbooks/:id/waypoints
func GetWaypoints(c *gin.Context) {
	// 获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")

	waypoints, err := waypointService.GetWaypointsByRoadbook(roadbookID, userID.(int64))
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权访问该路书")
			return
		}
		if err.Error() == "roadbook not found" {
			response.NotFound(c, "路书不存在")
			return
		}
		response.ServerError(c, "获取途经点失败: "+err.Error())
		return
	}

	response.Success(c, waypoints)
}

// GetWaypointsGroupedByDay 获取路书的途经点（按天分组）
// GET /api/v1/travel/roadbooks/:id/waypoints/grouped
func GetWaypointsGroupedByDay(c *gin.Context) {
	// 获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")

	grouped, err := waypointService.GetWaypointsGroupedByDay(roadbookID, userID.(int64))
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权访问该路书")
			return
		}
		if err.Error() == "roadbook not found" {
			response.NotFound(c, "路书不存在")
			return
		}
		response.ServerError(c, "获取途经点失败: "+err.Error())
		return
	}

	response.Success(c, grouped)
}

// UpdateWaypoint 更新途经点
// PUT /api/v1/travel/roadbooks/:id/waypoints/:waypointId
func UpdateWaypoint(c *gin.Context) {
	// 获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	// 获取途经点ID
	waypointIDStr := c.Param("waypointId")
	waypointID, err := strconv.ParseInt(waypointIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的途经点ID")
		return
	}

	var req travelSvc.UpdateWaypointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	req.ID = waypointID
	req.RoadbookID = roadbookID

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")
	req.UserID = userID.(int64)

	err = waypointService.UpdateWaypoint(&req)
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权操作该路书")
			return
		}
		if err.Error() == "roadbook not found" {
			response.NotFound(c, "路书不存在")
			return
		}
		if err.Error() == "waypoint not found" {
			response.NotFound(c, "途经点不存在")
			return
		}
		response.ServerError(c, "更新途经点失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchUpdateWaypoints 批量更新途经点
// PUT /api/v1/travel/roadbooks/:id/waypoints
func BatchUpdateWaypoints(c *gin.Context) {
	// 获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	var req travelSvc.BatchUpdateWaypointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	req.RoadbookID = roadbookID

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")
	req.UserID = userID.(int64)

	err = waypointService.BatchUpdateWaypoints(&req)
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权操作该路书")
			return
		}
		if err.Error() == "roadbook not found" {
			response.NotFound(c, "路书不存在")
			return
		}
		response.ServerError(c, "批量更新途经点失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateWaypointSortOrder 更新途经点排序
// PUT /api/v1/travel/roadbooks/:id/waypoints/sort
func UpdateWaypointSortOrder(c *gin.Context) {
	// 获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	var req travelSvc.UpdateSortOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	req.RoadbookID = roadbookID

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")
	req.UserID = userID.(int64)

	err = waypointService.UpdateSortOrder(&req)
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权操作该路书")
			return
		}
		if err.Error() == "roadbook not found" {
			response.NotFound(c, "路书不存在")
			return
		}
		response.ServerError(c, "更新排序失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteWaypoint 删除途经点
// DELETE /api/v1/travel/roadbooks/:id/waypoints/:waypointId
func DeleteWaypoint(c *gin.Context) {
	// 获取路书ID
	roadbookIDStr := c.Param("id")
	roadbookID, err := strconv.ParseInt(roadbookIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	// 获取途经点ID
	waypointIDStr := c.Param("waypointId")
	waypointID, err := strconv.ParseInt(waypointIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的途经点ID")
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")

	err = waypointService.DeleteWaypoint(roadbookID, waypointID, userID.(int64))
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权操作该路书")
			return
		}
		if err.Error() == "roadbook not found" {
			response.NotFound(c, "路书不存在")
			return
		}
		if err.Error() == "waypoint not found" {
			response.NotFound(c, "途经点不存在")
			return
		}
		response.ServerError(c, "删除途经点失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
