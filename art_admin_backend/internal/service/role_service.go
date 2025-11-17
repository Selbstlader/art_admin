package service

import (
	"errors"
	"time"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repository.NewRoleRepository(),
	}
}

// GetRoleList 获取角色列表
func (s *RoleService) GetRoleList(req *request.RoleListRequest) ([]response.RoleListItem, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.RoleID != nil {
		query["roleId"] = *req.RoleID
	}
	if req.RoleName != "" {
		query["roleName"] = req.RoleName
	}
	if req.RoleCode != "" {
		query["roleCode"] = req.RoleCode
	}
	if req.Description != "" {
		query["description"] = req.Description
	}
	if req.Enabled != nil {
		query["enabled"] = *req.Enabled
	}

	// 查询角色列表
	roles, total, err := s.roleRepo.FindWithPagination(query, req.Current, req.Size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应DTO
	result := make([]response.RoleListItem, 0, len(roles))
	for _, role := range roles {
		result = append(result, response.RoleListItem{
			RoleID:      role.RoleID,
			RoleName:    role.RoleName,
			RoleCode:    role.RoleCode,
			Description: role.Description,
			Enabled:     role.Enabled,
			CreateTime:  role.CreateTime,
		})
	}

	return result, total, nil
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req *request.CreateRoleRequest) error {
	// 检查角色编码是否已存在
	existRole, _ := s.roleRepo.FindByCode(req.RoleCode)
	if existRole != nil {
		return errors.New("角色编码已存在")
	}

	// 创建角色
	role := &model.Role{
		RoleName:    req.RoleName,
		RoleCode:    req.RoleCode,
		Description: req.Description,
		Enabled:     req.Enabled,
		CreateTime:  time.Now(),
	}

	if err := s.roleRepo.Create(role); err != nil {
		return err
	}

	// 分配菜单权限
	if len(req.MenuIDs) > 0 {
		if err := s.roleRepo.AssignMenus(role.RoleID, req.MenuIDs); err != nil {
			return err
		}
	}

	// 分配按钮权限
	if len(req.ButtonIDs) > 0 {
		if err := s.roleRepo.AssignButtons(role.RoleID, req.ButtonIDs); err != nil {
			return err
		}
	}

	return nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(req *request.UpdateRoleRequest) error {
	// 查询角色是否存在
	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 不允许修改超级管理员
	if role.RoleID == 1 {
		return errors.New("不能修改超级管理员角色")
	}

	// 检查角色编码是否被其他角色使用
	if role.RoleCode != req.RoleCode {
		existRole, _ := s.roleRepo.FindByCode(req.RoleCode)
		if existRole != nil && existRole.RoleID != req.RoleID {
			return errors.New("角色编码已被使用")
		}
	}

	// 更新角色信息
	role.RoleName = req.RoleName
	role.RoleCode = req.RoleCode
	role.Description = req.Description
	role.Enabled = req.Enabled

	if err := s.roleRepo.Update(role); err != nil {
		return err
	}

	// 更新菜单权限
	if err := s.roleRepo.AssignMenus(role.RoleID, req.MenuIDs); err != nil {
		return err
	}

	// 更新按钮权限
	if err := s.roleRepo.AssignButtons(role.RoleID, req.ButtonIDs); err != nil {
		return err
	}

	return nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(roleID int64) error {
	// 检查角色是否存在
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 不允许删除超级管理员
	if roleID == 1 {
		return errors.New("不能删除超级管理员角色")
	}

	// TODO: 检查是否有用户关联此角色

	return s.roleRepo.Delete(roleID)
}

// GetRolePermissions 获取角色权限
func (s *RoleService) GetRolePermissions(roleID int64) (menuIDs []int64, buttonIDs []int64, err error) {
	menuIDs, err = s.roleRepo.GetRoleMenus(roleID)
	if err != nil {
		return nil, nil, err
	}

	buttonIDs, err = s.roleRepo.GetRoleButtons(roleID)
	if err != nil {
		return nil, nil, err
	}

	return menuIDs, buttonIDs, nil
}

// UpdateRolePermissions 更新角色权限
func (s *RoleService) UpdateRolePermissions(roleID int64, req *request.UpdateRolePermissionsRequest) error {
	// 检查角色是否存在
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 不允许修改超级管理员权限
	if roleID == 1 {
		return errors.New("不能修改超级管理员角色权限")
	}

	// 更新菜单权限
	if err := s.roleRepo.AssignMenus(roleID, req.MenuIDs); err != nil {
		return err
	}

	// 更新按钮权限
	if err := s.roleRepo.AssignButtons(roleID, req.ButtonIDs); err != nil {
		return err
	}

	return nil
}
