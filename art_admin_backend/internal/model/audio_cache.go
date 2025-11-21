package model

import "time"

// AudioCache 音频缓存表
type AudioCache struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ContentHash  string    `gorm:"type:varchar(64);not null;uniqueIndex;comment:内容MD5哈希" json:"content_hash"`
	TextContent  string    `gorm:"type:text;comment:原始文本内容" json:"text_content"`
	AudioURL     string    `gorm:"type:varchar(500);not null;comment:音频URL" json:"audio_url"`
	VoiceType    string    `gorm:"type:varchar(50);default:xiaoyun;comment:音色类型" json:"voice_type"`
	Language     string    `gorm:"type:varchar(20);default:zh-CN;comment:语言" json:"language"`
	FileSize     int64     `gorm:"comment:文件大小(字节)" json:"file_size"`
	Duration     int       `gorm:"comment:时长(秒)" json:"duration"`
	AccessCount  int       `gorm:"default:0;comment:访问次数" json:"access_count"`
	LastAccessAt time.Time `gorm:"index;comment:最后访问时间" json:"last_access_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (AudioCache) TableName() string {
	return "audio_cache"
}
