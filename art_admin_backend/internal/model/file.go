package model

import (
	"time"

	"gorm.io/gorm"
)

// File 文件模型
type File struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	FileName     string         `gorm:"size:255;not null" json:"fileName"`         // 原始文件名
	StorageName  string         `gorm:"size:255;not null" json:"storageName"`      // 存储文件名
	FilePath     string         `gorm:"size:500;not null" json:"filePath"`         // 文件路径
	FileSize     int64          `gorm:"not null" json:"fileSize"`                  // 文件大小(字节)
	FileType     string         `gorm:"size:100" json:"fileType"`                  // 文件类型(MIME)
	FileExt      string         `gorm:"size:20" json:"fileExt"`                    // 文件扩展名
	Category     string         `gorm:"size:50;default:'default'" json:"category"` // 分类(image/document/video/audio/other)
	UploadUserID uint           `gorm:"index" json:"uploadUserId"`                 // 上传用户ID
	URL          string         `gorm:"size:500" json:"url"`                       // 访问URL
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (File) TableName() string {
	return "files"
}
