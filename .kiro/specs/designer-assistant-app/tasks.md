# Implementation Plan: Designer Assistant App

## Overview

本实现计划将工装设计师AI辅助系统移动端应用分解为可执行的开发任务。基于 UniApp (Vue 3 + TypeScript) 框架，实现 iOS、Android 和微信小程序三端应用，复用现有后端 API。

## Tasks

- [x] 1. 项目初始化与基础架构搭建
  - [x] 1.1 创建 UniApp 项目并配置开发环境
    - 使用 HBuilderX 或 CLI 创建 Vue 3 + TypeScript + Vite 项目
    - 配置 manifest.json（App 权限、小程序 AppID）
    - 配置 pages.json（页面路由、tabBar）
    - _Requirements: 全局配置_

  - [x] 1.2 搭建项目目录结构
    - 创建 pages/、components/、store/、api/、utils/ 目录
    - 配置 TypeScript 类型定义文件
    - _Requirements: 架构设计_

  - [x] 1.3 实现网络请求封装
    - 封装 uni.request 为统一请求工具
    - 实现请求拦截器（自动携带 token）
    - 实现响应拦截器（统一错误处理、token 过期处理）
    - _Requirements: 1.2, 1.3_

  - [ ]* 1.4 编写请求工具单元测试
    - 测试请求拦截器 token 注入
    - 测试响应拦截器错误处理
    - **Property 1: Token 管理一致性**
    - **Validates: Requirements 1.2, 1.4**

  - [x] 1.5 实现平台检测工具
    - 实现 isApp()、isMiniProgram()、canUseAIChat() 等方法
    - 使用条件编译处理平台差异
    - _Requirements: 4.1, 4.2, 7.1_

  - [ ]* 1.6 编写平台检测工具测试
    - **Property 19: 平台功能可用性**
    - **Validates: Requirements 4.1, 4.2, 7.1**

- [x] 2. Checkpoint - 确保基础架构完成
  - 确保所有测试通过，如有问题请询问用户

- [x] 3. 用户认证模块
  - [x] 3.1 实现 User Store (Pinia)
    - 定义 UserState 接口
    - 实现 login、logout、getUserInfo actions
    - 实现 token 本地存储管理
    - _Requirements: 1.2, 1.4_

  - [x] 3.2 实现登录页面 (pages/login/index.vue)
    - 实现账号密码登录表单
    - 实现登录按钮和加载状态
    - 登录成功后跳转首页
    - _Requirements: 1.1, 1.2_

  - [ ]* 3.3 编写 User Store 属性测试
    - **Property 1: Token 管理一致性**
    - **Validates: Requirements 1.2, 1.4**

- [x] 4. 项目列表模块
  - [x] 4.1 实现 Project Store (Pinia)
    - 定义 ProjectState 接口
    - 实现 fetchList、loadMore、search actions
    - 实现本地缓存管理
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 8.1_

  - [x] 4.2 实现缓存管理工具
    - 实现 setCache、getCache、isExpired 方法
    - 缓存过期时间设为 7 天
    - _Requirements: 8.1, 8.4_

  - [ ]* 4.3 编写缓存管理属性测试
    - **Property 14: 缓存数据一致性**
    - **Property 15: 缓存过期判断**
    - **Validates: Requirements 8.1, 8.4**

  - [x] 4.4 实现 ProjectCard 组件
    - 悬浮卡片设计，蓝色渐变背景
    - 显示项目名称、状态、面积、更新时间
    - 点击触发 @click 事件
    - _Requirements: 2.1_

  - [x] 4.5 实现首页项目列表 (pages/index/index.vue)
    - 实现搜索框和实时过滤
    - 实现下拉刷新和上拉加载
    - 实现空状态展示
    - 实现离线模式提示
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 8.2_

  - [ ]* 4.6 编写项目列表属性测试
    - **Property 2: 项目列表渲染正确性**
    - **Property 3: 分页加载数据增长**
    - **Property 4: 搜索过滤结果匹配**
    - **Validates: Requirements 2.1, 2.3, 2.4**

- [x] 5. Checkpoint - 确保项目列表模块完成
  - 确保所有测试通过，如有问题请询问用户

