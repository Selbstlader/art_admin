package repository

import (
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
)

type AppUserRepository struct{}

func NewAppUserRepository() *AppUserRepository {
	return &AppUserRepository{}
}

// FindByUserName 根据用户名查询APP用户
func (r *AppUserRepository) FindByUserName(userName string) (*model.AppUser, error) {
	var user model.AppUser
	err := database.DB.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone 根据手机号查询APP用户
func (r *AppUserRepository) FindByPhone(phone string) (*model.AppUser, error) {
	var user model.AppUser
	err := database.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查询APP用户
func (r *AppUserRepository) FindByID(id int64) (*model.AppUser, error) {
	var user model.AppUser
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建APP用户
func (r *AppUserRepository) Create(user *model.AppUser) error {
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	return database.DB.Create(user).Error
}

// Update 更新APP用户
func (r *AppUserRepository) Update(user *model.AppUser) error {
	user.UpdateTime = time.Now()
	return database.DB.Save(user).Error
}

// Delete 删除APP用户（软删除）
func (r *AppUserRepository) Delete(id int64) error {
	return database.DB.Delete(&model.AppUser{}, id).Error
}

// FindWithPagination 分页查询APP用户列表
func (r *AppUserRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.AppUser, int64, error) {
	var users []model.AppUser
	var total int64

	db := database.DB.Model(&model.AppUser{})

	// 动态条件查询
	if id, ok := query["id"]; ok && id != nil {
		db = db.Where("id = ?", id)
	}
	if userName, ok := query["userName"]; ok && userName != "" {
		db = db.Where("user_name LIKE ?", "%"+userName.(string)+"%")
	}
	if nickName, ok := query["nickName"]; ok && nickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+nickName.(string)+"%")
	}
	if phone, ok := query["phone"]; ok && phone != "" {
		db = db.Where("phone LIKE ?", "%"+phone.(string)+"%")
	}
	if userType, ok := query["userType"]; ok && userType != "" {
		db = db.Where("user_type = ?", userType)
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

// UpdatePassword 更新APP用户密码
func (r *AppUserRepository) UpdatePassword(id int64, password string) error {
	return database.DB.Model(&model.AppUser{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password":    password,
			"update_time": time.Now(),
		}).Error
}

// UpdateLastLogin 更新最后登录时间和IP
func (r *AppUserRepository) UpdateLastLogin(id int64, ip string) error {
	now := time.Now()
	return database.DB.Model(&model.AppUser{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_time": now,
			"last_login_ip":   ip,
			"update_time":     now,
		}).Error
}
