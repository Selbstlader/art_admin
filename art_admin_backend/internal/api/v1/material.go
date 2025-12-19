package v1

import (
	"strconv"
	"sync"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service/material"

	"github.com/gin-gonic/gin"
)

// 使用延迟初始化避免数据库未初始化时的空指针问题
// Use lazy initialization to avoid nil pointer when database is not initialized
var (
	materialService     material.MaterialService
	materialServiceOnce sync.Once
)

// getMaterialService 获取材料服务实例（延迟初始化）
// Get material service instance (lazy initialization)
func getMaterialService() material.MaterialService {
	materialServiceOnce.Do(func() {
		materialService = material.NewMaterialService()
	})
	return materialService
}

// CreateMaterial 创建材料
// @Summary 创建材料
// @Description Create a new material
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param request body request.CreateMaterialRequest true "创建材料请求"
// @Success 200 {object} response.Response{data=response.MaterialResponse}
// @Router /api/designer/materials [post]
func CreateMaterial(c *gin.Context) {
	var req request.CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := getMaterialService().CreateMaterial(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", result)
}

// UpdateMaterial 更新材料
// @Summary 更新材料
// @Description Update an existing material
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param request body request.UpdateMaterialRequest true "更新材料请求"
// @Success 200 {object} response.Response{data=response.MaterialResponse}
// @Router /api/designer/materials [put]
func UpdateMaterial(c *gin.Context) {
	var req request.UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := getMaterialService().UpdateMaterial(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", result)
}

// DeleteMaterial 删除材料
// @Summary 删除材料
// @Description Delete a material by ID
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param id path int true "材料ID"
// @Success 200 {object} response.Response
// @Router /api/designer/materials/{id} [delete]
func DeleteMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的材料ID")
		return
	}

	if err := getMaterialService().DeleteMaterial(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// BatchDeleteMaterials 批量删除材料
// @Summary 批量删除材料
// @Description Batch delete materials by IDs
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param request body request.BatchDeleteMaterialRequest true "批量删除请求"
// @Success 200 {object} response.Response
// @Router /api/designer/materials/batch-delete [post]
func BatchDeleteMaterials(c *gin.Context) {
	var req request.BatchDeleteMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := getMaterialService().BatchDeleteMaterials(req.IDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "批量删除成功", nil)
}

// GetMaterialDetail 获取材料详情
// @Summary 获取材料详情
// @Description Get material detail by ID
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param id path int true "材料ID"
// @Success 200 {object} response.Response{data=response.MaterialDetailResponse}
// @Router /api/designer/materials/{id} [get]
func GetMaterialDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的材料ID")
		return
	}

	result, err := getMaterialService().GetMaterialByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetMaterialList 获取材料列表
// @Summary 获取材料列表
// @Description Get material list with pagination and filters
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param current query int true "当前页码"
// @Param size query int true "每页条数"
// @Param name query string false "材料名称"
// @Param category query string false "分类"
// @Param brand query string false "品牌"
// @Param minPrice query number false "最低单价"
// @Param maxPrice query number false "最高单价"
// @Param status query string false "状态"
// @Param keyword query string false "关键字"
// @Success 200 {object} response.Response{data=response.MaterialListResponse}
// @Router /api/designer/materials [get]
func GetMaterialList(c *gin.Context) {
	var req request.MaterialListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := getMaterialService().GetMaterialList(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetMaterialCategories 获取材料分类列表
// @Summary 获取材料分类列表
// @Description Get all material categories with count
// @Tags 材料管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.MaterialCategoryResponse}
// @Router /api/designer/materials/categories [get]
func GetMaterialCategories(c *gin.Context) {
	result, err := getMaterialService().GetCategories()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetMaterialBrands 获取材料品牌列表
// @Summary 获取材料品牌列表
// @Description Get all material brands with count
// @Tags 材料管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.MaterialBrandResponse}
// @Router /api/designer/materials/brands [get]
func GetMaterialBrands(c *gin.Context) {
	result, err := getMaterialService().GetBrands()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetMaterialStats 获取材料统计信息
// @Summary 获取材料统计信息
// @Description Get material statistics
// @Tags 材料管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.MaterialStatsResponse}
// @Router /api/designer/materials/stats [get]
func GetMaterialStats(c *gin.Context) {
	result, err := getMaterialService().GetStats()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// BatchImportMaterials 批量导入材料
// @Summary 批量导入材料
// @Description Batch import materials
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param request body request.BatchImportMaterialRequest true "批量导入请求"
// @Success 200 {object} response.Response{data=response.BatchImportMaterialResponse}
// @Router /api/designer/materials/batch-import [post]
func BatchImportMaterials(c *gin.Context) {
	var req request.BatchImportMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := getMaterialService().BatchImportMaterials(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// CalculateMaterialUsage 计算材料用量
// @Summary 计算材料用量
// @Description Calculate material usage based on area
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param request body request.MaterialCalculateRequest true "计算请求"
// @Success 200 {object} response.Response{data=response.MaterialCalculateResponse}
// @Router /api/designer/materials/calculate [post]
func CalculateMaterialUsage(c *gin.Context) {
	var req request.MaterialCalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := getMaterialService().CalculateMaterialUsage(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// BatchCalculateMaterialUsage 批量计算材料用量
// @Summary 批量计算材料用量
// @Description Batch calculate material usage
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param request body request.MaterialBatchCalculateRequest true "批量计算请求"
// @Success 200 {object} response.Response{data=response.MaterialBatchCalculateResponse}
// @Router /api/designer/materials/batch-calculate [post]
func BatchCalculateMaterialUsage(c *gin.Context) {
	var req request.MaterialBatchCalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := getMaterialService().BatchCalculateMaterialUsage(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// RecommendMaterials 智能推荐材料
// @Summary 智能推荐材料
// @Description Intelligently recommend materials based on project requirements
// @Tags 材料管理
// @Accept json
// @Produce json
// @Param request body request.MaterialRecommendRequest true "推荐请求"
// @Success 200 {object} response.Response{data=response.MaterialRecommendResponse}
// @Router /api/designer/materials/recommend [post]
func RecommendMaterials(c *gin.Context) {
	var req request.MaterialRecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := getMaterialService().RecommendMaterials(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}
