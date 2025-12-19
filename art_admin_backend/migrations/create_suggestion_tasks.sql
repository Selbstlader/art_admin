-- 创建建议生成任务表 / Create suggestion tasks table
CREATE TABLE IF NOT EXISTS `suggestion_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint unsigned NOT NULL COMMENT '关联项目ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `category` varchar(50) DEFAULT '' COMMENT '指定类别',
  `count` int DEFAULT 5 COMMENT '请求生成数量',
  `generated_count` int DEFAULT 0 COMMENT '实际生成数量',
  `status` varchar(20) DEFAULT 'pending' COMMENT '状态: pending/processing/completed/failed',
  `error_message` varchar(500) DEFAULT '' COMMENT '错误信息',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_suggestion_tasks_project_id` (`project_id`),
  KEY `idx_suggestion_tasks_user_id` (`user_id`),
  KEY `idx_suggestion_tasks_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计建议生成任务表';
