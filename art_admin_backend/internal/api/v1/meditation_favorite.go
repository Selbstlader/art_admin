package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddMeditationFavorite 添加冥想收藏
// @Summary 添加冥想收藏
// @Description 收藏冥想内容
// @Tags 冥想收藏
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "冥想内容ID"
// @Success 200 {object} response.Response "收藏成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/{id}/favorite [post]
func AddMeditationFavorite(c *gin.Context) {
	idStr := c.Param("id")
	contentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := meditationFavoriteService.AddFavorite(userID, contentID); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "收藏成功", nil)
}

// RemoveMeditationFavorite 取消冥想收藏
// @Summary 取消冥想收藏
// @Description 取消收藏冥想内容
// @Tags 冥想收藏
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "冥想内容ID"
// @Success 200 {object} response.Response "取消收藏成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/{id}/favorite [delete]
func RemoveMeditationFavorite(c *gin.Context) {
	idStr := c.Param("id")
	contentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := meditationFavoriteService.RemoveFavorite(userID, contentID); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "取消收藏成功", nil)
}

// CheckMeditationFavorite 检查是否已收藏
// @Summary 检查是否已收藏
// @Description 检查冥想内容是否已收藏
// @Tags 冥想收藏
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "冥想内容ID"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/{id}/favorite/check [get]
func CheckMeditationFavorite(c *gin.Context) {
	idStr := c.Param("id")
	contentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	isFavorite, err := meditationFavoriteService.IsFavorite(userID, contentID)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, map[string]bool{"is_favorite": isFavorite})
}

// GetUserFavorites 获取用户收藏列表
// @Summary 获取用户收藏列表
// @Description 获取用户收藏的冥想内容列表
// @Tags 冥想收藏
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "当前页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/favorites [get]
func GetUserFavorites(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	items, total, err := meditationFavoriteService.GetUserFavorites(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取收藏列表失败")
		return
	}

	response.SuccessWithPagination(c, items, page, pageSize, total)
}

// UpdatePlayRecord 更新播放记录
// @Summary 更新播放记录
// @Description 更新冥想内容的播放记录
// @Tags 冥想播放记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdatePlayRecordRequest true "播放记录"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/play-record [post]
func UpdatePlayRecord(c *gin.Context) {
	var req request.UpdatePlayRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := meditationFavoriteService.UpdatePlayRecord(userID, req.ContentID, req.Duration, req.Progress, req.Completed); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新播放记录成功", nil)
}

// GetPlayRecord 获取播放记录
// @Summary 获取播放记录
// @Description 获取指定冥想内容的播放记录
// @Tags 冥想播放记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "冥想内容ID"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/{id}/play-record [get]
func GetPlayRecord(c *gin.Context) {
	idStr := c.Param("id")
	contentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	record, err := meditationFavoriteService.GetPlayRecord(userID, contentID)
	if err != nil {
		response.Success(c, nil) // 没有记录返回空
		return
	}

	response.Success(c, record)
}

// GetUserPlayHistory 获取用户播放历史
// @Summary 获取用户播放历史
// @Description 获取用户的冥想播放历史
// @Tags 冥想播放记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "当前页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/play-history [get]
func GetUserPlayHistory(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	items, total, err := meditationFavoriteService.GetUserPlayHistory(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取播放历史失败")
		return
	}

	response.SuccessWithPagination(c, items, page, pageSize, total)
}

// GetRecentlyPlayed 获取最近播放
// @Summary 获取最近播放
// @Description 获取用户最近播放的冥想内容
// @Tags 冥想播放记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/recently-played [get]
func GetRecentlyPlayed(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	items, err := meditationFavoriteService.GetRecentlyPlayed(userID, limit)
	if err != nil {
		response.ServerError(c, "获取最近播放失败")
		return
	}

	response.Success(c, items)
}

// GetContinuePlaying 获取继续播放列表
// @Summary 获取继续播放列表
// @Description 获取用户未完成的冥想内容
// @Tags 冥想播放记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "返回数量" default(5)
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/continue-playing [get]
func GetContinuePlaying(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 20 {
		limit = 5
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	items, err := meditationFavoriteService.GetContinuePlaying(userID, limit)
	if err != nil {
		response.ServerError(c, "获取继续播放列表失败")
		return
	}

	response.Success(c, items)
}

// GetPlayStats 获取播放统计
// @Summary 获取播放统计
// @Description 获取用户的冥想播放统计
// @Tags 冥想播放记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/play-stats [get]
func GetPlayStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	stats, err := meditationFavoriteService.GetPlayStats(userID)
	if err != nil {
		response.ServerError(c, "获取播放统计失败")
		return
	}

	response.Success(c, stats)
}
