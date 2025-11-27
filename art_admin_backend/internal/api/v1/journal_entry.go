package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateJournalEntry 创建日记条目
// @Summary 创建日记条目
// @Description 创建新的日记条目
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateJournalEntryRequest true "日记条目信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/create [post]
func CreateJournalEntry(c *gin.Context) {
	var req request.CreateJournalEntryRequest
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

	if err := journalEntryService.CreateJournalEntry(userID, &req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	// 触发成就检查（异步执行，不阻塞用户操作）
	if achievementService != nil {
		go achievementService.ProcessUserEvent(userID, "journal_created", map[string]interface{}{
			"title":   req.Title,
			"content": req.Content,
		})
	}

	response.SuccessWithMsg(c, "创建日记条目成功", nil)
}

// UpdateJournalEntry 更新日记条目
// @Summary 更新日记条目
// @Description 更新日记条目信息
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateJournalEntryRequest true "日记条目信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/update [put]
func UpdateJournalEntry(c *gin.Context) {
	var req request.UpdateJournalEntryRequest
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

	if err := journalEntryService.UpdateJournalEntry(userID, &req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新日记条目成功", nil)
}

// DeleteJournalEntry 删除日记条目
// @Summary 删除日记条目
// @Description 删除日记条目
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记条目ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/delete/{id} [delete]
func DeleteJournalEntry(c *gin.Context) {
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

	if err := journalEntryService.DeleteJournalEntry(userID, id); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除日记条目成功", nil)
}

// GetJournalEntryList 获取日记条目列表
// @Summary 获取日记条目列表
// @Description 分页查询日记条目列表
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int true "当前页码" minimum(1)
// @Param pageSize query int true "每页条数" minimum(1) maximum(100)
// @Param moodRecordId query int false "关联情绪记录ID"
// @Param isPublic query bool false "是否公开"
// @Param startDate query string false "开始日期(YYYY-MM-DD)"
// @Param endDate query string false "结束日期(YYYY-MM-DD)"
// @Param hasAttachment query bool false "是否有附件"
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/list [get]
func GetJournalEntryList(c *gin.Context) {
	var req request.JournalEntryListRequest
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

	entries, total, err := journalEntryService.GetJournalEntryList(userID, &req)
	if err != nil {
		response.ServerError(c, "查询日记条目列表失败")
		return
	}

	response.SuccessWithPagination(c, entries, req.Page, req.PageSize, total)
}

// GetJournalEntryDetail 获取日记条目详情
// @Summary 获取日记条目详情
// @Description 获取日记条目详细信息
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记条目ID"
// @Success 200 {object} response.Response{data=response.JournalEntryDetail} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/{id} [get]
func GetJournalEntryDetail(c *gin.Context) {
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

	detail, err := journalEntryService.GetJournalEntryDetail(userID, id)
	if err != nil {
		response.ServerError(c, "获取日记条目详情失败")
		return
	}

	response.Success(c, detail)
}

// GetJournalEntriesByDateRange 根据日期范围获取日记条目
// @Summary 根据日期范围获取日记条目
// @Description 根据日期范围获取日记条目列表
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param startDate query string true "开始日期(YYYY-MM-DD)"
// @Param endDate query string true "结束日期(YYYY-MM-DD)"
// @Param page query int false "当前页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/date-range [get]
func GetJournalEntriesByDateRange(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		response.BadRequest(c, "开始日期和结束日期不能为空")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.BadRequest(c, "开始日期格式错误")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.BadRequest(c, "结束日期格式错误")
		return
	}

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

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	entries, total, err := journalEntryService.GetJournalEntriesByDateRange(userID, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取日期范围日记失败")
		return
	}

	response.SuccessWithPagination(c, entries, page, pageSize, total)
}

// GetRecentJournalEntries 获取最近的日记条目
// @Summary 获取最近的日记条目
// @Description 获取最近的日记条目列表
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "返回数量限制" default(10)
// @Success 200 {object} response.Response{data=[]response.JournalEntryItem} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/recent [get]
func GetRecentJournalEntries(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	entries, err := journalEntryService.GetRecentJournalEntries(userID, limit)
	if err != nil {
		response.ServerError(c, "获取最近日记失败")
		return
	}

	response.Success(c, entries)
}

// SearchJournalEntries 搜索日记条目
// @Summary 搜索日记条目
// @Description 搜索日记条目
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyword query string true "搜索关键词"
// @Param page query int true "当前页码" minimum(1)
// @Param pageSize query int true "每页条数" minimum(1) maximum(100)
// @Success 200 {object} response.Response{data=response.PaginatedData} "搜索成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/search [get]
func SearchJournalEntries(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "搜索关键词不能为空")
		return
	}

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

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	entries, total, err := journalEntryService.SearchJournalEntries(userID, keyword, page, pageSize)
	if err != nil {
		response.ServerError(c, "搜索日记条目失败")
		return
	}

	response.SuccessWithPagination(c, entries, page, pageSize, total)
}

