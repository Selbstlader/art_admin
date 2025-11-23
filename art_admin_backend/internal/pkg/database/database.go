package database

import (
	"fmt"
	"log"
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := cfg.GetDSN()

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
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
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	log.Println("开始自动迁移数据库表结构...")

	// 自动迁移所有模型
	err := DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Button{},
		&model.DictionaryType{},
		&model.Dictionary{},
		&model.OperationLog{},
		// 项目管理模块
		&model.ProjectTemplate{},
		&model.TemplateTask{},
		&model.TemplateTaskDependency{},
		&model.Project{},
		&model.ProjectMember{},
		&model.Task{},
		&model.TaskDependency{},
		&model.TaskComment{},
		&model.TaskAttachment{},
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
