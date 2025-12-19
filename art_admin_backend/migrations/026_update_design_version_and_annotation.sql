-- 026_update_design_version_and_annotation.sql
-- 更新设计版本表和施工图标注表 / Update design_versions and construction_annotations tables

-- 为construction_annotations表添加version_id字段
ALTER TABLE construction_annotations 
ADD COLUMN version_id BIGINT UNSIGNED COMMENT '关联设计版本ID',
ADD COLUMN created_by BIGINT UNSIGNED COMMENT '创建者ID';

-- 添加外键约束
ALTER TABLE construction_annotations 
ADD CONSTRAINT fk_construction_annotations_version_id 
FOREIGN KEY (version_id) REFERENCES design_versions(id) ON DELETE SET NULL;

-- 添加索引
CREATE INDEX idx_construction_annotations_version_id ON construction_annotations(version_id);
CREATE INDEX idx_construction_annotations_created_by ON construction_annotations(created_by);
