package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	pkgResponse "art_admin_backend/internal/pkg/response"
	aiTagSvc "art_admin_backend/internal/service/ai_tag"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*** AI Tag API Handler ***/

var aiTagService = aiTagSvc.NewAITagService()

// GetAITagList retrieves paginated AI tags with optional keyword search
// @Summary Get AI tag list
// @Description Get paginated list of AI tags with optional keyword search
// @Tags AI标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param current query int true "Current page" minimum(1)
// @Param size query int true "Page size" minimum(1) maximum(100)
// @Param keyword query string false "Search keyword"
// @Success 200 {object} pkgResponse.Response{data=response.AITagListResponse} "Success"
// @Failure 400 {object} pkgResponse.Response "Bad request"
// @Failure 401 {object} pkgResponse.Response "Unauthorized"
// @Router /api/v1/ai-tags [get]
func GetAITagList(c *gin.Context) {
	var req request.AITagListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := aiTagService.ListTags(req.Current, req.Size, req.Keyword)
	if err != nil {
		pkgResponse.ServerError(c, "查询AI标签列表失败: "+err.Error())
		return
	}

	// Convert to response DTO
	listResp := convertToAITagListResponse(result)
	pkgResponse.SuccessWithPagination(c, listResp.List, req.Current, req.Size, listResp.Total)
}

// GetAITagByID retrieves an AI tag by ID
// @Summary Get AI tag by ID
// @Description Get AI tag details by ID
// @Tags AI标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "AI Tag ID"
// @Success 200 {object} pkgResponse.Response{data=response.AITagResponse} "Success"
// @Failure 400 {object} pkgResponse.Response "Bad request"
// @Failure 401 {object} pkgResponse.Response "Unauthorized"
// @Failure 404 {object} pkgResponse.Response "Not found"
// @Router /api/v1/ai-tags/{id} [get]
func GetAITagByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID")
		return
	}

	tag, err := aiTagService.GetTag(id)
	if err != nil {
		if err.Error() == "tag not found" {
			pkgResponse.NotFound(c, "AI标签不存在")
		} else {
			pkgResponse.ServerError(c, "查询AI标签失败: "+err.Error())
		}
		return
	}

	pkgResponse.Success(c, convertToAITagResponse(tag))
}

// CreateAITag creates a new AI tag
// @Summary Create AI tag
// @Description Create a new AI tag
// @Tags AI标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateAITagRequest true "AI Tag data"
// @Success 200 {object} pkgResponse.Response{data=response.AITagResponse} "Success"
// @Failure 400 {object} pkgResponse.Response "Bad request"
// @Failure 401 {object} pkgResponse.Response "Unauthorized"
// @Failure 500 {object} pkgResponse.Response "Server error"
// @Router /api/v1/ai-tags [post]
func CreateAITag(c *gin.Context) {
	var req request.CreateAITagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// Convert request DTO to service request
	svcReq := &aiTagSvc.CreateAITagRequest{
		Name:              req.Name,
		Description:       req.Description,
		KnowledgeBaseID:   req.KnowledgeBaseID,
		KnowledgeBaseName: req.KnowledgeBaseName,
		SystemPrompt:      req.SystemPrompt,
		ChatAPIKey:        req.ChatAPIKey,
	}

	tag, err := aiTagService.CreateTag(svcReq)
	if err != nil {
		if err.Error() == "tag name already exists" {
			pkgResponse.BadRequest(c, "标签名称已存在")
		} else if err.Error() == "tag name cannot be empty or whitespace only" ||
			err.Error() == "system prompt cannot be empty or whitespace only" {
			pkgResponse.BadRequest(c, err.Error())
		} else {
			pkgResponse.ServerError(c, "创建AI标签失败: "+err.Error())
		}
		return
	}

	pkgResponse.Success(c, convertToAITagResponse(tag))
}

// UpdateAITag updates an existing AI tag
// @Summary Update AI tag
// @Description Update an existing AI tag
// @Tags AI标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "AI Tag ID"
// @Param data body request.UpdateAITagRequest true "AI Tag data"
// @Success 200 {object} pkgResponse.Response{data=response.AITagResponse} "Success"
// @Failure 400 {object} pkgResponse.Response "Bad request"
// @Failure 401 {object} pkgResponse.Response "Unauthorized"
// @Failure 404 {object} pkgResponse.Response "Not found"
// @Failure 500 {object} pkgResponse.Response "Server error"
// @Router /api/v1/ai-tags/{id} [put]
func UpdateAITag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID")
		return
	}

	var req request.UpdateAITagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// Convert request DTO to service request
	svcReq := &aiTagSvc.UpdateAITagRequest{
		Name:              req.Name,
		Description:       req.Description,
		KnowledgeBaseID:   req.KnowledgeBaseID,
		KnowledgeBaseName: req.KnowledgeBaseName,
		SystemPrompt:      req.SystemPrompt,
		ChatAPIKey:        req.ChatAPIKey,
		Status:            req.Status,
	}

	tag, err := aiTagService.UpdateTag(id, svcReq)
	if err != nil {
		switch err.Error() {
		case "tag not found":
			pkgResponse.NotFound(c, "AI标签不存在")
		case "tag name already exists":
			pkgResponse.BadRequest(c, "标签名称已存在")
		case "tag name cannot be empty or whitespace only",
			"system prompt cannot be empty or whitespace only":
			pkgResponse.BadRequest(c, err.Error())
		default:
			pkgResponse.ServerError(c, "更新AI标签失败: "+err.Error())
		}
		return
	}

	pkgResponse.Success(c, convertToAITagResponse(tag))
}

