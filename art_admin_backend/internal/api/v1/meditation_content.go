package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var meditationContentService *service.MeditationContentService

// getMeditationContentService 获取冥想内容服务实例（懒加载）
func getMeditationContentService() *service.MeditationContentService {
	if meditationContentService == nil {
		meditationContentService = service.NewMeditationContentService()
	}
	return meditationContentService
}

// CreateMeditationContent 创建冥想内容
// @Summary 创建冥想内容
// @Description 创建新的冥想内容
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateMeditationContentRequest true "冥想内容信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/create [post]
func CreateMeditationContent(c *gin.Context) {
	var req request.CreateMeditationContentRequest
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

	// 获取服务实例（懒加载）
	svc := getMeditationContentService()

	if err := svc.CreateMeditationContent(userID, &req); err != nil {
		fmt.Printf("创建冥想内容失败: %v\n", err)
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	// 触发成就检查（异步执行，不阻塞用户操作）
	if achievementService != nil {
		go achievementService.ProcessUserEvent(userID, "meditation_completed", map[string]interface{}{
			"title":     req.Title,
			"duration":  req.Duration,
			"audio_url": req.AudioURL,
		})
	}

	response.SuccessWithMsg(c, "创建冥想内容成功", nil)
}

// UpdateMeditationContent 更新冥想内容
// @Summary 更新冥想内容
// @Description 更新冥想内容信息
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateMeditationContentRequest true "冥想内容信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/update [put]
func UpdateMeditationContent(c *gin.Context) {
	var req request.UpdateMeditationContentRequest
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

	if err := getMeditationContentService().UpdateMeditationContent(userID, &req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新冥想内容成功", nil)
}

// DeleteMeditationContent 删除冥想内容
// @Summary 删除冥想内容
// @Description 删除冥想内容
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "冥想内容ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/delete/{id} [delete]
func DeleteMeditationContent(c *gin.Context) {
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

	if err := getMeditationContentService().DeleteMeditationContent(userID, id); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除冥想内容成功", nil)
}

// GetMeditationContentList 获取冥想内容列表
// @Summary 获取冥想内容列表
// @Description 分页查询冥想内容列表
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int true "当前页码" minimum(1)
// @Param pageSize query int true "每页条数" minimum(1) maximum(100)
// @Param category query string false "冥想分类"
// @Param difficultyLevel query int false "难度等级"
// @Param isPublic query bool false "是否公开"
// @Param minDuration query int false "最小时长(秒)"
// @Param maxDuration query int false "最大时长(秒)"
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/list [get]
func GetMeditationContentList(c *gin.Context) {
	var req request.MeditationContentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	contents, total, err := getMeditationContentService().GetMeditationContentList(&req)
	if err != nil {
		response.ServerError(c, "查询冥想内容列表失败")
		return
	}

	response.SuccessWithPagination(c, contents, req.Page, req.PageSize, total)
}

// GetMeditationContentDetail 获取冥想内容详情
// @Summary 获取冥想内容详情
// @Description 获取冥想内容详细信息
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "冥想内容ID"
// @Success 200 {object} response.Response{data=response.MeditationContentDetail} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/{id} [get]
func GetMeditationContentDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	detail, err := getMeditationContentService().GetMeditationContentDetail(id)
	if err != nil {
		response.ServerError(c, "获取冥想内容详情失败")
		return
	}

	response.Success(c, detail)
}

// GetMeditationContentByCategory 根据分类获取冥想内容
// @Summary 根据分类获取冥想内容
// @Description 根据分类获取冥想内容列表
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category path string true "冥想分类"
// @Param limit query int false "返回数量限制" default(10)
// @Success 200 {object} response.Response{data=[]response.MeditationContentItem} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/category/{category} [get]
func GetMeditationContentByCategory(c *gin.Context) {
	category := c.Param("category")
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	contents, err := getMeditationContentService().GetMeditationContentByCategory(category, limit)
	if err != nil {
		response.ServerError(c, "获取分类冥想内容失败")
		return
	}

	response.Success(c, contents)
}

// GetMeditationContentByDifficulty 根据难度等级获取冥想内容
// @Summary 根据难度等级获取冥想内容
// @Description 根据难度等级获取冥想内容列表
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param difficulty path int true "难度等级"
// @Param limit query int false "返回数量限制" default(10)
// @Success 200 {object} response.Response{data=[]response.MeditationContentItem} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/difficulty/{difficulty} [get]
func GetMeditationContentByDifficulty(c *gin.Context) {
	difficultyStr := c.Param("difficulty")
	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil || difficulty < 1 || difficulty > 5 {
		response.BadRequest(c, "难度等级必须在1-5之间")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	contents, err := getMeditationContentService().GetMeditationContentByDifficulty(difficulty, limit)
	if err != nil {
		response.ServerError(c, "获取难度冥想内容失败")
		return
	}

	response.Success(c, contents)
}

// GetPopularMeditationContent 获取热门冥想内容
// @Summary 获取热门冥想内容
// @Description 获取热门冥想内容列表
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "返回数量限制" default(10)
// @Success 200 {object} response.Response{data=[]response.MeditationContentItem} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/popular [get]
func GetPopularMeditationContent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	contents, err := getMeditationContentService().GetPopularMeditationContent(limit)
	if err != nil {
		response.ServerError(c, "获取热门冥想内容失败")
		return
	}

	response.Success(c, contents)
}

// LikeMeditationContent 点赞冥想内容
// @Summary 点赞冥想内容
// @Description 点赞冥想内容
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "冥想内容ID"
// @Success 200 {object} response.Response "点赞成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/like/{id} [post]
func LikeMeditationContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := getMeditationContentService().LikeMeditationContent(id); err != nil {
		response.ServerError(c, "点赞失败")
		return
	}

	response.SuccessWithMsg(c, "点赞成功", nil)
}

