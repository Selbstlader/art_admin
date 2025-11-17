USE gin_admin;

-- 更新 sys_button 表结构
-- 将字段名从旧的改为新的

-- 1. 重命名字段
ALTER TABLE sys_button CHANGE COLUMN title auth_name VARCHAR(50) NOT NULL;
ALTER TABLE sys_button CHANGE COLUMN auth_mark auth_label VARCHAR(100) NOT NULL;
ALTER TABLE sys_button CHANGE COLUMN enabled is_enable TINYINT(1) DEFAULT 1;

-- 2. 添加新字段
ALTER TABLE sys_button ADD COLUMN auth_icon VARCHAR(50) DEFAULT '' AFTER auth_label;
ALTER TABLE sys_button ADD COLUMN auth_sort INT DEFAULT 0 AFTER auth_icon;

-- 3. 更新现有数据的排序
UPDATE sys_button SET auth_sort = 1 WHERE auth_label = 'user:add';
UPDATE sys_button SET auth_sort = 2 WHERE auth_label = 'user:edit';
UPDATE sys_button SET auth_sort = 3 WHERE auth_label = 'user:delete';
UPDATE sys_button SET auth_sort = 4 WHERE auth_label = 'user:export';

UPDATE sys_button SET auth_sort = 1 WHERE auth_label = 'role:add';
UPDATE sys_button SET auth_sort = 2 WHERE auth_label = 'role:edit';
UPDATE sys_button SET auth_sort = 3 WHERE auth_label = 'role:delete';
UPDATE sys_button SET auth_sort = 4 WHERE auth_label = 'role:assign';

UPDATE sys_button SET auth_sort = 1 WHERE auth_label = 'menu:add';
UPDATE sys_button SET auth_sort = 2 WHERE auth_label = 'menu:edit';
UPDATE sys_button SET auth_sort = 3 WHERE auth_label = 'menu:delete';

-- 验证更新结果
SELECT * FROM sys_button ORDER BY menu_id, auth_sort;








