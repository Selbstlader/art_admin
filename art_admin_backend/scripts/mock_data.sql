USE gin_admin;

-- 插入测试角色
INSERT INTO `sys_role` (`role_id`, `role_name`, `role_code`, `description`, `enabled`) VALUES
(1, '超级管理员', 'admin', '拥有所有权限', 1),
(2, '普通用户', 'user', '普通用户权限', 1),
(3, '访客', 'guest', '只读权限', 1);

-- 插入测试用户
-- 密码都是 123456 (bcrypt加密后)
INSERT INTO `sys_user` (`id`, `user_name`, `nick_name`, `password`, `email`, `user_phone`, `user_gender`, `avatar`, `status`, `create_by`) VALUES
(1, 'admin', '超级管理员', '$2a$10$hW6cCywvV3KqFpep3Qd.buD/04YsVcb6OsL4Bmh3DyvVNVn..AN8S', 'admin@example.com', '13800138000', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin', '1', 'system'),
(2, 'user', '普通用户', '$2a$10$hW6cCywvV3KqFpep3Qd.buD/04YsVcb6OsL4Bmh3DyvVNVn..AN8S', 'user@example.com', '13800138001', 'female', 'https://api.dicebear.com/7.x/avataaars/svg?seed=user', '1', 'system'),
(3, 'guest', '访客', '$2a$10$hW6cCywvV3KqFpep3Qd.buD/04YsVcb6OsL4Bmh3DyvVNVn..AN8S', 'guest@example.com', '13800138002', 'unknown', 'https://api.dicebear.com/7.x/avataaars/svg?seed=guest', '1', 'system');

-- 用户角色关联
INSERT INTO `sys_user_role` (`user_id`, `role_id`) VALUES
(1, 1), -- admin -> 超级管理员
(2, 2), -- user -> 普通用户
(3, 3); -- guest -> 访客

-- 插入测试菜单
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `path`, `component`, `redirect`, `title`, `icon`, `sort`, `enabled`) VALUES
-- 一级菜单
(1, 0, 'Dashboard', '/dashboard', 'Layout', '/dashboard/console', '仪表盘', 'dashboard', 1, 1),
(2, 0, 'System', '/system', 'Layout', '/system/user', '系统管理', 'setting', 2, 1),

-- Dashboard 子菜单
(11, 1, 'DashboardConsole', '/dashboard/console', '/views/dashboard/console/index.vue', NULL, '主控台', 'console', 1, 1),
(12, 1, 'DashboardAnalysis', '/dashboard/analysis', '/views/dashboard/analysis/index.vue', NULL, '分析页', 'analysis', 2, 1),

-- System 子菜单
(21, 2, 'SystemUser', '/system/user', '/views/system/user/index.vue', NULL, '用户管理', 'user', 1, 1),
(22, 2, 'SystemRole', '/system/role', '/views/system/role/index.vue', NULL, '角色管理', 'role', 2, 1),
(23, 2, 'SystemMenu', '/system/menu', '/views/system/menu/index.vue', NULL, '菜单管理', 'menu', 3, 1);

-- 角色菜单关联（admin拥有所有菜单）
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES
-- admin 角色
(1, 1), (1, 11), (1, 12),
(1, 2), (1, 21), (1, 22), (1, 23),
-- user 角色
(2, 1), (2, 11), (2, 12),
-- guest 角色
(3, 1), (3, 11);

-- 插入按钮权限
INSERT INTO `sys_button` (`menu_id`, `title`, `auth_mark`, `enabled`) VALUES
-- 用户管理按钮
(21, '新增用户', 'user:add', 1),
(21, '编辑用户', 'user:edit', 1),
(21, '删除用户', 'user:delete', 1),
(21, '导出用户', 'user:export', 1),

-- 角色管理按钮
(22, '新增角色', 'role:add', 1),
(22, '编辑角色', 'role:edit', 1),
(22, '删除角色', 'role:delete', 1),
(22, '分配权限', 'role:assign', 1),

-- 菜单管理按钮
(23, '新增菜单', 'menu:add', 1),
(23, '编辑菜单', 'menu:edit', 1),
(23, '删除菜单', 'menu:delete', 1);

-- 角色按钮关联（admin拥有所有按钮权限）
INSERT INTO `sys_role_button` (`role_id`, `button_id`)
SELECT 1, id FROM `sys_button`;

-- user 角色只有部分按钮权限
INSERT INTO `sys_role_button` (`role_id`, `button_id`)
SELECT 2, id FROM `sys_button` WHERE auth_mark IN ('user:edit', 'user:export');

