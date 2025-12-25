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

	// 检查是否已有用户数据
	var userCount int64
	DB.Model(&model.User{}).Count(&userCount)

	// 如果没有用户数据，初始化基础数据
	if userCount == 0 {
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
	} else {
		log.Println("用户数据已存在，跳过基础数据初始化")
	}

	// 6. 初始化部门（独立检查）
	if err := seedDepartments(); err != nil {
		log.Printf("初始化部门失败（可能已存在）: %v", err)
	}

	// 7. 初始化字典（独立检查）
	if err := seedDictionaries(); err != nil {
		log.Printf("初始化字典失败（可能已存在）: %v", err)
	}

	// 8. 初始化设计规范（独立检查）
	if err := seedDesignStandards(); err != nil {
		log.Printf("初始化设计规范失败（可能已存在）: %v", err)
	}

	log.Println("种子数据初始化完成!")
	return nil
}

// seedRoles 创建角色
func seedRoles() error {
	roles := []model.Role{
		{RoleID: 1, RoleName: "超级管理员", RoleCode: "admin", Description: "拥有所有权限", Enabled: true, CreateTime: time.Now()},
		{RoleID: 2, RoleName: "普通用户", RoleCode: "user", Description: "普通用户权限", Enabled: true, CreateTime: time.Now()},
		{RoleID: 3, RoleName: "访客", RoleCode: "guest", Description: "只读权限", Enabled: true, CreateTime: time.Now()},
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
	hashedPassword, _ := utils.HashPassword("123456")

	users := []model.User{
		{ID: 1, UserName: "admin", NickName: "超级管理员", Password: hashedPassword, Email: "admin@example.com", UserPhone: "13800138000", UserGender: "male", Avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=admin", Status: "1", CreateBy: "system", CreateTime: time.Now()},
		{ID: 2, UserName: "user", NickName: "普通用户", Password: hashedPassword, Email: "user@example.com", UserPhone: "13800138001", UserGender: "female", Avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=user", Status: "1", CreateBy: "system", CreateTime: time.Now()},
		{ID: 3, UserName: "guest", NickName: "访客", Password: hashedPassword, Email: "guest@example.com", UserPhone: "13800138002", UserGender: "unknown", Avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=guest", Status: "1", CreateBy: "system", CreateTime: time.Now()},
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("创建用户失败 %s: %v", user.UserName, err)
			return err
		}
	}

	// 分配用户角色
	userRoles := []struct{ UserID, RoleID int64 }{{1, 1}, {2, 2}, {3, 3}}
	for _, ur := range userRoles {
		DB.Exec("INSERT INTO sys_user_role (user_id, role_id, create_time) VALUES (?, ?, ?)", ur.UserID, ur.RoleID, time.Now())
	}

	log.Println("✓ 用户创建成功")
	return nil
}

// seedMenus 创建菜单（完整菜单数据）
func seedMenus() error {
	menus := []model.Menu{
		// 仪表盘
		{ID: 1, ParentID: 0, Name: "Dashboard", Path: "/dashboard", Component: "/index/index", Redirect: "/dashboard/console", Title: "仪表盘", Icon: "dashboard", Sort: 1, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 11, ParentID: 1, Name: "DashboardConsole", Path: "/dashboard/console", Component: "/dashboard/console", Title: "主控台", Icon: "console", Sort: 1, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 12, ParentID: 1, Name: "DashboardAnalysis", Path: "/dashboard/analysis", Component: "/dashboard/analysis", Title: "分析页", Icon: "analysis", Sort: 2, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},

		// 系统管理
		{ID: 2, ParentID: 0, Name: "System", Path: "/system", Component: "/index/index", Redirect: "/system/user", Title: "系统管理", Icon: "setting", Sort: 2, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 21, ParentID: 2, Name: "SystemUser", Path: "/system/user", Component: "/system/user", Title: "用户管理", Icon: "user", Sort: 1, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 22, ParentID: 2, Name: "SystemRole", Path: "/system/role", Component: "/system/role", Title: "角色管理", Icon: "role", Sort: 2, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 23, ParentID: 2, Name: "SystemMenu", Path: "/system/menu", Component: "/system/menu", Title: "菜单管理", Icon: "menu", Sort: 3, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 28, ParentID: 2, Name: "SystemDictionary", Path: "/system/dictionary", Component: "/system/dictionary", Title: "字典管理", Sort: 4, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 29, ParentID: 2, Name: "SystemDepartment", Path: "/system/department", Component: "/system/department", Title: "部门管理", Sort: 5, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 30, ParentID: 2, Name: "SystemOperationLog", Path: "/system/operation-log", Component: "/system/operation-log", Title: "日志管理", Sort: 6, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 39, ParentID: 2, Name: "SystemAppUser", Path: "/system/app-user", Component: "/system/app-user", Title: "app用户管理", Sort: 7, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 49, ParentID: 2, Name: "SystemUserCenter", Path: "/system/user-center", Component: "/system/user-center", Title: "个人中心", Sort: 8, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},

		// Dify管理
		{ID: 36, ParentID: 0, Name: "Dify", Path: "/dify", Component: "/index/index", Title: "dify管理", Sort: 4, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 37, ParentID: 36, Name: "DifyChat", Path: "/dify/chat", Component: "/dify/chat", Title: "ai对话", Sort: 1, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 38, ParentID: 36, Name: "DifyKnowledge", Path: "/dify/knowledge", Component: "/dify/knowledge", Title: "知识库管理", Sort: 2, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},

		// 设计师助手
		{ID: 1023, ParentID: 0, Name: "Designer", Path: "/designer", Component: "/index/index", Title: "设计师助手", Sort: 3, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 1024, ParentID: 1023, Name: "DesignerProjectList", Path: "/designer-assistant/project/ProjectList", Component: "/designer-assistant/project/ProjectList", Title: "项目管理", Sort: 1, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 1025, ParentID: 1023, Name: "DesignerProjectCreate", Path: "/designer-assistant/project/ProjectCreate", Component: "/designer-assistant/project/ProjectCreate", Title: "创建项目", Sort: 2, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1026, ParentID: 1023, Name: "DesignerProjectDetail", Path: "/designer-assistant/project/ProjectDetail", Component: "/designer-assistant/project/ProjectDetail", Title: "项目详情", Sort: 3, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1027, ParentID: 1023, Name: "DesignerDocumentUpload", Path: "/designer-assistant/document/DocumentUpload", Component: "/designer-assistant/document/DocumentUpload", Title: "文档分析", Sort: 4, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1028, ParentID: 1023, Name: "DesignerCompareUpload", Path: "/designer-assistant/design-compare/CompareUpload", Component: "/designer-assistant/design-compare/CompareUpload", Title: "设计比对", Sort: 5, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1029, ParentID: 1023, Name: "DesignerCadUpload", Path: "/designer-assistant/cad-viewer/CadUpload", Component: "/designer-assistant/cad-viewer/CadUpload", Title: "CAD预览", Sort: 6, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1030, ParentID: 1023, Name: "DesignerMaterialList", Path: "/designer-assistant/material/MaterialList", Component: "/designer-assistant/material/MaterialList", Title: "材料库", Sort: 7, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 1031, ParentID: 1023, Name: "DesignerCostConfig", Path: "/designer-assistant/cost/CostConfig", Component: "/designer-assistant/cost/CostConfig", Title: "成本估算", Sort: 8, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1032, ParentID: 1023, Name: "DesignerChat", Path: "/designer-assistant/chat/DesignerChat", Component: "/designer-assistant/chat/DesignerChat", Title: "AI对话", Sort: 9, IsEnable: true, IsMenu: true, KeepAlive: true, CreateTime: time.Now()},
		{ID: 1033, ParentID: 1023, Name: "DesignerComplianceCheck", Path: "/designer-assistant/compliance/ComplianceCheck", Component: "/designer-assistant/compliance/ComplianceCheck", Title: "合规检查", Sort: 10, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1034, ParentID: 1023, Name: "DesignerVersionList", Path: "/designer-assistant/version-compare/VersionList", Component: "/designer-assistant/version-compare/VersionList", Title: "版本对比", Sort: 11, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1035, ParentID: 1023, Name: "DesignerVersionDiff", Path: "/designer-assistant/version-compare/VersionDiff", Component: "/designer-assistant/version-compare/VersionDiff", Title: "版本差异", Sort: 12, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1037, ParentID: 1023, Name: "DesignerCadGeneration", Path: "/designer-assistant/cad-generation/CadGenerationList", Component: "/designer-assistant/cad-generation/CadGenerationList", Title: "CAD生成", Sort: 13, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
		{ID: 1039, ParentID: 1023, Name: "DesignerCadViewer", Path: "/designer-assistant/cad-viewer/CadViewer", Component: "/designer-assistant/cad-viewer/CadViewer", Title: "CAD查看", Sort: 14, IsEnable: true, IsMenu: true, KeepAlive: true, IsHide: true, CreateTime: time.Now()},
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
		DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 1, menu.ID, time.Now())
	}

	// 管理员拥有所有按钮
	var buttons []model.Button
	DB.Find(&buttons)
	for _, button := range buttons {
		DB.Exec("INSERT INTO sys_role_button (role_id, button_id, create_time) VALUES (?, ?, ?)", 1, button.ID, time.Now())
	}

	// 普通用户拥有Dashboard和设计师助手菜单
	userMenuIDs := []int64{1, 11, 12, 1023, 1024, 1030, 1032}
	for _, menuID := range userMenuIDs {
		DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 2, menuID, time.Now())
	}

	// 访客只有Dashboard查看权限
	guestMenuIDs := []int64{1, 11, 12}
	for _, menuID := range guestMenuIDs {
		DB.Exec("INSERT INTO sys_role_menu (role_id, menu_id, create_time) VALUES (?, ?, ?)", 3, menuID, time.Now())
	}

	log.Println("✓ 角色权限分配成功")
	return nil
}

// seedDepartments 初始化部门
func seedDepartments() error {
	var count int64
	DB.Model(&model.Department{}).Count(&count)
	if count > 0 {
		return nil
	}

	depts := []model.Department{
		{DeptID: 1, ParentID: nil, DeptName: "总公司", DeptCode: "HQ", OrderNum: 1, Status: 1, CreateTime: time.Now()},
		{DeptID: 2, ParentID: ptrInt64(1), DeptName: "技术部", DeptCode: "TECH", OrderNum: 1, Status: 1, CreateTime: time.Now()},
		{DeptID: 3, ParentID: ptrInt64(1), DeptName: "设计部", DeptCode: "DESIGN", OrderNum: 2, Status: 1, CreateTime: time.Now()},
		{DeptID: 4, ParentID: ptrInt64(1), DeptName: "市场部", DeptCode: "MARKET", OrderNum: 3, Status: 1, CreateTime: time.Now()},
	}

	for _, dept := range depts {
		if err := DB.Create(&dept).Error; err != nil {
			return err
		}
	}

	log.Println("✓ 部门创建成功")
	return nil
}

// seedDictionaries 初始化字典
func seedDictionaries() error {
	// 字典类型
	dictTypes := []model.DictionaryType{
		{ID: 1, TypeName: "用户状态", TypeCode: "user_status", Description: "用户状态", Enabled: true, CreateTime: time.Now()},
		{ID: 2, TypeName: "性别", TypeCode: "gender", Description: "性别", Enabled: true, CreateTime: time.Now()},
		{ID: 3, TypeName: "项目状态", TypeCode: "project_status", Description: "设计项目状态", Enabled: true, CreateTime: time.Now()},
		{ID: 4, TypeName: "设计风格", TypeCode: "design_style", Description: "设计风格类型", Enabled: true, CreateTime: time.Now()},
	}

	for _, dt := range dictTypes {
		var count int64
		DB.Model(&model.DictionaryType{}).Where("type_code = ?", dt.TypeCode).Count(&count)
		if count == 0 {
			DB.Create(&dt)
		}
	}

	// 字典数据
	dicts := []model.Dictionary{
		{TypeCode: "user_status", Label: "正常", Value: "1", OrderNum: 1, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "user_status", Label: "禁用", Value: "0", OrderNum: 2, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "gender", Label: "男", Value: "male", OrderNum: 1, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "gender", Label: "女", Value: "female", OrderNum: 2, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "gender", Label: "未知", Value: "unknown", OrderNum: 3, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "project_status", Label: "草稿", Value: "draft", OrderNum: 1, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "project_status", Label: "进行中", Value: "in_progress", OrderNum: 2, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "project_status", Label: "已完成", Value: "completed", OrderNum: 3, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "project_status", Label: "已归档", Value: "archived", OrderNum: 4, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "design_style", Label: "现代简约", Value: "modern", OrderNum: 1, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "design_style", Label: "中式风格", Value: "chinese", OrderNum: 2, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "design_style", Label: "欧式风格", Value: "european", OrderNum: 3, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "design_style", Label: "工业风格", Value: "industrial", OrderNum: 4, Enabled: true, CreateTime: time.Now()},
		{TypeCode: "design_style", Label: "北欧风格", Value: "nordic", OrderNum: 5, Enabled: true, CreateTime: time.Now()},
	}

	for _, d := range dicts {
		var count int64
		DB.Model(&model.Dictionary{}).Where("type_code = ? AND value = ?", d.TypeCode, d.Value).Count(&count)
		if count == 0 {
			DB.Create(&d)
		}
	}

	log.Println("✓ 字典数据创建成功")
	return nil
}

// seedDesignStandards 初始化设计规范
func seedDesignStandards() error {
	var count int64
	DB.Model(&model.DesignStandard{}).Count(&count)
	if count > 0 {
		return nil
	}

	standards := []model.DesignStandard{
		{Code: "GB50016-2014", Name: "建筑设计防火规范", Category: "fire", Content: "建筑物的耐火等级、防火分区、安全疏散、消防设施等应符合本规范要求", ApplicableTypes: `["office","commercial","industrial","public"]`, Version: "2018修订版", Source: "国家标准", Interpretation: "本规范是建筑防火设计的基本依据", Status: "active"},
		{Code: "GB50016-5.5.17", Name: "疏散走道宽度要求", Category: "fire", Content: "公共建筑内疏散走道的净宽度不应小于1.1m；人员密集场所不应小于1.4m", ApplicableTypes: `["commercial","public"]`, Version: "2018修订版", Source: "国家标准", Interpretation: "疏散走道是人员疏散的主要通道", Status: "active"},
		{Code: "GB50016-5.5.21", Name: "安全出口数量要求", Category: "fire", Content: "公共建筑每个防火分区或一个防火分区的每个楼层，其安全出口的数量应经计算确定，且不应少于2个", ApplicableTypes: `["commercial","public"]`, Version: "2018修订版", Source: "国家标准", Interpretation: "多个安全出口可确保在一个出口被封堵时仍有其他疏散途径", Status: "active"},
		{Code: "GB50763-2012", Name: "无障碍设计规范", Category: "accessibility", Content: "公共建筑应设置无障碍通道、无障碍电梯、无障碍卫生间等设施", ApplicableTypes: `["office","commercial","public"]`, Version: "2012版", Source: "国家标准", Interpretation: "保障残疾人、老年人等特殊群体的通行和使用需求", Status: "active"},
		{Code: "GB50763-3.3.1", Name: "无障碍坡道坡度要求", Category: "accessibility", Content: "无障碍坡道的坡度不应大于1:12，困难情况下不应大于1:8", ApplicableTypes: `["office","commercial","public"]`, Version: "2012版", Source: "国家标准", Interpretation: "坡度过大会增加轮椅使用者的通行难度和安全风险", Status: "active"},
		{Code: "GB50325-2020", Name: "民用建筑工程室内环境污染控制标准", Category: "environmental", Content: "室内装饰装修材料的甲醛、苯、TVOC等有害物质释放量应符合本标准限值", ApplicableTypes: `["office","commercial","public"]`, Version: "2020版", Source: "国家标准", Interpretation: "控制室内空气污染，保障人员健康", Status: "active"},
		{Code: "GB50352-2019", Name: "民用建筑设计统一标准", Category: "safety", Content: "建筑设计应满足结构安全、使用安全、防火安全等基本要求", ApplicableTypes: `["office","commercial","industrial","public"]`, Version: "2019版", Source: "国家标准", Interpretation: "民用建筑设计的基本准则和通用要求", Status: "active"},
		{Code: "GB50352-6.6.3", Name: "栏杆高度要求", Category: "safety", Content: "临空高度在24m以下时，栏杆高度不应低于1.05m；临空高度在24m及以上时，栏杆高度不应低于1.10m", ApplicableTypes: `["office","commercial","public"]`, Version: "2019版", Source: "国家标准", Interpretation: "防止人员坠落，保障使用安全", Status: "active"},
	}

	for _, s := range standards {
		DB.Create(&s)
	}

	log.Println("✓ 设计规范创建成功")
	return nil
}

// ptrInt64 返回 int64 指针
func ptrInt64(v int64) *int64 {
	return &v
}
