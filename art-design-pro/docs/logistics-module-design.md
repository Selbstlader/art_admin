模块业务背景+模块是企业管理的重要组成部分涵盖储调度地图定位、核心本模块采用微服务架构设计，支持高并发、高可用的业务场景。特性支持多仓多，实时库存跟踪智能调度径规划实时跟踪集成地理服务引擎多场景分析成本效率统计、报表生

## 2. 系统架构

### 2.1 整体架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    前端应用层 (Vue3 + TypeScript)            │
├─────────────────────────────────────────────────────────────┤
│  仓储管理  │  运输调度  │  地图服务  │  审批流程  │  数据分析  │
├─────────────────────────────────────────────────────────────┤
│                     API网关层                              │
├─────────────────────────────────────────────────────────────┤
│  物流心服务  │  审批引擎服务  │  地图服务  │  通知服务  │
├─────────────────────────────────────────────────────────────┤
│                   数据存储层                               │
│  MySQL(业务数据)  │  Redis(缓存)  │  MongoDB(日志)         │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 模块划分#2.2. 物流核心模块 (Logistics Core)
- **仓储管理子模块**
  - 仓库信息管理
  - 库位管理
  - 入库/出库流程
  - 库存盘点
  - 库存预警

- **运输管理子模块**
  - 运输计划管理
  - 车辆调度
  - 司机管理
  - 路线规划
  - 运输跟踪

#### 2.2.2 通用审批引擎 (OA Approval Engine)
- **流程配置管理**
  - 审批模板设计
  - 审批节点配置
  - 条件分支设置
  - 审批人员分配

- **审批执行引擎**
  - 流程实例管理
  - 任务分发
  - 审批记录
  - 状态流转

- **业务集成接口**
  - 多态关联支持
  - 事件回调机制
  - 自定义表单
  - API集成接口

#### 2.2 地图服务模块 (Map Service)
- **高德地图集成**
  - 地理编码服务
  - 路径规划API
  - 距离计算
  -地图可视化
位置服务
  - 地址解析
  - 坐标转换
  - 地理围栏
  - 位置跟踪

## 3. 数据库设计

### 3.1 物流核心表结构

#### 3.1.1 仓库管理相关表

```sql
-- 仓库信息表
CREATE TABLE warehouse (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    warehouse_code VARCHAR(50) NOT NULL UNIQUE COMMENT '仓库编码',
    warehouse_name VARCHAR(200) NOT NULL COMMENT '仓库名称',
    warehouse_type TINYINT DEFAULT 1 COMMENT '仓库类型1-自营，2-第三方',
    address ARCHAR(500) COMMENT '仓库地址',
    longitde DECIMAL(10,7) COMMENT '经度',
    latitude DECIMAL(10,7) COMMENT '纬度',
    area DECIMAL(10,2) COMMENT '仓库面积(平方米)',
    managr_idBIGINT COMMENT '负责人ID',
    contact_phone VARCHAR(20) COMMENT '联系电话',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-停用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_warehouse_code (warehouse_code),
    INDEX idx_status (status)
) COMMENT='仓库信息表';

-- 库位表
CREATE TABLE storage_location (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT NOT NULL COMMENT '所属仓库ID',
    location_code VARCHAR(50) NOT NULL COMMENT '库位编码',
    location_name VARCHAR(200) COMMENT '库位名称',
    location_type TINYINT DEFAULT 1 COMMENT '库位类型：1-普通，2-冷藏，-危险品',
    area_code VARCHAR(50) COMMENT '区域编码',
    shelf_code VARCHAR(50) COMMENT '货架编码',
    layer_num INT COMMENT '层数',
    position_num INT COMMENT '位置号',
    max_weight DECIMAL(10,2) COMMENT '最大承重(kg)',
    max_volume DECIMAL(10,2) COMMENT '最大体积(立方米)',
    status TINYINT DEFAULT 1 COMMENT'状态：1-可用，2-占用，3-维修',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(id),
    UNIQUE KEY uk_warehouse_location (warehouse_id, location_code),
    INDEX idx_warehouse_id (warehouse_id)
) COMMENT='库位表';
```

#### 3.1.2 库存管理相关表

```sql
-- 商品表
CREATE TABLE product (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_code VARCHAR(50) NOT NULL UNIQUE COMMENT '商品编码',
    product_name VARCHAR(200) NOT NULLCOMMEN '商品名称',
    categor_id BIGINT COMMENT '分类ID',
    unit VARCHAR(20) COMMENT '计量单位',
    weight DECIMAL(10,3) COMMENT '重量(kg)',
    volume DECIMAL(10,3) COMMENT '体积(立方米)',
    ric DECIMAL(10,2) COMMENT '单价',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-停用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMETAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_produt_code (product_code),
    INDEX idx_category_id (category_id)
) COMMENT='商品表';

-- 库存表
CREATE TABLE inventory (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    waehouse_d BIGINT NOT NULL COMMENT '仓库ID',
    location_id BIGINT COMMENT '库位ID',
    roduc_idBIGINTNOT NULL COMMENT '商品ID',
    batch_no ARCHAR(50) COMMENT '批次号',
    quantity DECIMAL(12,3) NOT NULL DEFAULT 0 COMMENT '库存数量',
    available_quantity DECIMAL(12,3) NOT NULL DEFAULT 0 COMMENT '可用数量',
    locked_quantity DECIMAL(12,3) NOT NULL DEFAULT 0 COMMENT '锁定数量',
    production_date DATE COMMENT '生产日期',
    expiry_date DATE COMMENT '过期日期',
    last_check_tme TIMESTAMP COMMENT '最后盘点时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updaed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (warehous_id) REFERENCES warehouse(id),    FOREIGN KEY (location_id) REFERENCES storage_location(id),
    FOREIGN KEY (product_id) REFERENCES product(id),
    UNIQUE KEY uk_location_product_batch (location_id, product_id, batch_no),
    INDEX idx_warehouse_product (warehouse_id, product_id),
    INDEX idx_batch_no (batch_no)
) COMMENT='库存表';
```

#### 3.1.3 入库出库记录表

```sql
- 入库单表
CREATE TABLE inbound_order (
    id BIGINT PRIMARY KEYAUTO_INCREMENT,
    order_no VARCHAR(50) NOT NULL UNIQUE COMMENT '入库单号',
    warehouse_id BIGINT NOT NLL COMMENT '入库仓库ID',
    supplier_id BIGINT COMMENT '供应商ID',
    inbound_type TNYINT DEFAULT 1 COMMENT '入类型1-采购入库，2-退货入库，3-调拨入库',
    total_quantity DECIML(12,3) DEFAULT 0 COMMENT '总数量',
    total_amout DECIMAL(12,2) DEFAULT 0 COMMENT '总金额',
    operator_id BIGINT COMMENT '操作员ID',
    status TINYINT DEFAULT 0 COMMENT '状态：0-待审核，1-已审核，2-已入库，3-已取消',
    remark TEXT COMMENT '备注',
    created_aTIMESTAMP EFAULT CURRENT_TIMESTAMP,
    updatd_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (warehouse_id) REFERENCES warehoue(id),
    INDEX idx_order_no (order_no),
    INDEX dx_warehouse_id (warehouse_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) COMMENT='入库单表';

-- 入库单明细表
CREATE TABLE inbound_order_item (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    inbound_order_id BIGINT NOT NULL COMMENT '入库单ID',
    product_id BIGINT NOT NULL COMMENT '商品ID',
    location_id BIGINT COMMENT '入库库位ID',
    batch_no VARCHAR(50) COMMENT '批次号',
    quantity DECIMAL(12,3) NOT NULL COMMENT '入库数量',
    unit_price DECIMAL(10,2) COMMENT '单价',
    total_price DECIMAL(12,2) COMMENT '总价',
    production_date DATE COMMENT '生产日期',
    expiry_date DATE COMMENT '过期日期',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inbound_order_id) REFERENCES inbound_order(id),
    FOREIGN KEY (product_id) REFERENCES product(id),
    FOREIGN KEY (location_id) REFERENCES storae_location(id),
    INDEX idx_inbound_order_id (inbound_order_id),
    INDEX idx_product_id (product_id)
) COMMENT='入库单明细表';

-- 出库单表
CREATE TABLE outbound_order (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_oARCHAR(50) NOT NULL UNIQUE COMMENT '出库单号',
    warehouse_id BIGINT NOT NULL COMMENT '出库仓库ID',
    customer_id BIGINT COMMENT '客户ID',
    outbound_type TINYINT DEFAULT 1 COMMENT '出库类型：1-销售出库，2-调拨出库，3-报废出库',
    total_quantity DECIMAL(12,3) DEFAULT 0 COMMENT '总数量',
    total_amount DECIMAL(12,2) DEFAULT 0 COMMENT '总金额',
    operator_id BIGINT COMMENT '操作员ID',
    status TINYINT DEFAULT 0 COMMENT '状态：0-待审核，1-已审核，2-已出库，3-已取消',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(id),
    INDEX idx_order_no (order_no),
    INDEX idx_warehos_id (warehouse_id),    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) COMMENT='出库单表';
```

