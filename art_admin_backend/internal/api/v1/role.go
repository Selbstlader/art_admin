package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var roleService = service.NewRoleService()

// GetRoleList 获取角色列表
// @Summary 获取角色列表
// @Description 分页查询角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param current query int true "当前页码" minimum(1)
// @Param size query int true "每页条数" minimum(1) maximum(100)
// @Param roleId query int false "角色ID"
// @Param roleName query string false "角色名"
// @Param roleCode query string false "角色编码"
// @Param description query string false "描述"
// @Param enabled query boolean false "是否启用"
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/role/list [get]
func GetRoleList(c *gin.Context) {
	var req request.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	roles, total, err := roleService.GetRoleList(&req)
	if err != nil {
		response.ServerError(c, "查询角色列表失败")
		return
	}

	response.SuccessWithPagination(c, roles, req.Current, req.Size, total)
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateRoleRequest true "角色信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/role/create [post]
func CreateRole(c *gin.Context) {
	var req request.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := roleService.CreateRole(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建角色成功", nil)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateRoleRequest true "角色信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/role/update [put]
func UpdateRole(c *gin.Context) {
	var req request.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := roleService.UpdateRole(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新角色成功", nil)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/role/delete/{id} [delete]
func DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := roleService.DeleteRole(id); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除角色成功", nil)
}

// GetRolePermissions 获取角色权限
// @Summary 获取角色权限
// @Description 获取角色的菜单和按钮权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/role/permissions/{id} [get]
func GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	menuIDs, buttonIDs, err := roleService.GetRolePermissions(id)
	if err != nil {
		response.ServerError(c, "查询角色权限失败")
		return
	}

	response.Success(c, map[string]interface{}{
		"menuIds":   menuIDs,
		"buttonIds": buttonIDs,
	})
}

// UpdateRolePermissions 更新角色权限
// @Summary 更新角色权限
// @Description 更新角色的菜单和按钮权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param request body request.UpdateRolePermissionsRequest true "权限信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/role/permissions/{id} [put]
func UpdateRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := roleService.UpdateRolePermissions(id, &req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新角色权限成功", nil)
}
