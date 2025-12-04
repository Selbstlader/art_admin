package travel

import (
	"art_admin_backend/internal/pkg/response"
	travelSvc "art_admin_backend/internal/service/travel"
	"strconv"

	"github.com/gin-gonic/gin"
)

var roadbookService = travelSvc.NewRoadbookService()

// GetRoadbookList 获取路书列表
func GetRoadbookList(c *gin.Context) {
	var req travelSvc.GetRoadbookListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")
	req.UserID = userID.(int64)

	roadbooks, total, err := roadbookService.GetRoadbookList(&req)
	if err != nil {
		response.ServerError(c, "查询路书列表失败")
		return
	}

	response.SuccessWithPagination(c, roadbooks, req.Page, req.PageSize, total)
}

// GetRoadbookDetail 获取路书详情
func GetRoadbookDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")

	roadbook, err := roadbookService.GetRoadbookDetail(id, userID.(int64))
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权访问该路书")
			return
		}
		response.ServerError(c, "查询路书详情失败")
		return
	}

	response.Success(c, roadbook)
}

// CreateRoadbook 创建路书
func CreateRoadbook(c *gin.Context) {
	var req travelSvc.CreateRoadbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户信息
	userID, _ := c.Get("userId")
	req.UserID = userID.(int64)

	roadbook, err := roadbookService.CreateRoadbook(&req)
	if err != nil {
		response.ServerError(c, "创建路书失败: "+err.Error())
		return
	}

	response.Success(c, roadbook)
}

// UpdateRoadbook 更新路书
func UpdateRoadbook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	var req travelSvc.UpdateRoadbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	req.ID = id

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")
	req.UserID = userID.(int64)

	err = roadbookService.UpdateRoadbook(&req)
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权修改该路书")
			return
		}
		response.ServerError(c, "更新路书失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteRoadbook 删除路书
func DeleteRoadbook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的路书ID")
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")

	err = roadbookService.DeleteRoadbook(id, userID.(int64))
	if err != nil {
		if err.Error() == "forbidden" {
			response.Forbidden(c, "无权删除该路书")
			return
		}
		response.ServerError(c, "删除路书失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
