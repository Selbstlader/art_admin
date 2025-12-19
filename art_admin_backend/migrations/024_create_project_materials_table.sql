-- 创建项目材料清单表 / Create project materials table
-- Migration: 024_create_project_materials_table.sql

CREATE TABLE IF NOT EXISTS project_materials (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL,
    material_id BIGINT UNSIGNED DEFAULT NULL,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(100) DEFAULT NULL,
    specification VARCHAR(200) DEFAULT NULL,
    unit VARCHAR(50) NOT NULL,
    unit_price DECIMAL(12, 2) DEFAULT 0,
    quantity DECIMAL(12, 2) DEFAULT 0,
    total_price DECIMAL(12, 2) DEFAULT 0,
    brand VARCHAR(100) DEFAULT NULL,
    supplier VARCHAR(200) DEFAULT NULL,
    remark VARCHAR(500) DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_project_materials_project_id (project_id),
    INDEX idx_project_materials_material_id (material_id),
    INDEX idx_project_materials_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目材料清单表';
