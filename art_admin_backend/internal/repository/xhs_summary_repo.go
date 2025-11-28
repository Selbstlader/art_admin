package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// XHSSummaryRepository 小红书总结仓库
type XHSSummaryRepository struct{}

// NewXHSSummaryRepository 创建仓库
func NewXHSSummaryRepository() *XHSSummaryRepository {
	return &XHSSummaryRepository{}
}

// Create 创建总结记录
func (r *XHSSummaryRepository) Create(summary *model.XHSSummary) error {
	db := database.GetDB()
	return db.Create(summary).Error
}

// Update 更新总结记录
func (r *XHSSummaryRepository) Update(summary *model.XHSSummary) error {
	db := database.GetDB()
	return db.Save(summary).Error
}

// Delete 删除总结记录
func (r *XHSSummaryRepository) Delete(id, userID int64) error {
	db := database.GetDB()
	return db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.XHSSummary{}).Error
}

// FindByID 根据ID查找
func (r *XHSSummaryRepository) FindByID(id int64) (*model.XHSSummary, error) {
	db := database.GetDB()
	var summary model.XHSSummary
	err := db.First(&summary, id).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// FindByUserID 查找用户的所有总结
func (r *XHSSummaryRepository) FindByUserID(userID int64, page, pageSize int) ([]model.XHSSummary, int64, error) {
	db := database.GetDB()
	var summaries []model.XHSSummary
	var total int64

	query := db.Model(&model.XHSSummary{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&summaries).Error; err != nil {
		return nil, 0, err
	}

	return summaries, total, nil
}

// FindFavorites 查找收藏的总结
func (r *XHSSummaryRepository) FindFavorites(userID int64, folderID int64, page, pageSize int) ([]model.XHSSummary, int64, error) {
	db := database.GetDB()
	var summaries []model.XHSSummary
	var total int64

	query := db.Model(&model.XHSSummary{}).Where("user_id = ? AND is_favorite = ?", userID, true)
	if folderID > 0 {
		query = query.Where("folder_id = ?", folderID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&summaries).Error; err != nil {
		return nil, 0, err
	}

	return summaries, total, nil
}

// UpdateFavorite 更新收藏状态
func (r *XHSSummaryRepository) UpdateFavorite(id, userID int64, isFavorite bool, folderID int64) error {
	db := database.GetDB()
	return db.Model(&model.XHSSummary{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"is_favorite": isFavorite,
			"folder_id":   folderID,
		}).Error
}

// Search 搜索总结
func (r *XHSSummaryRepository) Search(userID int64, keyword string, page, pageSize int) ([]model.XHSSummary, int64, error) {
	db := database.GetDB()
	var summaries []model.XHSSummary
	var total int64

	query := db.Model(&model.XHSSummary{}).Where("user_id = ?", userID)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		query = query.Where("note_title LIKE ? OR summary_title LIKE ? OR summary_content LIKE ?",
			likeKeyword, likeKeyword, likeKeyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&summaries).Error; err != nil {
		return nil, 0, err
	}

	return summaries, total, nil
}

// XHSFolderRepository 收藏夹仓库
type XHSFolderRepository struct{}

// NewXHSFolderRepository 创建收藏夹仓库
func NewXHSFolderRepository() *XHSFolderRepository {
	return &XHSFolderRepository{}
}

// Create 创建收藏夹
func (r *XHSFolderRepository) Create(folder *model.XHSFolder) error {
	db := database.GetDB()
	return db.Create(folder).Error
}

// Update 更新收藏夹
func (r *XHSFolderRepository) Update(folder *model.XHSFolder) error {
	db := database.GetDB()
	return db.Save(folder).Error
}

// Delete 删除收藏夹
func (r *XHSFolderRepository) Delete(id, userID int64) error {
	db := database.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		// 将该收藏夹下的总结移出收藏夹
		if err := tx.Model(&model.XHSSummary{}).
			Where("folder_id = ? AND user_id = ?", id, userID).
			Update("folder_id", 0).Error; err != nil {
			return err
		}
		// 删除收藏夹
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.XHSFolder{}).Error
	})
}

// FindByUserID 查找用户的所有收藏夹
func (r *XHSFolderRepository) FindByUserID(userID int64) ([]model.XHSFolder, error) {
	db := database.GetDB()
	var folders []model.XHSFolder
	err := db.Where("user_id = ?", userID).Order("sort_order ASC, created_at DESC").Find(&folders).Error
	return folders, err
}

// FindByID 根据ID查找收藏夹
func (r *XHSFolderRepository) FindByID(id int64) (*model.XHSFolder, error) {
	db := database.GetDB()
	var folder model.XHSFolder
	err := db.First(&folder, id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}
