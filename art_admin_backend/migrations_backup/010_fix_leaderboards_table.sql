-- 修复leaderboards表结构
-- 移除不必要的user_id字段，因为user_ranking JSON字段已包含用户信息

-- 先删除现有表（如果有数据的话需要备份）
DROP TABLE IF EXISTS leaderboards;

-- 重新创建表，移除user_id字段
CREATE TABLE leaderboards (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    type VARCHAR(20) NOT NULL COMMENT '排行榜类型',
    period VARCHAR(20) NOT NULL COMMENT '周期',
    user_ranking JSON COMMENT '用户排名数据 JSON',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_type_period (type, period),
    INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='排行榜缓存表';
