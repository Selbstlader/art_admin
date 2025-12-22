package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/service/design_version"
	"strconv"

	"github.com/gin-gonic/gin"
)

var designVersionService *design_version.DesignVersionService
var designVersionAIClient *volcengine.Client

// SetDesignVersionAIClient 设置设计版本服务的AI客户端
// Set AI client for design version service
func SetDesignVersionAIClient(client *volcengine.Client) {
	designVersionAIClient = client
}

// initDesignVersionService 初始化设计版本服务
// Initialize design version service
func initDesignVersionService() {
	if designVersionService == nil {
		db := database.GetDB()
		designVersionService = design_version.NewDesignVersionService(db)
		// 注入AI客户端 / Inject AI client
		if designVersionAIClient != nil {
			designVersionService.SetAIClient(designVersionAIClient)
		}
	}
}

// CreateDesignVersion 创建设计版本
// @Summary 创建设计版本
// @Description 创建新的设计版本
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param data body request.CreateDesignVersionRequest true "创建请求"
// @Success 200 {object} response.DesignVersionResponse
// @Router /api/designer/versions [post]
func CreateDesignVersion(c *gin.Context) {
	initDesignVersionService()

	var req request.CreateDesignVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getVersionUserID(c)

	resp, err := designVersionService.Create(&req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetDesignVersionList 获取设计版本列表
// @Summary 获取设计版本列表
// @Description 分页获取项目的设计版本列表
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param projectId query int true "项目ID"
// @Param current query int true "当前页"
// @Param size query int true "每页数量"
// @Param status query string false "状态筛选"
// @Success 200 {object} response.DesignVersionListResponse
// @Router /api/designer/versions [get]
func GetDesignVersionList(c *gin.Context) {
	initDesignVersionService()

	var req request.DesignVersionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	resp, err := designVersionService.List(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetDesignVersionDetail 获取设计版本详情
// @Summary 获取设计版本详情
// @Description 根据ID获取设计版本详情
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param id path int true "版本ID"
// @Success 200 {object} response.DesignVersionResponse
// @Router /api/designer/versions/{id} [get]
func GetDesignVersionDetail(c *gin.Context) {
	initDesignVersionService()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的ID")
		return
	}

	resp, err := designVersionService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// UpdateDesignVersion 更新设计版本
// @Summary 更新设计版本
// @Description 更新设计版本信息
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param id path int true "版本ID"
// @Param data body request.UpdateDesignVersionRequest true "更新请求"
// @Success 200 {object} response.DesignVersionResponse
// @Router /api/designer/versions/{id} [put]
func UpdateDesignVersion(c *gin.Context) {
	initDesignVersionService()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的ID")
		return
	}

	var req request.UpdateDesignVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	req.ID = uint(id)

	resp, err := designVersionService.Update(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// DeleteDesignVersion 删除设计版本
// @Summary 删除设计版本
// @Description 根据ID删除设计版本
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param id path int true "版本ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/designer/versions/{id} [delete]
func DeleteDesignVersion(c *gin.Context) {
	initDesignVersionService()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的ID")
		return
	}

	if err := designVersionService.Delete(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, "删除成功")
}

// CompareDesignVersions 对比设计版本
// @Summary 对比设计版本
// @Description 对比两个设计版本的差异
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param data body request.CompareDesignVersionsRequest true "对比请求"
// @Success 200 {object} response.VersionCompareResponse
// @Router /api/designer/versions/compare [post]
func CompareDesignVersions(c *gin.Context) {
	initDesignVersionService()

	var req request.CompareDesignVersionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	resp, err := designVersionService.Compare(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetVersionCompareList 获取版本对比列表
// @Summary 获取版本对比列表
// @Description 分页获取版本对比记录列表
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param projectId query int true "项目ID"
// @Param current query int true "当前页"
// @Param size query int true "每页数量"
// @Param compareStatus query string false "对比状态筛选"
// @Success 200 {object} response.VersionCompareListResponse
// @Router /api/designer/versions/compares [get]
func GetVersionCompareList(c *gin.Context) {
	initDesignVersionService()

	var req request.VersionCompareListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	resp, err := designVersionService.GetCompareList(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetVersionCompareDetail 获取版本对比详情
// @Summary 获取版本对比详情
// @Description 根据ID获取版本对比详情
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param id path int true "对比记录ID"
// @Success 200 {object} response.VersionCompareResponse
// @Router /api/designer/versions/compares/{id} [get]
func GetVersionCompareDetail(c *gin.Context) {
	initDesignVersionService()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的ID")
		return
	}

	resp, err := designVersionService.GetCompareDetail(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetVersionDiff 获取版本差异
// @Summary 获取版本差异
// @Description 获取两个版本的差异（用于左右分屏展示）
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param versionAId query int true "版本A ID"
// @Param versionBId query int true "版本B ID"
// @Success 200 {object} response.VersionDiffResponse
// @Router /api/designer/versions/diff [get]
func GetVersionDiff(c *gin.Context) {
	initDesignVersionService()

	versionAID, err := strconv.ParseUint(c.Query("versionAId"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的版本A ID")
		return
	}
	versionBID, err := strconv.ParseUint(c.Query("versionBId"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的版本B ID")
		return
	}

	resp, err := designVersionService.GetDiff(uint(versionAID), uint(versionBID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// DeleteVersionCompare 删除版本对比记录
// @Summary 删除版本对比记录
// @Description 根据ID删除版本对比记录
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param id path int true "对比记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/designer/versions/compares/{id} [delete]
func DeleteVersionCompare(c *gin.Context) {
	initDesignVersionService()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的ID")
		return
	}

	if err := designVersionService.DeleteCompare(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, "删除成功")
}

// StartAIVersionAnalysis 启动AI版本对比分析
// @Summary 启动AI版本对比分析
// @Description 异步启动AI版本对比分析，完成后通过通知告知用户
// @Tags 设计版本
// @Accept json
// @Produce json
// @Param data body request.AIAnalyzeVersionDiffRequest true "AI分析请求"
// @Success 200 {object} response.AIAnalyzeVersionDiffResponse
// @Router /api/designer/versions/ai-analyze [post]
func StartAIVersionAnalysis(c *gin.Context) {
	initDesignVersionService()

	var req request.AIAnalyzeVersionDiffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getVersionUserID(c)

	resp, err := designVersionService.StartAIAnalysis(&req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// getVersionUserID 从上下文获取用户ID
// Get user ID from context
func getVersionUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}
