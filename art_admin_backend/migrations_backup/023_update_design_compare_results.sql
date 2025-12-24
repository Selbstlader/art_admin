-- 更新设计比对结果表，支持多图片、多文档
-- Update design_compare_results table to support multiple images and documents

-- 添加新字段
ALTER TABLE `design_compare_results` 
ADD COLUMN `name` VARCHAR(200) DEFAULT '' COMMENT '比对任务名称' AFTER `project_id`,
ADD COLUMN `document_ids` JSON DEFAULT NULL COMMENT '关联需求文档ID列表(JSON数组)' AFTER `name`,
ADD COLUMN `design_images` JSON DEFAULT NULL COMMENT '设计图列表(JSON数组)' AFTER `document_ids`;

-- 迁移旧数据：将单个 document_id 转为数组
UPDATE `design_compare_results` 
SET `document_ids` = JSON_ARRAY(`document_id`)
WHERE `document_id` IS NOT NULL AND `document_id` > 0;

-- 迁移旧数据：将单个 design_image_path 转为数组
UPDATE `design_compare_results` 
SET `design_images` = JSON_ARRAY(JSON_OBJECT('fileName', '', 'filePath', `design_image_path`, 'fileType', 'image', 'fileSize', 0))
WHERE `design_image_path` IS NOT NULL AND `design_image_path` != '';

-- 扩大错误信息字段长度
ALTER TABLE `design_compare_results` 
MODIFY COLUMN `error_message` VARCHAR(1000) DEFAULT '' COMMENT '错误信息';

-- 注意：保留旧字段 document_id 和 design_image_path 以兼容旧数据
-- 后续可根据需要删除这些字段
