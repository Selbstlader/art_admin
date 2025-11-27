-- 创建APP用户表
CREATE TABLE IF NOT EXISTS `app_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_name` varchar(50) NOT NULL COMMENT '用户账号',
  `nick_name` varchar(50) NOT NULL COMMENT '用户昵称',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像',
  `user_type` char(1) NOT NULL DEFAULT '2' COMMENT '用户类型 1-管理员 2-普通用户',
  `status` char(1) NOT NULL DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `create_by` varchar(50) DEFAULT NULL COMMENT '创建人',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_by` varchar(50) DEFAULT NULL COMMENT '更新人',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_name` (`user_name`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_user_type` (`user_type`),
  KEY `idx_status` (`status`),
  KEY `idx_last_login_time` (`last_login_time`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='APP用户表';

-- 插入测试数据（可选，使用INSERT IGNORE避免重复数据）
-- 密码为: 123456 (已加密)
INSERT IGNORE INTO `app_user` (`user_name`, `nick_name`, `password`, `phone`, `email`, `user_type`, `status`) VALUES
('admin_app', 'APP管理员', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '13800138000', 'admin@app.com', '1', '1'),
('user001', '张三', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '13800138001', 'user001@app.com', '2', '1'),
('user002', '李四', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '13800138002', 'user002@app.com', '2', '1');
