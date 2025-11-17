package service

import (
	"errors"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/pkg/jwt"
	"art_admin_backend/internal/pkg/utils"
	"art_admin_backend/internal/repository"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Login 用户登录
func (s *AuthService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	// 查询用户
	user, err := s.userRepo.FindByUserName(req.UserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证用户状态
	if user.Status != "1" {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成令牌
	accessToken, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID int64) (*response.UserInfoResponse, error) {
	// 查询用户基本信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 查询用户角色
	roles, err := s.userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.RoleCode)
	}

	// 查询用户按钮权限
	buttons, err := s.userRepo.GetUserButtons(userID)
	if err != nil {
		return nil, err
	}

	authMarks := make([]string, 0, len(buttons))
	for _, button := range buttons {
		authMarks = append(authMarks, button.AuthLabel)  // Button模型使用AuthLabel
	}

	return &response.UserInfoResponse{
		UserID:   user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Roles:    roleCodes,
		Buttons:  authMarks,
	}, nil
}

