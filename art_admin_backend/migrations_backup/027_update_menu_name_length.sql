-- 更新菜单表name字段长度
-- 将name字段从varchar(50)扩展到varchar(100)以支持更长的路由名称

ALTER TABLE sys_menu MODIFY COLUMN name varchar(100) NOT NULL COMMENT '路由名称/权限标识';