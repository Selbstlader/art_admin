package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"errors"

	"gorm.io/gorm"
)

type DepartmentRepository struct{}

func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{}
}

// FindByID 根据ID查询部门
func (r *DepartmentRepository) FindByID(deptID int64) (*model.Department, error) {
	var dept model.Department
	err := database.DB.Where("dept_id = ?", deptID).First(&dept).Error
	if err != nil {
		return nil, err // 找不到记录时返回错误，因为调用方需要区分记录不存在和数据库错误
	}
	return &dept, nil
}

// FindByDeptCode 根据部门编码查询
func (r *DepartmentRepository) FindByDeptCode(deptCode string) (*model.Department, error) {
	var dept model.Department
	err := database.DB.Where("dept_code = ?", deptCode).First(&dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 记录不存在是正常情况
		}
		return nil, err
	}
	return &dept, nil
}

// FindAll 查询所有部门
func (r *DepartmentRepository) FindAll(query map[string]interface{}) ([]model.Department, error) {
	var departments []model.Department
	db := database.DB

	// 应用查询条件
	for key, value := range query {
		switch key {
		case "deptName":
			if name, ok := value.(string); ok && name != "" {
				db = db.Where("dept_name LIKE ?", "%"+name+"%")
			}
		case "deptCode":
			if code, ok := value.(string); ok && code != "" {
				db = db.Where("dept_code LIKE ?", "%"+code+"%")
			}
		case "status":
			if status, ok := value.(int); ok {
				db = db.Where("status = ?", status)
			}
		}
	}

	err := db.Order("order_num ASC, dept_id ASC").Find(&departments).Error
	return departments, err
}

// FindByParentID 根据父部门ID查询子部门
func (r *DepartmentRepository) FindByParentID(parentID int64) ([]model.Department, error) {
	var departments []model.Department
	err := database.DB.Where("parent_id = ?", parentID).Order("order_num ASC, dept_id ASC").Find(&departments).Error
	return departments, err
}

// Create 创建部门
func (r *DepartmentRepository) Create(department *model.Department) error {
	return database.DB.Create(department).Error
}

// Update 更新部门
func (r *DepartmentRepository) Update(department *model.Department) error {
	return database.DB.Save(department).Error
}

// Delete 删除部门
func (r *DepartmentRepository) Delete(deptID int64) error {
	return database.DB.Delete(&model.Department{}, deptID).Error
}

// FindUsersByDeptID 根据部门ID查询用户
func (r *DepartmentRepository) FindUsersByDeptID(deptID int64) ([]model.User, error) {
	var users []model.User
	err := database.DB.Where("dept_id = ?", deptID).Find(&users).Error
	return users, err
}

// CountByParentID 统计指定父部门下的子部门数量
func (r *DepartmentRepository) CountByParentID(parentID int64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Department{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}

// ExistsByDeptCode 检查部门编码是否存在（排除指定ID）
func (r *DepartmentRepository) ExistsByDeptCode(deptCode string, excludeID int64) (bool, error) {
	var count int64
	query := database.DB.Model(&model.Department{}).Where("dept_code = ?", deptCode)
	if excludeID > 0 {
		query = query.Where("dept_id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
