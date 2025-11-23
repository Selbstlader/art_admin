package database

import (
	"art_admin_backend/internal/model"
	"log"
)

// MigrateAll 完整的数据库迁移（包括中间表）
func MigrateAll() error {
	if DB == nil {
		return nil
	}

	log.Println("开始完整的数据库迁移...")

	// 1. 创建主表
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
		return err
	}

	// 2. 创建中间表（many2many关系表）
	// GORM 会自动创建这些表，但我们需要确保它们存在

	// 用户-角色关联表
	if !DB.Migrator().HasTable("sys_user_role") {
		err = DB.Exec(`
			CREATE TABLE IF NOT EXISTS sys_user_role (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				user_id BIGINT NOT NULL,
				role_id BIGINT NOT NULL,
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE KEY uk_user_role (user_id, role_id),
				KEY idx_user_id (user_id),
				KEY idx_role_id (role_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表'
		`).Error
		if err != nil {
			log.Printf("创建 sys_user_role 表失败: %v", err)
		} else {
			log.Println("创建 sys_user_role 表成功")
		}
	}

	// 角色-菜单关联表
	if !DB.Migrator().HasTable("sys_role_menu") {
		err = DB.Exec(`
			CREATE TABLE IF NOT EXISTS sys_role_menu (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				role_id BIGINT NOT NULL,
				menu_id BIGINT NOT NULL,
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE KEY uk_role_menu (role_id, menu_id),
				KEY idx_role_id (role_id),
				KEY idx_menu_id (menu_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联表'
		`).Error
		if err != nil {
			log.Printf("创建 sys_role_menu 表失败: %v", err)
		} else {
			log.Println("创建 sys_role_menu 表成功")
		}
	}

	// 角色-按钮关联表
	if !DB.Migrator().HasTable("sys_role_button") {
		err = DB.Exec(`
			CREATE TABLE IF NOT EXISTS sys_role_button (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				role_id BIGINT NOT NULL,
				button_id BIGINT NOT NULL,
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE KEY uk_role_button (role_id, button_id),
				KEY idx_role_id (role_id),
				KEY idx_button_id (button_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色按钮关联表'
		`).Error
		if err != nil {
			log.Printf("创建 sys_role_button 表失败: %v", err)
		} else {
			log.Println("创建 sys_role_button 表成功")
		}
	}

	log.Println("数据库迁移完成!")
	return nil
}

// DropAllTables 删除所有表（危险操作，仅用于开发环境）
func DropAllTables() error {
	if DB == nil {
		return nil
	}

	log.Println("警告：正在删除所有表...")

	// 删除中间表
	DB.Exec("DROP TABLE IF EXISTS sys_role_button")
	DB.Exec("DROP TABLE IF EXISTS sys_role_menu")
	DB.Exec("DROP TABLE IF EXISTS sys_user_role")

	// 删除主表
	DB.Migrator().DropTable(
		&model.Button{},
		&model.Menu{},
		&model.Role{},
		&model.User{},
	)

	log.Println("所有表已删除!")
	return nil
}
