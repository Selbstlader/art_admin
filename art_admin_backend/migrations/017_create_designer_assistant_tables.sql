-- 工装设计师AI辅助系统数据表
-- Designer AI Assistant System Database Tables
-- Requirements: 4.1, 4.2, 6.5, 8.1

-- 设计项目表 / Designer Projects Table
CREATE TABLE IF NOT EXISTS designer_projects (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL COMMENT '项目名称',
    description TEXT COMMENT '项目描述',
    area DECIMAL(15,2) DEFAULT 0 COMMENT '面积(平方米)',
    budget DECIMAL(15,2) DEFAULT 0 COMMENT '预算',
    style VARCHAR(100) COMMENT '设计风格',
    status VARCHAR(50) DEFAULT 'draft' COMMENT '状态(draft/in_progress/completed/archived)',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工装设计项目表';

-- 项目文档表 / Project Documents Table
CREATE TABLE IF NOT EXISTS project_documents (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    file_path VARCHAR(500) NOT NULL COMMENT '存储路径',
    file_type VARCHAR(50) COMMENT '文件类型(pdf/word/image)',
    file_size BIGINT DEFAULT 0 COMMENT '文件大小(字节)',
    analysis_status VARCHAR(50) DEFAULT 'pending' COMMENT '分析状态(pending/processing/completed/failed)',
    keywords JSON COMMENT '提取的关键字',
    summary TEXT COMMENT '文档摘要',
    project_name VARCHAR(200) COMMENT '提取的项目名称',
    extracted_area DECIMAL(15,2) DEFAULT 0 COMMENT '提取的面积',
    extracted_budget DECIMAL(15,2) DEFAULT 0 COMMENT '提取的预算',
    extracted_style VARCHAR(100) COMMENT '提取的风格',
    functional_zones JSON COMMENT '功能分区',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_analysis_status (analysis_status),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_documents_project FOREIGN KEY (project_id) REFERENCES designer_projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目文档表';

-- CAD文件表 / CAD Files Table
CREATE TABLE IF NOT EXISTS cad_files (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    original_path VARCHAR(500) NOT NULL COMMENT '原始文件路径',
    parsed_path VARCHAR(500) COMMENT '解析后数据路径',
    file_format VARCHAR(20) COMMENT '文件格式(dwg/dxf)',
    parse_status VARCHAR(50) DEFAULT 'pending' COMMENT '解析状态(pending/processing/completed/failed)',
    layer_count INT DEFAULT 0 COMMENT '图层数量',
    layers JSON COMMENT '图层列表',
    has_3d BOOLEAN DEFAULT FALSE COMMENT '是否包含3D信息',
    file_size BIGINT DEFAULT 0 COMMENT '文件大小(字节)',
    error_message VARCHAR(500) COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_parse_status (parse_status),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_cad_files_project FOREIGN KEY (project_id) REFERENCES designer_projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CAD文件表';

-- 材料库表 / Materials Table
CREATE TABLE IF NOT EXISTS materials (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL COMMENT '材料名称',
    category VARCHAR(100) COMMENT '分类',
    specification VARCHAR(200) COMMENT '规格',
    unit VARCHAR(50) COMMENT '单位(m²/m/个/kg)',
    unit_price DECIMAL(15,2) DEFAULT 0 COMMENT '单价',
    brand VARCHAR(100) COMMENT '品牌',
    supplier VARCHAR(200) COMMENT '供应商',
    description TEXT COMMENT '描述',
    image_url VARCHAR(500) COMMENT '图片URL',
    applicable_scenes JSON COMMENT '适用场景',
    status VARCHAR(50) DEFAULT 'active' COMMENT '状态(active/inactive)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_name (name),
    INDEX idx_category (category),
    INDEX idx_brand (brand),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='材料库表';

-- 成本估算表 / Cost Estimates Table
CREATE TABLE IF NOT EXISTS cost_estimates (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    material_cost DECIMAL(15,2) DEFAULT 0 COMMENT '材料费',
    labor_cost DECIMAL(15,2) DEFAULT 0 COMMENT '人工费',
    equipment_cost DECIMAL(15,2) DEFAULT 0 COMMENT '设备费',
    management_cost DECIMAL(15,2) DEFAULT 0 COMMENT '管理费',
    total_cost DECIMAL(15,2) DEFAULT 0 COMMENT '总计',
    budget_limit DECIMAL(15,2) DEFAULT 0 COMMENT '预算上限',
    items JSON COMMENT '明细项',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE INDEX idx_project_id (project_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_cost_estimates_project FOREIGN KEY (project_id) REFERENCES designer_projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成本估算表';

-- 设计比对结果表 / Design Compare Results Table
CREATE TABLE IF NOT EXISTS design_compare_results (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    document_id BIGINT UNSIGNED COMMENT '关联需求文档ID',
    design_image_path VARCHAR(500) NOT NULL COMMENT '设计图路径',
    match_items JSON COMMENT '匹配项',
    deviation_items JSON COMMENT '偏差项',
    suggestions JSON COMMENT '建议项',
    overall_score DECIMAL(5,2) DEFAULT 0 COMMENT '整体匹配度(0-100)',
    analysis_status VARCHAR(50) DEFAULT 'pending' COMMENT '分析状态',
    error_message VARCHAR(500) COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_document_id (document_id),
    INDEX idx_analysis_status (analysis_status),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_compare_results_project FOREIGN KEY (project_id) REFERENCES designer_projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_compare_results_document FOREIGN KEY (document_id) REFERENCES project_documents(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计比对结果表';
