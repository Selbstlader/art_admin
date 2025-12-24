-- 创建学科表
CREATE TABLE IF NOT EXISTS `subjects` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) NOT NULL COMMENT '学科名称',
  `icon` varchar(200) DEFAULT NULL COMMENT '学科图标',
  `description` text COMMENT '学科描述',
  `grade_levels` json DEFAULT NULL COMMENT '适用年级 ["小学","初中","高中"]',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态 1-启用 0-禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学科表';

-- 插入初始学科数据
INSERT INTO `subjects` (`name`, `icon`, `description`, `grade_levels`, `sort`, `status`) VALUES
('语文', 'book', '语文学科，包括阅读理解、作文写作等', '["小学","初中","高中"]', 1, 1),
('数学', 'calculator', '数学学科，包括算术、代数、几何等', '["小学","初中","高中"]', 2, 1),
('英语', 'language', '英语学科，包括听说读写等语言技能', '["小学","初中","高中"]', 3, 1),
('物理', 'atom', '物理学科，包括力学、电学、光学等', '["初中","高中"]', 4, 1),
('化学', 'flask', '化学学科，包括无机化学、有机化学等', '["初中","高中"]', 5, 1),
('生物', 'leaf', '生物学科，包括植物学、动物学、微生物学等', '["初中","高中"]', 6, 1),
('历史', 'clock', '历史学科，包括中国历史、世界历史等', '["初中","高中"]', 7, 1),
('地理', 'globe', '地理学科，包括自然地理、人文地理等', '["初中","高中"]', 8, 1),
('政治', 'balance', '政治学科，包括思想品德、政治理论等', '["初中","高中"]', 9, 1);
