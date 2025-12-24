-- 设计师AI助手系统 - AI分析相关表
-- Designer AI Assistant System - AI Analysis Tables

-- =====================================================
-- 3. AI分析与设计版本表 / AI Analysis & Design Version Tables
-- =====================================================

-- 设计版本表
CREATE TABLE IF NOT EXISTS design_versions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    version_number INT NOT NULL COMMENT '版本号',
    version_name VARCHAR(100) COMMENT '版本名称',
    description TEXT COMMENT '版本说明',
    design_images JSON COMMENT '设计图列表',
    cad_file_ids JSON COMMENT '关联CAD文件ID列表',
    layout_info JSON COMMENT '布局信息',
    area_info JSON COMMENT '面积信息',
    style_info JSON COMMENT '风格信息',
    material_info JSON COMMENT '材料信息',
    status VARCHAR(50) DEFAULT 'draft' COMMENT '状态',
    created_by BIGINT UNSIGNED COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计版本表';

-- 设计版本对比表
CREATE TABLE IF NOT EXISTS design_version_compares (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    version_a_id BIGINT UNSIGNED NOT NULL COMMENT '版本A ID',
    version_b_id BIGINT UNSIGNED NOT NULL COMMENT '版本B ID',
    layout_changes JSON COMMENT '布局变化',
    area_changes JSON COMMENT '面积调整',
    element_changes JSON COMMENT '元素增减',
    style_changes JSON COMMENT '风格变化',
    material_changes JSON COMMENT '材料变化',
    summary TEXT COMMENT '对比摘要',
    compare_status VARCHAR(50) DEFAULT 'pending' COMMENT '对比状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_version_a_id (version_a_id),
    INDEX idx_version_b_id (version_b_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计版本对比表';

-- 设计比对结果表
CREATE TABLE IF NOT EXISTS design_compare_results (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    name VARCHAR(200) DEFAULT '' COMMENT '比对任务名称',
    document_ids JSON COMMENT '关联需求文档ID列表',
    design_images JSON COMMENT '设计图列表',
    document_id BIGINT UNSIGNED COMMENT '关联需求文档ID(兼容旧数据)',
    design_image_path VARCHAR(500) COMMENT '设计图路径(兼容旧数据)',
    match_items JSON COMMENT '匹配项',
    deviation_items JSON COMMENT '偏差项',
    suggestions JSON COMMENT '建议项',
    overall_score DECIMAL(5,2) DEFAULT 0 COMMENT '整体匹配度(0-100)',
    analysis_status VARCHAR(50) DEFAULT 'pending' COMMENT '分析状态',
    error_message VARCHAR(1000) COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_document_id (document_id),
    INDEX idx_analysis_status (analysis_status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计比对结果表';

-- 施工图标注表
CREATE TABLE IF NOT EXISTS construction_annotations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    cad_file_id BIGINT UNSIGNED NOT NULL COMMENT '关联CAD文件ID',
    version_id BIGINT UNSIGNED COMMENT '关联设计版本ID',
    image_path VARCHAR(500) NOT NULL COMMENT '施工图图片路径',
    analysis_status VARCHAR(50) DEFAULT 'pending' COMMENT '分析状态: pending/processing/completed/failed',
    elements_detected JSON COMMENT '检测到的元素数据',
    annotations JSON COMMENT '标注数据',
    error_message VARCHAR(500) COMMENT '错误信息',
    created_by BIGINT UNSIGNED COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_cad_file_id (cad_file_id),
    INDEX idx_version_id (version_id),
    INDEX idx_analysis_status (analysis_status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='施工图标注表';

-- AI生成CAD任务表
CREATE TABLE IF NOT EXISTS cad_generations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    document_ids JSON COMMENT '关联文档ID列表',
    generation_type VARCHAR(50) NOT NULL COMMENT '生成类型: floor_plan/elevation/detail/ceiling/electric',
    prompt TEXT COMMENT '生成提示词',
    parameters JSON COMMENT '生成参数JSON',
    status VARCHAR(50) DEFAULT 'pending' COMMENT '状态: pending/processing/completed/failed',
    progress INT DEFAULT 0 COMMENT '进度百分比0-100',
    result_file_id BIGINT UNSIGNED COMMENT '生成的CAD文件ID',
    result_file_path VARCHAR(500) COMMENT '生成的文件路径',
    preview_image_url VARCHAR(500) COMMENT '预览图URL',
    error_message VARCHAR(1000) COMMENT '错误信息',
    ai_model VARCHAR(100) COMMENT '使用的AI模型',
    processing_time INT DEFAULT 0 COMMENT '处理耗时(秒)',
    created_by BIGINT UNSIGNED COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI生成CAD任务表';

-- 设计建议表
CREATE TABLE IF NOT EXISTS design_suggestions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    content TEXT NOT NULL COMMENT '建议内容',
    applicable_scene TEXT COMMENT '适用场景',
    cost_impact VARCHAR(200) COMMENT '成本影响',
    category VARCHAR(100) COMMENT '建议类别',
    priority INT DEFAULT 0 COMMENT '优先级',
    status VARCHAR(50) DEFAULT 'pending' COMMENT '状态: pending/adopted/ignored',
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

-- 建议生成任务表
CREATE TABLE IF NOT EXISTS suggestion_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    category VARCHAR(50) DEFAULT '' COMMENT '指定类别',
    count INT DEFAULT 5 COMMENT '请求生成数量',
    generated_count INT DEFAULT 0 COMMENT '实际生成数量',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态: pending/processing/completed/failed',
    error_message VARCHAR(500) DEFAULT '' COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='建议生成任务表';

-- 设计师AI对话消息表
CREATE TABLE IF NOT EXISTS designer_chat_messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(100) NOT NULL COMMENT '会话ID',
    project_id BIGINT UNSIGNED COMMENT '关联项目ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    role VARCHAR(20) NOT NULL COMMENT '角色: user/assistant',
    content TEXT NOT NULL COMMENT '消息内容',
    tokens_used INT DEFAULT 0 COMMENT 'Token消耗',
    response_time BIGINT DEFAULT 0 COMMENT '响应时间(毫秒)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_session_id (session_id),
    INDEX idx_project_id (project_id),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计师AI对话消息表';

-- 效果图生成记录表
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

-- 用户效果图配额表
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

-- 用户通知表
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    title VARCHAR(200) NOT NULL COMMENT '通知标题',
    content VARCHAR(500) COMMENT '通知内容',
    type VARCHAR(20) DEFAULT 'notice' COMMENT '类型: notice/message/email',
    is_read TINYINT(1) DEFAULT 0 COMMENT '是否已读',
    related_id BIGINT UNSIGNED COMMENT '关联ID',
    related_type VARCHAR(50) DEFAULT 'general' COMMENT '关联类型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户通知表';

-- 设计规范表
CREATE TABLE IF NOT EXISTS design_standards (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(100) NOT NULL COMMENT '规范编号',
    name VARCHAR(255) NOT NULL COMMENT '规范名称',
    category VARCHAR(100) NOT NULL COMMENT '规范类别: fire/accessibility/environmental/safety/other',
    content TEXT NOT NULL COMMENT '规范内容',
    applicable_types JSON COMMENT '适用项目类型',
    version VARCHAR(50) COMMENT '规范版本',
    effective_date DATE COMMENT '生效日期',
    source VARCHAR(255) COMMENT '来源',
    interpretation TEXT COMMENT '规范解读',
    status VARCHAR(50) DEFAULT 'active' COMMENT '状态: active/deprecated',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_code (code),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计规范表';

-- 聊天室表
CREATE TABLE IF NOT EXISTS chat_rooms (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '聊天室名称',
    type VARCHAR(20) DEFAULT 'group' COMMENT '类型: private/group',
    creator_id BIGINT UNSIGNED NOT NULL COMMENT '创建者ID',
    avatar VARCHAR(500) COMMENT '头像',
    description VARCHAR(500) COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_creator_id (creator_id),
    INDEX idx_type (type),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天室表';

-- 聊天消息表
CREATE TABLE IF NOT EXISTS chat_messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    room_id BIGINT UNSIGNED NOT NULL COMMENT '聊天室ID',
    sender_id BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
    content TEXT NOT NULL COMMENT '消息内容',
    type VARCHAR(20) DEFAULT 'text' COMMENT '消息类型: text/image/file',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_room_id (room_id),
    INDEX idx_sender_id (sender_id),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表';

-- 聊天室成员表
CREATE TABLE IF NOT EXISTS chat_room_members (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    room_id BIGINT UNSIGNED NOT NULL COMMENT '聊天室ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    role VARCHAR(20) DEFAULT 'member' COMMENT '角色: owner/admin/member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE INDEX idx_room_user (room_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天室成员表';
