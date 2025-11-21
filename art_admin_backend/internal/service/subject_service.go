package service

import (
	"context"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
)

// SubjectService 学科服务
type SubjectService struct {
	repo *repository.SubjectRepository
}

// NewSubjectService 创建学科服务
func NewSubjectService(repo *repository.SubjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

// ListSubjects 获取学科列表
func (s *SubjectService) ListSubjects(ctx context.Context) ([]*model.Subject, error) {
	return s.repo.List(ctx)
}

// GetSubjectByID 根据 ID 获取学科
func (s *SubjectService) GetSubjectByID(ctx context.Context, id int64) (*model.Subject, error) {
	return s.repo.GetByID(ctx, id)
}
