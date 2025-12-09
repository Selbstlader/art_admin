package v1

import (
	"context"
	"strconv"

	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	pkgResponse "art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/repository"
	"art_admin_backend/internal/service/workflow"

	"github.com/gin-gonic/gin"
)

// formEngineService 表单引擎服务实例
var formEngineService workflow.FormEngineService

// initFormEngineService 初始化表单引擎服务
func initFormEngineService() {
	if formEngineService == nil {
		repo := repository.NewFormTemplateRepository()
		serializer := workflow.NewSerializer()
		formEngineService = workflow.NewFormEngineService(repo, serializer)
	}
}

// GetFormTemplateList 获取表单模板列表
// @Summary 获取表单模板列表
// @Description 分页查询表单模板列表
// @Tags 工作流-表单模板
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Param name query string false "模板名称"
// @Param code query string false "模板编码"
// @Param status query string false "状态"
// @Success 200 {object} pkgResponse.Response{data=response.FormTemplateListResponse}
// @Router /api/workflow/form-template/list [get]
func GetFormTemplateList(c *gin.Context) {
	initFormEngineService()

	var req request.ListFormTemplateRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 转换为服务层请求
	svcReq := &workflow.ListFormTemplateRequest{
		Name:    req.Name,
		Code:    req.Code,
		Status:  req.Status,
		Current: req.Page,
		Size:    req.PageSize,
	}

	templates, total, err := formEngineService.ListTemplates(context.Background(), svcReq)
	if err != nil {
		pkgResponse.ServerError(c, "查询表单模板列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]response.FormTemplateResponse, 0, len(templates))
	for _, t := range templates {
		list = append(list, toFormTemplateResponse(t))
	}

	pkgResponse.SuccessWithPagination(c, list, req.Page, req.PageSize, total)
}

// GetFormTemplateDetail 获取表单模板详情
// @Summary 获取表单模板详情
// @Description 根据ID获取表单模板详情
// @Tags 工作流-表单模板
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "表单模板ID"
// @Success 200 {object} pkgResponse.Response{data=response.FormTemplateDetailResponse}
// @Router /api/workflow/form-template/{id} [get]
func GetFormTemplateDetail(c *gin.Context) {
	initFormEngineService()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	template, err := formEngineService.GetTemplate(context.Background(), id)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeNotFound, err.Error())
		return
	}

	// 获取表单结构
	schema, err := formEngineService.GetSchema(context.Background(), id)
	if err != nil {
		pkgResponse.ServerError(c, "获取表单结构失败: "+err.Error())
		return
	}

	resp := toFormTemplateDetailResponse(template, schema)
	pkgResponse.Success(c, resp)
}

// CreateFormTemplate 创建表单模板
// @Summary 创建表单模板
// @Description 创建新的表单模板
// @Tags 工作流-表单模板
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateFormTemplateRequest true "表单模板信息"
// @Success 200 {object} pkgResponse.Response{data=response.FormTemplateResponse}
// @Router /api/workflow/form-template [post]
func CreateFormTemplate(c *gin.Context) {
	initFormEngineService()

	var req request.CreateFormTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetUserID(c)

	// 转换为服务层请求
	svcReq := &workflow.CreateFormTemplateRequest{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Schema:      toModelFormSchema(req.Schema),
		CreatedBy:   userID,
	}

	template, err := formEngineService.CreateTemplate(context.Background(), svcReq)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "创建表单模板成功", toFormTemplateResponse(template))
}

// UpdateFormTemplate 更新表单模板
// @Summary 更新表单模板
// @Description 更新表单模板
// @Tags 工作流-表单模板
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "表单模板ID"
// @Param request body request.UpdateFormTemplateRequest true "表单模板信息"
// @Success 200 {object} pkgResponse.Response{data=response.FormTemplateResponse}
// @Router /api/workflow/form-template/{id} [put]
func UpdateFormTemplate(c *gin.Context) {
	initFormEngineService()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	var req request.UpdateFormTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 转换为服务层请求
	svcReq := &workflow.UpdateFormTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		Schema:      toModelFormSchema(req.Schema),
	}

	template, err := formEngineService.UpdateTemplate(context.Background(), id, svcReq)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "更新表单模板成功", toFormTemplateResponse(template))
}

