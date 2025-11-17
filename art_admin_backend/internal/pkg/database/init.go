package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"art_admin_backend/internal/pkg/config"
)

// AutoCreateDatabase 自动创建数据库（如果不存在）
func AutoCreateDatabase(cfg *config.DatabaseConfig) error {
	// 连接到MySQL服务器（不指定数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %w", err)
	}

	// 创建数据库
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET %s COLLATE %s",
		cfg.Database,
		cfg.Charset,
		"utf8mb4_general_ci",
	)

	if err := db.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("创建数据库失败: %w", err)
	}

	log.Printf("数据库 %s 已就绪", cfg.Database)
	return nil
}

// InitializeDatabase 完整初始化数据库（创建库→连接→迁移→初始化数据）
func InitializeDatabase(cfg *config.DatabaseConfig) error {
	// 1. 自动创建数据库
	if err := AutoCreateDatabase(cfg); err != nil {
		return err
	}

	// 2. 连接数据库
	if err := InitDB(cfg); err != nil {
		return err
	}

	// 3. 自动迁移表结构
	if err := AutoMigrate(); err != nil {
		return err
	}

	// 4. 初始化种子数据
	if err := SeedData(); err != nil {
		log.Printf("警告：初始化数据失败: %v", err)
		// 不返回错误，因为数据可能已经存在
	}

	return nil
}

