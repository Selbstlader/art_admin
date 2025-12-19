-- 025_1_create_design_version_tables.sql
-- 创建设计版本相关表 / Create design version related tables

-- 1. 创建设计版本表 / Create design versions table
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
    
    INDEX idx_design_versions_project_id (project_id),
    INDEX idx_design_versions_status (status),
    INDEX idx_design_versions_created_by (created_by),
    INDEX idx_design_versions_deleted_at (deleted_at),
    
    FOREIGN KEY (project_id) REFERENCES designer_projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设计版本表';

-- 2. 创建版本对比表 / Create design version compare table
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
    
    INDEX idx_design_version_compares_project_id (project_id),
    INDEX idx_design_version_compares_version_a_id (version_a_id),
    INDEX idx_design_version_compares_version_b_id (version_b_id),
    INDEX idx_design_version_compares_deleted_at (deleted_at),
    
    FOREIGN KEY (project_id) REFERENCES designer_projects(id) ON DELETE CASCADE,
    FOREIGN KEY (version_a_id) REFERENCES design_versions(id) ON DELETE CASCADE,
    FOREIGN KEY (version_b_id) REFERENCES design_versions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设计版本对比表';