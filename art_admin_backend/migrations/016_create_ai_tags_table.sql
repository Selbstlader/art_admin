-- AI标签管理数据库迁移
-- 创建时间: 2024-12-09
-- 功能: 用于管理Dify知识库配置和AI提示词的标签系统

-- AI标签表
CREATE TABLE IF NOT EXISTS ai_tags (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT 'AI标签ID',
    name VARCHAR(100) NOT NULL UNIQUE COMMENT '标签名称',
    description VARCHAR(500) DEFAULT '' COMMENT '标签描述',
    knowledge_base_id VARCHAR(100) DEFAULT '' COMMENT 'Dify知识库ID（从API获取）',
    knowledge_base_name VARCHAR(200) DEFAULT '' COMMENT '知识库名称（冗余存储）',
    system_prompt TEXT NOT NULL COMMENT '系统提示词',
    chat_api_key VARCHAR(200) DEFAULT '' COMMENT 'Dify Chat App API Key（手动输入）',
    status TINYINT DEFAULT 1 COMMENT '状态: 1=active, 0=inactive',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI标签表';
