-- 修复user_points表缺少分类积分字段的问题
-- Migration: 008_add_category_points_to_user_points

-- 添加缺失的分类积分字段
ALTER TABLE user_points ADD COLUMN mood_points INT DEFAULT 0 COMMENT '情绪积分';
ALTER TABLE user_points ADD COLUMN meditation_points INT DEFAULT 0 COMMENT '冥想积分';
ALTER TABLE user_points ADD COLUMN journal_points INT DEFAULT 0 COMMENT '日记积分';
ALTER TABLE user_points ADD COLUMN analysis_points INT DEFAULT 0 COMMENT 'AI分析积分';
ALTER TABLE user_points ADD COLUMN bonus_points INT DEFAULT 0 COMMENT '额外积分';