// DeleteFormTemplate 删除表单模板
// @Summary 删除表单模板
// @Description 删除表单模板
// @Tags 工作流-表单模板
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "表单模板ID"
// @Success 200 {object} pkgResponse.Response
// @Router /api/workflow/form-template/{id} [delete]
func DeleteFormTemplate(c *gin.Context) {
	initFormEngineService()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	if err := formEngineService.DeleteTemplate(context.Background(), id); err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "删除表单模板成功", nil)
}

// ============================================================================
// 辅助函数：DTO转换
// ============================================================================

// toFormTemplateResponse 转换为表单模板响应
func toFormTemplateResponse(t *model.FormTemplate) response.FormTemplateResponse {
	return response.FormTemplateResponse{
		ID:          t.ID,
		Name:        t.Name,
		Code:        t.Code,
		Description: t.Description,
		Status:      t.Status,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// toFormTemplateDetailResponse 转换为表单模板详情响应
func toFormTemplateDetailResponse(t *model.FormTemplate, schema *model.FormSchema) response.FormTemplateDetailResponse {
	resp := response.FormTemplateDetailResponse{
		FormTemplateResponse: toFormTemplateResponse(t),
	}

	if schema != nil {
		resp.Schema = toFormSchemaResponse(schema)
	}

	return resp
}

// toFormSchemaResponse 转换为表单结构响应
func toFormSchemaResponse(schema *model.FormSchema) *response.FormSchemaResponse {
	if schema == nil {
		return nil
	}

	resp := &response.FormSchemaResponse{
		Fields: make([]response.FormFieldResponse, 0, len(schema.Fields)),
	}

	for _, field := range schema.Fields {
		resp.Fields = append(resp.Fields, toFormFieldResponse(&field))
	}

	return resp
}

// toFormFieldResponse 转换为表单字段响应
func toFormFieldResponse(field *model.FormField) response.FormFieldResponse {
	resp := response.FormFieldResponse{
		Key:          field.Key,
		Label:        field.Label,
		Type:         string(field.Type),
		Required:     field.Required,
		Placeholder:  field.Placeholder,
		DefaultValue: field.DefaultValue,
	}

	if len(field.Options) > 0 {
		resp.Options = make([]response.SelectOptionResponse, 0, len(field.Options))
		for _, opt := range field.Options {
			resp.Options = append(resp.Options, response.SelectOptionResponse{
				Label: opt.Label,
				Value: opt.Value,
			})
		}
	}

	if field.Validation != nil {
		resp.Validation = &response.FieldValidationResponse{
			MinLength: field.Validation.MinLength,
			MaxLength: field.Validation.MaxLength,
			Min:       field.Validation.Min,
			Max:       field.Validation.Max,
			Pattern:   field.Validation.Pattern,
		}
	}

	return resp
}

// toModelFormSchema 转换请求DTO为模型
func toModelFormSchema(dto *request.FormSchemaDTO) *model.FormSchema {
	if dto == nil {
		return nil
	}

	schema := &model.FormSchema{
		Fields: make([]model.FormField, 0, len(dto.Fields)),
	}

	for _, fieldDTO := range dto.Fields {
		schema.Fields = append(schema.Fields, toModelFormField(&fieldDTO))
	}

	return schema
}

// toModelFormField 转换字段DTO为模型
func toModelFormField(dto *request.FormFieldDTO) model.FormField {
	field := model.FormField{
		Key:          dto.Key,
		Label:        dto.Label,
		Type:         model.FieldType(dto.Type),
		Required:     dto.Required,
		Placeholder:  dto.Placeholder,
		DefaultValue: dto.DefaultValue,
	}

	if len(dto.Options) > 0 {
		field.Options = make([]model.SelectOption, 0, len(dto.Options))
		for _, opt := range dto.Options {
			field.Options = append(field.Options, model.SelectOption{
				Label: opt.Label,
				Value: opt.Value,
			})
		}
	}

	if dto.Validation != nil {
		field.Validation = &model.FieldValidation{
			MinLength: dto.Validation.MinLength,
			MaxLength: dto.Validation.MaxLength,
			Min:       dto.Validation.Min,
			Max:       dto.Validation.Max,
			Pattern:   dto.Validation.Pattern,
		}
	}

	return field
}
