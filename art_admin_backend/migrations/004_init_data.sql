-- 设计师AI助手系统 - 初始化数据
-- Designer AI Assistant System - Initial Data

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =====================================================
-- 1. 菜单数据 / Menu Data
-- =====================================================

-- 清空菜单表重新初始化
TRUNCATE TABLE sys_menu;

INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `path`, `component`, `redirect`, `title`, `icon`, `is_hide`, `is_hide_tab`, `is_iframe`, `link`, `keep_alive`, `fixed_tab`, `is_full_page`, `sort`, `is_enable`, `create_time`, `update_time`, `deleted_at`, `is_menu`, `show_badge`, `show_text_badge`, `active_path`) VALUES 
-- 仪表盘
(1, 0, 'Dashboard', '/dashboard', '/index/index', '/dashboard/console', '仪表盘', 'dashboard', 0, 0, 0, '', 1, 0, 0, 1, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(11, 1, 'DashboardConsole', '/dashboard/console', '/dashboard/console', '', '主控台', 'console', 0, 0, 0, '', 1, 0, 0, 1, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(12, 1, 'DashboardAnalysis', '/dashboard/analysis', '/dashboard/analysis', '', '分析页', 'analysis', 0, 0, 0, '', 1, 0, 0, 2, 1, NOW(), NOW(), NULL, 1, 0, '', ''),

-- 系统管理
(2, 0, 'System', '/system', '/index/index', '/system/user', '系统管理', 'setting', 0, 0, 0, '', 1, 0, 0, 99, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(21, 2, 'SystemUser', '/system/user', '/system/user', '', '用户管理', 'user', 0, 0, 0, '', 1, 0, 0, 1, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(22, 2, 'SystemRole', '/system/role', '/system/role', '', '角色管理', 'role', 0, 0, 0, '', 1, 0, 0, 2, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(23, 2, 'SystemMenu', '/system/menu', '/system/menu', '', '菜单管理', 'menu', 0, 0, 0, '', 1, 0, 0, 3, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(28, 2, 'SystemDictionary', '/system/dictionary', '/system/dictionary', '', '字典管理', '', 0, 0, 0, '', 1, 0, 0, 4, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(29, 2, 'SystemDepartment', '/system/department', '/system/department', '', '部门管理', '', 0, 0, 0, '', 1, 0, 0, 5, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(30, 2, 'SystemOperationLog', '/system/operation-log', '/system/operation-log', '', '日志管理', '', 0, 0, 0, '', 1, 0, 0, 6, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(39, 2, 'SystemAppUser', '/system/app-user', '/system/app-user', '', 'APP用户管理', '', 0, 0, 0, '', 1, 0, 0, 7, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(49, 2, 'SystemUserCenter', '/system/user-center', '/system/user-center', '', '个人中心', '', 1, 0, 0, '', 1, 0, 0, 8, 1, NOW(), NOW(), NULL, 1, 0, '', ''),

-- 设计师助手（核心功能）
(1023, 0, 'Designer', '/designer', '/index/index', '', '设计师助手', 'design', 0, 0, 0, '', 1, 0, 0, 2, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1024, 1023, 'DesignerProjectList', '/designer-assistant/project/ProjectList', '/designer-assistant/project/ProjectList', '', '项目管理', '', 0, 0, 0, '', 1, 0, 0, 1, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1025, 1023, 'DesignerProjectCreate', '/designer-assistant/project/ProjectCreate', '/designer-assistant/project/ProjectCreate', '', '创建项目', '', 1, 0, 0, '', 1, 0, 0, 2, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1026, 1023, 'DesignerProjectDetail', '/designer-assistant/project/ProjectDetail', '/designer-assistant/project/ProjectDetail', '', '项目详情', '', 1, 0, 0, '', 1, 0, 0, 3, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1027, 1023, 'DesignerDocumentUpload', '/designer-assistant/document/DocumentUpload', '/designer-assistant/document/DocumentUpload', '', '文档分析', '', 1, 0, 0, '', 1, 0, 0, 4, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1028, 1023, 'DesignerCompareUpload', '/designer-assistant/design-compare/CompareUpload', '/designer-assistant/design-compare/CompareUpload', '', '设计比对', '', 1, 0, 0, '', 1, 0, 0, 5, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1029, 1023, 'DesignerCadUpload', '/designer-assistant/cad-viewer/CadUpload', '/designer-assistant/cad-viewer/CadUpload', '', 'CAD预览', '', 1, 0, 0, '', 1, 0, 0, 6, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1030, 1023, 'DesignerMaterialList', '/designer-assistant/material/MaterialList', '/designer-assistant/material/MaterialList', '', '材料库', '', 0, 0, 0, '', 1, 0, 0, 7, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1031, 1023, 'DesignerCostConfig', '/designer-assistant/cost/CostConfig', '/designer-assistant/cost/CostConfig', '', '成本估算', '', 1, 0, 0, '', 1, 0, 0, 8, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1032, 1023, 'DesignerChat', '/designer-assistant/chat/DesignerChat', '/designer-assistant/chat/DesignerChat', '', 'AI对话', '', 0, 0, 0, '', 1, 0, 0, 9, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1033, 1023, 'DesignerComplianceCheck', '/designer-assistant/compliance/ComplianceCheck', '/designer-assistant/compliance/ComplianceCheck', '', '合规检查', '', 1, 0, 0, '', 1, 0, 0, 10, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1034, 1023, 'DesignerVersionList', '/designer-assistant/version-compare/VersionList', '/designer-assistant/version-compare/VersionList', '', '版本对比', '', 1, 0, 0, '', 1, 0, 0, 11, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1035, 1023, 'DesignerVersionDiff', '/designer-assistant/version-compare/VersionDiff', '/designer-assistant/version-compare/VersionDiff', '', '版本差异', '', 1, 0, 0, '', 1, 0, 0, 12, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1037, 1023, 'DesignerCadGeneration', '/designer-assistant/cad-generation/CadGenerationList', '/designer-assistant/cad-generation/CadGenerationList', '', 'CAD生成', '', 1, 0, 0, '', 1, 0, 0, 13, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(1039, 1023, 'DesignerCadViewer', '/designer-assistant/cad-viewer/CadViewer', '/designer-assistant/cad-viewer/CadViewer', '', 'CAD查看', '', 1, 0, 0, '', 1, 0, 0, 14, 1, NOW(), NOW(), NULL, 1, 0, '', ''),

-- Dify管理
(36, 0, 'Dify', '/dify', '/index/index', '', 'Dify管理', '', 0, 0, 0, '', 1, 0, 0, 3, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(37, 36, 'DifyChat', '/dify/chat', '/dify/chat', '', 'AI对话', '', 0, 0, 0, '', 1, 0, 0, 1, 1, NOW(), NOW(), NULL, 1, 0, '', ''),
(38, 36, 'DifyKnowledge', '/dify/knowledge', '/dify/knowledge', '', '知识库管理', '', 0, 0, 0, '', 1, 0, 0, 2, 1, NOW(), NOW(), NULL, 1, 0, '', '');

-- =====================================================
-- 2. 角色数据 / Role Data
-- =====================================================

INSERT INTO sys_role (name, code, description, status, sort) VALUES
('超级管理员', 'super_admin', '拥有所有权限', 1, 1),
('设计师', 'designer', '设计师角色', 1, 2),
('项目经理', 'project_manager', '项目管理角色', 1, 3)
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- =====================================================
-- 3. 管理员用户 / Admin User (密码: 123456)
-- =====================================================

INSERT INTO sys_user (username, password, nickname, email, status) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '系统管理员', 'admin@example.com', 1)
ON DUPLICATE KEY UPDATE nickname = VALUES(nickname);

-- =====================================================
-- 4. APP测试用户 / APP Test Users (密码: 123456)
-- =====================================================

INSERT INTO app_user (user_name, nick_name, password, phone, email, user_type, status) VALUES
('designer001', '设计师小王', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '13800138001', 'designer001@example.com', '2', '1'),
('designer002', '设计师小李', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '13800138002', 'designer002@example.com', '2', '1')
ON DUPLICATE KEY UPDATE nick_name = VALUES(nick_name);

-- =====================================================
-- 5. 设计规范数据 / Design Standards Data
-- =====================================================

INSERT INTO design_standards (code, name, category, content, applicable_types, version, effective_date, source, interpretation, status) VALUES
-- 消防规范
('GB50016-2014', '建筑设计防火规范', 'fire', '建筑物的耐火等级、防火分区、安全疏散、消防设施等应符合本规范要求', '["office","commercial","industrial","public"]', '2018修订版', '2018-10-01', '国家标准', '本规范是建筑防火设计的基本依据', 'active'),
('GB50016-5.5.17', '疏散走道宽度要求', 'fire', '公共建筑内疏散走道的净宽度不应小于1.1m；人员密集场所不应小于1.4m', '["commercial","public"]', '2018修订版', '2018-10-01', '国家标准', '疏散走道是人员疏散的主要通道', 'active'),
('GB50016-5.5.21', '安全出口数量要求', 'fire', '公共建筑每个防火分区或一个防火分区的每个楼层，其安全出口的数量应经计算确定，且不应少于2个', '["commercial","public"]', '2018修订版', '2018-10-01', '国家标准', '多个安全出口可确保在一个出口被封堵时仍有其他疏散途径', 'active'),
-- 无障碍规范
('GB50763-2012', '无障碍设计规范', 'accessibility', '公共建筑应设置无障碍通道、无障碍电梯、无障碍卫生间等设施', '["office","commercial","public"]', '2012版', '2012-09-01', '国家标准', '保障残疾人、老年人等特殊群体的通行和使用需求', 'active'),
('GB50763-3.3.1', '无障碍坡道坡度要求', 'accessibility', '无障碍坡道的坡度不应大于1:12，困难情况下不应大于1:8', '["office","commercial","public"]', '2012版', '2012-09-01', '国家标准', '坡度过大会增加轮椅使用者的通行难度和安全风险', 'active'),
-- 环保规范
('GB50325-2020', '民用建筑工程室内环境污染控制标准', 'environmental', '室内装饰装修材料的甲醛、苯、TVOC等有害物质释放量应符合本标准限值', '["office","commercial","public"]', '2020版', '2020-08-01', '国家标准', '控制室内空气污染，保障人员健康', 'active'),
-- 安全规范
('GB50352-2019', '民用建筑设计统一标准', 'safety', '建筑设计应满足结构安全、使用安全、防火安全等基本要求', '["office","commercial","industrial","public"]', '2019版', '2019-10-01', '国家标准', '民用建筑设计的基本准则和通用要求', 'active'),
('GB50352-6.6.3', '栏杆高度要求', 'safety', '临空高度在24m以下时，栏杆高度不应低于1.05m；临空高度在24m及以上时，栏杆高度不应低于1.10m', '["office","commercial","public"]', '2019版', '2019-10-01', '国家标准', '防止人员坠落，保障使用安全', 'active')
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- =====================================================
-- 6. 材料库示例数据 / Sample Materials Data
-- =====================================================

INSERT INTO materials (name, category, specification, unit, unit_price, brand, description, status) VALUES
('乳胶漆', '涂料', '内墙乳胶漆', 'm²', 35.00, '立邦', '环保内墙涂料，VOC含量低', 'active'),
('强化复合地板', '地板', '12mm厚', 'm²', 120.00, '圣象', 'E0级环保强化地板', 'active'),
('石膏板', '板材', '9.5mm厚', 'm²', 28.00, '可耐福', '纸面石膏板，用于吊顶隔墙', 'active'),
('轻钢龙骨', '龙骨', '50型', 'm', 8.50, '龙牌', '镀锌轻钢龙骨', 'active'),
('PVC地板', '地板', '2.0mm厚', 'm²', 65.00, 'LG', '商用PVC卷材地板', 'active'),
('铝扣板', '吊顶', '600x600mm', 'm²', 85.00, '欧普', '集成吊顶铝扣板', 'active'),
('玻璃隔断', '隔断', '12mm钢化玻璃', 'm²', 380.00, '信义', '办公室玻璃隔断', 'active'),
('办公地毯', '地毯', '方块地毯', 'm²', 95.00, 'Interface', '商用方块地毯', 'active')
ON DUPLICATE KEY UPDATE name = VALUES(name);

SET FOREIGN_KEY_CHECKS = 1;
