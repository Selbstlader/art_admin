package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// WaypointRepository 途经点数据访问层
type WaypointRepository struct {
	db *gorm.DB
}

// NewWaypointRepository 创建途经点仓库实例
func NewWaypointRepository() *WaypointRepository {
	return &WaypointRepository{
		db: database.DB,
	}
}

// NewWaypointRepositoryWithDB 创建带自定义DB的途经点仓库实例（用于测试）
func NewWaypointRepositoryWithDB(db *gorm.DB) *WaypointRepository {
	return &WaypointRepository{
		db: db,
	}
}

// getDB 获取数据库连接
func (r *WaypointRepository) getDB() *gorm.DB {
	if r.db != nil {
		return r.db
	}
	return database.DB
}

// Create 创建途经点
func (r *WaypointRepository) Create(waypoint *model.RoadbookWaypoint) error {
	return r.getDB().Create(waypoint).Error
}

// GetByID 根据ID获取途经点
func (r *WaypointRepository) GetByID(id int64) (*model.RoadbookWaypoint, error) {
	var waypoint model.RoadbookWaypoint
	if err := r.getDB().First(&waypoint, id).Error; err != nil {
		return nil, err
	}
	return &waypoint, nil
}

// Update 更新途经点
func (r *WaypointRepository) Update(waypoint *model.RoadbookWaypoint) error {
	return r.getDB().Save(waypoint).Error
}

// Delete 删除途经点
func (r *WaypointRepository) Delete(id int64) error {
	return r.getDB().Delete(&model.RoadbookWaypoint{}, id).Error
}

// ListByRoadbookID 根据路书ID获取所有途经点（按排序顺序）
func (r *WaypointRepository) ListByRoadbookID(roadbookID int64) ([]model.RoadbookWaypoint, error) {
	var waypoints []model.RoadbookWaypoint
	if err := r.getDB().Where("roadbook_id = ?", roadbookID).
		Order("day_index ASC, sort_order ASC").
		Find(&waypoints).Error; err != nil {
		return nil, err
	}
	return waypoints, nil
}

// ListByRoadbookIDGroupedByDay 根据路书ID获取途经点并按天分组
func (r *WaypointRepository) ListByRoadbookIDGroupedByDay(roadbookID int64) (map[int][]model.RoadbookWaypoint, error) {
	waypoints, err := r.ListByRoadbookID(roadbookID)
	if err != nil {
		return nil, err
	}

	grouped := make(map[int][]model.RoadbookWaypoint)
	for _, wp := range waypoints {
		grouped[wp.DayIndex] = append(grouped[wp.DayIndex], wp)
	}
	return grouped, nil
}

// BatchCreate 批量创建途经点
func (r *WaypointRepository) BatchCreate(waypoints []model.RoadbookWaypoint) error {
	if len(waypoints) == 0 {
		return nil
	}
	return r.getDB().Create(&waypoints).Error
}

// BatchUpdate 批量更新途经点（事务）
func (r *WaypointRepository) BatchUpdate(waypoints []model.RoadbookWaypoint) error {
	if len(waypoints) == 0 {
		return nil
	}

	tx := r.getDB().Begin()
	for _, wp := range waypoints {
		if err := tx.Save(&wp).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// UpdateSortOrder 更新途经点排序
func (r *WaypointRepository) UpdateSortOrder(id int64, sortOrder int) error {
	return r.getDB().Model(&model.RoadbookWaypoint{}).
		Where("id = ?", id).
		Update("sort_order", sortOrder).Error
}

// BatchUpdateSortOrder 批量更新途经点排序（事务）
func (r *WaypointRepository) BatchUpdateSortOrder(updates map[int64]int) error {
	if len(updates) == 0 {
		return nil
	}

	tx := r.getDB().Begin()
	for id, sortOrder := range updates {
		if err := tx.Model(&model.RoadbookWaypoint{}).
			Where("id = ?", id).
			Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// DeleteByRoadbookID 删除路书的所有途经点
func (r *WaypointRepository) DeleteByRoadbookID(roadbookID int64) error {
	return r.getDB().Where("roadbook_id = ?", roadbookID).Delete(&model.RoadbookWaypoint{}).Error
}

// CountByRoadbookID 统计路书的途经点数量
func (r *WaypointRepository) CountByRoadbookID(roadbookID int64) (int64, error) {
	var count int64
	if err := r.getDB().Model(&model.RoadbookWaypoint{}).
		Where("roadbook_id = ?", roadbookID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByID 检查途经点是否存在
func (r *WaypointRepository) ExistsByID(id int64) (bool, error) {
	var count int64
	if err := r.getDB().Model(&model.RoadbookWaypoint{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetMaxSortOrder 获取路书中某天的最大排序值
func (r *WaypointRepository) GetMaxSortOrder(roadbookID int64, dayIndex int) (int, error) {
	var maxOrder int
	err := r.getDB().Model(&model.RoadbookWaypoint{}).
		Where("roadbook_id = ? AND day_index = ?", roadbookID, dayIndex).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}

// ReplaceWaypoints 替换路书的所有途经点（事务）
func (r *WaypointRepository) ReplaceWaypoints(roadbookID int64, waypoints []model.RoadbookWaypoint) error {
	tx := r.getDB().Begin()

	// 删除旧的途经点
	if err := tx.Where("roadbook_id = ?", roadbookID).Delete(&model.RoadbookWaypoint{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建新的途经点
	for i := range waypoints {
		waypoints[i].RoadbookID = roadbookID
		if err := tx.Create(&waypoints[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
