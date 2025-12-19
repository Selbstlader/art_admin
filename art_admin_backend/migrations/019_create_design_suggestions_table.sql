-- 创建设计建议表 / Create design suggestions table
CREATE TABLE IF NOT EXISTS design_suggestions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    content TEXT NOT NULL COMMENT '建议内容',
    applicable_scene TEXT COMMENT '适用场景',
    cost_impact VARCHAR(200) COMMENT '成本影响',
    category VARCHAR(100) COMMENT '建议类别',
    priority INT DEFAULT 0 COMMENT '优先级',
    status VARCHAR(50) DEFAULT 'pending' COMMENT '状态(pending/adopted/ignored)',
    detail_info JSON COMMENT '详细信息',
    reference_images JSON COMMENT '参考图片URL列表',
    user_id BIGINT UNSIGNED COMMENT '创建用户ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计建议表';

-- 创建建议偏好记录表 / Create suggestion preferences table
CREATE TABLE IF NOT EXISTS suggestion_preferences (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    suggestion_id BIGINT UNSIGNED NOT NULL COMMENT '关联建议ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型(adopted/ignored/viewed)',
    feedback TEXT COMMENT '用户反馈',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_suggestion_id (suggestion_id),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='建议偏好记录表';

-- 创建AI使用日志表（如果不存在）/ Create AI usage log table if not exists
CREATE TABLE IF NOT EXISTS ai_usage_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED COMMENT '用户ID',
    project_id BIGINT UNSIGNED COMMENT '项目ID',
    service_type VARCHAR(100) COMMENT '服务类型',
    tokens_used INT DEFAULT 0 COMMENT 'Token使用量',
    response_time BIGINT DEFAULT 0 COMMENT '响应时间(毫秒)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_project_id (project_id),
    INDEX idx_service_type (service_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI使用日志表';
