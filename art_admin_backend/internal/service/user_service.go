package service

import (
	"errors"
	"time"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/utils"
	"art_admin_backend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(req *request.UserListRequest) ([]response.UserListItem, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.ID != nil {
		query["id"] = *req.ID
	}
	if req.UserName != "" {
		query["userName"] = req.UserName
	}
	if req.UserGender != "" {
		query["userGender"] = req.UserGender
	}
	if req.UserPhone != "" {
		query["userPhone"] = req.UserPhone
	}
	if req.UserEmail != "" {
		query["userEmail"] = req.UserEmail
	}
	if req.Status != "" {
		query["status"] = req.Status
	}

	// 查询用户列表
	users, total, err := s.userRepo.FindWithPagination(query, req.Current, req.Size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应DTO
	result := make([]response.UserListItem, 0, len(users))
	for _, user := range users {
		// 查询用户角色
		roleCodes, _ := s.userRepo.GetUserRolesByUserID(user.ID)

		result = append(result, response.UserListItem{
			ID:         user.ID,
			UserName:   user.UserName,
			NickName:   user.NickName,
			Avatar:     user.Avatar,
			Status:     user.Status,
			UserGender: user.UserGender,
			UserPhone:  user.UserPhone,
			UserEmail:  user.Email,
			UserRoles:  roleCodes,
			CreateBy:   user.CreateBy,
			CreateTime: user.CreateTime,
			UpdateBy:   user.UpdateBy,
			UpdateTime: user.UpdateTime,
		})
	}

	return result, total, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *request.CreateUserRequest) error {
	// 检查用户名是否已存在
	existUser, _ := s.userRepo.FindByUserName(req.UserName)
	if existUser != nil {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		UserName:   req.UserName,
		NickName:   req.NickName,
		Password:   hashedPassword,
		Email:      req.Email,
		UserPhone:  req.UserPhone,
		UserGender: req.UserGender,
		Avatar:     req.Avatar,
		Status:     req.Status,
		CreateBy:   "admin", // TODO: 从上下文获取当前用户
		CreateTime: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.AssignRoles(user.ID, req.RoleIDs); err != nil {
			return err
		}
	}

	return nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(req *request.UpdateUserRequest) error {
	// 查询用户是否存在
	user, err := s.userRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 更新用户信息
	user.NickName = req.NickName
	user.Email = req.Email
	user.UserPhone = req.UserPhone
	user.UserGender = req.UserGender
	user.Avatar = req.Avatar
	user.Status = req.Status
	user.UpdateBy = "admin" // TODO: 从上下文获取当前用户

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 更新角色
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.AssignRoles(user.ID, req.RoleIDs); err != nil {
			return err
		}
	}

	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int64) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 不允许删除超级管理员
	if id == 1 {
		return errors.New("不能删除超级管理员")
	}

	return s.userRepo.Delete(id)
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(req *request.ResetPasswordRequest) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	return s.userRepo.UpdatePassword(req.ID, hashedPassword)
}
