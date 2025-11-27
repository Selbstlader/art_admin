package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	appUserSvc "art_admin_backend/internal/service/app_user"
	"strconv"

	"github.com/gin-gonic/gin"
)

var appUserService = appUserSvc.NewAppUserService()

// GetAppUserList 获取APP用户列表
// @Summary 获取APP用户列表
// @Description 分页查询APP用户列表
// @Tags APP用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param current query int true "当前页码" minimum(1)
// @Param size query int true "每页条数" minimum(1) maximum(100)
// @Param id query int false "用户ID"
// @Param userName query string false "用户账号"
// @Param nickName query string false "用户昵称"
// @Param phone query string false "手机号"
// @Param userType query string false "用户类型"
// @Param status query string false "状态"
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/app-user/list [get]
func GetAppUserList(c *gin.Context) {
	var req request.AppUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	users, total, err := appUserService.GetAppUserList(&req)
	if err != nil {
		response.ServerError(c, "查询APP用户列表失败")
		return
	}

	response.SuccessWithPagination(c, users, req.Current, req.Size, total)
}

// CreateAppUser 创建APP用户
// @Summary 创建APP用户
// @Description 创建新APP用户
// @Tags APP用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateAppUserRequest true "用户信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/app-user/create [post]
func CreateAppUser(c *gin.Context) {
	var req request.CreateAppUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := appUserService.CreateAppUser(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建APP用户成功", nil)
}

// UpdateAppUser 更新APP用户
// @Summary 更新APP用户
// @Description 更新APP用户信息
// @Tags APP用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateAppUserRequest true "用户信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/app-user/update [put]
func UpdateAppUser(c *gin.Context) {
	var req request.UpdateAppUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := appUserService.UpdateAppUser(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新APP用户成功", nil)
}

// DeleteAppUser 删除APP用户
// @Summary 删除APP用户
// @Description 删除APP用户
// @Tags APP用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/app-user/delete/{id} [delete]
func DeleteAppUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := appUserService.DeleteAppUser(id); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除APP用户成功", nil)
}

// ResetAppUserPassword 重置APP用户密码
// @Summary 重置APP用户密码
// @Description 重置APP用户密码
// @Tags APP用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.ResetAppUserPasswordRequest true "重置密码信息"
// @Success 200 {object} response.Response "重置成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/app-user/reset-password [post]
func ResetAppUserPassword(c *gin.Context) {
	var req request.ResetAppUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := appUserService.ResetAppUserPassword(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "重置密码成功", nil)
}
