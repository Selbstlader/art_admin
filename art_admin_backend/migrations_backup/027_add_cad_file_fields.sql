-- 027_add_cad_file_fields.sql
-- 为CAD文件表添加新字段 / Add new fields to cad_files table

-- 添加文件路径字段
ALTER TABLE cad_files 
ADD COLUMN file_path VARCHAR(500) COMMENT '文件路径';

-- 添加AI生成标记字段
ALTER TABLE cad_files 
ADD COLUMN is_ai_generated BOOLEAN DEFAULT FALSE COMMENT '是否AI生成';