#### 3.1.4 运输管理相关表

```sql
- 运输计划表
CREATE TABLE transport_plan (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    plan_no VARCHAR(50) NOT NULL UNIQUE COMMENT '运输计划编号',
    plan_name VARCHAR(200) COMMENT '计划名称',
    transport_type TINYINT DEFAULT 1 COMMENT '运输类型：1-陆运，2-空运，3-海运',
    from_warehouse_id BIGINT COMMENT '起始仓库ID',
    to_warehouse_idBIGINT COMMENT '目标仓库ID',
    from_address VARCHAR(500) COMMENT '起始址',
    to_address VARCHAR(500) COMMENT '目标地址',
    from_longitude DECIMAL(10,7) COMMENT '起始经度',
    from_latitude DECIMAL(10,7) COMMENT '起始纬度',
    to_longitude DECIMAL(10,7) COMMENT '目标经度',
    to_latitude DECIMAL(10,7) COMMENT '目标纬度',
    estimated_distance DECIMAL(10,2) COMMENT '预估距离(km)',
    estimated_time INT COMMENT '预估时间(分钟)',
    actual_distance DECIMAL(10,2) COMMENT '实际距离(km)',
    actual_time INT COMMENT '实际时间(分钟)',
    vehicle_id BIGINT COMMENT '车辆ID',
    driver_id BIGINT COMMENT '司机ID',
    status TINYINT DEFAULT 1 COMMENT '状态：1-待调度，2-已调度，3-运输中，4-已完成，5-已取消',
    start_time TIMESTAMP COMMENT '开始时间',
    end_time TIMESTAMP COMMENT '结束时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_plan_no (plan_no),
    INDEX idx_status (status),
    INDEX idx_vehicle_id (vehicle_id),
    INDEX idx_driver_id (driver_id)
) COMMENT='运输计划表';

-- 车辆信息表
CREATE TABLE vehicle (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    vehicle_no VARCHAR(50) NOT NULL UNIQUE COMMENT '车牌号',
    vehicle_type VARCHAR(50) COMMENT '车辆类型',
    brand VARCHAR(100) COMMENT '品牌',
    model VARCHAR(100) COMMENT '型号',
    load_capacity DECIMAL(8,2) COMMENT '载重(吨)',
    volume_capacity DECIMAL(8,2) COMMENT '容积(立方米)',
    driver_id BIGINT COMMENT '默认司机ID',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-维修，3-停用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_vehicle_no (vehicle_no),
    INDEX idx_driver_id (driver_id)
) COMMENT='车辆信息表';

-- 司机信息表
CREATE TABLE driver (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    driver_name VARCHAR(100) NOT NULL COMMENT '司机姓名',
    phone VARCHAR(20) COMMENT '联系电话',
    license_no VARCHAR(50) COMMENT '驾驶证号',
    license_type VARCHAR(20) COMMENT '驾照类型',
    experience_years INT COMMENT '驾龄',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-请假，3-离职',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_driver_name (driver_name),
    INDEX idx_phone (phone)
) COMMENT='司机信息表';

-- 地理编码缓存表
CREATE TABLE geocode_cache (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    address VARCHAR(500) NOT NULL COMMENT '地址',
    longitude DECIMAL(10,7) NOT NULL COMMENT '经度',
    latitude DECIMAL(10,7) NOT NULL COMMENT '纬度',
    province VARCHAR(50) COMMENT '省份',
    city VARCHAR(50) COMMENT '城市',
    district VARCHAR(50) COMMENT '区县',
    formatted_address VARCHAR(500) COMMENT '格式化地址',
    accuracy TINYINT DEFAULT 1 COMMENT '精度级别：1-精确，2-近似',
    source VARCHAR(20) DEFAULT 'amap' COMMENT '数据源',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_address (address),
    INDEX idx_longitude_latitude (longitude, latitude),
    INDEX idx_city_district (city, district)
) COMMENT='地理编码缓存表';

-- 库存事务日志表
CREATE TABLE inventory_transaction_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT NOT NULL COMMENT '仓库ID',
    location_id BIGINT COMMENT '库位ID',
    product_id BIGINT NOT NULL COMMENT '商品ID',
    batch_no VARCHAR(50) COMMENT '批次号',
    transaction_type TINYINT NOT NULL COMMENT '事务类型：1-入库，2-出库，3-调拨，4-盘点调整',
    transaction_no VARCHAR(50) NOT NULL COMMENT '事务单号',
    quantity_before DECIMAL(12,3) NOT NULL COMMENT '变更前数量',
    quantity_change DECIMAL(12,3) NOT NULL COMMENT '变更数量(正数为增加，负数为减少)',
    quantity_after DECIMAL(12,3) NOT NULL COMMENT '变更后数量',
    unit_price DECIMAL(10,2) COMMENT '单价',
    total_amount DECIMAL(12,2) COMMENT '总金额',
    operator_id BIGINT COMMENT '操作员ID',
    operator_name VARCHAR(100) COMMENT '操作员姓名',
    remark VARCHAR(500) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(id),
    FOREIGN KEY (location_id) REFERENCES storage_location(id),
    FOREIGN KEY (product_id) REFERENCES product(id),
    INDEX idx_transaction_no (transaction_no),
    INDEX idx_warehouse_product (warehouse_id, product_id),
    INDEX idx_warehouse_product_time (warehouse_id, product_id, created_at),
    INDEX idx_transaction_type (transaction_type),
    INDEX idx_created_at (created_at)
) COMMENT='库存事务日志表';
```

### 3.2 通用审批引擎表结构

