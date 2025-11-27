-- 修复achievement_events表缺少processed_at列的问题
-- Migration: 006_add_processed_at_to_achievement_events

-- 添加缺失的processed_at列
ALTER TABLE achievement_events ADD COLUMN processed_at TIMESTAMP NULL COMMENT '处理时间';
