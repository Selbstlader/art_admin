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

// processDefService 流程定义服务实例
var processDefService workflow.ProcessDefinitionService

// initProcessDefService 初始化流程定义服务
func initProcessDefService() {
	if processDefService == nil {
		repo := repository.NewProcessDefinitionRepository()
		serializer := workflow.NewSerializer()
		validator := workflow.NewValidator()
		processDefService = workflow.NewProcessDefinitionService(repo, serializer, validator)
	}
}

// GetProcessDefList 获取流程定义列表
// @Summary 获取流程定义列表
// @Description 分页查询流程定义列表
// @Tags 工作流-流程定义
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Param name query string false "流程名称"
// @Param code query string false "流程编码"
// @Param category query string false "分类"
// @Param status query string false "状态"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessDefListResponse}
// @Router /api/workflow/process-def/list [get]
func GetProcessDefList(c *gin.Context) {
	initProcessDefService()

	var req request.ListProcessDefRequest
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
	svcReq := &workflow.ListProcessDefRequest{
		Name:     req.Name,
		Code:     req.Code,
		Category: req.Category,
		Status:   req.Status,
		Current:  req.Page,
		Size:     req.PageSize,
	}

	defs, total, err := processDefService.List(context.Background(), svcReq)
	if err != nil {
		pkgResponse.ServerError(c, "查询流程定义列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]response.ProcessDefResponse, 0, len(defs))
	for _, def := range defs {
		list = append(list, toProcessDefResponse(def))
	}

	pkgResponse.SuccessWithPagination(c, list, req.Page, req.PageSize, total)
}

// GetProcessDefDetail 获取流程定义详情
// @Summary 获取流程定义详情
// @Description 根据ID获取流程定义详情
// @Tags 工作流-流程定义
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流程定义ID"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessDefDetailResponse}
// @Router /api/workflow/process-def/{id} [get]
func GetProcessDefDetail(c *gin.Context) {
	initProcessDefService()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	def, err := processDefService.GetByID(context.Background(), id)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeNotFound, err.Error())
		return
	}

	// 获取流程图
	graph, err := processDefService.GetGraph(context.Background(), id)
	if err != nil {
		pkgResponse.ServerError(c, "获取流程图失败: "+err.Error())
		return
	}

	resp := toProcessDefDetailResponse(def, graph)
	pkgResponse.Success(c, resp)
}

// CreateProcessDef 创建流程定义
// @Summary 创建流程定义
// @Description 创建新的流程定义
// @Tags 工作流-流程定义
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateProcessDefRequest true "流程定义信息"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessDefResponse}
// @Router /api/workflow/process-def [post]
func CreateProcessDef(c *gin.Context) {
	initProcessDefService()

	var req request.CreateProcessDefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetUserID(c)

	// 转换为服务层请求
	svcReq := &workflow.CreateProcessDefRequest{
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		Category:       req.Category,
		FormTemplateID: req.FormTemplateID,
		Graph:          toModelProcessGraph(req.Graph),
		CreatedBy:      userID,
	}

	def, err := processDefService.Create(context.Background(), svcReq)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "创建流程定义成功", toProcessDefResponse(def))
}

// UpdateProcessDef 更新流程定义
// @Summary 更新流程定义
// @Description 更新流程定义（仅草稿状态可更新）
// @Tags 工作流-流程定义
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流程定义ID"
// @Param request body request.UpdateProcessDefRequest true "流程定义信息"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessDefResponse}
// @Router /api/workflow/process-def/{id} [put]
func UpdateProcessDef(c *gin.Context) {
	initProcessDefService()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	var req request.UpdateProcessDefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 转换为服务层请求
	svcReq := &workflow.UpdateProcessDefRequest{
		Name:           req.Name,
		Description:    req.Description,
		Category:       req.Category,
		FormTemplateID: req.FormTemplateID,
		Graph:          toModelProcessGraph(req.Graph),
	}

	def, err := processDefService.Update(context.Background(), id, svcReq)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "更新流程定义成功", toProcessDefResponse(def))
}

// DeleteProcessDef 删除流程定义
// @Summary 删除流程定义
// @Description 删除流程定义（仅草稿状态可删除）
// @Tags 工作流-流程定义
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流程定义ID"
// @Success 200 {object} pkgResponse.Response
// @Router /api/workflow/process-def/{id} [delete]
func DeleteProcessDef(c *gin.Context) {
	initProcessDefService()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	if err := processDefService.Delete(context.Background(), id); err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "删除流程定义成功", nil)
}

// PublishProcessDef 发布流程定义
// @Summary 发布流程定义
// @Description 发布流程定义，生成新版本
// @Tags 工作流-流程定义
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流程定义ID"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessDefResponse}
// @Router /api/workflow/process-def/{id}/publish [post]
func PublishProcessDef(c *gin.Context) {
	initProcessDefService()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	def, err := processDefService.Publish(context.Background(), id)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "发布流程定义成功", toProcessDefResponse(def))
}

// ============================================================================
// 辅助函数：DTO转换
// ============================================================================