```sql
-- 审批模板表
CREATE TABLE approval_template (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    template_name VARCHAR(200) NOT NULL COMMENT '模板名称',
    template_code VARCHAR(50) NOT NULL UNIQUE COMMENT '模板编码',
    business_type VARCHAR(50) COMMENT '业类型',
    description TEXT COMMENT '模板描述',
    version INT DEFAULT 1 COMMENT '版本号',
    status TINYINT DEFAULT 1 COMMENT '状态1-启用，0-禁用',
   created_by BIGINT COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_template_code (template_code),
    INDEX idx_business_typ (usiness_type)
) COMMENT='审批模板表';

-- 审批节点配置表
CREATETBLE approval_node_config (
    id BIGINT RIMARY KEY AUTO_INCREMENT,
    template_id BIGINT NOT NULL COMMENT '模板D',
    node_nameVARCHAR200) NOT NULL COMMENT '节点名称',
    node_type TINYINT DEFAULT 1 COMMENT '节点类型：1-审批，2-抄送，3-条件分支',
    node_order INT NOT NULL COMMENT '节点顺序',
    approval_type TINYINT DEFAULT 1 COMMENT '审批类型：1-单人审批，2-多人会签，3-多人或签',
    approval_users JSON COMMENT '审批人员配置',
    approval_roles JSON COMMENT '审批角色配置',
    conditions SON COMMENT '节点条件配置',
    time_limit INT COMMENT '超时时间(小时)',
    is_required TINYINT DEFAULT 1 COMMENT '是否必须：1-是，0-否',
    created_at TIMESTAMP DEFAULT CURRENT_TIMETM,
    FOREGN KEY (template_id) REFERENCESapproval_template(id),
    INDEX idx_template_id (template_id),
    INDEX idx_node_order (node_order)
) COMMENT='审批节点配置表';

-- 审批实例表
CREATE TABLE approval_instance (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    instance_no VARCHAR(50) NOT NULL UNIQUE COMMENT '审批实例编号',
    template_id BIGINT NOT NULL COMMENT '模板ID',
    business_type VARCHAR(50) COMMENT '业务类型',
    business_id BIGINT COMMENT '业务ID',
    title VARCHAR(500) COMMENT '审批标题',
    content TEXT COMMENT '审批内容',
    form_data JSON COMMENT '表单数据',
    applicant_id BIGINT COMMENT '申请人ID',
    current_node_id BIGINT COMMENT '当前节点ID',
    status TINYINT DEFAULT 1 COMMENT '状态：1-进行中，2-已通过，3-已拒绝，4-已撤销',
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP COMMENT '结束时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES approval_template(id),
    INDEX idx_instance_no (instance_no),
    INDEX idx_business (business_type, business_id),
    INDEX idx_applicant_id (applicant_id),
    INDEX idx_status (status)
) COMMENT='审批实例表';

-- 审批记录表
CREATE TABLE approval_record (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    instance_id BIGINT NOT NULL COMMENT '审批实例ID',
    node_id BIGINT NOT NULL COMMENT '节点ID',
    node_name VARCHAR(200) COMMENT '节点名称',
    approver_id BIGINT COMMENT '审批人ID',
    approver_name VARCHAR(100) COMMENT '审批人姓名',
    action TINYINT COMMENT '操作：1-同意，2-拒绝，3-转交，4-撤销',
    comment TEXT COMMENT '审批意见',
    approval_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (instance_id) REFERENCES approval_instance(id),
    INDEX idx_instance_id (instance_id),
    INDEX idx_approver_id (approver_id),
    INDEX idx_approval_time (approval_time)
) COMMENT='审批记录表';

-- 审批实例节点状态表
CREATE TABLE approval_instance_node (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    instance_id BIGINT NOT NULL COMMENT '审批实例ID',
    node_id BIGINT NOT NULL COMMENT '节点配置ID',
    node_name VARCHAR(200) COMMENT '节点名称',
    status TINYINT DEFAULT 1 COMMENT '状态：1-待处理，2-处理中，3-已完成，4-已跳过',
    start_time TIMESTAMP COMMENT '开始时间',
    end_time TIMESTAMP COMMENT '结束时间',
    approvers JSON COMMENT '节点审批人员列表',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (instance_id) REFERENCES approval_instance(id),
    UNIQUE KEY uk_instance_node (instance_id, node_id),
    INDEX idx_instance_id (instance_id),
    INDEX idx_status (status)
) COMMENT='审批实例节点状态表';
```

## 4. API接口设计

### 4.1 仓储管理API

#### 4.1.1 仓库管理接口

```typescript
// 仓库列表查询
GET /api/logistics/warehouse/list
{
  "page": 1,
  "size": 20,
  "warehouse_code": "WH001",
  "warehouse_name": "主仓库",
  "status": 1
}

// 仓库详情
GET /api/logistics/warehouse/:id

// 新增仓库
POST /api/logistics/warehouse
{
  "warehouse_code": "WH002",
  "warehouse_name": "分仓库",
  "warehouse_type": 1,
  "address": "北京市朝阳区xxx街道",
  "longitude": 116.397128,
  "latitude": 39.916527,
  "area": 1000.50,
  "manager_id": 1001,
  "contact_phone": "13800138000"
}

// 更新仓库
PUT /api/logistics/warehouse/:id

// 删除仓库
DELETE /api/logistics/warehouse/:id
```

#### 4.1.2 入库管理接口

```typescript
// 入库单列表
GET /api/logistics/inbound/list
{
  "page": 1,
  "size": 20,
  "order_no": "RK20231201001",
  "warehouse_id": 1,
  "status": 0,
  "start_date": "2023-12-01",
  "end_date": "2023-12-31"
}

// 创建入库单
POST /api/logistics/inbound
{
  "warehouse_id": 1,
  "supplier_id": 1001,
  "inbound_type": 1,
  "remark": "采购入库",
  "items": [
    {
      "product_id": 2001,
      "location_id": 3001,
      "batch_no": "B0231201001",
      "quantity": 10000,
      "unit_price": 50.00,
      "production_date": "2023-11-01",
      "expiry_date": "2024-11-01"
    }
  ]
}

// 入库单审核
POST /api/logistics/inbound/:id/approve
{
  "action": 1, // 1-同意，2-拒绝
  "comment": "审核通过"
}

// 确认入库
POST /api/logistics/inbound/:id/confirm
```

### 4.2 运输管理API

#### 4.2.1 运输计划接口

```typescript
// 运输计划列表
GET /api/logistics/transport/list
{
  "page": 1,
  "size": 20,
  "plan_no": "TP20231201001",
  "status": 1,
  "vehicle_id": 1001,
  "driver_id": 2001
}

// 创建运输计划
POST /api/logistics/transport
{
  "plan_name": "北京到上海运输",
  "transport_type": 1,
  "from_warehouse_id": 1,
  "to_warehouse_id": 2,
  "from_address": "北京市朝阳区xxx",
  "to_address": "上海市浦东新区xxx",
  "vehicle_id": 1001,  "driver_id": 2001
}

// 路径规划
POST /api/logistics/transport/routeplan
{
 "from_address": "北京市朝阳区",
  "to_address": "上海市浦东新区",
  "waypoints": ["天津市河西区"]
}

// 更新运输
PUT /api/logistics/transport/:id/status
{
  "status": 3, // 运输中
  "location": {
    "longitude": 116.397128,
    "latitude": 39.916527,
    "address": "北京市朝阳区xxx"
  }
}
```

### 4.3 审批引擎API

#### 4.3.1 模板接口

```typescript
// 审批模板列表
GET /api/approval/template/list
{
  "page": 1,
  "size": 20,
  "template_name": "入库审批",
  "business_type": "inbound",
  "status": 1
}

// 创建审批模板
POST /api/approval/template
{
  "template_name": "入库审批流程",
  "template_code": "INBOUND_APPROVAL",
  "business_type": "inbound",
  "description": "用于入库单的审批流程",
  "nodes": [
    {
      "node_name": "主管审批",
      "node_type": 1,
      "node_order": 1,
      "approval_type": 1,
      "approval_roles": ["warehouse_manager"],
      "time_limit": 24,
      "is_required": 1
    },
    {
      "node_name": "财务审批",
      "node_type": 1,
      "node_order": 2,
      "approval_type": 1,
      "approval_roles": ["finance_manager"],
      "time_limit": 48,
      "is_required": 1
    }
  ]
}

// 启动审批流程
POST /api/approval/instance
{
  "template_code": "INBOUND_APPROVAL",
  "business_type": "inbound",
  "business_id": 1001,
  "title": "入库单RK20231201001审批",
  "content": "请审核入库单信息",
  "form_data": {
    "order_no": "RK20231201001",
    "total_amount": 5000.00,
    "supplier_name": "xxx供应商"
  }
}

// 审批操作
OST /api/approval/instance/:id/approve
{
  "action": 1, // 1-同意，2-拒绝
  "comment": "审批通过，信息核对无误"
}
```

## 5. 前端页面设计

### 5.1 页面结构

