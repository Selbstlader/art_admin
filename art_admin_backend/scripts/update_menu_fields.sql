USE gin_admin;

-- 更新 sys_menu 表结构
-- 将字段名从旧的改为新的

-- 1. 添加新字段（如果不存在）
ALTER TABLE sys_menu ADD COLUMN IF NOT EXISTS label VARCHAR(100) DEFAULT '' AFTER path;
ALTER TABLE sys_menu ADD COLUMN IF NOT EXISTS is_menu TINYINT(1) DEFAULT 1 AFTER sort;
ALTER TABLE sys_menu ADD COLUMN IF NOT EXISTS show_badge TINYINT(1) DEFAULT 0 AFTER is_iframe;
ALTER TABLE sys_menu ADD COLUMN IF NOT EXISTS show_text_badge VARCHAR(50) DEFAULT '' AFTER show_badge;
ALTER TABLE sys_menu ADD COLUMN IF NOT EXISTS active_path VARCHAR(200) DEFAULT '' AFTER fixed_tab;

-- 2. 重命名字段（如果使用的是旧字段名）
-- 如果字段名是 enabled，改为 is_enable
-- ALTER TABLE sys_menu CHANGE COLUMN enabled is_enable TINYINT(1) DEFAULT 1;

-- 3. 确保所有现有菜单都有 is_menu = true
UPDATE sys_menu SET is_menu = true WHERE is_menu IS NULL OR is_menu = false;

-- 验证更新结果
SELECT id, name, title, path, is_enable, is_menu FROM sys_menu ORDER BY sort;








