-- 创建教材生成任务表
CREATE TABLE IF NOT EXISTS `material_generate_tasks` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `subject_id` bigint(20) NOT NULL COMMENT '学科ID',
  `grade` varchar(50) NOT NULL COMMENT '年级',
  `topic` varchar(200) NOT NULL COMMENT '学习题材',
  `difficulty` tinyint(4) NOT NULL DEFAULT '1' COMMENT '难度 1-基础 2-进阶 3-高级',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0-待处理 1-生成中 2-已完成 3-失败',
  `material_id` bigint(20) DEFAULT NULL COMMENT '生成的教材ID',
  `error_msg` text COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教材生成任务表';