```
/src/views/logistics/
├── warehouse/           # 仓储管理
│   ├── index.vue       # 仓库列表
│   ├── detail.vue      # 仓库详情
│   ├── location/       # 库位管理
│   └── inventory/      # 库存管理
├── inbound/            # 入库管理
│   ├── index.vue       # 入库单列表
│   ├── create.vue      # 创建入库单
│   └── detail.vue      # 入库单详情
├── outbound/           # 出库管理
│   ├── index.vue       # 出库单列表
│   ├── create.vue      # 创建出库单
│   └── detail.vue      # 出库单详情
├── transport/          # 运输管理
│   ├── dex.vue       # 运输计划列表
│   ├── create.vue      # 创建运输计划
│   ├── detail.vue      # 运输计划详情
│   ├── map.vue         # 地图跟踪
│   ├── vehcle/        # 车辆管理
│   └── driver/         # 司机管理
└── pproval/           # 审批管理    ├── template/       # 模板管理
    ├── instance/       # 审批实例
    └── todo/           # 我的待办
```

### 5.2 核心页面功能设计

#### 5.2.1 仓储管理页面

**仓库列表页面特性：**
支持仓库编码、名称、状态筛选
- 地图展示仓库分布位置
- 快速查看库存概览
- 支持仓库导入导出功能

库存管理页面特性：**
- 实时库存数据展示
- 支持多维度库存查询
- 库存预警提示
- 批量库存调整功能

#### 5.2.2 运输管理页面

**运输计划页面特性：**
- 日程视图展示运输计划
- 车辆和司机资源调度
- 径规划和距离计算
- 运输状态实时更新

**地图跟踪页面特性：**
- 高德地图集成
- 实时位置跟踪
- 历史轨迹回放
- 地理围栏设置

#### 5.2.3 审批页面

待办审批页面特性**
- 个人待办事项列表
- 快速审批操作
- 审批历史查看
- 消息通知提醒

**模板配置页面特性：**
- 可视化流程设计器
- 拖拽式节点配置
- 条件分支设置
- 模板版本管理

## 6. 技术实现要点

### 6.1 高德地图集成（Redis缓存优化）

```typescript
// Redis缓存服务
class RedisCacheService {
  private redis: Redis;
  
  constructor() {
    this.redis = new Redis({
      host: process.env.REDIS_HOST || 'localhost',
      port: parseInt(process.env.REDIS_PORT || '6379'),
      password: process.env.REDIS_PASSWORD,
      db: parseInt(process.env.REDIS_DB || '0')
    });
  }
  
  async get<T>(key: string): Promise<T | null> {
    const value = await this.redis.get(key);
    return value ? JSON.parse(value) : null;
  }
  
  async set(key: string, value: any, ttl: number = 24 * 60 * 60): Promise<void> {
    await this.redis.setex(key, ttl, JSON.stringify(value));
  }
  
  async del(key: string): Promise<void> {
    await this.redis.del(key);
  }
}

// 地图服务封装（Redis缓存 + 容错机制）
class MapService {
  private amap: any;
  private cache: RedisCacheService;
  private rateLimiter: { [key: string]: number } = {};
  private readonly RATE_LIMIT = 100; // 每分钟最大调用次数
  private readonly CACHE_TTL = 24 * 60 * 60; // 缓存24小时
  
  constructor() {
    this.cache = new RedisCacheService();
  }
  
  // 初始化地图
  initMap(container: string, center: [number, number]) {
    this.amap = new AMap.Map(container, {
      center,
      zoom: 10,
      mapStyle: 'amap://styles/normal'
    });
  }
  
  // 地址转坐标（Redis缓存）
  async geocode(address: string): Promise<[number, number]> {
    const cacheKey = `geocode:${Buffer.from(address).toString('base64')}`;
    
    // 检查Redis缓存
    const cached = await this.cache.get<[number, number]>(cacheKey);
    if (cached) {
      return cached;
    }
    
    // 检查频率限制
    if (this.checkRateLimit('geocode')) {
      throw new Error('地图API调用频率超限，请稍后重试');
    }
    
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error('地址解析超时'));
      }, 10000);
      
      AMap.plugin('AMap.Geocoder', () => {
        const geocoder = new AMap.Geocoder();
        geocoder.getLocation(address, async (status: string, result: any) => {
          clearTimeout(timeout);
          if (status === 'complete' && result.geocodes.length > 0) {
            const location = result.geocodes[0].location;
            const coords = [location.lng, location.lat];
            
            // 缓存到Redis
            await this.cache.set(cacheKey, coords, this.CACHE_TTL);
            
            resolve(coords);
          } else {
            reject(new Error('地址解析失败'));
          }
        });
      });
    });
  }
  
  // 路径规划（Redis缓存）
  async routePlanning(
    origin: [number, number], 
    destination: [number, number],
    waypoints?: [number, number][]
  ): Promise<any> {
    const routeKey = `route:${origin.join(',')}_${destination.join(',')}_${waypoints?.join(',') || ''}`;
    
    // 检查Redis缓存
    const cached = await this.cache.get<any>(routeKey);
    if (cached) {
      return cached;
    }
    
    if (this.checkRateLimit('route')) {
      throw new Error('路径规划API调用频率超限，请稍后重试');
    }
    
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error('路径规划超时'));
      }, 15000);
      
      AMap.plugin('AMap.Driving', () => {
        const driving = new AMap.Driving({
          map: this.amap,
          panel: 'panel'
        });
        
        driving.search(origin, destination, {
          waypoints
        }, async (status: string, result: any) => {
          clearTimeout(timeout);
          if (status === 'complete') {
            // 缓存到Redis
            await this.cache.set(routeKey, result, this.CACHE_TTL);
            resolve(result);
          } else {
            reject(new Error('路径规划失败'));
          }
        });
      });
    });
  }
  
  // 频率限制检查
  private checkRateLimit(apiType: string): boolean {
    const now = Date.now();
    const key = `${apiType}_${Math.floor(now / 60000)}`; // 按分钟分组
    
    if (!this.rateLimiter[key]) {
      this.rateLimiter[key] = 0;
    }
    
    this.rateLimiter[key]++;
    return this.rateLimiter[key] > this.RATE_LIMIT;
  }
}
```

### 6.2 审批引擎核心实现（Webhook回调机制）

