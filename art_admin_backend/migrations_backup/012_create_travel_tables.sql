-- 旅游规划模块数据库迁移
-- 创建时间: 2024-12-04

-- 路书表
CREATE TABLE IF NOT EXISTS roadbooks (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    title VARCHAR(100) NOT NULL COMMENT '路书标题',
    description TEXT COMMENT '路书描述',
    cover_url VARCHAR(500) COMMENT '封面图片URL',
    start_date DATE COMMENT '开始日期',
    end_date DATE COMMENT '结束日期',
    visibility TINYINT DEFAULT 1 COMMENT '可见性: 1-公开 2-私有 3-指定用户',
    travel_mode VARCHAR(20) DEFAULT 'driving' COMMENT '出行方式: driving/walking/transit/riding',
    total_distance INT DEFAULT 0 COMMENT '总距离(米)',
    total_duration INT DEFAULT 0 COMMENT '总时长(分钟)',
    total_budget DECIMAL(10,2) DEFAULT 0 COMMENT '总预算',
    view_count INT UNSIGNED DEFAULT 0 COMMENT '浏览量',
    favorite_count INT UNSIGNED DEFAULT 0 COMMENT '收藏量',
    like_count INT UNSIGNED DEFAULT 0 COMMENT '点赞量',
    comment_count INT UNSIGNED DEFAULT 0 COMMENT '评论量',
    share_count INT UNSIGNED DEFAULT 0 COMMENT '分享量',
    status TINYINT DEFAULT 1 COMMENT '状态: 1-草稿 2-已发布 3-已下架',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_visibility_status (visibility, status),
    INDEX idx_created_at (created_at),
    FULLTEXT INDEX ft_title_desc (title, description)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路书表';

-- 途经点表
CREATE TABLE IF NOT EXISTS roadbook_waypoints (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    roadbook_id BIGINT UNSIGNED NOT NULL COMMENT '路书ID',
    name VARCHAR(100) NOT NULL COMMENT '地点名称',
    address VARCHAR(300) COMMENT '详细地址',
    longitude DECIMAL(10, 7) NOT NULL COMMENT '经度',
    latitude DECIMAL(10, 7) NOT NULL COMMENT '纬度',
    poi_id VARCHAR(50) COMMENT '高德POI ID',
    poi_type VARCHAR(50) COMMENT 'POI类型',
    day_index INT DEFAULT 1 COMMENT '第几天',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    stay_duration INT DEFAULT 60 COMMENT '停留时长(分钟)',
    budget DECIMAL(10,2) DEFAULT 0 COMMENT '预算',
    notes TEXT COMMENT '备注',
    images JSON COMMENT '图片列表',
    waypoint_type TINYINT DEFAULT 2 COMMENT '类型: 1-起点 2-途经点 3-终点',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_roadbook_id (roadbook_id),
    INDEX idx_day_index (day_index),
    FOREIGN KEY (roadbook_id) REFERENCES roadbooks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='途经点表';

-- 标签表
CREATE TABLE IF NOT EXISTS roadbook_tags (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
    color VARCHAR(20) DEFAULT '#409EFF' COMMENT '标签颜色',
    use_count INT UNSIGNED DEFAULT 0 COMMENT '使用次数',
    is_system TINYINT DEFAULT 0 COMMENT '是否系统标签: 0-否 1-是',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_is_system (is_system)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- 标签关联表
CREATE TABLE IF NOT EXISTS roadbook_tag_relations (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    roadbook_id BIGINT UNSIGNED NOT NULL COMMENT '路书ID',
    tag_id BIGINT UNSIGNED NOT NULL COMMENT '标签ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_roadbook_tag (roadbook_id, tag_id),
    FOREIGN KEY (roadbook_id) REFERENCES roadbooks(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES roadbook_tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路书标签关联表';

-- 评论表
CREATE TABLE IF NOT EXISTS roadbook_comments (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    roadbook_id BIGINT UNSIGNED NOT NULL COMMENT '路书ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    parent_id BIGINT UNSIGNED DEFAULT 0 COMMENT '父评论ID',
    content VARCHAR(500) NOT NULL COMMENT '评论内容',
    like_count INT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_roadbook_id (roadbook_id),
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (roadbook_id) REFERENCES roadbooks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 收藏表
CREATE TABLE IF NOT EXISTS roadbook_favorites (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    roadbook_id BIGINT UNSIGNED NOT NULL COMMENT '路书ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_roadbook (user_id, roadbook_id),
    INDEX idx_user_id (user_id),
    INDEX idx_roadbook_id (roadbook_id),
    FOREIGN KEY (roadbook_id) REFERENCES roadbooks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

-- 分享记录表
CREATE TABLE IF NOT EXISTS roadbook_shares (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    roadbook_id BIGINT UNSIGNED NOT NULL COMMENT '路书ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    share_code VARCHAR(32) NOT NULL UNIQUE COMMENT '分享码',
    share_type TINYINT DEFAULT 1 COMMENT '分享类型: 1-链接 2-二维码',
    visit_count INT UNSIGNED DEFAULT 0 COMMENT '访问次数',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_roadbook_id (roadbook_id),
    INDEX idx_user_id (user_id),
    INDEX idx_share_code (share_code),
    FOREIGN KEY (roadbook_id) REFERENCES roadbooks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分享记录表';

-- 路书模板表
CREATE TABLE IF NOT EXISTS roadbook_templates (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    roadbook_id BIGINT UNSIGNED NOT NULL COMMENT '来源路书ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '创建者ID',
    title VARCHAR(100) NOT NULL COMMENT '模板标题',
    description TEXT COMMENT '模板描述',
    cover_url VARCHAR(500) COMMENT '封面图片URL',
    destination VARCHAR(100) COMMENT '目的地',
    days INT NOT NULL COMMENT '天数',
    category TINYINT DEFAULT 1 COMMENT '分类: 1-自然风光 2-人文历史 3-美食之旅 4-亲子游 5-其他',
    use_count INT UNSIGNED DEFAULT 0 COMMENT '使用次数',
    status TINYINT DEFAULT 0 COMMENT '状态: 0-待审核 1-已通过 2-已拒绝',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_category (category),
    INDEX idx_status (status),
    FOREIGN KEY (roadbook_id) REFERENCES roadbooks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路书模板表';

-- 统计日表
CREATE TABLE IF NOT EXISTS roadbook_stats_daily (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    roadbook_id BIGINT UNSIGNED NOT NULL COMMENT '路书ID',
    stat_date DATE NOT NULL COMMENT '统计日期',
    view_count INT UNSIGNED DEFAULT 0 COMMENT '浏览量',
    favorite_count INT UNSIGNED DEFAULT 0 COMMENT '收藏量',
    share_count INT UNSIGNED DEFAULT 0 COMMENT '分享量',
    UNIQUE KEY uk_roadbook_date (roadbook_id, stat_date),
    INDEX idx_stat_date (stat_date),
    FOREIGN KEY (roadbook_id) REFERENCES roadbooks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路书每日统计表';

-- 插入系统预设标签
INSERT INTO roadbook_tags (name, color, is_system) VALUES
('自驾游', '#409EFF', 1),
('徒步', '#67C23A', 1),
('骑行', '#E6A23C', 1),
('亲子游', '#F56C6C', 1),
('情侣游', '#FF69B4', 1),
('美食之旅', '#FF8C00', 1),
('自然风光', '#32CD32', 1),
('人文历史', '#8B4513', 1),
('海岛度假', '#00CED1', 1),
('城市探索', '#9370DB', 1);
