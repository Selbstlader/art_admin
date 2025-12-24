-- 创建设计规范表 / Create design standards table
-- 存储各类设计规范条款（消防、无障碍、环保等）
CREATE TABLE IF NOT EXISTS design_standards (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(100) NOT NULL COMMENT '规范编号',
    name VARCHAR(255) NOT NULL COMMENT '规范名称',
    category VARCHAR(100) NOT NULL COMMENT '规范类别(fire/accessibility/environmental/safety/other)',
    content TEXT NOT NULL COMMENT '规范内容',
    applicable_types JSON COMMENT '适用项目类型(office/commercial/industrial/public)',
    version VARCHAR(50) COMMENT '规范版本',
    effective_date DATE COMMENT '生效日期',
    source VARCHAR(255) COMMENT '来源(国家标准/行业标准/地方标准)',
    interpretation TEXT COMMENT '规范解读',
    status VARCHAR(50) DEFAULT 'active' COMMENT '状态(active/deprecated)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_code (code),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设计规范表';

-- 创建合规检查结果表 / Create compliance check results table
CREATE TABLE IF NOT EXISTS compliance_check_results (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '关联项目ID',
    check_type VARCHAR(100) COMMENT '检查类型(full/partial)',
    passed_items JSON COMMENT '通过项列表',
    failed_items JSON COMMENT '不合规项列表',
    suggestions JSON COMMENT '改进建议列表',
    overall_score DECIMAL(5,2) DEFAULT 0 COMMENT '整体合规分数(0-100)',
    check_status VARCHAR(50) DEFAULT 'pending' COMMENT '检查状态(pending/processing/completed/failed)',
    error_message VARCHAR(500) COMMENT '错误信息',
    checked_standards JSON COMMENT '已检查的规范ID列表',
    user_id BIGINT UNSIGNED COMMENT '执行检查的用户ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_project_id (project_id),
    INDEX idx_check_status (check_status),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合规检查结果表';

-- 插入示例设计规范数据 / Insert sample design standards data
INSERT INTO design_standards (code, name, category, content, applicable_types, version, effective_date, source, interpretation, status) VALUES
-- 消防规范 / Fire safety standards
('GB50016-2014', '建筑设计防火规范', 'fire', '建筑物的耐火等级、防火分区、安全疏散、消防设施等应符合本规范要求', '["office","commercial","industrial","public"]', '2018修订版', '2018-10-01', '国家标准', '本规范是建筑防火设计的基本依据，规定了各类建筑的防火设计要求', 'active'),
('GB50016-5.5.17', '疏散走道宽度要求', 'fire', '公共建筑内疏散走道的净宽度不应小于1.1m；人员密集场所不应小于1.4m', '["commercial","public"]', '2018修订版', '2018-10-01', '国家标准', '疏散走道是人员疏散的主要通道，宽度直接影响疏散效率', 'active'),
('GB50016-5.5.21', '安全出口数量要求', 'fire', '公共建筑每个防火分区或一个防火分区的每个楼层，其安全出口的数量应经计算确定，且不应少于2个', '["commercial","public"]', '2018修订版', '2018-10-01', '国家标准', '多个安全出口可确保在一个出口被封堵时仍有其他疏散途径', 'active'),

-- 无障碍规范 / Accessibility standards
('GB50763-2012', '无障碍设计规范', 'accessibility', '公共建筑应设置无障碍通道、无障碍电梯、无障碍卫生间等设施', '["office","commercial","public"]', '2012版', '2012-09-01', '国家标准', '保障残疾人、老年人等特殊群体的通行和使用需求', 'active'),
('GB50763-3.3.1', '无障碍坡道坡度要求', 'accessibility', '无障碍坡道的坡度不应大于1:12，困难情况下不应大于1:8', '["office","commercial","public"]', '2012版', '2012-09-01', '国家标准', '坡度过大会增加轮椅使用者的通行难度和安全风险', 'active'),
('GB50763-3.7.1', '无障碍卫生间要求', 'accessibility', '公共建筑的公共卫生间应设置无障碍厕位或无障碍卫生间', '["office","commercial","public"]', '2012版', '2012-09-01', '国家标准', '无障碍卫生间应满足轮椅回转空间和扶手设置要求', 'active'),

-- 环保规范 / Environmental standards
('GB50325-2020', '民用建筑工程室内环境污染控制标准', 'environmental', '室内装饰装修材料的甲醛、苯、TVOC等有害物质释放量应符合本标准限值', '["office","commercial","public"]', '2020版', '2020-08-01', '国家标准', '控制室内空气污染，保障人员健康', 'active'),
('GB50325-4.2.1', '甲醛浓度限值', 'environmental', 'I类民用建筑工程室内空气中甲醛浓度限值为0.07mg/m³，II类为0.08mg/m³', '["office","commercial","public"]', '2020版', '2020-08-01', '国家标准', 'I类包括住宅、医院、学校等；II类包括办公楼、商店等', 'active'),

-- 安全规范 / Safety standards
('GB50352-2019', '民用建筑设计统一标准', 'safety', '建筑设计应满足结构安全、使用安全、防火安全等基本要求', '["office","commercial","industrial","public"]', '2019版', '2019-10-01', '国家标准', '民用建筑设计的基本准则和通用要求', 'active'),
('GB50352-6.6.3', '栏杆高度要求', 'safety', '临空高度在24m以下时，栏杆高度不应低于1.05m；临空高度在24m及以上时，栏杆高度不应低于1.10m', '["office","commercial","public"]', '2019版', '2019-10-01', '国家标准', '防止人员坠落，保障使用安全', 'active');