```typescript
// Webhook回调接口定义
interface WebhookCallback {
  url: string;
  event: string; // 'approval.completed' | 'approval.rejected' | 'approval.cancelled'
  headers?: { [key: string]: string };
  retryCount?: number;
  timeout?: number;
}

// 审批引擎服务
class ApprovalEngine {
  private webhookService: WebhookService;
  
  constructor() {
    this.webhookService = new WebhookService();
  }
  
  // 启动审批流程
  async startApproval(params: StartApprovalParams) {
    // 1. 获取审批模板
    const template = await this.getTemplate(params.templateCode);
    
    // 2. 创建审批实例
    const instance = await this.createInstance({
      templateId: template.id,
      businessType: params.businessType,
      businessId: params.businessId,
      title: params.title,
      content: params.content,
      formData: params.formData,
      applicantId: params.applicantId,
      webhookUrl: params.webhookUrl // 业务模块提供的回调地址
    });
    
    // 3. 初始化第一个节点
    const firstNode = await this.getFirstNode(template.id);
    if (firstNode) {
      await this.createApprovalTask({
        instanceId: instance.id,
        nodeId: firstNode.id,
        approvers: await this.resolveApprovers(firstNode)
      });
    }
    
    // 4. 发送通知
    await this.sendNotification(instance);
    
    return instance;
  }
  
  // 处理审批
  async processApproval(params: ProcessApprovalParams) {
    // 1. 验证权限
    await this.validateApprovalPermission(params);
    
    // 2. 记录审批结果
    await this.recordApproval({
      instanceId: params.instanceId,
      nodeId: params.nodeId,
      approverId: params.approverId,
      action: params.action,
      comment: params.comment
    });
    
    // 3. 检查节点是否完成
    const nodeCompleted = await this.checkNodeCompletion(params.nodeId);
    if (nodeCompleted) {
      // 4. 流转到下一个节点
      const nextNode = await this.moveToNextNode(params.instanceId, params.nodeId);
      
      // 5. 如果流程完成，触发Webhook回调
      if (!nextNode) {
        await this.triggerApprovalCompleted(params.instanceId);
      }
    }
    
    // 6. 更新实例状态
    await this.updateInstanceStatus(params.instanceId);
    
    // 7. 发送通知
    await this.sendApprovalNotification(params);
  }
  
  // 审批完成回调
  private async triggerApprovalCompleted(instanceId: number) {
    const instance = await this.getApprovalInstance(instanceId);
    
    const callbackData = {
      instanceId: instance.id,
      instanceNo: instance.instance_no,
      businessType: instance.business_type,
      businessId: instance.business_id,
      status: instance.status,
      completedAt: new Date().toISOString(),
      formData: instance.form_data
    };
    
    // 触发Webhook回调
    await this.webhookService.sendCallback({
      url: instance.webhook_url,
      event: 'approval.completed',
      data: callbackData
    });
  }
}

// Webhook服务
class WebhookService {
  async sendCallback(params: { url: string; event: string; data: any }) {
    const payload = {
      event: params.event,
      data: params.data,
      timestamp: new Date().toISOString(),
      signature: this.generateSignature(params.data)
    };
    
    try {
      const response = await fetch(params.url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'User-Agent': 'Logistics-Approval-Webhook/1.0',
          ...this.getAuthHeaders()
        },
        body: JSON.stringify(payload),
        timeout: 30000
      });
      
      if (!response.ok) {
        throw new Error(`Webhook调用失败: ${response.status}`);
      }
      
      console.log(`Webhook回调成功: ${params.event} -> ${params.url}`);
    } catch (error) {
      console.error(`Webhook回调失败: ${error.message}`);
      // 可以加入重试机制
      await this.retryCallback(params, 1);
    }
  }
  
  private async retryCallback(params: { url: string; event: string; data: any }, retryCount: number) {
    if (retryCount > 3) {
      console.error('Webhook回调重试次数超限，放弃重试');
      return;
    }
    
    // 延迟重试
    setTimeout(async () => {
      try {
        await this.sendCallback(params);
      } catch (error) {
        await this.retryCallback(params, retryCount + 1);
      }
    }, Math.pow(2, retryCount) * 1000); // 指数退避
  }
  
  private generateSignature(data: any): string {
    // 生成签名用于验证回调的真实性
    const crypto = require('crypto');
    const secret = process.env.WEBHOOK_SECRET || 'default-secret';
    return crypto.createHmac('sha256', secret).update(JSON.stringify(data)).digest('hex');
  }
  
  private getAuthHeaders(): { [key: string]: string } {
    return {
      'Authorization': `Bearer ${process.env.WEBHOOK_TOKEN || ''}`
    };
  }
}
```

### 6.3 状态机设计

```typescript
// 入库单状态机
enum InboundOrderStatus {
  DRAFT = 0,        // 草稿
  PENDING_APPROVE = 1,  // 待审核
  APPROVED = 2,     // 已审核
  RECEIVED = 3,     // 已入库
  CANCELLED = 4     // 已取消
}

class InboundOrderStateMachine {
  private transitions = {
    [InboundOrderStatus.DRAFT]: [InboundOrderStatus.PENDING_APPROVE, InboundOrderStatus.CANCELLED],
    [InboundOrderStatus.PENDING_APPROVE]: [InboundOrderStatus.APPROVED, InboundOrderStatus.CANCELLED],
    [InboundOrderStatus.APPROVED]: [InboundOrderStatus.RECEIVED, InboundOrderStatus.CANCELLED],
    [InboundOrderStatus.RECEIVED]: [], // 终态
    [InboundOrderStatus.CANCELLED]: [] // 终态
  };
  
  canTransition(from: InboundOrderStatus, to: InboundOrderStatus): boolean {
    return this.transitions[from]?.includes(to) || false;
  }
  
  transition(order: InboundOrder, newStatus: InboundOrderStatus): boolean {
    if (this.canTransition(order.status, newStatus)) {
      order.status = newStatus;
      return true;
    }
    throw new Error(`无法从状态 ${order.status} 转换到 ${newStatus}`);
  }
}
```

## 7. 数据库优化配置

### 7.1 连接池配置

```yaml
# MySQL连接池配置
database:
  mysql:
    host: localhost
    port: 3306
    database: logistics_db
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}
    pool:
      minimum: 5
      maximum: 20
      acquire: 30000
      idle: 10000
      evict: 1000
    
  redis:
    host: ${REDIS_HOST}
    port: ${REDIS_PORT}
    password: ${REDIS_PASSWORD}
    db: 0
    pool:
      min: 2
      max: 10
      idle_timeout: 30000
```

### 7.2 查询优化建议

```sql
-- 常用查询优化
-- 1. 库存查询优化（使用复合索引）
EXPLAIN SELECT * FROM inventory 
WHERE warehouse_id = 1 AND product_id = 1001 
ORDER BY created_at DESC LIMIT 10;

-- 2. 审批实例查询优化
EXPLAIN SELECT i.*, t.template_name 
FROM approval_instance i 
LEFT JOIN approval_template t ON i.template_id = t.id 
WHERE i.business_type = 'inbound' AND i.status = 1 
ORDER BY i.created_at DESC;

-- 3. 地理编码查询优化
EXPLAIN SELECT * FROM geocode_cache 
WHERE city = '北京市' AND district = '朝阳区' 
ORDER BY accuracy DESC, created_at DESC;
```

## 8. 数据保留策略

### 8.1 日志数据归档

```sql
-- 库存事务日志归档策略
-- 1. 创建归档表
CREATE TABLE inventory_transaction_log_archive LIKE inventory_transaction_log;

-- 2. 定期归档脚本（保留2年数据）
DELIMITER //
CREATE PROCEDURE archive_inventory_logs()
BEGIN
    DECLARE archive_date DATE;
    SET archive_date = DATE_SUB(CURRENT_DATE, INTERVAL 2 YEAR);
    
    -- 移动数据到归档表
    INSERT INTO inventory_transaction_log_archive 
    SELECT * FROM inventory_transaction_log 
    WHERE created_at < archive_date;
    
    -- 删除原表数据
    DELETE FROM inventory_transaction_log 
    WHERE created_at < archive_date;
    
    -- 优化表
    OPTIMIZE TABLE inventory_transaction_log;
END //
DELIMITER ;

-- 3. 定时任务（每月执行一次）
-- 添加到crontab: 0 0 1 * * mysql -u user -p -e "CALL archive_inventory_logs();"
```

### 8.2 审批历史清理

```sql
-- 审批记录保留策略
DELIMITER //
CREATE PROCEDURE cleanup_approval_history()
BEGIN
    DECLARE cleanup_date DATE;
    SET cleanup_date = DATE_SUB(CURRENT_DATE, INTERVAL 5 YEAR);
    
    -- 删除5年前的已完成审批记录
    DELETE FROM approval_record 
    WHERE instance_id IN (
        SELECT id FROM approval_instance 
        WHERE status IN (2, 3, 4) AND end_time < cleanup_date
    );
    
    -- 删除对应的审批实例
    DELETE FROM approval_instance 
    WHERE status IN (2, 3, 4) AND end_time < cleanup_date;
    
    -- 清理审批实例节点状态
    DELETE FROM approval_instance_node 
    WHERE instance_id NOT IN (SELECT id FROM approval_instance);
END //
DELIMITER ;
```

## 9. 部署和运维

