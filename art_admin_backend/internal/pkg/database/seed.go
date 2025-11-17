package database

import (
	"log"
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/utils"
)

// SeedData 初始化种子数据
func SeedData() error {
	if DB == nil {
		return nil
	}

	log.Println("开始初始化种子数据...")

	// 检查是否已有数据
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		log.Println("数据已存在，跳过初始化")
		return nil
	}

	// 1. 创建角色
	if err := seedRoles(); err != nil {
		return err
	}

	// 2. 创建用户
	if err := seedUsers(); err != nil {
		return err
	}

	// 3. 创建菜单
	if err := seedMenus(); err != nil {
		return err
	}

	// 4. 创建按钮权限
	if err := seedButtons(); err != nil {
		return err
	}

	// 5. 分配角色权限
	if err := seedRolePermissions(); err != nil {
		return err
	}

	log.Println("种子数据初始化完成!")
	return nil
}

// seedRoles 创建角色
func seedRoles() error {
	roles := []model.Role{
		{
			RoleID:      1,
			RoleName:    "超级管理员",
			RoleCode:    "admin",
			Description: "拥有所有权限",
			Enabled:     true,
			CreateTime:  time.Now(),
		},
		{
			RoleID:      2,
			RoleName:    "普通用户",
			RoleCode:    "user",
			Description: "普通用户权限",
			Enabled:     true,
			CreateTime:  time.Now(),
		},
		{
			RoleID:      3,
			RoleName:    "访客",
			RoleCode:    "guest",
			Description: "只读权限",
			Enabled:     true,
			CreateTime:  time.Now(),
		},
	}

	for _, role := range roles {
		if err := DB.Create(&role).Error; err != nil {
			log.Printf("创建角色失败 %s: %v", role.RoleName, err)
			return err
		}
	}

	log.Println("✓ 角色创建成功")
	return nil
}

// seedUsers 创建用户
func seedUsers() error {
	// 密码: 123456
	hashedPassword, _ := utils.HashPassword("123456")

	users := []model.User{
		{
			ID:         1,
			UserName:   "admin",
			NickName:   "超级管理员",
			Password:   hashedPassword,
			Email:      "admin@example.com",
			UserPhone:  "13800138000",
			UserGender: "male",
			Avatar:     "https://api.dicebear.com/7.x/avataaars/svg?seed=admin",
			Status:     "1",
			CreateBy:   "system",
			CreateTime: time.Now(),
		},
		{
			ID:         2,
			UserName:   "user",
			NickName:   "普通用户",
			Password:   hashedPassword,
			Email:      "user@example.com",
			UserPhone:  "13800138001",
			UserGender: "female",
			Avatar:     "https://api.dicebear.com/7.x/avataaars/svg?seed=user",
			Status:     "1",
			CreateBy:   "system",
			CreateTime: time.Now(),
		},
		{
			ID:         3,
			UserName:   "guest",
			NickName:   "访客",
			Password:   hashedPassword,
			Email:      "guest@example.com",
			UserPhone:  "13800138002",
			UserGender: "unknown",
			Avatar:     "https://api.dicebear.com/7.x/avataaars/svg?seed=guest",
			Status:     "1",
			CreateBy:   "system",
			CreateTime: time.Now(),
		},
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("创建用户失败 %s: %v", user.UserName, err)
			return err
		}
	}

	// 分配用户角色
	userRoles := []struct {
		UserID int64
		RoleID int64
	}{
		{1, 1}, // admin -> 超级管理员
		{2, 2}, // user -> 普通用户
		{3, 3}, // guest -> 访客
	}

	for _, ur := range userRoles {
		if err := DB.Exec("INSERT INTO sys_user_role (user_id, role_id, create_time) VALUES (?, ?, ?)",
			ur.UserID, ur.RoleID, time.Now()).Error; err != nil {
			log.Printf("分配用户角色失败: %v", err)
			return err
		}
	}

	log.Println("✓ 用户创建成功")
	return nil
}

