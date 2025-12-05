-- 创建文件管理表
CREATE TABLE IF NOT EXISTS `files` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `file_name` VARCHAR(255) NOT NULL COMMENT '原始文件名',
    `storage_name` VARCHAR(255) NOT NULL COMMENT '存储文件名',
    `file_path` VARCHAR(500) NOT NULL COMMENT '文件路径',
    `file_size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
    `file_type` VARCHAR(100) DEFAULT '' COMMENT '文件类型(MIME)',
    `file_ext` VARCHAR(20) DEFAULT '' COMMENT '文件扩展名',
    `category` VARCHAR(50) DEFAULT 'default' COMMENT '分类(image/document/video/audio/other)',
    `upload_user_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '上传用户ID',
    `url` VARCHAR(500) DEFAULT '' COMMENT '访问URL',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    INDEX `idx_category` (`category`),
    INDEX `idx_upload_user_id` (`upload_user_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件管理表';
