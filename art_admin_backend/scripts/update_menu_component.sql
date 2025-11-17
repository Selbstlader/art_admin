-- 更新菜单组件路径，使其符合前端期望的格式
USE gin_admin;

-- 更新子菜单的 component 路径（去掉 /views 前缀和 /index.vue 后缀）
UPDATE sys_menu SET component = '/dashboard/console' WHERE id = 11;
UPDATE sys_menu SET component = '/dashboard/analysis' WHERE id = 12;
UPDATE sys_menu SET component = '/system/user' WHERE id = 21;
UPDATE sys_menu SET component = '/system/role' WHERE id = 22;
UPDATE sys_menu SET component = '/system/menu' WHERE id = 23;

-- 更新系统管理菜单的图标
UPDATE sys_menu SET icon = 'setting' WHERE id = 2;

SELECT '菜单组件路径已更新' as message;

-- 查看更新后的菜单
SELECT id, name, path, component, icon FROM sys_menu ORDER BY sort, id;