// toProcessDefResponse 转换为流程定义响应
func toProcessDefResponse(def *model.ProcessDefinition) response.ProcessDefResponse {
	return response.ProcessDefResponse{
		ID:             def.ID,
		Name:           def.Name,
		Code:           def.Code,
		Description:    def.Description,
		Category:       def.Category,
		FormTemplateID: def.FormTemplateID,
		Version:        def.Version,
		Status:         def.Status,
		CreatedBy:      def.CreatedBy,
		CreatedAt:      def.CreatedAt,
		UpdatedAt:      def.UpdatedAt,
	}
}

// toProcessDefDetailResponse 转换为流程定义详情响应
func toProcessDefDetailResponse(def *model.ProcessDefinition, graph *model.ProcessGraph) response.ProcessDefDetailResponse {
	resp := response.ProcessDefDetailResponse{
		ProcessDefResponse: toProcessDefResponse(def),
	}

	if graph != nil {
		resp.Graph = toProcessGraphResponse(graph)
	}

	return resp
}

// toProcessGraphResponse 转换为流程图响应
func toProcessGraphResponse(graph *model.ProcessGraph) *response.ProcessGraphResponse {
	if graph == nil {
		return nil
	}

	resp := &response.ProcessGraphResponse{
		Nodes: make([]response.ProcessNodeResponse, 0, len(graph.Nodes)),
		Edges: make([]response.ProcessEdgeResponse, 0, len(graph.Edges)),
	}

	for _, node := range graph.Nodes {
		resp.Nodes = append(resp.Nodes, toProcessNodeResponse(&node))
	}

	for _, edge := range graph.Edges {
		resp.Edges = append(resp.Edges, toProcessEdgeResponse(&edge))
	}

	if graph.GlobalProps != nil {
		resp.GlobalProps = &response.GlobalPropertiesResponse{
			AllowWithdraw: graph.GlobalProps.AllowWithdraw,
		}
	}

	return resp
}

// toProcessNodeResponse 转换为流程节点响应
func toProcessNodeResponse(node *model.ProcessNode) response.ProcessNodeResponse {
	resp := response.ProcessNodeResponse{
		ID:   node.ID,
		Type: string(node.Type),
		Name: node.Name,
		Position: response.PositionResponse{
			X: node.Position.X,
			Y: node.Position.Y,
		},
	}

	if node.Properties != nil {
		resp.Properties = &response.NodePropertiesResponse{
			ApprovalMode:     string(node.Properties.ApprovalMode),
			TimeoutHours:     node.Properties.TimeoutHours,
			RejectAction:     string(node.Properties.RejectAction),
			FieldPermissions: node.Properties.FieldPermissions,
		}

		if node.Properties.AssigneeRule != nil {
			resp.Properties.AssigneeRule = &response.AssigneeRuleResponse{
				Type:   string(node.Properties.AssigneeRule.Type),
				Values: node.Properties.AssigneeRule.Values,
			}
		}
	}

	return resp
}

// toProcessEdgeResponse 转换为流程边响应
func toProcessEdgeResponse(edge *model.ProcessEdge) response.ProcessEdgeResponse {
	return response.ProcessEdgeResponse{
		ID:        edge.ID,
		Source:    edge.Source,
		Target:    edge.Target,
		Condition: edge.Condition,
	}
}

// toModelProcessGraph 转换请求DTO为模型
func toModelProcessGraph(dto *request.ProcessGraphDTO) *model.ProcessGraph {
	if dto == nil {
		return nil
	}

	graph := &model.ProcessGraph{
		Nodes: make([]model.ProcessNode, 0, len(dto.Nodes)),
		Edges: make([]model.ProcessEdge, 0, len(dto.Edges)),
	}

	for _, nodeDTO := range dto.Nodes {
		graph.Nodes = append(graph.Nodes, toModelProcessNode(&nodeDTO))
	}

	for _, edgeDTO := range dto.Edges {
		graph.Edges = append(graph.Edges, toModelProcessEdge(&edgeDTO))
	}

	if dto.GlobalProps != nil {
		graph.GlobalProps = &model.GlobalProperties{
			AllowWithdraw: dto.GlobalProps.AllowWithdraw,
		}
	}

	return graph
}

// toModelProcessNode 转换节点DTO为模型
func toModelProcessNode(dto *request.ProcessNodeDTO) model.ProcessNode {
	node := model.ProcessNode{
		ID:   dto.ID,
		Type: model.NodeType(dto.Type),
		Name: dto.Name,
		Position: model.Position{
			X: dto.Position.X,
			Y: dto.Position.Y,
		},
	}

	if dto.Properties != nil {
		node.Properties = &model.NodeProperties{
			ApprovalMode:     model.ApprovalMode(dto.Properties.ApprovalMode),
			TimeoutHours:     dto.Properties.TimeoutHours,
			RejectAction:     model.RejectAction(dto.Properties.RejectAction),
			FieldPermissions: dto.Properties.FieldPermissions,
		}

		if dto.Properties.AssigneeRule != nil {
			node.Properties.AssigneeRule = &model.AssigneeRule{
				Type:   model.AssigneeType(dto.Properties.AssigneeRule.Type),
				Values: dto.Properties.AssigneeRule.Values,
			}
		}
	}

	return node
}

// toModelProcessEdge 转换边DTO为模型
func toModelProcessEdge(dto *request.ProcessEdgeDTO) model.ProcessEdge {
	return model.ProcessEdge{
		ID:        dto.ID,
		Source:    dto.Source,
		Target:    dto.Target,
		Condition: dto.Condition,
	}
}
