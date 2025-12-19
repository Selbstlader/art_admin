-- 创建施工图标注表
-- Create construction annotations table

CREATE TABLE IF NOT EXISTS construction_annotations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    cad_file_id BIGINT UNSIGNED NOT NULL COMMENT '关联CAD文件ID',
    image_path VARCHAR(500) NOT NULL COMMENT '施工图图片路径',
    analysis_status VARCHAR(50) NOT NULL DEFAULT 'pending' COMMENT '分析状态: pending/processing/completed/failed',
    elements_detected JSON COMMENT '检测到的元素数据',
    annotations JSON COMMENT '标注数据',
    error_message VARCHAR(500) COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_project_id (project_id),
    INDEX idx_cad_file_id (cad_file_id),
    INDEX idx_analysis_status (analysis_status),
    INDEX idx_deleted_at (deleted_at),
    
    FOREIGN KEY (project_id) REFERENCES designer_projects(id) ON DELETE CASCADE,
    FOREIGN KEY (cad_file_id) REFERENCES cad_files(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='施工图标注表';