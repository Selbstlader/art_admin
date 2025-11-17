package repository

import (
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
)

type MenuRepository struct{}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{}
}

// FindByID 根据ID查询菜单
func (r *MenuRepository) FindByID(id int64) (*model.Menu, error) {
	var menu model.Menu
	err := database.DB.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// FindAll 查询所有菜单
func (r *MenuRepository) FindAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *model.Menu) error {
	menu.CreateTime = time.Now()
	menu.UpdateTime = time.Now()
	return database.DB.Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *model.Menu) error {
	menu.UpdateTime = time.Now()
	return database.DB.Save(menu).Error
}

// Delete 删除菜单（软删除）
func (r *MenuRepository) Delete(id int64) error {
	// 先检查是否有子菜单
	var count int64
	if err := database.DB.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return database.DB.Error  // 可以返回自定义错误
	}
	
	return database.DB.Delete(&model.Menu{}, id).Error
}

// HasChildren 检查是否有子菜单
func (r *MenuRepository) HasChildren(id int64) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetMenusByUserID 根据用户ID获取菜单列表
func (r *MenuRepository) GetMenusByUserID(userID int64) ([]model.Menu, error) {
	var menus []model.Menu

	err := database.DB.
		Distinct().
		Select("sys_menu.*").
		Table("sys_menu").
		Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_menu.role_id").
		Where("sys_user_role.user_id = ? AND sys_menu.is_enable = ?", userID, true).
		Order("sys_menu.sort ASC").
		Find(&menus).Error

	return menus, err
}

// GetButtonsByMenuID 根据菜单ID获取按钮权限
func (r *MenuRepository) GetButtonsByMenuID(menuID int64, userID int64) ([]model.Button, error) {
	var buttons []model.Button

	err := database.DB.
		Distinct().
		Select("sys_button.*").
		Table("sys_button").
		Joins("JOIN sys_role_button ON sys_role_button.button_id = sys_button.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_button.role_id").
		Where("sys_button.menu_id = ? AND sys_user_role.user_id = ? AND sys_button.is_enable = ?", menuID, userID, true).
		Find(&buttons).Error

	return buttons, err
}
