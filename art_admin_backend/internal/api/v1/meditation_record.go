package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateMeditationRecord 创建冥想记录
func CreateMeditationRecord(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req request.CreateMeditationRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	// 调用服务创建记录
	err := meditationRecordService.CreateMeditationRecord(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "冥想记录创建成功",
		Data:    nil,
	})
}

// UpdateMeditationRecord 更新冥想记录
func UpdateMeditationRecord(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req request.UpdateMeditationRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	// 调用服务更新记录
	err := meditationRecordService.UpdateMeditationRecord(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "冥想记录更新成功",
		Data:    nil,
	})
}

// DeleteMeditationRecord 删除冥想记录
func DeleteMeditationRecord(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// 获取记录ID
	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	// 调用服务删除记录
	err = meditationRecordService.DeleteMeditationRecord(userID, recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "冥想记录删除成功",
		Data:    nil,
	})
}

// GetMeditationRecordList 获取冥想记录列表
func GetMeditationRecordList(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 获取筛选参数
	meditationType := c.Query("meditation_type")
	difficultyLevelStr := c.Query("difficulty_level")

	var difficultyLevel *int
	if difficultyLevelStr != "" {
		level, err := strconv.Atoi(difficultyLevelStr)
		if err == nil && level >= 1 && level <= 5 {
			difficultyLevel = &level
		}
	}

	// 构建请求结构体
	req := &request.MeditationRecordListRequest{
		Page:            page,
		PageSize:        pageSize,
		MeditationType:  meditationType,
		DifficultyLevel: nil, // 设置默认值
	}
	if difficultyLevel != nil {
		req.DifficultyLevel = difficultyLevel
	}

	// 调用服务获取列表
	records, total, err := meditationRecordService.GetMeditationRecordList(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建响应数据
	result := map[string]interface{}{
		"records":     records,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取冥想记录列表成功",
		Data:    result,
	})
}

// GetMeditationRecordDetail 获取冥想记录详情
func GetMeditationRecordDetail(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// 获取记录ID
	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	// 调用服务获取详情
	record, err := meditationRecordService.GetMeditationRecord(userID, recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取冥想记录详情成功",
		Data:    record,
	})
}

// GetMeditationStatistics 获取冥想统计信息
func GetMeditationStatistics(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// 获取天数参数
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		days = 30
	}

	// 调用服务获取统计信息
	stats, err := meditationRecordService.GetMeditationStatistics(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取冥想统计信息成功",
		Data:    stats,
	})
}
