package repository

import (
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

// FindByID 根据ID查询角色
func (r *RoleRepository) FindByID(roleID int64) (*model.Role, error) {
	var role model.Role
	err := database.DB.Where("role_id = ?", roleID).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByCode 根据角色编码查询
func (r *RoleRepository) FindByCode(code string) (*model.Role, error) {
	var role model.Role
	err := database.DB.Where("role_code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Create 创建角色
func (r *RoleRepository) Create(role *model.Role) error {
	role.CreateTime = time.Now()
	role.UpdateTime = time.Now()
	return database.DB.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *model.Role) error {
	role.UpdateTime = time.Now()
	return database.DB.Save(role).Error
}

// Delete 删除角色（软删除）
func (r *RoleRepository) Delete(roleID int64) error {
	return database.DB.Where("role_id = ?", roleID).Delete(&model.Role{}).Error
}

// AssignMenus 分配菜单给角色
func (r *RoleRepository) AssignMenus(roleID int64, menuIDs []int64) error {
	tx := database.DB.Begin()
	
	// 删除现有菜单关联
	if err := tx.Exec("DELETE FROM sys_role_menu WHERE role_id = ?", roleID).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 添加新的菜单关联
	for _, menuID := range menuIDs {
		if err := tx.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 
			roleID, menuID, time.Now()).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit().Error
}

// AssignButtons 分配按钮权限给角色
func (r *RoleRepository) AssignButtons(roleID int64, buttonIDs []int64) error {
	tx := database.DB.Begin()
	
	// 删除现有按钮关联
	if err := tx.Exec("DELETE FROM sys_role_button WHERE role_id = ?", roleID).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 添加新的按钮关联
	for _, buttonID := range buttonIDs {
		if err := tx.Exec("INSERT INTO sys_role_button (role_id, button_id, create_time) VALUES (?, ?, ?)", 
			roleID, buttonID, time.Now()).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit().Error
}

// GetRoleMenus 获取角色的菜单列表
func (r *RoleRepository) GetRoleMenus(roleID int64) ([]int64, error) {
	var menuIDs []int64
	err := database.DB.Table("sys_role_menu").
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

// GetRoleButtons 获取角色的按钮权限列表
func (r *RoleRepository) GetRoleButtons(roleID int64) ([]int64, error) {
	var buttonIDs []int64
	err := database.DB.Table("sys_role_button").
		Where("role_id = ?", roleID).
		Pluck("button_id", &buttonIDs).Error
	return buttonIDs, err
}

// FindWithPagination 分页查询角色列表
func (r *RoleRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	db := database.DB.Model(&model.Role{})

	// 动态条件查询
	if roleID, ok := query["roleId"]; ok && roleID != nil {
		db = db.Where("role_id = ?", roleID)
	}
	if roleName, ok := query["roleName"]; ok && roleName != "" {
		db = db.Where("role_name LIKE ?", "%"+roleName.(string)+"%")
	}
	if roleCode, ok := query["roleCode"]; ok && roleCode != "" {
		db = db.Where("role_code LIKE ?", "%"+roleCode.(string)+"%")
	}
	if description, ok := query["description"]; ok && description != "" {
		db = db.Where("description LIKE ?", "%"+description.(string)+"%")
	}
	if enabled, ok := query["enabled"]; ok && enabled != nil {
		db = db.Where("enabled = ?", enabled)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err := db.Offset(offset).Limit(size).Order("role_id ASC").Find(&roles).Error

	return roles, total, err
}