### 9.1 服务部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                    负载均衡器 (Nginx)                        │
├─────────────────────────────────────────────────────────────┤
│  Web服务1  │  Web服务2  │  Web服务3  │  Admin服务           │
├─────────────────────────────────────────────────────────────┤
│  物流服务1  │  物流服务2  │  审批服务1  │  审批服务2          │
├─────────────────────────────────────────────────────────────┤
│  MySQL主库  │  MySQL从库  │  Redis集群  │  MongoDB集群       │
└─────────────────────────────────────────────────────────────┘
```

### 9.2 监控指标

**业务监控指标：**
- 入库单处理时效
- 库存周转率
- 运输准时率
- 审批流程耗时

**技术监控指标：**
- API响应时间
- 数据库连接池使用率
- Redis缓存命中率
- 地图API调用次数和成功率

### 9.3 安全考虑

- 地图API密钥加密存储
- 审批操作日志记录
- 敏感数据脱敏处理
- 权限控制和数据隔离
- Webhook签名验证

### 9.4 API限流策略

```typescript
// API限流中间件配置
class RateLimitMiddleware {
  private redis: Redis;
  private readonly limits = {
    // 物流API限流配置
    '/api/logistics/warehouse': { requests: 100, window: 60000 }, // 每分钟100次
    '/api/logistics/inbound': { requests: 50, window: 60000 },   // 每分钟50次
    '/api/logistics/transport': { requests: 30, window: 60000 }, // 每分钟30次
    
    // 审批API限流配置
    '/api/approval/instance': { requests: 20, window: 60000 },   // 每分钟20次
    '/api/approval/template': { requests: 10, window: 60000 },   // 每分钟10次
    
    // 地图相关API限流（更严格）
    '/api/logistics/map': { requests: 10, window: 60000 },       // 每分钟10次
  };
  
  async checkLimit(userId: string, endpoint: string): Promise<boolean> {
    const limit = this.limits[endpoint] || { requests: 100, window: 60000 };
    const key = `rate_limit:${userId}:${endpoint}`;
    
    const current = await this.redis.incr(key);
    if (current === 1) {
      await this.redis.expire(key, Math.ceil(limit.window / 1000));
    }
    
    return current <= limit.requests;
  }
}

// 限流异常处理
export class RateLimitExceededError extends Error {
  constructor(endpoint: string, limit: number) {
    super(`API调用频率超限: ${endpoint} 每分钟最多${limit}次`);
    this.name = 'RateLimitExceededError';
  }
}
```

### 9.5 灾备和恢复策略

**数据库备份策略：**
```bash
#!/bin/bash
# 数据库备份脚本
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="logistics_db"

# 全量备份（每日凌晨2点）
mysqldump -u $DB_USER -p$DB_PASS \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  $DB_NAME | gzip > $BACKUP_DIR/full_backup_$DATE.sql.gz

# 增量备份（每4小时）
mysqlbinlog --read-from-remote-server \
  --host=$DB_HOST \
  --raw \
  --stop-never \
  mysql-bin > $BACKUP_DIR/binlog_$DATE.log

# 保留30天备份，删除过期文件
find $BACKUP_DIR -name "*.gz" -mtime +30 -delete
find $BACKUP_DIR -name "*.log" -mtime +7 -delete
```

**Redis数据备份：**
```bash
#!/bin/bash
# Redis备份脚本
REDIS_DIR="/backup/redis"
DATE=$(date +%Y%m%d_%H%M%S)

# RDB快照备份
redis-cli --rdb $REDIS_DIR/dump_$DATE.rdb

# AOF备份
cp /var/lib/redis/appendonly.aof $REDIS_DIR/aof_$DATE.aof

# 保留7天备份
find $REDIS_DIR -name "*.rdb" -mtime +7 -delete
find $REDIS_DIR -name "*.aof" -mtime +7 -delete
```

**故障恢复流程：**
1. **数据库故障恢复**
   - 停止应用服务
   - 恢复最近的完整备份
   - 应用增量日志
   - 验证数据完整性
   - 重启应用服务

2. **Redis故障恢复**
   - 切换到备用Redis实例
   - 恢复RDB或AOF文件
   - 重建缓存数据
   - 验证缓存一致性

### 9.6 故障排查指南

**常见问题处理：**

1. **地图API调用失败**
   - 检查API密钥是否有效
   - 验证调用频率是否超限
   - 查看Redis缓存是否正常工作

2. **审批流程卡住**
   - 检查当前节点状态
   - 验证审批人员权限
   - 查看Webhook回调日志

3. **库存数据不一致**
   - 检查事务日志完整性
   - 验证并发控制机制
   - 执行库存数据修复脚本

## 10. 后续扩展规划

### 10.1 功能扩展
- 智能仓储机器人集成
- 无人机配送支持
- 区块链溯源功能
- 大数据分析预测

### 10.2 技术升级
- 微服务架构演进
- 容器化部署
- 云原生改造
- AI算法优化

## 11. 总结

物流仓储+运输物流模块采用模块化、可扩展的架构设计，通过通用审批引擎支持多业务场景的审批需求，集成高德地图提供丰富的地理位置服务。系统具备高性能、高可用、易维护的特点，能够满足企业物流管理的核心需求，并为未来的功能扩展预留了充分的空间。

### 11.1 核心优势
1. **模块化设计**：物流核心、审批引擎、地图服务三大模块独立部署和维护
2. **通用审批引擎**：支持多业务场景，具备Webhook回调机制
3. **性能优化**：Redis缓存、数据库索引优化、连接池配置
4. **成本控制**：地理编码缓存减少API调用成本
5. **数据安全**：完整的审计追踪和数据保留策略

### 11.2 技术亮点
- 多态关联的审批引擎设计
- Redis分布式缓存架构
- Webhook异步回调机制
- 完善的状态机管理
- 自动化数据归档策略

该设计文档为物流系统的开发提供了完整的技术方案和实施指导，确保系统的可靠性、可扩展性和可维护性。

## 12. 开发者快速入门指南

### 12.1 环境准备

**开发环境要求：**
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+
- MongoDB 4.4+

**本地开发配置：**
```bash
# 1. 克隆项目
git clone <repository-url>
cd logistics-module

# 2. 安装依赖
npm install

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息

# 4. 初始化数据库
npm run db:migrate
npm run db:seed

# 5. 启动开发服务
npm run dev
```

### 12.2 项目结构

```
logistics-module/
├── src/
│   ├── modules/
│   │   ├── logistics/          # 物流核心模块
│   │   ├── approval/           # 审批引擎模块
│   │   └── map/               # 地图服务模块
│   ├── common/
│   │   ├── cache/             # Redis缓存服务
│   │   ├── database/          # 数据库连接
│   │   └── middleware/        # 中间件
│   ├── api/                   # API路由
│   ├── types/                 # TypeScript类型定义
│   └── utils/                 # 工具函数
├── docs/                      # 文档
├── tests/                     # 测试文件
├── scripts/                   # 脚本文件
└── docker/                    # Docker配置
```

### 12.3 开发规范

**代码规范：**
- 使用TypeScript严格模式
- 遵循ESLint和Prettier配置
- 所有API必须有完整的类型定义
- 数据库操作必须使用事务

**Git提交规范：**
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 代码重构
test: 测试相关
chore: 构建过程或辅助工具的变动
```

## 13. 测试策略

### 13.1 单元测试

**测试框架：** Jest + Supertest

**覆盖范围：**
- 业务逻辑层：≥90%
- 数据访问层：≥85%
- 工具函数：≥95%

**示例测试：**
```typescript
// tests/unit/logistics/inbound.test.ts
describe('InboundOrderService', () => {
  test('should create inbound order successfully', async () => {
    const orderData = {
      warehouseId: 1,
      products: [{ productId: 1001, quantity: 10 }],
      applicantId: 1001
    };
    
    const result = await inboundService.createOrder(orderData);
    expect(result.status).toBe(InboundOrderStatus.DRAFT);
    expect(result.orderNo).toMatch(/^IB\d{12}$/);
  });
  
  test('should throw error for insufficient inventory', async () => {
    // 测试库存不足场景
  });
});
```

### 13.2 集成测试

**测试场景：**
- 审批流程端到端测试
- 地图API集成测试
- 数据库事务测试
- Redis缓存一致性测试

