package repository

import (
	"context"

	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// SubjectRepository 学科仓储
type SubjectRepository struct {
	db *gorm.DB
}

// NewSubjectRepository 创建学科仓储
func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

// GetByID 根据 ID 获取学科
func (r *SubjectRepository) GetByID(ctx context.Context, id int64) (*model.Subject, error) {
	var subject model.Subject
	if err := r.db.WithContext(ctx).Where("id = ? AND status = 1", id).First(&subject).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

// List 获取学科列表
func (r *SubjectRepository) List(ctx context.Context) ([]*model.Subject, error) {
	var subjects []*model.Subject
	if err := r.db.WithContext(ctx).
		Where("status = 1").
		Order("sort ASC, id ASC").
		Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

// Create 创建学科
func (r *SubjectRepository) Create(ctx context.Context, subject *model.Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

// Update 更新学科
func (r *SubjectRepository) Update(ctx context.Context, subject *model.Subject) error {
	return r.db.WithContext(ctx).Save(subject).Error
}

// Delete 删除学科（软删除）
func (r *SubjectRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.Subject{}).
		Where("id = ?", id).
		Update("status", 0).Error
}