- [x] 6. 项目详情模块
  - [x] 6.1 实现项目详情页 (pages/project/detail.vue)
    - 展示项目基本信息（名称、描述、面积、预算、风格、状态）
    - 展示文档数量、设计图数量、最近分析时间
    - 提供"查看文档分析"、"查看成本预算"入口
    - 实现加载失败重试
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

  - [ ]* 6.2 编写项目详情属性测试
    - **Property 20: 错误状态显示**
    - **Validates: Requirements 3.5**

- [x] 7. 文档分析模块
  - [x] 7.1 实现文档列表页 (pages/project/documents.vue)
    - 展示文档列表（文件名、分析状态、分析时间）
    - 点击已分析文档跳转详情
    - 未分析文档显示提示状态
    - _Requirements: 5.1, 5.4_

  - [x] 7.2 实现文档分析结果展示
    - 展示关键字提取结果
    - 按"项目概述、核心需求、特殊要求"三部分展示摘要
    - _Requirements: 5.2, 5.3_

  - [ ]* 7.3 编写文档摘要属性测试
    - **Property 8: 文档摘要结构完整性**
    - **Validates: Requirements 5.3**

- [x] 8. 成本预算模块
  - [x] 8.1 实现 CostSummaryCard 组件
    - 3D 立体数据图表设计
    - 显示总成本和预算上限
    - 超支时显示红色警告样式
    - _Requirements: 6.1, 6.2_

  - [x] 8.2 实现成本报告页 (pages/project/cost.vue)
    - 展示成本汇总卡片
    - 按类别（材料费、人工费、设备费、管理费）分项展示
    - 点击类别展开明细
    - _Requirements: 6.1, 6.2, 6.3, 6.4_

  - [ ]* 8.3 编写成本模块属性测试
    - **Property 9: 成本超支警告显示**
    - **Property 10: 成本明细分类正确性**
    - **Validates: Requirements 6.2, 6.3**

- [ ] 9. Checkpoint - 确保项目相关模块完成
  - 确保所有测试通过，如有问题请询问用户

- [x] 10. AI 对话模块 (仅 App)
  - [x] 10.1 实现 Chat Store (Pinia)
    - 定义 ChatState 接口
    - 实现 sendMessage、loadHistory、clearSession actions
    - 实现流式消息处理
    - _Requirements: 4.4, 4.5, 4.7_

  - [x] 10.2 实现 ChatBubble 组件
    - 玻璃拟态设计
    - 区分用户消息和 AI 消息样式
    - 支持流式输出动画
    - _Requirements: 4.4, 4.5_

  - [x] 10.3 实现 AI 对话页面 (pages/chat/index.vue)
    - 使用条件编译仅在 App 端显示
    - 实现对话界面和输入框
    - 实现流式输出展示
    - 实现历史消息上拉加载
    - 实现项目上下文传递
    - 实现超时错误处理
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 4.8_

  - [ ]* 10.4 编写 AI 对话属性测试
    - **Property 5: 消息发送即时显示**
    - **Property 6: 流式输出内容累积**
    - **Property 7: 项目上下文传递**
    - **Validates: Requirements 4.4, 4.5, 4.6**

- [x] 11. 设计图预览模块
  - [x] 11.1 实现设计图列表页 (pages/project/designs.vue)
    - 展示设计图缩略图网格
    - 显示加载进度指示器
    - 加载失败显示占位图
    - _Requirements: 9.1, 9.4, 9.5_

  - [x] 11.2 实现设计图全屏预览
    - 支持双指缩放和拖动
    - 支持左右滑动切换
    - _Requirements: 9.2, 9.3_

  - [ ]* 11.3 编写设计图切换属性测试
    - **Property 16: 设计图切换边界处理**
    - **Validates: Requirements 9.3**

- [ ] 12. Checkpoint - 确保 AI 对话和设计图模块完成
  - 确保所有测试通过，如有问题请询问用户

