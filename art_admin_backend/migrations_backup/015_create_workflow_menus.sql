-- 工作流菜单配置
-- 用于后端控制模式下的菜单数据
-- 注意：执行前请确保ID不与现有菜单冲突，可根据实际情况调整ID值

-- 工作流主菜单 (ID: 1000)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1000, 0, 'Workflow', '/workflow', '工作流', '/index/index', 'icon-workflow', 1, 50, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 审批工作台 (ID: 1001, parent: 1000)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1001, 1000, 'WorkflowWorkbench', 'workbench', '审批工作台', '/index/index', '', 1, 1, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 我的待办 (ID: 1002, parent: 1001)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1002, 1001, 'WorkflowTodo', 'todo', '我的待办', '/workflow/workbench/todo', '', 1, 1, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 我的已办 (ID: 1003, parent: 1001)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1003, 1001, 'WorkflowDone', 'done', '我的已办', '/workflow/workbench/done', '', 1, 2, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 我发起的 (ID: 1004, parent: 1001)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1004, 1001, 'WorkflowInitiated', 'initiated', '我发起的', '/workflow/workbench/initiated', '', 1, 3, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 审批详情（隐藏菜单）(ID: 1005, parent: 1001)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1005, 1001, 'WorkflowDetail', 'detail/:id', '审批详情', '/workflow/workbench/detail', '', 1, 4, 1, 0, 1)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 发起申请 (ID: 1006, parent: 1000)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1006, 1000, 'WorkflowApply', 'apply', '发起申请', '/workflow/apply/index', '', 1, 2, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 填写申请（隐藏菜单）(ID: 1007, parent: 1000)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1007, 1000, 'WorkflowApplyForm', 'apply/form/:defId', '填写申请', '/workflow/apply/form', '', 1, 3, 1, 0, 1)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 流程管理 (ID: 1008, parent: 1000)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1008, 1000, 'WorkflowManagement', 'management', '流程管理', '/index/index', '', 1, 4, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 流程定义 (ID: 1009, parent: 1008)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1009, 1008, 'ProcessDefList', 'process-def', '流程定义', '/workflow/process-def/index', '', 1, 1, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 流程设计（隐藏菜单）(ID: 1010, parent: 1008)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1010, 1008, 'ProcessDefEdit', 'process-def/edit/:id?', '流程设计', '/workflow/process-def/edit', '', 1, 2, 1, 0, 1)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 表单模板 (ID: 1011, parent: 1008)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1011, 1008, 'FormTemplateList', 'form-template', '表单模板', '/workflow/form-template/index', '', 1, 3, 1, 1, 0)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 表单设计（隐藏菜单）(ID: 1012, parent: 1008)
INSERT INTO sys_menu (id, parent_id, name, path, title, component, icon, is_enable, sort, is_menu, keep_alive, is_hide)
VALUES (1012, 1008, 'FormTemplateEdit', 'form-template/edit/:id?', '表单设计', '/workflow/form-template/edit', '', 1, 4, 1, 0, 1)
ON DUPLICATE KEY UPDATE title = VALUES(title);
