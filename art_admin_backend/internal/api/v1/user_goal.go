package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 目标服务实例
var userGoalService = service.NewUserGoalService()

// CreateGoal 创建目标
// @Summary 创建目标
// @Description 创建用户目标
// @Tags 目标管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateGoalRequest true "目标信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/goals [post]
func CreateGoal(c *gin.Context) {
	var req request.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := userGoalService.CreateGoal(userID, req.GoalType, req.TargetValue, req.Period); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建目标成功", nil)
}

// UpdateGoal 更新目标
// @Summary 更新目标
// @Description 更新用户目标
// @Tags 目标管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "目标ID"
// @Param request body request.UpdateGoalRequest true "目标信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/goals/{id} [put]
func UpdateGoal(c *gin.Context) {
	idStr := c.Param("id")
	goalID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := userGoalService.UpdateGoal(userID, goalID, req.TargetValue, req.IsActive); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新目标成功", nil)
}

// DeleteGoal 删除目标
// @Summary 删除目标
// @Description 删除用户目标
// @Tags 目标管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "目标ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/goals/{id} [delete]
func DeleteGoal(c *gin.Context) {
	idStr := c.Param("id")
	goalID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := userGoalService.DeleteGoal(userID, goalID); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除目标成功", nil)
}

// GetUserGoals 获取用户目标列表
// @Summary 获取用户目标列表
// @Description 获取用户的所有目标及进度
// @Tags 目标管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/goals [get]
func GetUserGoals(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	goals, err := userGoalService.GetUserGoals(userID)
	if err != nil {
		response.ServerError(c, "获取目标列表失败")
		return
	}

	response.Success(c, goals)
}

// Checkin 打卡
// @Summary 打卡
// @Description 用户打卡
// @Tags 打卡管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CheckinRequest true "打卡信息"
// @Success 200 {object} response.Response "打卡成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/checkin [post]
func Checkin(c *gin.Context) {
	var req request.CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := userGoalService.Checkin(userID, req.CheckinType, req.Note); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "打卡成功", nil)
}

// GetTodayCheckin 获取今日打卡状态
// @Summary 获取今日打卡状态
// @Description 获取用户今日打卡状态
// @Tags 打卡管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/checkin/today [get]
func GetTodayCheckin(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	status, err := userGoalService.GetTodayCheckin(userID)
	if err != nil {
		response.ServerError(c, "获取打卡状态失败")
		return
	}

	response.Success(c, status)
}

// GetCheckinHistory 获取打卡历史
// @Summary 获取打卡历史
// @Description 获取用户打卡历史
// @Tags 打卡管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "查询天数" default(30)
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/checkin/history [get]
func GetCheckinHistory(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	history, err := userGoalService.GetCheckinHistory(userID, days)
	if err != nil {
		response.ServerError(c, "获取打卡历史失败")
		return
	}

	response.Success(c, history)
}

// GetCheckinStreak 获取连续打卡天数
// @Summary 获取连续打卡天数
// @Description 获取用户连续打卡天数
// @Tags 打卡管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "打卡类型"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/checkin/streak [get]
func GetCheckinStreak(c *gin.Context) {
	checkinType := c.Query("type")

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	streak, err := userGoalService.GetCheckinStreak(userID, checkinType)
	if err != nil {
		response.ServerError(c, "获取连续打卡天数失败")
		return
	}

	response.Success(c, map[string]int{"streak": streak})
}

// GetCheckinCalendar 获取打卡日历
// @Summary 获取打卡日历
// @Description 获取用户打卡日历数据
// @Tags 打卡管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int true "年份"
// @Param month query int true "月份"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/checkin/calendar [get]
func GetCheckinCalendar(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.BadRequest(c, "无效的年份参数")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		response.BadRequest(c, "无效的月份参数")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	calendar, err := userGoalService.GetCheckinCalendar(userID, year, month)
	if err != nil {
		response.ServerError(c, "获取打卡日历失败")
		return
	}

	response.Success(c, calendar)
}

// GetCheckinStats 获取打卡统计
// @Summary 获取打卡统计
// @Description 获取用户打卡统计信息
// @Tags 打卡管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/checkin/stats [get]
func GetCheckinStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	stats, err := userGoalService.GetCheckinStats(userID)
	if err != nil {
		response.ServerError(c, "获取打卡统计失败")
		return
	}

	response.Success(c, stats)
}

// UpdateGoalProgress 更新目标进度
// @Summary 更新目标进度
// @Description 更新用户目标进度
// @Tags 目标管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "目标ID"
// @Param request body request.UpdateGoalProgressRequest true "进度信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/goals/{id}/progress [put]
func UpdateGoalProgress(c *gin.Context) {
	idStr := c.Param("id")
	goalID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateGoalProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := userGoalService.UpdateGoalProgress(userID, goalID, req.Value); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新进度成功", nil)
}

// GetGoalProgress 获取目标进度历史
// @Summary 获取目标进度历史
// @Description 获取用户目标进度历史
// @Tags 目标管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "目标ID"
// @Param days query int false "查询天数" default(30)
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/goals/{id}/progress [get]
func GetGoalProgress(c *gin.Context) {
	idStr := c.Param("id")
	goalID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	progress, err := userGoalService.GetGoalProgress(userID, goalID, days)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, progress)
}
