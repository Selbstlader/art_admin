package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	dictionarySvc "art_admin_backend/internal/service/dictionary"
	"strconv"

	"github.com/gin-gonic/gin"
)

var dictionaryService = dictionarySvc.NewDictionaryService()

// ========== 字典类型相关 ==========

// GetDictionaryTypeList 获取字典类型列表
// @Summary 获取字典类型列表
// @Description 分页查询字典类型列表
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param current query int true "当前页码" minimum(1)
// @Param size query int true "每页条数" minimum(1) maximum(100)
// @Param id query int false "字典类型ID"
// @Param typeName query string false "字典类型名称"
// @Param typeCode query string false "字典类型编码"
// @Param description query string false "描述"
// @Param enabled query boolean false "是否启用"
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dictionary/type/list [get]
func GetDictionaryTypeList(c *gin.Context) {
	var req request.DictionaryTypeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	dictTypes, total, err := dictionaryService.GetDictionaryTypeList(&req)
	if err != nil {
		response.ServerError(c, "查询字典类型列表失败")
		return
	}

	response.SuccessWithPagination(c, dictTypes, req.Current, req.Size, total)
}

// GetDictionaryTypeByID 根据ID获取字典类型
// @Summary 根据ID获取字典类型
// @Description 根据ID获取字典类型详情
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "字典类型ID"
// @Success 200 {object} response.Response{data=response.DictionaryTypeResponse} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "字典类型不存在"
// @Router /api/dictionary/type/{id} [get]
func GetDictionaryTypeByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	dictType, err := dictionaryService.GetDictionaryTypeByID(id)
	if err != nil {
		response.NotFound(c, "字典类型不存在")
		return
	}

	response.Success(c, dictType)
}

// CreateDictionaryType 创建字典类型
// @Summary 创建字典类型
// @Description 创建新的字典类型
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateDictionaryTypeRequest true "字典类型信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "创建失败"
// @Router /api/dictionary/type [post]
func CreateDictionaryType(c *gin.Context) {
	var req request.CreateDictionaryTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户
	createBy := c.GetString("userName")
	if createBy == "" {
		createBy = "system"
	}

	err := dictionaryService.CreateDictionaryType(&req, createBy)
	if err != nil {
		response.ServerError(c, "创建字典类型失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateDictionaryType 更新字典类型
// @Summary 更新字典类型
// @Description 更新字典类型信息
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.UpdateDictionaryTypeRequest true "字典类型信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "字典类型不存在"
// @Failure 500 {object} response.Response "更新失败"
// @Router /api/dictionary/type [put]
func UpdateDictionaryType(c *gin.Context) {
	var req request.UpdateDictionaryTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户
	updateBy := c.GetString("userName")
	if updateBy == "" {
		updateBy = "system"
	}

	err := dictionaryService.UpdateDictionaryType(&req, updateBy)
	if err != nil {
		if err.Error() == "字典类型不存在" {
			response.NotFound(c, "字典类型不存在")
		} else {
			response.ServerError(c, "更新字典类型失败: "+err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// DeleteDictionaryType 删除字典类型
// @Summary 删除字典类型
// @Description 删除字典类型
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.DeleteDictionaryTypeRequest true "字典类型ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "字典类型不存在"
// @Failure 500 {object} response.Response "删除失败"
// @Router /api/dictionary/type [delete]
func DeleteDictionaryType(c *gin.Context) {
	var req request.DeleteDictionaryTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := dictionaryService.DeleteDictionaryType(req.ID)
	if err != nil {
		if err.Error() == "字典类型不存在" {
			response.NotFound(c, "字典类型不存在")
		} else {
			response.ServerError(c, "删除字典类型失败: "+err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// ========== 字典数据相关 ==========

// GetDictionaryList 获取字典数据列表
// @Summary 获取字典数据列表
// @Description 分页查询字典数据列表
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param current query int true "当前页码" minimum(1)
// @Param size query int true "每页条数" minimum(1) maximum(100)
// @Param id query int false "字典数据ID"
// @Param typeCode query string false "字典类型编码"
// @Param label query string false "字典标签"
// @Param value query string false "字典值"
// @Param enabled query boolean false "是否启用"
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dictionary/list [get]
func GetDictionaryList(c *gin.Context) {
	var req request.DictionaryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	dictionaries, total, err := dictionaryService.GetDictionaryList(&req)
	if err != nil {
		response.ServerError(c, "查询字典数据列表失败")
		return
	}

	response.SuccessWithPagination(c, dictionaries, req.Current, req.Size, total)
}

// GetDictionaryByID 根据ID获取字典数据
// @Summary 根据ID获取字典数据
// @Description 根据ID获取字典数据详情
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "字典数据ID"
// @Success 200 {object} response.Response{data=response.DictionaryResponse} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "字典数据不存在"
// @Router /api/dictionary/{id} [get]
func GetDictionaryByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	dictionary, err := dictionaryService.GetDictionaryByID(id)
	if err != nil {
		response.NotFound(c, "字典数据不存在")
		return
	}

	response.Success(c, dictionary)
}

// GetDictionaryByTypeCode 根据类型编码获取字典数据
// @Summary 根据类型编码获取字典数据
// @Description 根据类型编码获取字典数据列表（用于下拉框等）
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param typeCode query string true "字典类型编码"
// @Success 200 {object} response.Response{data=[]response.DictionaryResponse} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dictionary/by-type [get]
func GetDictionaryByTypeCode(c *gin.Context) {
	var req request.GetDictionaryByTypeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	dictionaries, err := dictionaryService.GetDictionaryByTypeCode(req.TypeCode)
	if err != nil {
		response.ServerError(c, "查询字典数据失败")
		return
	}

	response.Success(c, dictionaries)
}

// CreateDictionary 创建字典数据
// @Summary 创建字典数据
// @Description 创建新的字典数据
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateDictionaryRequest true "字典数据信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "创建失败"
// @Router /api/dictionary [post]
func CreateDictionary(c *gin.Context) {
	var req request.CreateDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户
	createBy := c.GetString("userName")
	if createBy == "" {
		createBy = "system"
	}

	err := dictionaryService.CreateDictionary(&req, createBy)
	if err != nil {
		response.ServerError(c, "创建字典数据失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateDictionary 更新字典数据
// @Summary 更新字典数据
// @Description 更新字典数据信息
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.UpdateDictionaryRequest true "字典数据信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "字典数据不存在"
// @Failure 500 {object} response.Response "更新失败"
// @Router /api/dictionary [put]
func UpdateDictionary(c *gin.Context) {
	var req request.UpdateDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户
	updateBy := c.GetString("userName")
	if updateBy == "" {
		updateBy = "system"
	}

	err := dictionaryService.UpdateDictionary(&req, updateBy)
	if err != nil {
		if err.Error() == "字典数据不存在" {
			response.NotFound(c, "字典数据不存在")
		} else {
			response.ServerError(c, "更新字典数据失败: "+err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// DeleteDictionary 删除字典数据
// @Summary 删除字典数据
// @Description 删除字典数据
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.DeleteDictionaryRequest true "字典数据ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "字典数据不存在"
// @Failure 500 {object} response.Response "删除失败"
// @Router /api/dictionary [delete]
func DeleteDictionary(c *gin.Context) {
	var req request.DeleteDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := dictionaryService.DeleteDictionary(req.ID)
	if err != nil {
		if err.Error() == "字典数据不存在" {
			response.NotFound(c, "字典数据不存在")
		} else {
			response.ServerError(c, "删除字典数据失败: "+err.Error())
		}
		return
	}

	response.Success(c, nil)
}
