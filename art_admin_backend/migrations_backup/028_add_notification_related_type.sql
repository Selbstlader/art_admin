-- 添加通知关联类型字段
-- Add notification related type field for routing

-- 直接添加列（如果已存在会报错但不影响）
ALTER TABLE notifications ADD COLUMN related_type VARCHAR(50) DEFAULT 'general';

-- 更新现有记录的默认值
UPDATE notifications SET related_type = 'general' WHERE related_type IS NULL OR related_type = '';
