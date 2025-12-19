-- 添加错误信息字段到项目文档表
-- Add error_message field to project_documents table

ALTER TABLE project_documents 
ADD COLUMN error_message VARCHAR(500) DEFAULT '' COMMENT '分析错误信息 / Analysis error message' 
AFTER analysis_status;
