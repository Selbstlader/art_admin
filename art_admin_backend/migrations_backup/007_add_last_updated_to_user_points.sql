-- 修复user_points表缺少last_updated列的问题
-- Migration: 007_add_last_updated_to_user_points

-- 添加缺失的last_updated列
ALTER TABLE user_points ADD COLUMN last_updated TIMESTAMP NULL COMMENT '最后更新时间';

-- 添加索引
ALTER TABLE user_points ADD INDEX idx_last_updated (last_updated);
