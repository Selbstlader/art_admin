-- AI 学习系统数据库表结构
-- 创建时间: 2024-11-21

-- 学科表
CREATE TABLE IF NOT EXISTS `subjects` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL COMMENT '学科名称',
  `icon` VARCHAR(200) COMMENT '学科图标',
  `description` TEXT COMMENT '学科描述',
  `grade_levels` JSON COMMENT '适用年级 ["小学","初中","高中"]',
  `sort` INT DEFAULT 0 COMMENT '排序',
  `status` TINYINT DEFAULT 1 COMMENT '状态 1-启用 0-禁用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_status` (`status`),
  INDEX `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学科表';

-- 教材内容表
CREATE TABLE IF NOT EXISTS `learning_materials` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '创建用户ID',
  `subject_id` BIGINT NOT NULL COMMENT '学科ID',
  `title` VARCHAR(200) NOT NULL COMMENT '教材标题',
  `topic` VARCHAR(200) NOT NULL COMMENT '学习题材',
  `grade` VARCHAR(50) NOT NULL COMMENT '年级',
  `difficulty` TINYINT DEFAULT 1 COMMENT '难度 1-基础 2-进阶 3-高级',
  `summary` TEXT COMMENT '内容概要',
  `content` LONGTEXT NOT NULL COMMENT '教材内容(JSON格式)',
  `audio_url` VARCHAR(500) COMMENT '音频URL',
  `audio_duration` INT COMMENT '音频时长(秒)',
  `total_time` INT COMMENT '预计学习时长(分钟)',
  `view_count` INT DEFAULT 0 COMMENT '浏览次数',
  `favorite_count` INT DEFAULT 0 COMMENT '收藏次数',
  `status` TINYINT DEFAULT 1 COMMENT '状态 1-正常 0-已删除',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_subject_id` (`subject_id`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教材内容表';

-- 教材章节表
CREATE TABLE IF NOT EXISTS `material_sections` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `material_id` BIGINT NOT NULL COMMENT '教材ID',
  `type` VARCHAR(20) NOT NULL COMMENT '类型 knowledge|example|exercise',
  `title` VARCHAR(200) COMMENT '章节标题',
  `content` LONGTEXT COMMENT '章节内容(JSON格式)',
  `sort` INT DEFAULT 0 COMMENT '排序',
  `audio_url` VARCHAR(500) COMMENT '章节音频URL',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_material_id` (`material_id`),
  INDEX `idx_type` (`type`),
  INDEX `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教材章节表';

-- 音频缓存表
CREATE TABLE IF NOT EXISTS `audio_cache` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `content_hash` VARCHAR(64) NOT NULL COMMENT '内容MD5哈希',
  `text_content` TEXT COMMENT '原始文本内容',
  `audio_url` VARCHAR(500) NOT NULL COMMENT '音频URL',
  `voice_type` VARCHAR(50) DEFAULT 'xiaoyun' COMMENT '音色类型',
  `language` VARCHAR(20) DEFAULT 'zh-CN' COMMENT '语言',
  `file_size` BIGINT COMMENT '文件大小(字节)',
  `duration` INT COMMENT '时长(秒)',
  `access_count` INT DEFAULT 0 COMMENT '访问次数',
  `last_access_at` TIMESTAMP COMMENT '最后访问时间',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_content_hash` (`content_hash`),
  INDEX `idx_last_access` (`last_access_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='音频缓存表';

-- 插入默认学科数据
INSERT INTO `subjects` (`name`, `icon`, `description`, `grade_levels`, `sort`, `status`) VALUES
('语文', '📚', '中文语言文学学科', '["小学","初中","高中"]', 1, 1),
('数学', '🔢', '数学逻辑与计算学科', '["小学","初中","高中"]', 2, 1),
('英语', '🔤', '英语语言学科', '["小学","初中","高中"]', 3, 1),
('物理', '⚛️', '物理学科', '["初中","高中"]', 4, 1),
('化学', '🧪', '化学学科', '["初中","高中"]', 5, 1),
('生物', '🧬', '生物学科', '["初中","高中"]', 6, 1),
('历史', '📜', '历史学科', '["初中","高中"]', 7, 1),
('地理', '🌍', '地理学科', '["初中","高中"]', 8, 1),
('政治', '⚖️', '政治学科', '["初中","高中"]', 9, 1);
