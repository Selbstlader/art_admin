-- 成就系统表结构创建
-- 创建成就系统的完整表结构

-- 删除可能冲突的旧user_achievements表
DROP TABLE IF EXISTS user_achievements;

-- 成就模板表
CREATE TABLE IF NOT EXISTS achievements (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '成就ID',
    type VARCHAR(50) NOT NULL COMMENT '成就类型',
    name VARCHAR(100) NOT NULL COMMENT '成就名称',
    description TEXT COMMENT '成就描述',
    icon VARCHAR(100) COMMENT '图标',
    tier VARCHAR(20) NOT NULL COMMENT '等级: bronze, silver, gold, platinum',
    target_value INT NOT NULL COMMENT '目标值',
    points INT DEFAULT 0 COMMENT '奖励积分',
    is_hidden BOOLEAN DEFAULT FALSE COMMENT '是否隐藏成就',
    criteria JSON COMMENT '触发条件 JSON',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_type (type),
    INDEX idx_tier (tier),
    INDEX idx_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成就模板表';

-- 用户成就记录表
CREATE TABLE IF NOT EXISTS user_achievements (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    achievement_id BIGINT NOT NULL COMMENT '成就ID',
    current_value INT DEFAULT 0 COMMENT '当前进度',
    is_unlocked BOOLEAN DEFAULT FALSE COMMENT '是否已解锁',
    unlocked_at TIMESTAMP NULL COMMENT '解锁时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_achievement_id (achievement_id),
    INDEX idx_user_unlocked (user_id, is_unlocked),
    UNIQUE KEY uk_user_achievement (user_id, achievement_id),
    FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户成就记录表';

-- 用户积分记录表
CREATE TABLE IF NOT EXISTS user_points (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    total_points INT DEFAULT 0 COMMENT '总积分',
    level INT DEFAULT 1 COMMENT '等级',
    experience_points INT DEFAULT 0 COMMENT '经验值',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    UNIQUE KEY uk_user_points (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户积分记录表';

-- 成就事件表
CREATE TABLE IF NOT EXISTS achievement_events (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '事件ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    event_type VARCHAR(50) NOT NULL COMMENT '事件类型',
    event_data JSON COMMENT '事件数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_event_type (event_type),
    INDEX idx_user_created (user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成就事件表';

-- 排行榜表
CREATE TABLE IF NOT EXISTS leaderboards (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    type VARCHAR(20) NOT NULL COMMENT '排行榜类型',
    period VARCHAR(20) NOT NULL COMMENT '周期',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    total_points INT DEFAULT 0 COMMENT '总积分',
    `rank` INT DEFAULT 0 COMMENT '排名',
    level INT DEFAULT 1 COMMENT '等级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_type_period (type, period),
    INDEX idx_user_id (user_id),
    INDEX idx_rank (`rank`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='排行榜表';

-- 添加额外的性能优化索引（避免重复创建表内已定义的索引）
-- CREATE INDEX idx_user_achievements_user_unlocked ON user_achievements(user_id, is_unlocked); -- 已在表内定义
-- CREATE INDEX idx_achievements_type_active ON achievements(type, is_active); -- 已在表内定义
-- CREATE INDEX idx_user_points_user_id ON user_points(user_id); -- 已在表内定义
-- CREATE INDEX idx_achievement_events_user_created ON achievement_events(user_id, created_at); -- 已在表内定义
-- CREATE INDEX idx_leaderboards_type_period ON leaderboards(type, period); -- 已在表内定义

-- 情绪记录成就
INSERT INTO achievements (type, name, description, icon, tier, target_value, points, is_hidden, criteria, is_active, created_at, updated_at) VALUES
-- 铜牌成就（可见）
('mood_record', '初次记录', '记录你的第一个情绪状态', 'mood_1', 'bronze', 1, 10, false, '{"type": "count", "value": 1}', true, NOW(), NOW()),
('mood_record', '情绪新手', '累计记录3个情绪状态', 'mood_3', 'bronze', 3, 20, false, '{"type": "count", "value": 3}', true, NOW(), NOW()),
('mood_record', '情绪观察者', '累计记录10个情绪状态', 'mood_10', 'bronze', 10, 30, false, '{"type": "count", "value": 10}', true, NOW(), NOW()),

-- 银牌成就（可见）
('mood_record', '情绪记录者', '累计记录25个情绪状态', 'mood_25', 'silver', 25, 50, false, '{"type": "count", "value": 25}', true, NOW(), NOW()),
('mood_record', '连续记录3天', '连续3天记录情绪状态', 'streak_3', 'silver', 3, 60, false, '{"type": "streak", "value": 3}', true, NOW(), NOW()),
('mood_record', '连续记录7天', '连续7天记录情绪状态', 'streak_7', 'silver', 7, 80, false, '{"type": "streak", "value": 7}', true, NOW(), NOW()),

-- 金牌成就（隐藏，直到银牌成就解锁）
('mood_record', '情绪大师', '累计记录100个情绪状态', 'mood_100', 'gold', 100, 120, true, '{"type": "count", "value": 100}', true, NOW(), NOW()),
('mood_record', '连续记录30天', '连续30天记录情绪状态', 'streak_30', 'gold', 30, 150, true, '{"type": "streak", "value": 30}', true, NOW(), NOW()),

-- 铂金成就（隐藏，直到金牌成就解锁）
('mood_record', '情绪传奇', '累计记录365个情绪状态', 'mood_365', 'platinum', 365, 300, true, '{"type": "count", "value": 365}', true, NOW(), NOW()),
('mood_record', '年度记录者', '连续365天记录情绪状态', 'streak_365', 'platinum', 365, 500, true, '{"type": "streak", "value": 365}', true, NOW(), NOW()),

-- 日记成就
-- 铜牌成就（可见）
('journal', '日记初体验', '写下第一篇日记', 'journal_1', 'bronze', 1, 15, false, '{"type": "count", "value": 1}', true, NOW(), NOW()),
('journal', '日记新手', '累计写3篇日记', 'journal_3', 'bronze', 3, 25, false, '{"type": "count", "value": 3}', true, NOW(), NOW()),
('journal', '坚持记录', '累计写10篇日记', 'journal_10', 'bronze', 10, 40, false, '{"type": "count", "value": 10}', true, NOW(), NOW()),

-- 银牌成就（可见）
('journal', '日记达人', '累计写25篇日记', 'journal_25', 'silver', 25, 70, false, '{"type": "count", "value": 25}', true, NOW(), NOW()),
('journal', '深度思考者', '累计写50篇日记', 'journal_50', 'silver', 50, 90, false, '{"type": "count", "value": 50}', true, NOW(), NOW()),

-- 金牌成就（隐藏）
('journal', '日记作家', '累计写100篇日记', 'journal_100', 'gold', 100, 140, true, '{"type": "count", "value": 100}', true, NOW(), NOW()),
('journal', '持续记录者', '累计写200篇日记', 'journal_200', 'gold', 200, 200, true, '{"type": "count", "value": 200}', true, NOW(), NOW()),

-- 铂金成就（隐藏）
('journal', '日记大师', '累计写365篇日记', 'journal_365', 'platinum', 365, 350, true, '{"type": "count", "value": 365}', true, NOW(), NOW()),

-- 冥想成就
-- 铜牌成就（可见）
('meditation', '冥想初体验', '完成10分钟冥想', 'meditation_10', 'bronze', 10, 20, false, '{"type": "minutes", "value": 10}', true, NOW(), NOW()),
('meditation', '冥想练习者', '累计冥想30分钟', 'meditation_30', 'bronze', 30, 35, false, '{"type": "minutes", "value": 30}', true, NOW(), NOW()),

-- 银牌成就（可见）
('meditation', '冥想爱好者', '累计冥想60分钟', 'meditation_60', 'silver', 60, 60, false, '{"type": "minutes", "value": 60}', true, NOW(), NOW()),
('meditation', '冥想达人', '累计冥想150分钟', 'meditation_150', 'silver', 150, 100, false, '{"type": "minutes", "value": 150}', true, NOW(), NOW()),

-- 金牌成就（隐藏）
('meditation', '冥想大师', '累计冥想300分钟', 'meditation_300', 'gold', 300, 180, true, '{"type": "minutes", "value": 300}', true, NOW(), NOW()),
('meditation', '冥想专家', '累计冥想600分钟', 'meditation_600', 'gold', 600, 250, true, '{"type": "minutes", "value": 600}', true, NOW(), NOW()),

-- 铂金成就（隐藏）
('meditation', '冥想传奇', '累计冥想1000分钟', 'meditation_1000', 'platinum', 1000, 400, true, '{"type": "minutes", "value": 1000}', true, NOW(), NOW()),

-- AI分析成就
-- 铜牌成就（可见）
('analysis', 'AI分析初体验', '完成第一次AI情绪分析', 'analysis_1', 'bronze', 1, 25, false, '{"type": "count", "value": 1}', true, NOW(), NOW()),
('analysis', '分析探索者', '完成3次AI情绪分析', 'analysis_3', 'bronze', 3, 40, false, '{"type": "count", "value": 3}', true, NOW(), NOW()),

-- 银牌成就（可见）
('analysis', '分析爱好者', '完成5次AI情绪分析', 'analysis_5', 'silver', 5, 60, false, '{"type": "count", "value": 5}', true, NOW(), NOW()),
('analysis', '分析达人', '完成10次AI情绪分析', 'analysis_10', 'silver', 10, 100, false, '{"type": "count", "value": 10}', true, NOW(), NOW()),

-- 金牌成就（隐藏）
('analysis', '分析专家', '完成20次AI情绪分析', 'analysis_20', 'gold', 20, 160, true, '{"type": "count", "value": 20}', true, NOW(), NOW()),
('analysis', '分析大师', '完成35次AI情绪分析', 'analysis_35', 'gold', 35, 220, true, '{"type": "count", "value": 35}', true, NOW(), NOW()),

-- 铂金成就（隐藏）
('analysis', '分析传奇', '完成50次AI情绪分析', 'analysis_50', 'platinum', 50, 350, true, '{"type": "count", "value": 50}', true, NOW(), NOW()),

-- 质量成就（隐藏，鼓励深度使用）
('quality', '深度思考', '写10篇超过500字的日记', 'quality_long', 'gold', 10, 120, true, '{"type": "long_journal", "value": 10, "min_words": 500}', true, NOW(), NOW()),
('quality', '积极心态', '连续7天保持平均情绪分数7分以上', 'quality_positive', 'gold', 7, 150, true, '{"type": "positive_streak", "value": 7, "min_score": 7}', true, NOW(), NOW()),
('quality', '情绪稳定', '连续30天记录情绪且分数波动不超过2分', 'quality_stable', 'platinum', 30, 300, true, '{"type": "stable_mood", "value": 30, "max_variance": 2}', true, NOW(), NOW()),

-- 特殊成就（隐藏，里程碑式奖励）
('growth', '新手毕业', '解锁所有铜牌成就', 'graduation_bronze', 'silver', 0, 200, true, '{"type": "tier_complete", "tier": "bronze"}', true, NOW(), NOW()),
('growth', '银牌达人', '解锁所有银牌成就', 'graduation_silver', 'gold', 0, 400, true, '{"type": "tier_complete", "tier": "silver"}', true, NOW(), NOW()),
('growth', '金牌大师', '解锁所有金牌成就', 'graduation_gold', 'platinum', 0, 800, true, '{"type": "tier_complete", "tier": "gold"}', true, NOW(), NOW()),

-- 参与度成就（鼓励多样化使用）
('engagement', '全面发展', '在一天内同时记录情绪、写日记、完成冥想', 'engagement_daily', 'silver', 3, 100, false, '{"type": "daily_all", "value": 3}', true, NOW(), NOW()),
('engagement', '周活跃', '一周内至少有5天记录任何内容', 'engagement_weekly', 'silver', 5, 80, false, '{"type": "weekly_active", "value": 5}', true, NOW(), NOW()),
('engagement', '月活跃', '一个月内至少有20天记录任何内容', 'engagement_monthly', 'gold', 20, 200, true, '{"type": "monthly_active", "value": 20}', true, NOW(), NOW());
