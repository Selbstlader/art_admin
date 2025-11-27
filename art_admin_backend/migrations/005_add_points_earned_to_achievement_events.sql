-- 修复achievement_events表缺少points_earned列的问题
-- Migration: 005_add_points_earned_to_achievement_events

-- 添加缺失的points_earned列
ALTER TABLE achievement_events ADD COLUMN points_earned INT DEFAULT 0 COMMENT '获得积分';
