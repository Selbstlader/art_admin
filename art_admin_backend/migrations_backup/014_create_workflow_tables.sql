-- OA工作流引擎数据库迁移
-- 创建时间: 2024-12-05

-- 表单模板表
CREATE TABLE IF NOT EXISTS wf_form_template (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '模板名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '模板编码',
    description VARCHAR(500) COMMENT '模板描述',
    schema_json TEXT NOT NULL COMMENT '表单结构JSON',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态: active/disabled',
    created_by BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_status (status),
    INDEX idx_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='表单模板表';

-- 流程定义表
CREATE TABLE IF NOT EXISTS wf_process_definition (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '流程名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '流程编码',
    description VARCHAR(500) COMMENT '流程描述',
    category VARCHAR(50) COMMENT '流程分类',
    form_template_id BIGINT UNSIGNED COMMENT '关联的表单模板ID',
    graph_json TEXT NOT NULL COMMENT '流程图JSON',
    version INT DEFAULT 1 COMMENT '版本号',
    status VARCHAR(20) DEFAULT 'draft' COMMENT '状态: draft/published/disabled',
    created_by BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by),
    INDEX idx_form_template_id (form_template_id),
    FOREIGN KEY (form_template_id) REFERENCES wf_form_template(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='流程定义表';

-- 流程实例表
CREATE TABLE IF NOT EXISTS wf_process_instance (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    process_def_id BIGINT UNSIGNED NOT NULL COMMENT '流程定义ID',
    process_def_version INT NOT NULL COMMENT '流程定义版本',
    title VARCHAR(200) NOT NULL COMMENT '流程标题',
    initiator_id BIGINT UNSIGNED NOT NULL COMMENT '发起人ID',
    initiator_name VARCHAR(50) COMMENT '发起人名称',
    status VARCHAR(20) DEFAULT 'running' COMMENT '状态: running/completed/rejected/withdrawn',
    current_node_id VARCHAR(50) COMMENT '当前节点ID',
    form_data TEXT COMMENT '表单数据JSON',
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    INDEX idx_process_def_id (process_def_id),
    INDEX idx_initiator_id (initiator_id),
    INDEX idx_status (status),
    INDEX idx_started_at (started_at),
    FOREIGN KEY (process_def_id) REFERENCES wf_process_definition(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='流程实例表';

-- 工作流任务表
CREATE TABLE IF NOT EXISTS wf_task (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    process_inst_id BIGINT UNSIGNED NOT NULL COMMENT '流程实例ID',
    node_id VARCHAR(50) NOT NULL COMMENT '节点ID',
    node_name VARCHAR(100) COMMENT '节点名称',
    assignee_id BIGINT UNSIGNED NOT NULL COMMENT '审批人ID',
    assignee_name VARCHAR(50) COMMENT '审批人名称',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态: pending/approved/rejected/delegated/transferred',
    comment VARCHAR(500) COMMENT '审批意见',
    delegated_from BIGINT UNSIGNED COMMENT '委托来源用户ID',
    transferred_from BIGINT UNSIGNED COMMENT '转办来源用户ID',
    due_at TIMESTAMP NULL COMMENT '截止时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    INDEX idx_process_inst_id (process_inst_id),
    INDEX idx_assignee_id (assignee_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_due_at (due_at),
    FOREIGN KEY (process_inst_id) REFERENCES wf_process_instance(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流任务表';

-- 任务日志表
CREATE TABLE IF NOT EXISTS wf_task_log (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    process_inst_id BIGINT UNSIGNED NOT NULL COMMENT '流程实例ID',
    operator_id BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
    operator_name VARCHAR(50) COMMENT '操作人名称',
    action VARCHAR(20) NOT NULL COMMENT '操作类型: approve/reject/delegate/transfer/withdraw',
    comment VARCHAR(500) COMMENT '操作备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_task_id (task_id),
    INDEX idx_process_inst_id (process_inst_id),
    INDEX idx_operator_id (operator_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (task_id) REFERENCES wf_task(id) ON DELETE CASCADE,
    FOREIGN KEY (process_inst_id) REFERENCES wf_process_instance(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务日志表';

-- 插入示例表单模板
INSERT INTO wf_form_template (name, code, description, schema_json, status, created_by) VALUES
('请假申请表单', 'leave_form', '员工请假申请表单模板', '{"fields":[{"key":"leaveType","label":"请假类型","type":"select","required":true,"options":[{"label":"年假","value":"annual"},{"label":"事假","value":"personal"},{"label":"病假","value":"sick"}]},{"key":"startDate","label":"开始日期","type":"date","required":true},{"key":"endDate","label":"结束日期","type":"date","required":true},{"key":"days","label":"请假天数","type":"number","required":true,"validation":{"min":0.5,"max":30}},{"key":"reason","label":"请假原因","type":"textarea","required":true,"validation":{"minLength":10,"maxLength":500}}]}', 'active', 1),
('报销申请表单', 'expense_form', '费用报销申请表单模板', '{"fields":[{"key":"expenseType","label":"报销类型","type":"select","required":true,"options":[{"label":"差旅费","value":"travel"},{"label":"办公用品","value":"office"},{"label":"招待费","value":"entertainment"}]},{"key":"amount","label":"报销金额","type":"number","required":true,"validation":{"min":0.01}},{"key":"expenseDate","label":"费用发生日期","type":"date","required":true},{"key":"description","label":"费用说明","type":"textarea","required":true},{"key":"attachments","label":"附件","type":"file","required":false}]}', 'active', 1);
