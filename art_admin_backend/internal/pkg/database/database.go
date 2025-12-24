package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接，支持 MySQL 和 PostgreSQL
// Initialize database connection, supports MySQL and PostgreSQL
func InitDB(cfg *config.DatabaseConfig) error {
	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	var db *gorm.DB
	var err error

	// 优先使用 DATABASE_URL 环境变量（Render 等云平台）
	// Prefer DATABASE_URL env var for cloud platforms like Render
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Println("使用 DATABASE_URL 连接 PostgreSQL...")
		db, err = gorm.Open(postgres.Open(databaseURL), gormConfig)
	} else if cfg.Driver == "postgres" {
		// PostgreSQL 连接
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	} else {
		// MySQL 连接（默认）
		dsn := cfg.GetDSN()
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	}

	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	DB = db
	log.Println("数据库连接成功")
	return nil
}

// AutoMigrate 自动迁移数据库表结构
// 基础功能 + 系统功能 + Dify AI + 聊天室
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	log.Println("开始自动迁移数据库表结构...")

	// 自动迁移所有模型 - 基础功能 + 系统功能 + 聊天室
	err := DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Button{},
		&model.DictionaryType{},
		&model.Dictionary{},
		&model.OperationLog{},
		&model.Department{},
		&model.File{},
		&model.AppUser{},
		// 聊天室模块
		&model.ChatRoom{},
		&model.ChatMessage{},
		&model.ChatRoomMember{},
	)

	if err != nil {
		return fmt.Errorf("数据库表迁移失败: %w", err)
	}

	log.Println("数据库表迁移完成!")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
