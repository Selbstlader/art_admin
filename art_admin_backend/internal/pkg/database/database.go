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
		// 禁用外键约束，避免迁移时的依赖问题
		DisableForeignKeyConstraintWhenMigrating: true,
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
// 完整迁移所有模型，按依赖顺序执行
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	log.Println("开始自动迁移数据库表结构...")

	// 清理可能存在问题的约束（MySQL 9.x 兼容性处理）
	// Clean up potentially problematic constraints for MySQL 9.x compatibility
	cleanupConstraints()

	// ========== 第一批：基础系统表（无外键依赖）==========
	log.Println("迁移第一批：基础系统表...")
	if err := DB.AutoMigrate(
		&model.User{},           // sys_user - 用户表
		&model.Role{},           // sys_role - 角色表
		&model.Menu{},           // sys_menu - 菜单表
		&model.Button{},         // sys_button - 按钮权限表
		&model.Department{},     // sys_department - 部门表
		&model.DictionaryType{}, // sys_dictionary_type - 字典类型表
		&model.Dictionary{},     // sys_dictionary - 字典数据表
		&model.OperationLog{},   // sys_operation_log - 操作日志表
		&model.File{},           // files - 文件管理表
		&model.AppUser{},        // app_user - APP用户表
	); err != nil {
		return fmt.Errorf("迁移基础系统表失败: %w", err)
	}

	// ========== 第二批：设计师项目核心表 ==========
	log.Println("迁移第二批：设计师项目核心表...")
	if err := DB.AutoMigrate(
		&model.DesignerProject{}, // designer_projects - 设计项目表
		&model.Material{},        // materials - 材料库表
		&model.DesignStandard{},  // design_standards - 设计规范表
	); err != nil {
		return fmt.Errorf("迁移设计师项目核心表失败: %w", err)
	}

	// ========== 第三批：项目关联表（依赖 designer_projects）==========
	log.Println("迁移第三批：项目关联表...")
	if err := DB.AutoMigrate(
		&model.ProjectDocument{},       // project_documents - 项目文档表
		&model.CadFile{},               // cad_files - CAD文件表
		&model.CostEstimate{},          // cost_estimates - 成本估算表
		&model.DesignCompareResult{},   // design_compare_results - 设计比对结果表
		&model.DesignVersion{},         // design_versions - 设计版本表
		&model.DesignVersionCompare{},  // design_version_compares - 版本对比表
		&model.DesignSuggestion{},      // design_suggestions - 设计建议表
		&model.SuggestionPreference{},  // suggestion_preferences - 建议偏好记录表
		&model.SuggestionTask{},        // suggestion_tasks - 建议生成任务表
		&model.ComplianceCheckResult{}, // compliance_check_results - 合规检查结果表
		&model.ProjectMaterial{},       // project_materials - 项目材料清单表
	); err != nil {
		return fmt.Errorf("迁移项目关联表失败: %w", err)
	}

	// ========== 第四批：CAD相关表（依赖 cad_files）==========
	log.Println("迁移第四批：CAD相关表...")
	if err := DB.AutoMigrate(
		&model.CadGeneration{},          // cad_generations - AI生成CAD任务表
		&model.ConstructionAnnotation{}, // construction_annotations - 施工图标注表
	); err != nil {
		return fmt.Errorf("迁移CAD相关表失败: %w", err)
	}

	// ========== 第五批：效果图和通知表 ==========
	log.Println("迁移第五批：效果图和通知表...")
	if err := DB.AutoMigrate(
		&model.RenderRecord{},    // render_records - 效果图生成记录表
		&model.UserRenderQuota{}, // user_render_quotas - 用户效果图配额表
		&model.Notification{},    // notifications - 用户通知表
	); err != nil {
		return fmt.Errorf("迁移效果图和通知表失败: %w", err)
	}

	// ========== 第六批：聊天室模块 ==========
	log.Println("迁移第六批：聊天室模块...")
	if err := DB.AutoMigrate(
		&model.ChatRoom{},       // chat_rooms - 聊天室表
		&model.ChatMessage{},    // chat_messages - 聊天消息表
		&model.ChatRoomMember{}, // chat_room_members - 聊天室成员表
	); err != nil {
		return fmt.Errorf("迁移聊天室模块失败: %w", err)
	}

	// ========== 第七批：AI对话和日志表 ==========
	log.Println("迁移第七批：AI对话和日志表...")
	if err := DB.AutoMigrate(
		&model.DesignerChatMessage{}, // designer_chat_messages - 设计师AI对话消息表
		&model.AIUsageLog{},          // ai_usage_logs - AI使用日志表
	); err != nil {
		return fmt.Errorf("迁移AI对话和日志表失败: %w", err)
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

// cleanupConstraints 清理可能导致迁移失败的约束
// 处理 MySQL 9.x 中 GORM 迁移时的约束冲突问题
// Clean up constraints that may cause migration failures for MySQL 9.x compatibility
func cleanupConstraints() {
	log.Println("检查并清理可能冲突的数据库约束...")

	// 检查 sys_user 表是否存在
	var sysUserExists int64
	DB.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'sys_user'`).Scan(&sysUserExists)

	if sysUserExists > 0 {
		// 检查是否有用户数据
		var userDataCount int64
		DB.Raw(`SELECT COUNT(*) FROM sys_user`).Scan(&userDataCount)

		if userDataCount == 0 {
			// 数据库有表但没有数据，执行完全清理
			log.Println("检测到空数据库，执行完全清理以避免索引冲突...")
			dropAllTables()
			return
		} else {
			// 有数据，尝试修复索引问题
			log.Println("检测到有数据的数据库，尝试修复索引...")
			fixIndexNaming()
		}
	}

	log.Println("约束检查完成")
}

// fixIndexNaming 修复索引命名问题
// Fix index naming issues for existing tables with data
func fixIndexNaming() {
	// 需要检查和修复的索引映射：旧名称 -> 新名称
	indexFixes := map[string]struct {
		table    string
		oldIndex string
		newIndex string
		column   string
	}{
		"sys_user_user_name": {
			table:    "sys_user",
			oldIndex: "uni_sys_user_user_name",
			newIndex: "idx_user_name",
			column:   "user_name",
		},
		"sys_role_role_code": {
			table:    "sys_role",
			oldIndex: "uni_sys_role_role_code",
			newIndex: "idx_role_code",
			column:   "role_code",
		},
	}

	for name, fix := range indexFixes {
		// 检查旧索引是否存在
		var oldExists int64
		DB.Raw(`SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?`,
			fix.table, fix.oldIndex).Scan(&oldExists)

		if oldExists > 0 {
			log.Printf("修复索引 %s: %s -> %s", name, fix.oldIndex, fix.newIndex)
			// 删除旧索引
			DB.Exec(fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`", fix.table, fix.oldIndex))
			// 创建新索引
			DB.Exec(fmt.Sprintf("CREATE UNIQUE INDEX `%s` ON `%s` (`%s`)", fix.newIndex, fix.table, fix.column))
		}
	}
}

