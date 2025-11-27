-- 为user_points表添加复合索引优化排行榜查询性能
-- 索引：(total_points DESC, updated_at DESC) 用于按积分和时间排序的排行榜查询

CREATE INDEX idx_user_points_total_points_updated_at 
ON user_points (total_points DESC, updated_at DESC);

-- 为分类积分查询添加索引
CREATE INDEX idx_user_points_mood_points 
ON user_points (mood_points DESC);

CREATE INDEX idx_user_points_meditation_points 
ON user_points (meditation_points DESC);

CREATE INDEX idx_user_points_journal_points 
ON user_points (journal_points DESC);

CREATE INDEX idx_user_points_analysis_points 
ON user_points (analysis_points DESC);
