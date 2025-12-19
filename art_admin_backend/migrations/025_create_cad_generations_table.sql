-- 025_create_cad_generations_table.sql
-- AI生成CAD任务表 / CAD generation tasks table

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
    
    INDEX idx_cad_generations_project_id (project_id),
    INDEX idx_cad_generations_status (status),
    INDEX idx_cad_generations_created_by (created_by),
    INDEX idx_cad_generations_deleted_at (deleted_at),
    
    FOREIGN KEY (project_id) REFERENCES designer_projects(id) ON DELETE CASCADE,
    FOREIGN KEY (result_file_id) REFERENCES cad_files(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI生成CAD任务表';