// dropAllTables 删除所有业务表（按依赖顺序倒序删除）
// Drop all business tables in reverse dependency order
func dropAllTables() {
	log.Println("开始删除所有业务表...")

	// 禁用外键检查
	DB.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 按依赖关系倒序删除（先删除依赖表，再删除被依赖表）
	// Drop in reverse dependency order
	tablesToDrop := []string{
		// AI对话和日志表
		"ai_usage_logs", "designer_chat_messages",
		// 聊天室模块
		"chat_room_members", "chat_messages", "chat_rooms",
		// 效果图和通知表
		"notifications", "user_render_quotas", "render_records",
		// CAD相关表
		"construction_annotations", "cad_generations",
		// 项目关联表
		"project_materials", "compliance_check_results", "suggestion_tasks",
		"suggestion_preferences", "design_suggestions", "design_version_compares",
		"design_versions", "design_compare_results", "cost_estimates",
		"cad_files", "project_documents",
		// 设计师项目核心表
		"design_standards", "materials", "designer_projects",
		// 系统基础表（关联表先删）
		"sys_role_button", "sys_role_menu", "sys_user_role",
		"sys_operation_log", "files", "app_user",
		"sys_dictionary", "sys_dictionary_type",
		"sys_department", "sys_button", "sys_menu",
		"sys_role", "sys_user",
	}

	for _, tableName := range tablesToDrop {
		var tableExists int64
		DB.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`, tableName).Scan(&tableExists)

		if tableExists > 0 {
			log.Printf("删除表: %s", tableName)
			if err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tableName)).Error; err != nil {
				log.Printf("警告: 删除表 %s 失败: %v", tableName, err)
			}
		}
	}

	// 重新启用外键检查
	DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	log.Println("所有业务表已删除，将由 GORM 重新创建")
}
