package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// CommentRepository 评论数据访问层
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库实例
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		db: database.DB,
	}
}

// NewCommentRepositoryWithDB 创建带自定义DB的评论仓库实例（用于测试）
func NewCommentRepositoryWithDB(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

// getDB 获取数据库连接
func (r *CommentRepository) getDB() *gorm.DB {
	if r.db != nil {
		return r.db
	}
	return database.DB
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.RoadbookComment) error {
	return r.getDB().Create(comment).Error
}

// GetByID 根据ID获取评论
func (r *CommentRepository) GetByID(id int64) (*model.RoadbookComment, error) {
	var comment model.RoadbookComment
	if err := r.getDB().First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Update 更新评论
func (r *CommentRepository) Update(comment *model.RoadbookComment) error {
	return r.getDB().Save(comment).Error
}

// SoftDelete 软删除评论
func (r *CommentRepository) SoftDelete(id int64) error {
	return r.getDB().Delete(&model.RoadbookComment{}, id).Error
}

// ListByRoadbookID 根据路书ID获取评论列表（按时间倒序，分页）
func (r *CommentRepository) ListByRoadbookID(roadbookID int64, page, pageSize int) ([]model.RoadbookComment, int64, error) {
	var comments []model.RoadbookComment
	var total int64

	// 限制最大pageSize为100
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if page < 1 {
		page = 1
	}

	db := r.getDB().Model(&model.RoadbookComment{}).Where("roadbook_id = ?", roadbookID)

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListTopLevelByRoadbookID 获取路书的顶级评论（parentId=0）
func (r *CommentRepository) ListTopLevelByRoadbookID(roadbookID int64, page, pageSize int) ([]model.RoadbookComment, int64, error) {
	var comments []model.RoadbookComment
	var total int64

	// 限制最大pageSize为100
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if page < 1 {
		page = 1
	}

	db := r.getDB().Model(&model.RoadbookComment{}).Where("roadbook_id = ? AND parent_id = 0", roadbookID)

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListRepliesByParentID 获取某评论的回复列表
func (r *CommentRepository) ListRepliesByParentID(parentID int64) ([]model.RoadbookComment, error) {
	var replies []model.RoadbookComment
	if err := r.getDB().Where("parent_id = ?", parentID).Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}

// ExistsByID 检查评论是否存在
func (r *CommentRepository) ExistsByID(id int64) (bool, error) {
	var count int64
	if err := r.getDB().Model(&model.RoadbookComment{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByIDAndRoadbookID 检查评论是否存在且属于指定路书
func (r *CommentRepository) ExistsByIDAndRoadbookID(id, roadbookID int64) (bool, error) {
	var count int64
	if err := r.getDB().Model(&model.RoadbookComment{}).Where("id = ? AND roadbook_id = ?", id, roadbookID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IncrementLikeCount 增加点赞数
func (r *CommentRepository) IncrementLikeCount(id int64, delta int) error {
	return r.getDB().Model(&model.RoadbookComment{}).Where("id = ?", id).
		UpdateColumn("like_count", r.getDB().Raw("like_count + ?", delta)).Error
}

// GetCommentTree 获取评论树形结构
func (r *CommentRepository) GetCommentTree(roadbookID int64, page, pageSize int) ([]model.RoadbookComment, int64, error) {
	// 先获取顶级评论
	topComments, total, err := r.ListTopLevelByRoadbookID(roadbookID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 获取所有顶级评论的回复
	for i := range topComments {
		replies, err := r.ListRepliesByParentID(topComments[i].ID)
		if err != nil {
			return nil, 0, err
		}
		topComments[i].Replies = replies
	}

	return topComments, total, nil
}
