-- 效果图生成记录表（异步任务）/ Render records table (async task)
CREATE TABLE IF NOT EXISTS render_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    project_id BIGINT UNSIGNED COMMENT '项目ID',
    cad_file_id BIGINT UNSIGNED COMMENT 'CAD文件ID',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态: pending/processing/completed/failed',
    image_url VARCHAR(500) COMMENT '效果图URL',
    design_proposal TEXT COMMENT '设计方案',
    style VARCHAR(100) COMMENT '设计风格',
    room_type VARCHAR(100) COMMENT '空间类型',
    prompt TEXT COMMENT '生成提示词',
    error_message VARCHAR(500) COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_project_id (project_id),
    INDEX idx_cad_file_id (cad_file_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='效果图生成记录表';

-- 用户效果图配额表 / User render quota table
CREATE TABLE IF NOT EXISTS user_render_quotas (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    daily_limit INT DEFAULT 3 COMMENT '每日限制次数',
    used_today INT DEFAULT 0 COMMENT '今日已使用次数',
    last_reset_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '上次重置日期',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户效果图配额表';
