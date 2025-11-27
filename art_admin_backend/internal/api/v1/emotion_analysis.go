package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	moodAnalyticsSvc "art_admin_backend/internal/service/mood_analytics"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AnalyzeMoodPatterns 分析情绪模式
func AnalyzeMoodPatterns(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取请求参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数，范围1-365"})
		return
	}

	// 调用服务分析
	result, err := moodRecordService.AnalyzeMoodPatterns(c.Request.Context(), userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "情绪模式分析成功",
		Data:    result,
	})
}

// GetMoodTrends 获取情绪趋势
func GetMoodTrends(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取请求参数
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数，范围1-365"})
		return
	}

	// 调用服务获取趋势
	trends, err := GetMoodRecordService().GetMoodTrends(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取情绪趋势成功",
		Data:    trends,
	})
}

// AnalyzeMeditationSession 分析冥想会话
func AnalyzeMeditationSession(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取冥想记录ID
	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	// 调用服务分析
	result, err := GetMeditationRecordService().AnalyzeMeditationSession(c.Request.Context(), userID, recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "冥想会话分析成功",
		Data:    result,
	})
}

// GetMeditationTrends 获取冥想趋势
func GetMeditationTrends(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取请求参数
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数，范围1-365"})
		return
	}

	// 调用服务获取趋势
	trends, err := GetMeditationRecordService().GetMeditationTrends(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取冥想趋势成功",
		Data:    trends,
	})
}

// GetMeditationRecommendations 获取冥想推荐
func GetMeditationRecommendations(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 调用服务获取推荐
	recommendations, err := GetMeditationRecordService().GetMeditationRecommendations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取冥想推荐成功",
		Data:    recommendations,
	})
}

// AnalyzeJournalEntry 分析日记条目
func AnalyzeJournalEntry(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取日记条目ID
	entryIDStr := c.Param("id")
	entryID, err := strconv.ParseInt(entryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的条目ID"})
		return
	}

	// 调用服务分析
	result, err := GetJournalEntryService().AnalyzeJournalEntry(c.Request.Context(), userID, entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "日记条目分析成功",
		Data:    result,
	})
}

// GetJournalInsights 获取日记洞察
func GetJournalInsights(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取请求参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数，范围1-365"})
		return
	}

	// 调用服务获取洞察
	insights, err := journalEntryService.GetJournalInsights(c.Request.Context(), userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取日记洞察成功",
		Data:    insights,
	})
}

// GetJournalTrends 获取日记趋势
func GetJournalTrends(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取请求参数
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数，范围1-365"})
		return
	}

	// 调用服务获取趋势
	trends, err := journalEntryService.GetJournalTrends(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取日记趋势成功",
		Data:    trends,
	})
}

// BatchAnalyze 批量分析
func BatchAnalyze(c *gin.Context) {
	var req response.BatchAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	// 验证请求参数
	if len(req.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID列表不能为空"})
		return
	}
	if len(req.UserIDs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "批量分析最多支持100个用户"})
		return
	}

	// 调用AI分析服务进行批量分析
	aiService := moodAnalyticsSvc.NewAIAnalysisService()
	result := &response.BatchAnalysisResponse{
		TotalUsers:  len(req.UserIDs),
		Results:     make(map[int64]interface{}),
		Errors:      make(map[int64]string),
		ProcessedAt: time.Now(),
	}

	processedCount := 0
	failedCount := 0

	for _, userID := range req.UserIDs {
		// 创建分析请求
		dtoReq := &request.EmotionAnalysisRequest{
			UserID:       userID,
			StartDate:    req.StartDate,
			EndDate:      req.EndDate,
			AnalysisType: req.AnalysisType,
		}

		// 转换为service层请求类型
		serviceReq := &moodAnalyticsSvc.EmotionAnalysisRequest{
			UserID:       dtoReq.UserID,
			StartDate:    dtoReq.StartDate,
			EndDate:      dtoReq.EndDate,
			AnalysisType: dtoReq.AnalysisType,
		}

		// 执行分析
		analysisResult, err := aiService.AnalyzeUserEmotions(c.Request.Context(), serviceReq)
		if err != nil {
			result.Errors[userID] = err.Error()
			failedCount++
		} else {
			result.Results[userID] = analysisResult
			processedCount++
		}
	}

	result.ProcessedUsers = processedCount
	result.FailedUsers = failedCount

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "批量分析完成",
		Data:    result,
	})
}

// GetAnalysisStatus 获取分析状态
func GetAnalysisStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取记录类型和ID
	recordType := c.Query("type")
	recordIDStr := c.Query("id")

	if recordType == "" || recordIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数: type 和 id"})
		return
	}

	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	var status response.AnalysisStatus

	switch recordType {
	case "mood":
		// 获取情绪记录分析状态
		record, err := moodRecordService.GetMoodRecord(userID, recordID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "情绪记录不存在"})
			return
		}
		status = response.AnalysisStatus{
			Status:   record.AnalysisStatus,
			Progress: getAnalysisProgress(record.AnalysisStatus),
		}
	case "meditation":
		// 获取冥想记录分析状态
		record, err := meditationRecordService.GetMeditationRecord(userID, recordID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "冥想记录不存在"})
			return
		}
		status = response.AnalysisStatus{
			Status:   record.AnalysisStatus,
			Progress: getAnalysisProgress(record.AnalysisStatus),
		}
	case "journal":
		// 获取日记条目分析状态
		record, err := journalEntryService.GetJournalEntry(userID, recordID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "日记条目不存在"})
			return
		}
		status = response.AnalysisStatus{
			Status:   record.AnalysisStatus,
			Progress: getAnalysisProgress(record.AnalysisStatus),
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录类型"})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Code:    200,
		Message: "获取分析状态成功",
		Data:    status,
	})
}

// getAnalysisProgress 根据分析状态获取进度
func getAnalysisProgress(status string) float64 {
	switch status {
	case "pending":
		return 0.0
	case "analyzing":
		return 50.0
	case "completed":
		return 100.0
	case "failed":
		return 0.0
	default:
		return 0.0
	}
}
