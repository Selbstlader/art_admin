# Requirements Document

## Introduction

本功能模块是工装设计师AI辅助系统的移动端应用，基于 UniApp 开发，支持 iOS、Android 和微信小程序三端。作为 Web 端的轻量级辅助工具，主要面向设计师在外出、没带电脑时快速查看项目数据和进行简单操作的场景。App 复用现有后端 API，与 Web 端数据完全同步。

**平台差异说明：**
- iOS/Android App：功能完整，包含 AI 对话功能
- 微信小程序：受平台限制，AI 对话功能不可用，其他功能正常

## Glossary

- **Designer_App**: 工装设计师AI辅助系统移动端应用
- **Web_System**: 现有的 Web 端工装设计师AI辅助系统
- **Project**: 设计项目，包含项目基本信息、文档、设计图等
- **AI_Chat**: AI 对话功能，基于火山AI提供设计咨询服务
- **Push_Notification**: 推送通知，用于项目更新、审批提醒等
- **Offline_Cache**: 离线缓存，支持无网络时查看已缓存的数据
- **UniApp**: 基于 Vue.js 的跨平台开发框架，一套代码编译到多个平台
- **Mini_Program**: 微信小程序，受平台限制部分功能不可用
- **Native_App**: 原生 App（iOS/Android），功能完整无限制

## Requirements

### Requirement 1: 用户认证与账号同步

**User Story:** As a 工装设计师, I want to 使用与 Web 端相同的账号登录 App, so that 我可以在移动端访问我的所有项目数据。

#### Acceptance Criteria

1. WHEN 用户打开 App THEN Designer_App SHALL 展示登录页面，支持账号密码登录
2. WHEN 用户登录成功 THEN Designer_App SHALL 获取并缓存用户 token，后续请求自动携带认证信息
3. WHEN 用户 token 过期 THEN Designer_App SHALL 自动跳转登录页并提示"登录已过期，请重新登录"
4. WHEN 用户点击退出登录 THEN Designer_App SHALL 清除本地缓存的 token 和用户数据
5. WHILE 用户已登录 THEN Designer_App SHALL 在首页展示用户头像和名称

### Requirement 2: 项目列表查看

**User Story:** As a 工装设计师, I want to 在手机上快速浏览我的项目列表, so that 我可以随时了解项目概况。

#### Acceptance Criteria

1. WHEN 用户进入项目列表页 THEN Designer_App SHALL 展示项目卡片列表，每个卡片包含项目名称、状态、面积、最近更新时间
2. WHEN 用户下拉刷新 THEN Designer_App SHALL 重新请求项目列表并更新显示
3. WHEN 用户上拉加载 THEN Designer_App SHALL 分页加载更多项目（每页 10 条）
4. WHEN 用户在搜索框输入关键字 THEN Designer_App SHALL 按项目名称实时过滤列表
5. WHEN 项目列表为空 THEN Designer_App SHALL 展示空状态提示"暂无项目"

### Requirement 3: 项目详情查看

**User Story:** As a 工装设计师, I want to 查看项目的详细信息, so that 我可以在外出时快速了解项目情况。

#### Acceptance Criteria

1. WHEN 用户点击项目卡片 THEN Designer_App SHALL 跳转到项目详情页，展示项目基本信息（名称、描述、面积、预算、风格、状态）
2. WHEN 项目详情加载完成 THEN Designer_App SHALL 展示项目关联的文档数量、设计图数量、最近分析时间
3. WHEN 用户点击"查看文档分析" THEN Designer_App SHALL 展示该项目的文档分析结果摘要
4. WHEN 用户点击"查看成本预算" THEN Designer_App SHALL 展示该项目的成本估算报告
5. IF 项目数据加载失败 THEN Designer_App SHALL 展示错误提示并提供重试按钮

### Requirement 4: AI 对话咨询（仅 App）

**User Story:** As a 工装设计师, I want to 在手机上与 AI 进行设计相关对话, so that 我可以随时随地咨询设计问题。

#### Acceptance Criteria

1. WHERE 平台为 iOS 或 Android App THEN Designer_App SHALL 在底部导航栏展示"AI对话"入口
2. WHERE 平台为微信小程序 THEN Designer_App SHALL 隐藏"AI对话"入口，不展示该功能
3. WHEN 用户进入 AI 对话页面 THEN Designer_App SHALL 展示对话界面，底部有输入框和发送按钮
4. WHEN 用户发送消息 THEN Designer_App SHALL 立即显示用户消息气泡，并展示 AI 正在输入的加载状态
5. WHEN AI 返回回复 THEN Designer_App SHALL 以流式方式逐字展示 AI 回复内容
6. WHEN 用户在项目详情页发起对话 THEN Designer_App SHALL 自动携带项目上下文，AI 回复基于当前项目信息
7. WHEN 用户查看历史对话 THEN Designer_App SHALL 支持上拉加载更多历史消息
8. IF AI 服务响应超时 THEN Designer_App SHALL 展示"AI 服务暂时不可用，请稍后重试"

### Requirement 5: 文档分析结果查看