```typescript
// tests/integration/approval.test.ts
describe('Approval Flow Integration', () => {
  test('should complete approval workflow', async () => {
    // 1. 创建入库单
    const order = await createInboundOrder();
    
    // 2. 启动审批流程
    const instance = await approvalEngine.startApproval({
      businessType: 'inbound',
      businessId: order.id,
      templateCode: 'inbound_approval'
    });
    
    // 3. 模拟审批操作
    await approvalEngine.processApproval({
      instanceId: instance.id,
      approverId: 2001,
      action: 'approve'
    });
    
    // 4. 验证最终状态
    const finalOrder = await getInboundOrder(order.id);
    expect(finalOrder.status).toBe(InboundOrderStatus.APPROVED);
  });
});
```

### 13.3 性能测试

**测试工具：** Artillery + K6

**测试指标：**
- API响应时间：P95 < 500ms
- 并发用户数：支持1000+
- 数据库连接池：利用率 < 80%
- Redis缓存命中率：> 90%

**性能测试配置：**
```yaml
# artillery-config.yml
config:
  target: 'http://localhost:3000'
  phases:
    - duration: 60
      arrivalRate: 10
    - duration: 120
      arrivalRate: 50
    - duration: 60
      arrivalRate: 100

scenarios:
  - name: "物流API压力测试"
    weight: 70
    flow:
      - get:
          url: "/api/logistics/warehouse/list"
      - post:
          url: "/api/logistics/inbound/create"
          json:
            warehouseId: 1
            products: [{ productId: 1001, quantity: 5 }]
```

### 13.4 安全测试

**测试范围：**
- SQL注入防护
- XSS攻击防护
- CSRF令牌验证
- API权限验证
- 数据脱敏验证

## 14. API文档附录

### 14.1 核心API示例

**创建入库单：**
```http
POST /api/logistics/inbound/create
Content-Type: application/json
Authorization: Bearer <token>

{
  "warehouseId": 1,
  "supplierId": 1001,
  "expectedDate": "2024-01-15",
  "products": [
    {
      "productId": 1001,
      "quantity": 100,
      "batchNo": "B20240115001",
      "remark": "第一批次"
    }
  ],
  "remark": "常规入库"
}

Response:
{
  "code": 200,
  "data": {
    "id": 1001,
    "orderNo": "IB202401150001",
    "status": 0,
    "createdAt": "2024-01-15T10:30:00Z"
  }
}
```

**审批操作：**
```http
POST /api/approval/process
Content-Type: application/json
Authorization: Bearer <token>

{
  "instanceId": 1001,
  "nodeId": 2001,
  "action": "approve",
  "comment": "审核通过，质量符合要求"
}

Response:
{
  "code": 200,
  "data": {
    "instanceId": 1001,
    "status": 2,
    "nextNode": null,
    "completedAt": "2024-01-15T14:30:00Z"
  }
}
```

**地理编码查询：**
```http
GET /api/logistics/map/geocode?address=北京市朝阳区建国门外大街1号
Authorization: Bearer <token>

Response:
{
  "code": 200,
  "data": {
    "longitude": 116.4307,
    "latitude": 39.9155,
    "province": "北京市",
    "city": "北京市",
    "district": "朝阳区",
    "formattedAddress": "北京市朝阳区建国门外大街1号",
    "accuracy": 1
  }
}
```

### 14.2 错误码定义

| 错误码 | 描述 | 解决方案 |
|--------|------|----------|
| 1001 | 参数验证失败 | 检查请求参数格式 |
| 1002 | 权限不足 | 联系管理员分配权限 |
| 2001 | 仓库不存在 | 检查仓库ID是否正确 |
| 2002 | 库存不足 | 调整入库数量或联系采购 |
| 3001 | 审批实例不存在 | 检查审批ID是否正确 |
| 3002 | 无审批权限 | 联系流程管理员 |
| 4001 | 地图API调用失败 | 检查网络连接和API密钥 |
| 4002 | 地址解析失败 | 检查地址格式是否正确 |

### 14.3 Webhook回调格式

**审批完成回调：**
```json
{
  "event": "approval.completed",
  "data": {
    "instanceId": 1001,
    "instanceNo": "AP202401150001",
    "businessType": "inbound",
    "businessId": 1001,
    "status": 2,
    "completedAt": "2024-01-15T14:30:00Z",
    "formData": {
      "warehouseId": 1,
      "totalAmount": 10000.00
    }
  },
  "timestamp": "2024-01-15T14:30:01Z",
  "signature": "sha256=abcdef1234567890..."
}
```

---

## 15. 迁移策略

### 15.1 系统集成方案

**现有系统集成策略：**
```sql
-- 1. 数据迁移脚本
-- 创建迁移日志表
CREATE TABLE migration_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    migration_name VARCHAR(200) NOT NULL,
    status TINYINT DEFAULT 0 COMMENT '0-待执行，1-执行中，2-完成，3-失败',
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 仓库数据迁移
INSERT INTO warehouse (warehouse_code, warehouse_name, address, longitude, latitude, status)
SELECT 
    CONCAT('WH', LPAD(old.id, 6, '0')),
    old.name,
    old.address,
    COALESCE(gc.longitude, 0),
    COALESCE(gc.latitude, 0),
    1
FROM old_warehouse old
LEFT JOIN geocode_cache gc ON old.address = gc.address
WHERE old.is_active = 1;

-- 库存数据迁移
INSERT INTO inventory (warehouse_id, product_id, quantity, batch_no, created_at)
SELECT 
    w.id,
    p.id,
    old.stock_quantity,
    old.batch_number,
    old.created_time
FROM old_inventory old
JOIN warehouse w ON old.warehouse_code = w.warehouse_code
JOIN product p ON old.product_code = p.product_code;
```

**分阶段部署计划：**

**第一阶段（2周）：基础环境搭建**
- 数据库环境准备和迁移脚本开发
- Redis缓存环境配置
- 基础服务框架搭建
- CI/CD流水线配置

**第二阶段（3周）：核心功能开发**
- 仓储管理模块开发
- 库存管理功能实现
- 基础API接口开发
- 单元测试编写

**第三阶段（3周）：审批引擎开发**
- 通用审批引擎开发
- 审批模板配置功能
- Webhook回调机制实现
- 集成测试

**第四阶段（2周）：地图服务集成**
- 高德地图API集成
- 地理编码缓存实现
- 路径规划功能开发
- 性能优化

**第五阶段（2周）：前端开发**
- 管理界面开发
- 地图展示功能
- 审批流程界面
- 用户体验优化

**第六阶段（1周）：测试和上线**
- 端到端测试
- 性能测试
- 安全测试
- 生产环境部署

### 15.2 回滚计划

**回滚触发条件：**
- 系统性能下降超过30%
- 关键功能异常率超过5%
- 数据一致性问题
- 安全漏洞发现

**回滚步骤：**
```bash
#!/bin/bash
# 系统回滚脚本

# 1. 停止新服务
systemctl stop logistics-api
systemctl stop approval-engine

# 2. 恢复数据库
mysql -u root -p < backup_rollback.sql

# 3. 恢复Redis缓存
redis-cli FLUSHALL
redis-cli --rdb /backup/redis/rollback.rdb

# 4. 启动旧服务
systemctl start old-logistics-service

# 5. 验证系统状态
curl -f http://localhost:8080/health || exit 1

echo "系统回滚完成"
```

### 15.3 项目时间估算

**开发资源需求：**

| 模块 | 开发人员 | 测试人员 | 运维人员 | 预估工期 |
|------|----------|----------|----------|----------|
| 物流核心模块 | 3人 | 1人 | - | 3周 |
| 审批引擎 | 2人 | 1人 | - | 3周 |
| 地图服务 | 2人 | 1人 | - | 2周 |
| 前端开发 | 3人 | 1人 | - | 3周 |
| 数据库设计 | 1人 | - | 1人 | 1周 |
| 部署运维 | - | - | 2人 | 2周 |
| 测试验证 | - | 2人 | 1人 | 2周 |

