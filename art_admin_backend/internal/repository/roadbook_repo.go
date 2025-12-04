package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// RoadbookRepository 路书数据访问层
type RoadbookRepository struct {
	db *gorm.DB
}

// NewRoadbookRepository 创建路书仓库实例
func NewRoadbookRepository() *RoadbookRepository {
	return &RoadbookRepository{
		db: database.DB,
	}
}

// NewRoadbookRepositoryWithDB 创建带自定义DB的路书仓库实例（用于测试）
func NewRoadbookRepositoryWithDB(db *gorm.DB) *RoadbookRepository {
	return &RoadbookRepository{
		db: db,
	}
}

// Create 创建路书
func (r *RoadbookRepository) Create(roadbook *model.Roadbook) error {
	return r.getDB().Create(roadbook).Error
}

// GetByID 根据ID获取路书（不包含已删除的）
func (r *RoadbookRepository) GetByID(id int64) (*model.Roadbook, error) {
	var roadbook model.Roadbook
	if err := r.getDB().First(&roadbook, id).Error; err != nil {
		return nil, err
	}
	return &roadbook, nil
}

// Update 更新路书
func (r *RoadbookRepository) UpdateRoadbook(roadbook *model.Roadbook) error {
	return r.getDB().Save(roadbook).Error
}

// List 分页查询路书列表（按创建时间倒序）
func (r *RoadbookRepository) List(page, pageSize int) ([]model.Roadbook, int64, error) {
	var roadbooks []model.Roadbook
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

	db := r.getDB().Model(&model.Roadbook{})

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&roadbooks).Error; err != nil {
		return nil, 0, err
	}

	return roadbooks, total, nil
}

// ListByUserID 根据用户ID分页查询路书列表
func (r *RoadbookRepository) ListByUserID(userID int64, page, pageSize int) ([]model.Roadbook, int64, error) {
	var roadbooks []model.Roadbook
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

	db := r.getDB().Model(&model.Roadbook{}).Where("user_id = ?", userID)

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&roadbooks).Error; err != nil {
		return nil, 0, err
	}

	return roadbooks, total, nil
}

// SoftDelete 软删除路书
func (r *RoadbookRepository) SoftDelete(id int64) error {
	return r.getDB().Delete(&model.Roadbook{}, id).Error
}

// ExistsByID 检查路书是否存在（不包含已删除的）
func (r *RoadbookRepository) ExistsByID(id int64) (bool, error) {
	var count int64
	if err := r.getDB().Model(&model.Roadbook{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// getDB 获取数据库连接
func (r *RoadbookRepository) getDB() *gorm.DB {
	if r.db != nil {
		return r.db
	}
	return database.DB
}

// FindWithPagination 分页查询路书列表
func (r *RoadbookRepository) FindWithPagination(page, pageSize int, userID int64, keyword string, tagID int64, visibility *int, status *int) ([]model.Roadbook, int64, error) {
	var roadbooks []model.Roadbook
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

	db := r.getDB().Model(&model.Roadbook{})

	// 用户只能看到自己的路书或公开的路书
	db = db.Where("user_id = ? OR visibility = ?", userID, model.VisibilityPublic)

	// 条件查询
	if keyword != "" {
		db = db.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if visibility != nil {
		db = db.Where("visibility = ?", *visibility)
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	if tagID > 0 {
		db = db.Joins("JOIN roadbook_tag_relations ON roadbook_tag_relations.roadbook_id = roadbooks.id").
			Where("roadbook_tag_relations.tag_id = ?", tagID)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&roadbooks).Error; err != nil {
		return nil, 0, err
	}

	return roadbooks, total, nil
}

// FindByID 根据ID查询路书（包含途经点）
func (r *RoadbookRepository) FindByID(id int64) (*model.Roadbook, error) {
	var roadbook model.Roadbook
	if err := r.getDB().Preload("Waypoints").First(&roadbook, id).Error; err != nil {
		return nil, err
	}
	return &roadbook, nil
}

// CreateWithWaypoints 创建路书及途经点（事务）
func (r *RoadbookRepository) CreateWithWaypoints(roadbook *model.Roadbook, waypoints []model.RoadbookWaypoint, tagIDs []int64) error {
	tx := r.getDB().Begin()

	// 创建路书
	if err := tx.Create(roadbook).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建途经点
	for i := range waypoints {
		waypoints[i].RoadbookID = roadbook.ID
		if err := tx.Create(&waypoints[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 创建标签关联
	for _, tagID := range tagIDs {
		relation := &model.RoadbookTagRelation{
			RoadbookID: roadbook.ID,
			TagID:      tagID,
		}
		if err := tx.Create(relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 计算总距离和时间
	if len(waypoints) > 0 {
		var totalDistance, totalDuration int
		for _, wp := range waypoints {
			totalDuration += wp.StayDuration
		}
		roadbook.TotalDuration = totalDuration
		roadbook.TotalDistance = totalDistance
		if err := tx.Save(roadbook).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Update 更新路书（带标签）
func (r *RoadbookRepository) Update(roadbook *model.Roadbook, tagIDs []int64) error {
	tx := r.getDB().Begin()

	if err := tx.Save(roadbook).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新标签关联
	if tagIDs != nil {
		// 删除旧的标签关联
		if err := tx.Where("roadbook_id = ?", roadbook.ID).Delete(&model.RoadbookTagRelation{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		// 创建新的标签关联
		for _, tagID := range tagIDs {
			relation := &model.RoadbookTagRelation{
				RoadbookID: roadbook.ID,
				TagID:      tagID,
			}
			if err := tx.Create(relation).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// Delete 软删除路书
func (r *RoadbookRepository) Delete(id int64) error {
	return r.getDB().Delete(&model.Roadbook{}, id).Error
}

// IncrementViewCount 增加浏览量
func (r *RoadbookRepository) IncrementViewCount(id int64) error {
	return r.getDB().Model(&model.Roadbook{}).Where("id = ?", id).
		UpdateColumn("view_count", r.getDB().Raw("view_count + 1")).Error
}

// IncrementFavoriteCount 增加收藏量
func (r *RoadbookRepository) IncrementFavoriteCount(id int64, delta int) error {
	return r.getDB().Model(&model.Roadbook{}).Where("id = ?", id).
		UpdateColumn("favorite_count", r.getDB().Raw("favorite_count + ?", delta)).Error
}

// IncrementCommentCount 增加评论量
func (r *RoadbookRepository) IncrementCommentCount(id int64, delta int) error {
	return r.getDB().Model(&model.Roadbook{}).Where("id = ?", id).
		UpdateColumn("comment_count", r.getDB().Raw("comment_count + ?", delta)).Error
}

// IncrementShareCount 增加分享量
func (r *RoadbookRepository) IncrementShareCount(id int64) error {
	return r.getDB().Model(&model.Roadbook{}).Where("id = ?", id).
		UpdateColumn("share_count", r.getDB().Raw("share_count + 1")).Error
}
