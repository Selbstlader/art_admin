-- 用户通知表 / User notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    title VARCHAR(200) NOT NULL COMMENT '通知标题',
    content VARCHAR(500) COMMENT '通知内容',
    type VARCHAR(20) DEFAULT 'notice' COMMENT '类型: notice/message/email',
    is_read TINYINT(1) DEFAULT 0 COMMENT '是否已读',
    related_id BIGINT UNSIGNED COMMENT '关联ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户通知表';

-- 清理7天前的通知数据（可选，定期执行）
-- DELETE FROM notifications WHERE created_at < DATE_SUB(NOW(), INTERVAL 7 DAY);
