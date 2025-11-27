package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateMoodRecord 创建情绪记录
// @Summary 创建情绪记录
// @Description 创建新的情绪记录
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateMoodRecordRequest true "情绪记录信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/create [post]
func CreateMoodRecord(c *gin.Context) {
	var req request.CreateMoodRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := moodRecordService.CreateMoodRecord(userID, &req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	// 触发成就检查（异步执行，不阻塞用户操作）
	if achievementService != nil {
		go achievementService.ProcessUserEvent(userID, "mood_record_created", map[string]interface{}{
			"mood_score": req.Intensity,
			"mood_type":  req.MoodType,
		})
	}

	response.SuccessWithMsg(c, "创建情绪记录成功", nil)
}

// UpdateMoodRecord 更新情绪记录
// @Summary 更新情绪记录
// @Description 更新情绪记录信息
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateMoodRecordRequest true "情绪记录信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/update [put]
func UpdateMoodRecord(c *gin.Context) {
	var req request.UpdateMoodRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := moodRecordService.UpdateMoodRecord(userID, &req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新情绪记录成功", nil)
}

// DeleteMoodRecord 删除情绪记录
// @Summary 删除情绪记录
// @Description 删除情绪记录
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "记录ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/delete/{id} [delete]
func DeleteMoodRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := moodRecordService.DeleteMoodRecord(userID, id); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除情绪记录成功", nil)
}

// GetMoodRecordList 获取情绪记录列表
// @Summary 获取情绪记录列表
// @Description 分页查询情绪记录列表
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param current query int true "当前页码" minimum(1)
// @Param size query int true "每页条数" minimum(1) maximum(100)
// @Param user_id query int false "用户ID"
// @Param mood_type query string false "情绪类型"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param min_intensity query int false "最小强度" minimum(1) maximum(10)
// @Param max_intensity query int false "最大强度" minimum(1) maximum(10)
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/list [get]
func GetMoodRecordList(c *gin.Context) {
	var req request.MoodRecordListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	records, total, err := moodRecordService.GetMoodRecordList(userID, &req)
	if err != nil {
		response.ServerError(c, "查询情绪记录列表失败")
		return
	}

	response.SuccessWithPagination(c, records, req.Page, req.PageSize, total)
}

// GetMoodStatistics 获取情绪统计
// @Summary 获取情绪统计
// @Description 获取用户情绪统计信息
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query int true "用户ID"
// @Param days query int false "统计天数" default(30)
// @Success 200 {object} response.Response{data=response.MoodStatistics} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/statistics [get]
func GetMoodStatistics(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	statistics, err := moodRecordService.GetMoodStatistics(userID, days)
	if err != nil {
		response.ServerError(c, "获取情绪统计失败")
		return
	}

	response.Success(c, statistics)
}

// GetMoodRecordDetail 获取情绪记录详情
// @Summary 获取情绪记录详情
// @Description 获取单条情绪记录详情
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "记录ID"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/{id} [get]
func GetMoodRecordDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	record, err := moodRecordService.GetMoodRecord(userID, id)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, record)
}

// GetMoodRecordsByDate 按日期获取情绪记录
// @Summary 按日期获取情绪记录
// @Description 获取指定日期的所有情绪记录
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date path string true "日期(格式:2006-01-02)"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/date/{date} [get]
func GetMoodRecordsByDate(c *gin.Context) {
	dateStr := c.Param("date")
	if dateStr == "" {
		response.BadRequest(c, "日期参数不能为空")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	records, err := moodRecordService.GetMoodRecordsByDate(userID, dateStr)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, records)
}

// GetMoodCalendar 获取情绪日历数据
// @Summary 获取情绪日历数据
// @Description 获取指定月份的情绪日历视图数据
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int true "年份"
// @Param month query int true "月份(1-12)"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/calendar [get]
func GetMoodCalendar(c *gin.Context) {
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

	calendarData, err := moodRecordService.GetMoodCalendar(userID, year, month)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, calendarData)
}

// GetMoodAnalytics 获取情绪分析
// @Summary 获取情绪分析
// @Description 获取用户情绪深度分析
// @Tags 情绪管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.MoodAnalyticsRequest true "分析请求"
// @Success 200 {object} response.Response{data=response.MoodAnalyticsItem} "分析成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/mood-record/analytics [post]
func GetMoodAnalytics(c *gin.Context) {
	var req request.MoodAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	analytics, err := moodRecordService.GetMoodAnalytics(userID, &req)
	if err != nil {
		response.ServerError(c, "获取情绪分析失败")
		return
	}

	response.Success(c, analytics)
}
