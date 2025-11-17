-- 修复 Layout 组件路径，使其符合前端期望
USE gin_admin;

-- 将所有父级菜单的 component 从 "Layout" 改为 "/index/index"
UPDATE sys_menu SET component = '/index/index' WHERE component = 'Layout';

SELECT '父级菜单组件路径已修复' as message;

-- 查看更新后的结果
SELECT id, name, path, component, parent_id FROM sys_menu ORDER BY parent_id, sort, id;

