-- 小红书总结功能相关表

-- 小红书笔记总结表
CREATE TABLE IF NOT EXISTS xhs_summaries (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    original_url VARCHAR(500) DEFAULT '' COMMENT '原始小红书链接',
    note_title VARCHAR(200) DEFAULT '' COMMENT '笔记标题',
    note_content TEXT COMMENT '笔记原文内容',
    image_urls JSON COMMENT '图片URL列表',
    summary_title VARCHAR(100) DEFAULT '' COMMENT '总结标题',
    summary_content TEXT COMMENT '总结内容',
    key_points JSON COMMENT '关键信息列表',
    tags JSON COMMENT '标签列表',
    sentiment VARCHAR(20) DEFAULT 'neutral' COMMENT '情感倾向: positive, neutral, negative',
    ocr_texts JSON COMMENT 'OCR识别的文字',
    analysis JSON COMMENT '专业分析结果',
    style VARCHAR(20) DEFAULT 'concise' COMMENT '总结风格: concise, detailed, casual',
    is_favorite TINYINT(1) DEFAULT 0 COMMENT '是否收藏: 0-否, 1-是',
    folder_id BIGINT DEFAULT 0 COMMENT '收藏夹ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_folder_id (folder_id),
    INDEX idx_is_favorite (is_favorite),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小红书笔记总结表';

-- 小红书收藏夹表
CREATE TABLE IF NOT EXISTS xhs_folders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    name VARCHAR(50) NOT NULL COMMENT '收藏夹名称',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小红书收藏夹表';