**总体时间线：**
- **总开发周期：** 13周
- **测试周期：** 4周（与开发并行）
- **部署上线：** 1周
- **项目总时长：** 14周

**关键里程碑：**
- Week 2: 数据库设计完成
- Week 5: 核心物流功能完成
- Week 8: 审批引擎完成
- Week 10: 地图服务集成完成
- Week 13: 前端开发完成
- Week 14: 系统测试通过，准备上线

**风险评估和缓冲时间：**
- 技术风险缓冲：+2周
- 人员变动缓冲：+1周
- 需求变更缓冲：+2周
- **建议项目周期：** 19周（含缓冲）

### 15.4 成功标准

**技术指标：**
- API响应时间P95 < 500ms
- 系统可用性 > 99.9%
- 数据库查询优化率 > 95%
- 缓存命中率 > 90%

**业务指标：**
- 入库审批时效提升50%
- 库存准确率 > 99.5%
- 运输路径优化成本降低20%
- 用户满意度 > 90%

**质量指标：**
- 代码测试覆盖率 > 85%
- 安全漏洞数量 = 0
- 生产环境故障数 < 5次/月
- 文档完整性 > 95%

---

## 16. 术语表

### 16.1 技术术语

| 术语 | 中文解释 | 应用场景 |
|------|----------|----------|
| Redis | 内存数据库 | 用于缓存地理编码数据，提高查询速度 |
| Webhook | 网络回调机制 | 审批完成后自动通知业务系统 |
| State Machine | 状态机 | 管理入库单、运输单等业务状态流转 |
| API | 应用程序接口 | 系统间数据交互的标准接口 |
| Database Index | 数据库索引 | 加速数据库查询，提升系统性能 |
| Cache Hit Rate | 缓存命中率 | 衡量缓存效果的重要指标 |
| Load Balancer | 负载均衡器 | 分散系统负载，提高可用性 |
| Microservices | 微服务架构 | 将系统拆分为独立的小服务 |
| Docker | 容器化技术 | 简化部署和环境管理 |
| CI/CD | 持续集成/持续部署 | 自动化代码构建和部署流程 |

### 16.2 业务术语

| 术语 | 中文解释 | 业务含义 |
|------|----------|----------|
| 入库单 | Inbound Order | 货物进入仓库的凭证单据 |
| 出库单 | Outbound Order | 货物离开仓库的凭证单据 |
| 库存盘点 | Inventory Counting | 定期检查实际库存与系统记录是否一致 |
| 批次号 | Batch Number | 同一批次货物的唯一标识 |
| 库位 | Storage Location | 仓库内具体的存储位置 |
| 审批流程 | Approval Process | 需要多人确认的业务流程 |
| 审批节点 | Approval Node | 审批流程中的具体审批环节 |
| 地理编码 | Geocoding | 将地址转换为经纬度坐标的过程 |
| 路径规划 | Route Planning | 为运输车辆规划最优行驶路线 |
| 电子围栏 | Geofencing | 基于地理位置的虚拟边界 |

## 17. 常见问题解答（FAQ）

### 17.1 审批引擎相关问题

**Q: 审批引擎是否支持复杂的审批条件？**
A: 是的，审批引擎支持基于金额、部门、商品类型等多种条件的动态审批规则配置，可以通过可视化界面灵活设置。

**Q: 如果审批人不在公司，如何处理审批？**
A: 系统支持多种审批方式：
- 邮件审批：通过邮件链接直接审批
- 移动端审批：支持手机APP审批
- 代理审批：可设置临时代理人
- 自动转交：超时自动转交给上级

**Q: 审批流程是否可以中途修改？**
A: 可以。系统支持流程版本管理，正在进行的审批使用原版本，新创建的审批使用新版本，确保业务连续性。

### 17.2 地图服务相关问题

**Q: 高德地图API调用成本如何控制？**
A: 系统采用多重成本控制策略：
- Redis缓存：相同地址24小时内只调用一次API
- 频率限制：每分钟最多100次调用
- 智能去重：批量地址去重后统一查询
- 成本监控：实时统计API调用费用

**Q: 地图服务是否支持离线使用？**
A: 基础的地图展示功能支持离线使用，但地址解析和路径规划需要网络连接。建议在有网络时预先缓存常用地址。

**Q: 如何保证地址解析的准确性？**
A: 系统采用多级验证机制：
- 多源对比：结合多个地图服务验证
- 人工校验：重要地址支持人工确认
- 精度分级：根据业务需求选择不同精度
- 历史记录：保留解析历史便于追溯

### 17.3 数据迁移相关问题

**Q: 现有数据迁移是否会影响业务？**
A: 数据迁移采用增量迁移策略：
- 先迁移历史数据，不影响现有业务
- 业务切换时只迁移增量数据
- 支持双系统并行运行一段时间
- 提供完整的回滚方案

**Q: 迁移过程中数据安全如何保障？**
A: 多重安全保障措施：
- 数据加密传输和存储
- 迁移过程全程日志记录
- 数据一致性自动校验
- 备份数据多重保存

**Q: 迁移失败如何处理？**
A: 完善的失败处理机制：
- 自动检测数据一致性
- 支持断点续传
- 快速回滚到迁移前状态
- 详细的错误报告和修复建议

### 17.4 系统性能相关问题

**Q: 系统能支持多少并发用户？**
A: 根据性能测试结果：
- 支持1000+并发用户同时使用
- API响应时间P95 < 500ms
- 数据库连接池利用率 < 80%
- 可通过水平扩展进一步提升性能

**Q: 大量数据查询会不会影响系统性能？**
A: 系统采用多重优化策略：
- 数据库索引优化，加速查询
- Redis缓存热点数据
- 分页查询避免大量数据传输
- 读写分离分散数据库压力

**Q: 系统升级是否需要停机？**
A: 采用蓝绿部署策略：
- 零停机时间升级
- 自动健康检查
- 快速回滚机制
- 灰度发布降低风险

### 17.5 运维管理相关问题

**Q: 系统监控包含哪些内容？**
A: 全方位监控体系：
- 业务指标：订单处理时效、库存准确率
- 技术指标：API响应时间、系统资源使用率
- 安全指标：异常登录、权限变更
- 成本指标：API调用费用、服务器成本

**Q: 如何快速定位和解决问题？**
A: 完善的问题诊断体系：
- 详细的操作日志记录
- 自动异常检测和告警
- 标准化的故障处理流程
- 24小时技术支持团队

**Q: 系统备份策略是怎样的？**
A: 多层备份保障：
- 数据库每日全量备份
- 实时增量备份
- 异地备份存储
- 定期备份恢复演练

---

## 18. 联系方式和支持

### 18.1 技术支持团队

| 角色 | 联系方式 | 支持时间 |
|------|----------|----------|
| 项目经理 | project-manager@company.com | 工作日 9:00-18:00 |
| 技术负责人 | tech-lead@company.com | 7×24小时紧急支持 |
| 运维工程师 | ops@company.com | 7×24小时 |
| 业务顾问 | business@company.com | 工作日 9:00-18:00 |

### 18.2 紧急联系流程

1. **系统故障**：立即拨打运维热线 400-xxx-xxxx
2. **数据问题**：联系技术负责人，启动数据恢复流程
3. **安全问题**：立即联系安全团队，启动应急响应
4. **业务咨询**：通过企业微信群或邮件联系业务顾问

### 18.3 文档维护

本文档由物流系统开发组负责维护，如有更新建议或问题反馈，请通过以下方式联系：
- 邮箱：docs-feedback@company.com
- 内部工单系统：创建"文档更新"工单
- 版本控制：Git仓库提交Pull Request

---

**文档版本：** v1.0  
**最后更新：** 2024-01-15  
**维护团队：** 物流系统开发组  
**项目周期：** 14-19周（根据风险缓冲调整）  
**文档状态：** 最终版，已通过技术评审