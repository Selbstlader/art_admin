package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	authSvc "art_admin_backend/internal/service/auth"

	"github.com/gin-gonic/gin"
)

var authService = authSvc.NewAuthService()

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录参数"
// @Success 200 {object} response.Response{data=response.LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := authService.Login(&req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "登录成功", result)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.UserInfoResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/user/info [get]
func GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未授权")
		return
	}

	result, err := authService.GetUserInfo(userID)
	if err != nil {
		response.ServerError(c, "获取用户信息失败")
		return
	}

	response.Success(c, result)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 当前登录用户更新个人信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateUserInfoRequest true "用户信息参数"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/user/info [put]
func UpdateUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未授权")
		return
	}

	var req request.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := authService.UpdateUserInfo(userID, &req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", nil)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 当前登录用户修改密码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.ChangePasswordRequest true "密码参数"
// @Success 200 {object} response.Response "修改成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/user/change-password [post]
func ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未授权")
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证新密码和确认密码是否一致
	if req.NewPassword != req.ConfirmPassword {
		response.BadRequest(c, "新密码和确认密码不一致")
		return
	}

	err := authService.ChangePassword(userID, &req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "密码修改成功", nil)
}