**User Story:** As a 工装设计师, I want to 在手机上查看文档分析结果, so that 我可以快速回顾项目需求要点。

#### Acceptance Criteria

1. WHEN 用户进入文档分析页面 THEN Designer_App SHALL 展示文档列表，每项显示文件名、分析状态、分析时间
2. WHEN 用户点击已分析的文档 THEN Designer_App SHALL 展示分析结果，包含关键字提取和文档摘要
3. WHEN 展示文档摘要 THEN Designer_App SHALL 按"项目概述、核心需求、特殊要求"三部分结构化展示
4. WHEN 文档尚未分析 THEN Designer_App SHALL 展示"该文档尚未分析"状态，不提供分析操作（需在 Web 端进行）

### Requirement 6: 成本预算报告查看

**User Story:** As a 工装设计师, I want to 在手机上查看项目成本预算, so that 我可以在与客户沟通时快速提供预算参考。

#### Acceptance Criteria

1. WHEN 用户进入成本报告页面 THEN Designer_App SHALL 展示成本汇总卡片，显示总成本和预算上限
2. WHEN 总成本超过预算上限 THEN Designer_App SHALL 以红色警告样式展示超支金额
3. WHEN 用户查看成本明细 THEN Designer_App SHALL 按类别（材料费、人工费、设备费、管理费）分项展示
4. WHEN 用户点击某个成本类别 THEN Designer_App SHALL 展开显示该类别下的具体明细项

### Requirement 7: 消息通知中心（仅 App）

**User Story:** As a 工装设计师, I want to 接收项目相关的推送通知, so that 我可以及时了解项目动态。

#### Acceptance Criteria

1. WHERE 平台为 iOS 或 Android App THEN Designer_App SHALL 在底部导航栏展示"消息"入口
2. WHERE 平台为微信小程序 THEN Designer_App SHALL 使用微信订阅消息替代推送通知
3. WHEN App 启动时 THEN Designer_App SHALL 请求推送通知权限（iOS/Android）
4. WHEN 有新的项目更新 THEN Designer_App SHALL 发送推送通知，点击通知跳转到对应项目
5. WHEN 用户进入消息中心 THEN Designer_App SHALL 展示通知列表，按时间倒序排列
6. WHEN 用户点击某条通知 THEN Designer_App SHALL 标记为已读并跳转到相关页面
7. WHEN 用户有未读消息 THEN Designer_App SHALL 在底部导航栏消息图标上显示红点

### Requirement 8: 离线数据缓存

**User Story:** As a 工装设计师, I want to 在无网络时也能查看已加载的数据, so that 我在网络不好的环境下仍能使用 App。

#### Acceptance Criteria

1. WHEN 项目列表加载成功 THEN Designer_App SHALL 将数据缓存到本地存储
2. WHEN 网络不可用时打开 App THEN Designer_App SHALL 展示缓存的项目列表，并提示"当前为离线模式"
3. WHEN 网络恢复 THEN Designer_App SHALL 自动刷新数据并更新缓存
4. WHEN 缓存数据超过 7 天 THEN Designer_App SHALL 在下次联网时自动更新

### Requirement 9: 设计图预览

**User Story:** As a 工装设计师, I want to 在手机上预览设计图, so that 我可以快速向客户展示设计方案。

#### Acceptance Criteria

1. WHEN 用户进入设计图列表 THEN Designer_App SHALL 展示项目关联的设计图缩略图列表
2. WHEN 用户点击设计图 THEN Designer_App SHALL 全屏展示设计图，支持双指缩放和拖动查看
3. WHEN 用户左右滑动 THEN Designer_App SHALL 切换到上一张/下一张设计图
4. WHEN 设计图加载中 THEN Designer_App SHALL 展示加载进度指示器
5. IF 设计图加载失败 THEN Designer_App SHALL 展示占位图并提示"图片加载失败"

### Requirement 10: 材料库快速查询

**User Story:** As a 工装设计师, I want to 在手机上快速查询材料信息, so that 我可以在与客户或供应商沟通时提供材料参考。

#### Acceptance Criteria

1. WHEN 用户进入材料库页面 THEN Designer_App SHALL 展示材料分类列表（地板、墙面、天花、家具等）
2. WHEN 用户选择分类 THEN Designer_App SHALL 展示该分类下的材料列表
3. WHEN 用户搜索材料 THEN Designer_App SHALL 支持按材料名称、品牌进行模糊搜索
4. WHEN 用户点击材料 THEN Designer_App SHALL 展示材料详情（名称、规格、单价、供应商、适用场景）
5. WHEN 展示材料详情 THEN Designer_App SHALL 显示材料图片，支持点击放大查看

### Requirement 11: 我的页面

**User Story:** As a 工装设计师, I want to 查看我的基本信息和退出登录, so that 我可以确认当前登录状态。

#### Acceptance Criteria

1. WHEN 用户进入"我的"页面 THEN Designer_App SHALL 展示用户头像、名称、角色信息
2. WHEN 用户点击"退出登录" THEN Designer_App SHALL 弹出确认弹窗，确认后退出并跳转登录页
3. WHEN 用户点击"关于" THEN Designer_App SHALL 展示 App 版本号信息