// seedMenus 创建菜单
func seedMenus() error {
	menus := []model.Menu{
		{
			ID:         1,
			ParentID:   0,
			Name:       "Dashboard",         // 路由名称
			Title:      "仪表盘",            // 显示名称
			Path:       "/dashboard",
			Component:  "/index/index",
			Redirect:   "/dashboard/console",
			Icon:       "dashboard",
			Sort:       1,
			IsEnable:   true,
			IsMenu:     true,
			KeepAlive:  true,
			CreateTime: time.Now(),
		},
		{
			ID:         11,
			ParentID:   1,
			Name:       "DashboardConsole",
			Title:      "主控台",
			Path:       "/dashboard/console",
			Component:  "/dashboard/console",
			Icon:       "console",
			Sort:       1,
			IsEnable:   true,
			IsMenu:     true,
			KeepAlive:  true,
			CreateTime: time.Now(),
		},
		{
			ID:         12,
			ParentID:   1,
			Name:       "DashboardAnalysis",
			Title:      "分析页",
			Path:       "/dashboard/analysis",
			Component:  "/dashboard/analysis",
			Icon:       "analysis",
			Sort:       2,
			IsEnable:   true,
			IsMenu:     true,
			KeepAlive:  true,
			CreateTime: time.Now(),
		},
		{
			ID:         2,
			ParentID:   0,
			Name:       "System",
			Title:      "系统管理",
			Path:       "/system",
			Component:  "/index/index",
			Redirect:   "/system/user",
			Icon:       "setting",
			Sort:       10,
			IsEnable:   true,
			IsMenu:     true,
			KeepAlive:  true,
			CreateTime: time.Now(),
		},
		{
			ID:         21,
			ParentID:   2,
			Name:       "SystemUser",
			Title:      "用户管理",
			Path:       "/system/user",
			Component:  "/system/user",
			Icon:       "user",
			Sort:       1,
			IsEnable:   true,
			IsMenu:     true,
			KeepAlive:  true,
			CreateTime: time.Now(),
		},
		{
			ID:         22,
			ParentID:   2,
			Name:       "SystemRole",
			Title:      "角色管理",
			Path:       "/system/role",
			Component:  "/system/role",
			Icon:       "role",
			Sort:       2,
			IsEnable:   true,
			IsMenu:     true,
			KeepAlive:  true,
			CreateTime: time.Now(),
		},
		{
			ID:         23,
			ParentID:   2,
			Name:       "SystemMenu",
			Title:      "菜单管理",
			Path:       "/system/menu",
			Component:  "/system/menu",
			Icon:       "menu",
			Sort:       3,
			IsEnable:   true,
			IsMenu:     true,
			KeepAlive:  true,
			CreateTime: time.Now(),
		},
	}

	for _, menu := range menus {
		if err := DB.Create(&menu).Error; err != nil {
			log.Printf("创建菜单失败 %s: %v", menu.Name, err)
			return err
		}
	}

	log.Println("✓ 菜单创建成功")
	return nil
}

// seedButtons 创建按钮权限
func seedButtons() error {
	buttons := []model.Button{
		// 用户管理按钮
		{MenuID: 21, AuthName: "新增用户", AuthLabel: "user:add", AuthSort: 1, IsEnable: true, CreateTime: time.Now()},
		{MenuID: 21, AuthName: "编辑用户", AuthLabel: "user:edit", AuthSort: 2, IsEnable: true, CreateTime: time.Now()},
		{MenuID: 21, AuthName: "删除用户", AuthLabel: "user:delete", AuthSort: 3, IsEnable: true, CreateTime: time.Now()},
		{MenuID: 21, AuthName: "导出用户", AuthLabel: "user:export", AuthSort: 4, IsEnable: true, CreateTime: time.Now()},
		// 角色管理按钮
		{MenuID: 22, AuthName: "新增角色", AuthLabel: "role:add", AuthSort: 1, IsEnable: true, CreateTime: time.Now()},
		{MenuID: 22, AuthName: "编辑角色", AuthLabel: "role:edit", AuthSort: 2, IsEnable: true, CreateTime: time.Now()},
		{MenuID: 22, AuthName: "删除角色", AuthLabel: "role:delete", AuthSort: 3, IsEnable: true, CreateTime: time.Now()},
		{MenuID: 22, AuthName: "分配权限", AuthLabel: "role:assign", AuthSort: 4, IsEnable: true, CreateTime: time.Now()},
		// 菜单管理按钮
		{MenuID: 23, AuthName: "新增菜单", AuthLabel: "menu:add", AuthSort: 1, IsEnable: true, CreateTime: time.Now()},
		{MenuID: 23, AuthName: "编辑菜单", AuthLabel: "menu:edit", AuthSort: 2, IsEnable: true, CreateTime: time.Now()},
		{MenuID: 23, AuthName: "删除菜单", AuthLabel: "menu:delete", AuthSort: 3, IsEnable: true, CreateTime: time.Now()},
	}

	for _, button := range buttons {
		if err := DB.Create(&button).Error; err != nil {
			log.Printf("创建按钮失败 %s: %v", button.AuthName, err)
			return err
		}
	}

	log.Println("✓ 按钮权限创建成功")
	return nil
}

// seedRolePermissions 分配角色权限
func seedRolePermissions() error {
	// 管理员拥有所有菜单
	var menus []model.Menu
	DB.Find(&menus)
	for _, menu := range menus {
		DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)",
			1, menu.ID, time.Now())
	}

	// 管理员拥有所有按钮
	var buttons []model.Button
	DB.Find(&buttons)
	for _, button := range buttons {
		DB.Exec("INSERT INTO sys_role_button (role_id, button_id, create_time) VALUES (?, ?, ?)",
			1, button.ID, time.Now())
	}

	// 普通用户只有Dashboard菜单
	DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 2, 1, time.Now())
	DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 2, 11, time.Now())
	DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 2, 12, time.Now())

	// 访客只有Dashboard查看权限
	DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 3, 1, time.Now())
	DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 3, 11, time.Now())

	log.Println("✓ 角色权限分配成功")
	return nil
}