- [x] 13. 材料库模块
  - [x] 13.1 实现 MaterialCard 组件
    - 显示材料图片、名称、品牌、单价
    - 点击跳转详情
    - _Requirements: 10.1_

  - [x] 13.2 实现材料分类页 (pages/material/index.vue)
    - 展示材料分类列表（地板、墙面、天花、家具等）
    - 点击分类跳转材料列表
    - _Requirements: 10.1_

  - [x] 13.3 实现材料列表页 (pages/material/list.vue)
    - 展示分类下的材料列表
    - 实现搜索功能（名称、品牌模糊搜索）
    - _Requirements: 10.2, 10.3_

  - [x] 13.4 实现材料详情页 (pages/material/detail.vue)
    - 展示材料详情（名称、规格、单价、供应商、适用场景）
    - 图片支持点击放大
    - _Requirements: 10.4, 10.5_

  - [ ]* 13.5 编写材料搜索属性测试
    - **Property 4: 搜索过滤结果匹配**
    - **Property 17: 材料分类过滤正确性**
    - **Validates: Requirements 10.2, 10.3**

- [x] 14. 消息通知模块 (仅 App)
  - [x] 14.1 实现 Notification Store (Pinia)
    - 定义 NotificationState 接口
    - 实现 fetchList、markAsRead、getUnreadCount actions
    - _Requirements: 7.5, 7.6, 7.7_

  - [x] 14.2 实现消息中心页 (pages/message/index.vue)
    - 使用条件编译仅在 App 端显示
    - 展示通知列表（按时间倒序）
    - 点击标记已读并跳转
    - _Requirements: 7.1, 7.5, 7.6_

  - [x] 14.3 实现推送通知权限请求
    - App 启动时请求推送权限
    - 小程序使用订阅消息
    - _Requirements: 7.3, 7.2_

  - [x] 14.4 实现未读消息红点显示
    - 在 tabBar 消息图标上显示红点
    - _Requirements: 7.7_

  - [ ]* 14.5 编写通知模块属性测试
    - **Property 11: 通知列表时间排序**
    - **Property 12: 通知已读状态更新**
    - **Property 13: 未读消息红点显示**
    - **Validates: Requirements 7.5, 7.6, 7.7**

- [x] 15. 我的页面模块
  - [x] 15.1 实现我的页面 (pages/mine/index.vue)
    - 展示用户头像、名称、角色
    - 实现退出登录（确认弹窗）
    - 实现关于页面（版本号）
    - _Requirements: 11.1, 11.2, 11.3, 1.5_

  - [ ]* 15.2 编写用户信息展示属性测试
    - **Property 18: 用户信息展示完整性**
    - **Validates: Requirements 11.1**

- [ ] 16. Checkpoint - 确保所有模块完成
  - 确保所有测试通过，如有问题请询问用户

- [ ] 17. 离线模式与网络状态处理
  - [ ] 17.1 实现网络状态监听
    - 监听网络状态变化
    - 网络恢复时自动刷新数据
    - _Requirements: 8.2, 8.3_

  - [ ] 17.2 实现离线模式 UI 提示
    - 显示"当前为离线模式"提示
    - 网络恢复时自动隐藏
    - _Requirements: 8.2_

- [ ] 18. UI 样式优化
  - [ ] 18.1 实现全局样式变量
    - 定义蓝色渐变色系
    - 定义玻璃拟态样式
    - 定义悬浮卡片阴影
    - _Requirements: UI 设计要求_

  - [ ] 18.2 实现 tabBar 配置
    - 配置底部导航栏（首页、材料库、AI对话*、消息*、我的）
    - *仅 App 端显示
    - _Requirements: 4.1, 7.1_

- [ ] 19. 最终测试与优化
  - [ ] 19.1 运行所有属性测试
    - 确保 20 个正确性属性全部通过
    - _Requirements: 全部_

  - [ ] 19.2 平台兼容性测试
    - iOS 真机测试
    - Android 真机测试
    - 微信小程序测试
    - _Requirements: 全部_

- [ ] 20. Final Checkpoint - 项目完成
  - 确保所有测试通过，所有功能正常运行

## Notes

- 标记 `*` 的任务为可选测试任务，可跳过以加快 MVP 开发
- 每个 Checkpoint 用于验证阶段性成果
- AI 对话和消息通知模块使用条件编译，仅在 App 端可用
- 属性测试使用 fast-check 库，每个测试运行 100 次迭代