// GetJournalStatistics 获取日记统计信息
// @Summary 获取日记统计信息
// @Description 获取日记统计信息
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "统计天数" default(30)
// @Success 200 {object} response.Response{data=response.JournalStatistics} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/statistics [get]
func GetJournalStatistics(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	stats, err := journalEntryService.GetJournalStatistics(userID, days)
	if err != nil {
		response.ServerError(c, "获取日记统计失败")
		return
	}

	response.Success(c, stats)
}

// UpdateSentimentScore 更新情感得分
// @Summary 更新情感得分
// @Description 更新日记条目的情感得分
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记条目ID"
// @Param request body request.UpdateSentimentScoreRequest true "情感得分"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/sentiment/{id} [put]
func UpdateSentimentScore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateSentimentScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证情感得分范围
	if req.SentimentScore < -1 || req.SentimentScore > 1 {
		response.BadRequest(c, "情感得分必须在-1到1之间")
		return
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := journalEntryService.UpdateSentimentScore(userID, id, req.SentimentScore); err != nil {
		response.ServerError(c, "更新情感得分失败")
		return
	}

	response.SuccessWithMsg(c, "更新情感得分成功", nil)
}

// UpdateJournalTags 更新日记标签
// @Summary 更新日记标签
// @Description 更新日记条目的标签
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记条目ID"
// @Param request body request.UpdateJournalTagsRequest true "标签列表"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/{id}/tags [put]
func UpdateJournalTags(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateJournalTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := journalEntryService.UpdateJournalTags(userID, id, req.Tags); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新标签成功", nil)
}

// AddJournalImages 添加日记图片
// @Summary 添加日记图片
// @Description 为日记条目添加图片
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记条目ID"
// @Param request body request.AddJournalImagesRequest true "图片URL列表"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/{id}/images [post]
func AddJournalImages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.AddJournalImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := journalEntryService.AddJournalImages(userID, id, req.Images); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "添加图片成功", nil)
}

// RemoveJournalImage 删除日记图片
// @Summary 删除日记图片
// @Description 从日记条目中删除图片
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记条目ID"
// @Param request body request.RemoveJournalImageRequest true "图片URL"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/{id}/images [delete]
func RemoveJournalImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.RemoveJournalImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	if err := journalEntryService.RemoveJournalImage(userID, id, req.ImageURL); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除图片成功", nil)
}

// GetAllUserTags 获取用户所有标签
// @Summary 获取用户所有标签
// @Description 获取用户在日记中使用过的所有标签
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/tags [get]
func GetAllUserTags(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	tags, err := journalEntryService.GetAllUserTags(userID)
	if err != nil {
		response.ServerError(c, "获取标签失败")
		return
	}

	response.Success(c, tags)
}

// GetJournalsByTag 按标签获取日记
// @Summary 按标签获取日记
// @Description 获取包含指定标签的日记列表
// @Tags 日记管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tag query string true "标签名称"
// @Param page query int false "当前页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/journal-entry/by-tag [get]
func GetJournalsByTag(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		response.BadRequest(c, "标签参数不能为空")
		return
	}

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

	entries, total, err := journalEntryService.GetJournalsByTag(userID, tag, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询日记失败")
		return
	}

	response.SuccessWithPagination(c, entries, page, pageSize, total)
}