// DeleteAITag soft deletes an AI tag
// @Summary Delete AI tag
// @Description Soft delete an AI tag by setting status to inactive
// @Tags AI标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "AI Tag ID"
// @Success 200 {object} pkgResponse.Response "Success"
// @Failure 400 {object} pkgResponse.Response "Bad request"
// @Failure 401 {object} pkgResponse.Response "Unauthorized"
// @Failure 404 {object} pkgResponse.Response "Not found"
// @Failure 500 {object} pkgResponse.Response "Server error"
// @Router /api/v1/ai-tags/{id} [delete]
func DeleteAITag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID")
		return
	}

	err = aiTagService.DeleteTag(id)
	if err != nil {
		if err.Error() == "tag not found" {
			pkgResponse.NotFound(c, "AI标签不存在")
		} else {
			pkgResponse.ServerError(c, "删除AI标签失败: "+err.Error())
		}
		return
	}

	pkgResponse.Success(c, nil)
}

// TestAITag tests an AI tag configuration
// @Summary Test AI tag configuration
// @Description Test AI tag configuration by calling Dify API
// @Tags AI标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "AI Tag ID"
// @Success 200 {object} pkgResponse.Response{data=response.AITagTestResponse} "Success"
// @Failure 400 {object} pkgResponse.Response "Bad request"
// @Failure 401 {object} pkgResponse.Response "Unauthorized"
// @Failure 404 {object} pkgResponse.Response "Not found"
// @Failure 500 {object} pkgResponse.Response "Server error"
// @Router /api/v1/ai-tags/{id}/test [post]
func TestAITag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID")
		return
	}

	result, err := aiTagService.TestTag(id)
	if err != nil {
		pkgResponse.ServerError(c, "测试AI标签失败: "+err.Error())
		return
	}

	// Convert to response DTO
	testResp := &response.AITagTestResponse{
		Success:           result.Success,
		Response:          result.Response,
		KnowledgeAccessed: result.KnowledgeAccessed,
		ErrorMessage:      result.ErrorMessage,
	}

	pkgResponse.Success(c, testResp)
}

/*** Helper functions for converting model to response DTO ***/

// convertToAITagResponse converts model.AITag to response.AITagResponse
func convertToAITagResponse(tag *model.AITag) *response.AITagResponse {
	if tag == nil {
		return nil
	}
	return &response.AITagResponse{
		ID:                tag.ID,
		Name:              tag.Name,
		Description:       tag.Description,
		KnowledgeBaseID:   tag.KnowledgeBaseID,
		KnowledgeBaseName: tag.KnowledgeBaseName,
		SystemPrompt:      tag.SystemPrompt,
		ChatAPIKey:        tag.ChatAPIKey,
		Status:            tag.Status,
		CreatedAt:         tag.CreatedAt,
		UpdatedAt:         tag.UpdatedAt,
	}
}

// convertToAITagListResponse converts service response to DTO response
func convertToAITagListResponse(result *aiTagSvc.AITagListResponse) *response.AITagListResponse {
	if result == nil {
		return nil
	}

	list := make([]*response.AITagResponse, 0, len(result.List))
	for _, tag := range result.List {
		list = append(list, &response.AITagResponse{
			ID:                tag.ID,
			Name:              tag.Name,
			Description:       tag.Description,
			KnowledgeBaseID:   tag.KnowledgeBaseID,
			KnowledgeBaseName: tag.KnowledgeBaseName,
			SystemPrompt:      tag.SystemPrompt,
			ChatAPIKey:        tag.ChatAPIKey,
			Status:            tag.Status,
			CreatedAt:         tag.CreatedAt,
			UpdatedAt:         tag.UpdatedAt,
		})
	}

	return &response.AITagListResponse{
		List:  list,
		Total: result.Total,
		Page:  result.Page,
		Size:  result.Size,
	}
}
