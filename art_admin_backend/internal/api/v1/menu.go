package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	systemSvc "art_admin_backend/internal/service/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

var menuService = systemSvc.NewMenuService()

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取当前用户的菜单树
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]response.MenuResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/system/menus [get]
func GetMenuList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未授权")
		return
	}

	menus, err := menuService.GetMenuList(userID)
	if err != nil {
		response.ServerError(c, "获取菜单列表失败")
		return
	}

	response.Success(c, menus)
}

// GetAllMenus 获取所有菜单
// @Summary 获取所有菜单
// @Description 获取所有菜单（管理用）
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]response.MenuResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/menu/all [get]
func GetAllMenus(c *gin.Context) {
	menus, err := menuService.GetAllMenus()
	if err != nil {
		response.ServerError(c, "获取菜单列表失败")
		return
	}

	response.Success(c, menus)
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateMenuRequest true "菜单信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/menu/create [post]
func CreateMenu(c *gin.Context) {
	var req request.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := menuService.CreateMenu(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建菜单成功", nil)
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateMenuRequest true "菜单信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/menu/update [put]
func UpdateMenu(c *gin.Context) {
	var req request.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := menuService.UpdateMenu(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新菜单成功", nil)
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/menu/delete/{id} [delete]
func DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := menuService.DeleteMenu(id); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除菜单成功", nil)
}
