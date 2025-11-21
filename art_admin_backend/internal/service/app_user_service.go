package service

import (
	"errors"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/utils"
	"art_admin_backend/internal/repository"
)

type AppUserService struct {
	appUserRepo *repository.AppUserRepository
}

func NewAppUserService() *AppUserService {
	return &AppUserService{
		appUserRepo: repository.NewAppUserRepository(),
	}
}

// GetAppUserList 获取APP用户列表
func (s *AppUserService) GetAppUserList(req *request.AppUserListRequest) ([]response.AppUserListItem, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.ID != nil {
		query["id"] = *req.ID
	}
	if req.UserName != "" {
		query["userName"] = req.UserName
	}
	if req.NickName != "" {
		query["nickName"] = req.NickName
	}
	if req.Phone != "" {
		query["phone"] = req.Phone
	}
	if req.UserType != "" {
		query["userType"] = req.UserType
	}
	if req.Status != "" {
		query["status"] = req.Status
	}

	// 查询APP用户列表
	users, total, err := s.appUserRepo.FindWithPagination(query, req.Current, req.Size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应DTO
	result := make([]response.AppUserListItem, 0, len(users))
	for _, user := range users {
		result = append(result, response.AppUserListItem{
			ID:            user.ID,
			UserName:      user.UserName,
			NickName:      user.NickName,
			Phone:         user.Phone,
			Email:         user.Email,
			Avatar:        user.Avatar,
			UserType:      user.UserType,
			Status:        user.Status,
			LastLoginTime: user.LastLoginTime,
			LastLoginIP:   user.LastLoginIP,
			CreateBy:      user.CreateBy,
			CreateTime:    user.CreateTime,
			UpdateBy:      user.UpdateBy,
			UpdateTime:    user.UpdateTime,
		})
	}

	return result, total, nil
}

// CreateAppUser 创建APP用户
func (s *AppUserService) CreateAppUser(req *request.CreateAppUserRequest) error {
	// 检查用户名是否已存在
	existUser, _ := s.appUserRepo.FindByUserName(req.UserName)
	if existUser != nil {
		return errors.New("用户账号已存在")
	}

	// 检查手机号是否已存在
	existPhone, _ := s.appUserRepo.FindByPhone(req.Phone)
	if existPhone != nil {
		return errors.New("手机号已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建APP用户
	user := &model.AppUser{
		UserName: req.UserName,
		NickName: req.NickName,
		Password: hashedPassword,
		Phone:    req.Phone,
		Email:    req.Email,
		Avatar:   req.Avatar,
		UserType: req.UserType,
		Status:   req.Status,
	}

	if err := s.appUserRepo.Create(user); err != nil {
		return errors.New("创建APP用户失败")
	}

	return nil
}

// UpdateAppUser 更新APP用户
func (s *AppUserService) UpdateAppUser(req *request.UpdateAppUserRequest) error {
	// 检查用户是否存在
	user, err := s.appUserRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("APP用户不存在")
	}

	// 检查手机号是否被其他用户使用
	existPhone, _ := s.appUserRepo.FindByPhone(req.Phone)
	if existPhone != nil && existPhone.ID != req.ID {
		return errors.New("手机号已被其他用户使用")
	}

	// 更新用户信息
	user.NickName = req.NickName
	user.Phone = req.Phone
	user.Email = req.Email
	user.Avatar = req.Avatar
	user.UserType = req.UserType
	user.Status = req.Status

	if err := s.appUserRepo.Update(user); err != nil {
		return errors.New("更新APP用户失败")
	}

	return nil
}

// DeleteAppUser 删除APP用户
func (s *AppUserService) DeleteAppUser(id int64) error {
	// 检查用户是否存在
	_, err := s.appUserRepo.FindByID(id)
	if err != nil {
		return errors.New("APP用户不存在")
	}

	if err := s.appUserRepo.Delete(id); err != nil {
		return errors.New("删除APP用户失败")
	}

	return nil
}

// ResetAppUserPassword 重置APP用户密码
func (s *AppUserService) ResetAppUserPassword(req *request.ResetAppUserPasswordRequest) error {
	// 检查用户是否存在
	user, err := s.appUserRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("APP用户不存在")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.appUserRepo.Update(user); err != nil {
		return errors.New("重置密码失败")
	}

	return nil
}