// UnlikeMeditationContent 取消点赞冥想内容
// @Summary 取消点赞冥想内容
// @Description 取消点赞冥想内容
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "冥想内容ID"
// @Success 200 {object} response.Response "取消点赞成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/unlike/{id} [post]
func UnlikeMeditationContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := getMeditationContentService().UnlikeMeditationContent(id); err != nil {
		response.ServerError(c, "取消点赞失败")
		return
	}

	response.SuccessWithMsg(c, "取消点赞成功", nil)
}

// SearchMeditationContent 搜索冥想内容
// @Summary 搜索冥想内容
// @Description 搜索冥想内容
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyword query string true "搜索关键词"
// @Param page query int true "当前页码" minimum(1)
// @Param pageSize query int true "每页条数" minimum(1) maximum(100)
// @Success 200 {object} response.Response{data=response.PaginatedData} "搜索成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/meditation-content/search [get]
func SearchMeditationContent(c *gin.Context) {
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

	contents, total, err := getMeditationContentService().SearchMeditationContent(keyword, page, pageSize)
	if err != nil {
		response.ServerError(c, "搜索冥想内容失败")
		return
	}

	response.SuccessWithPagination(c, contents, page, pageSize, total)
}

// GetMeditationCategories 获取冥想分类列表
// @Summary 获取冥想分类列表
// @Description 获取所有冥想分类
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]string} "查询成功"
// @Router /api/meditation-content/categories [get]
func GetMeditationCategories(c *gin.Context) {
	categories, err := getMeditationContentService().GetCategories()
	if err != nil {
		response.ServerError(c, "获取冥想分类失败")
		return
	}

	response.Success(c, categories)
}

// GetMeditationDifficultyLevels 获取难度等级列表
// @Summary 获取难度等级列表
// @Description 获取所有难度等级
// @Tags 冥想管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]int} "查询成功"
// @Router /api/meditation-content/difficulty-levels [get]
func GetMeditationDifficultyLevels(c *gin.Context) {
	levels, err := getMeditationContentService().GetDifficultyLevels()
	if err != nil {
		response.ServerError(c, "获取难度等级失败")
		return
	}

	response.Success(c, levels)
}
