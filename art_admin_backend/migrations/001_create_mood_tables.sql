-- 创建情绪管理相关表
-- Migration: 001_create_mood_tables

-- 情绪记录表
CREATE TABLE IF NOT EXISTS mood_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    mood_type VARCHAR(50) NOT NULL COMMENT '情绪类型',
    intensity INT NOT NULL COMMENT '情绪强度 1-10',
    triggers JSON COMMENT '触发因素',
    activities JSON COMMENT '相关活动',
    note TEXT COMMENT '备注',
    location VARCHAR(255) COMMENT '地理位置',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_time (user_id, created_at),
    INDEX idx_mood_type (mood_type),
    INDEX idx_intensity (intensity)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='情绪记录表';

-- 冥想内容表
CREATE TABLE IF NOT EXISTS meditation_contents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '内容ID',
    title VARCHAR(255) NOT NULL COMMENT '标题',
    category VARCHAR(100) NOT NULL COMMENT '分类',
    duration INT NOT NULL COMMENT '时长(秒)',
    audio_url VARCHAR(500) COMMENT '音频文件URL',
    cover_image VARCHAR(500) COMMENT '封面图片URL',
    description TEXT COMMENT '描述',
    difficulty_level INT DEFAULT 1 COMMENT '难度等级 1-5',
    tags JSON COMMENT '标签',
    is_public BOOLEAN DEFAULT TRUE COMMENT '是否公开',
    view_count INT DEFAULT 0 COMMENT '查看次数',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    created_by BIGINT COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_category (category),
    INDEX idx_difficulty (difficulty_level),
    INDEX idx_public (is_public),
    INDEX idx_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='冥想内容表';

-- 用户冥想会话表
CREATE TABLE IF NOT EXISTS user_meditation_sessions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '会话ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    content_id BIGINT NOT NULL COMMENT '内容ID',
    completed_duration INT NOT NULL COMMENT '实际完成时长(秒)',
    is_completed BOOLEAN DEFAULT FALSE COMMENT '是否完成',
    started_at TIMESTAMP NULL COMMENT '开始时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_content_id (content_id),
    INDEX idx_completed (is_completed),
    INDEX idx_session_time (user_id, created_at),
    FOREIGN KEY (content_id) REFERENCES meditation_contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户冥想会话表';

-- 日记表
CREATE TABLE IF NOT EXISTS journal_entries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日记ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    title VARCHAR(255) COMMENT '标题',
    content TEXT NOT NULL COMMENT '内容',
    mood_record_id BIGINT COMMENT '关联的情绪记录ID',
    sentiment_score DECIMAL(3,2) COMMENT '情感分析得分 -1.0 到 1.0',
    is_private BOOLEAN DEFAULT TRUE COMMENT '是否私密',
    tags JSON COMMENT '标签',
    images JSON COMMENT '图片URLs',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_mood_record_id (mood_record_id),
    INDEX idx_private (is_private),
    INDEX idx_journal_time (user_id, created_at),
    FOREIGN KEY (mood_record_id) REFERENCES mood_records(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='日记表';

-- 用户成就表
CREATE TABLE IF NOT EXISTS user_achievements (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '成就ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    achievement_type VARCHAR(100) NOT NULL COMMENT '成就类型',
    target_value INT NOT NULL COMMENT '目标值',
    current_value INT DEFAULT 0 COMMENT '当前值',
    is_unlocked BOOLEAN DEFAULT FALSE COMMENT '是否已解锁',
    unlocked_at TIMESTAMP NULL COMMENT '解锁时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_user_achievement (user_id, achievement_type),
    INDEX idx_user_id (user_id),
    INDEX idx_unlocked (is_unlocked)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户成就表';

-- 情绪分析聚合表
CREATE TABLE IF NOT EXISTS mood_analytics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分析ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    analyze_type VARCHAR(20) NOT NULL COMMENT '分析类型 daily, weekly, monthly',
    analyze_date DATE NOT NULL COMMENT '分析日期',
    mood_distribution JSON COMMENT '情绪分布',
    average_intensity DECIMAL(3,1) COMMENT '平均强度',
    total_records INT DEFAULT 0 COMMENT '总记录数',
    most_common_mood VARCHAR(50) COMMENT '最常见情绪',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_user_analyze (user_id, analyze_type, analyze_date),
    INDEX idx_user_id (user_id),
    INDEX idx_analyze_date (analyze_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='情绪分析聚合表';

-- 插入一些示例冥想内容
INSERT INTO meditation_contents (title, category, duration, description, difficulty_level, is_public) VALUES
('深度睡眠冥想', 'sleep', 1800, '帮助你快速入睡的深度冥想练习', 1, TRUE),
('压力释放冥想', 'stress', 900, '缓解日常压力的放松冥想', 2, TRUE),
('专注力训练', 'focus', 600, '提升注意力和专注力的冥想练习', 3, TRUE),
('焦虑缓解冥想', 'anxiety', 1200, '减轻焦虑情绪的舒缓冥想', 2, TRUE),
('正念呼吸练习', 'breathing', 300, '基础的正念呼吸冥想', 1, TRUE),
('身体扫描冥想', 'body_scan', 1500, '全身放松的身体扫描练习', 2, TRUE);
