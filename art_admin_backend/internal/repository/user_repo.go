package repository

import (
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByUserName 根据用户名查询用户
func (r *UserRepository) FindByUserName(userName string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查询用户
func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	return database.DB.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	user.UpdateTime = time.Now()
	return database.DB.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(id int64) error {
	return database.DB.Delete(&model.User{}, id).Error
}

// AssignRoles 分配角色给用户
func (r *UserRepository) AssignRoles(userID int64, roleIDs []int64) error {
	// 开启事务
	tx := database.DB.Begin()
	
	// 删除现有角色关联
	if err := tx.Exec("DELETE FROM sys_user_role WHERE user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 添加新的角色关联
	for _, roleID := range roleIDs {
		if err := tx.Exec("INSERT INTO sys_user_role (user_id, role_id, create_time) VALUES (?, ?, ?)", 
			userID, roleID, time.Now()).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit().Error
}

// GetUserRoles 获取用户角色
func (r *UserRepository) GetUserRoles(userID int64) ([]model.Role, error) {
	var roles []model.Role
	err := database.DB.
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// GetUserButtons 获取用户按钮权限
func (r *UserRepository) GetUserButtons(userID int64) ([]model.Button, error) {
	var buttons []model.Button
	err := database.DB.
		Joins("JOIN sys_role_button ON sys_role_button.button_id = sys_button.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_button.role_id").
		Where("sys_user_role.user_id = ? AND sys_button.enabled = ?", userID, true).
		Distinct().
		Find(&buttons).Error
	return buttons, err
}

// FindWithPagination 分页查询用户列表
func (r *UserRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := database.DB.Model(&model.User{})

	// 动态条件查询
	if id, ok := query["id"]; ok && id != nil {
		db = db.Where("id = ?", id)
	}
	if userName, ok := query["userName"]; ok && userName != "" {
		db = db.Where("user_name LIKE ?", "%"+userName.(string)+"%")
	}
	if userGender, ok := query["userGender"]; ok && userGender != "" {
		db = db.Where("user_gender = ?", userGender)
	}
	if userPhone, ok := query["userPhone"]; ok && userPhone != "" {
		db = db.Where("user_phone LIKE ?", "%"+userPhone.(string)+"%")
	}
	if userEmail, ok := query["userEmail"]; ok && userEmail != "" {
		db = db.Where("email LIKE ?", "%"+userEmail.(string)+"%")
	}
	if status, ok := query["status"]; ok && status != "" {
		db = db.Where("status = ?", status)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err := db.Offset(offset).Limit(size).Order("id DESC").Find(&users).Error

	return users, total, err
}

// GetUserRolesByUserID 批量获取用户角色
func (r *UserRepository) GetUserRolesByUserID(userID int64) ([]string, error) {
	var roleCodes []string
	err := database.DB.
		Table("sys_role").
		Select("sys_role.role_code").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.role_id").
		Where("sys_user_role.user_id = ?", userID).
		Pluck("role_code", &roleCodes).Error
	return roleCodes, err
}

// UpdatePassword 更新用户密码
func (r *UserRepository) UpdatePassword(id int64, password string) error {
	return database.DB.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password":    password,
			"update_time": time.Now(),
		}).Error
}
