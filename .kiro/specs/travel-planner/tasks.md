# Implementation Plan

## Phase 1: 基础架构搭建

- [x] 1. 后端基础架构
  - [x] 1.1 创建旅游模块目录结构和基础文件
    - 创建 internal/api/travel、internal/model、internal/repository、internal/service/travel 目录
    - 创建基础的路由注册文件
    - _Requirements: 14.1_
  - [x] 1.2 创建数据库迁移文件
    - 创建 roadbooks、roadbook_waypoints、roadbook_tags、roadbook_tag_relations 表
    - 创建 roadbook_comments、roadbook_favorites、roadbook_shares 表
    - 创建 roadbook_templates、roadbook_stats_daily 表
    - _Requirements: 14.1_
  - [x] 1.3 创建数据模型 (Model)
    - 实现 Roadbook、RoadbookWaypoint、RoadbookTag 模型
    - 实现 RoadbookComment、RoadbookFavorite、RoadbookShare 模型
    - _Requirements: 14.1_
  - [x] 1.4 编写属性测试：JSON序列化Round-Trip
    - **Property 1: 路书数据持久化 Round-Trip**
    - **Validates: Requirements 3.2, 14.3**

- [x] 2. 前端基础架构
  - [x] 2.1 创建旅游模块目录结构
    - 创建 src/views/travel 及子目录
    - 创建 components、composables、store、api、types 目录
    - _Requirements: 1.1_
  - [x] 2.2 创建TypeScript类型定义
    - 定义 Roadbook、Waypoint、Comment、Tag 等接口
    - 定义 API 请求和响应类型
    - _Requirements: 1.1_
  - [x] 2.3 创建API接口文件
    - 实现路书相关API调用函数
    - 实现评论、收藏、分享API调用函数
    - _Requirements: 1.1_
  - [x] 2.4 创建Pinia状态管理
    - 创建 roadbook store 管理路书状态
    - 创建 map store 管理地图状态
    - _Requirements: 1.1_

- [x] 3. Checkpoint - 确保基础架构完成
  - 确保所有测试通过，如有问题请询问用户

## Phase 2: 路书核心功能

- [x] 4. 路书CRUD后端实现
  - [x] 4.1 实现路书Repository层
    - 实现 Create、GetByID、Update、Delete、List 方法
    - 实现软删除和分页查询
    - _Requirements: 3.2, 3.3, 3.4, 14.2_
  - [x] 4.2 编写属性测试：路书列表排序
    - **Property 2: 路书列表按时间倒序排列**
    - **Validates: Requirements 3.3**
  - [x] 4.3 编写属性测试：路书软删除
    - **Property 3: 路书软删除后不可查询**
    - **Validates: Requirements 3.4**
  - [x] 4.4 实现路书Service层
    - 实现创建、更新、删除路书业务逻辑
    - 实现路书列表查询和详情查询
    - 实现事务处理确保数据一致性
    - _Requirements: 3.2, 3.3, 3.4, 14.1_
  - [x] 4.5 编写属性测试：事务一致性
    - **Property 17: 事务一致性保证**
    - **Validates: Requirements 14.1**
  - [x] 4.6 实现路书API处理器
    - 实现 POST/GET/PUT/DELETE /api/v1/travel/roadbooks 接口
    - 实现参数验证和错误处理
    - _Requirements: 3.1, 3.2, 3.3, 3.4_
  - [x] 4.7 编写属性测试：分页查询限制
    - **Property 18: 分页查询数量限制**
    - **Validates: Requirements 14.2**

- [x] 5. 途经点功能实现
  - [x] 5.1 实现途经点Repository层
    - 实现途经点的CRUD操作
    - 实现批量更新和排序功能
    - _Requirements: 3.5_
  - [x] 5.2 编写属性测试：途经点数据完整性
    - **Property 4: 途经点数据完整性**
    - **Validates: Requirements 3.5**
  - [x] 5.3 实现途经点Service层
    - 实现添加、更新、删除途经点
    - 实现按天分组查询
    - _Requirements: 3.5, 4.3_
  - [x] 5.4 编写属性测试：途经点按天分组
    - **Property 5: 途经点按天分组正确性**
    - **Validates: Requirements 4.3**
  - [x] 5.5 实现途经点API处理器
    - 实现途经点相关API接口
    - _Requirements: 3.5_

- [x] 6. Checkpoint - 确保路书核心功能完成
  - 确保所有测试通过，如有问题请询问用户

## Phase 3: 前端路书功能

- [x] 7. 高德地图集成
  - [x] 7.1 创建高德地图容器组件
    - 实现 AMapContainer.vue 组件
    - 集成高德地图JS API 2.0
    - _Requirements: 1.1, 1.2_
  - [x] 7.2 创建地图Hook (useAMap)
    - 实现地图初始化、标记管理、路线绘制功能
    - 实现POI搜索功能
    - _Requirements: 1.3, 1.4, 2.1_
  - [x] 7.3 创建路线规划Hook (useRoute)
    - 实现路线规划API调用
    - 实现多种出行方式切换
    - _Requirements: 2.1, 2.2, 2.3, 2.5_

- [x] 8. 路书列表页面
  - [x] 8.1 创建路书列表页面
    - 实现路书卡片组件
    - 实现分页加载
    - _Requirements: 3.3_
  - [x] 8.2 创建路书筛选和搜索功能
    - 实现按标签筛选
    - 实现关键词搜索
    - _Requirements: 8.2, 8.3_

- [x] 9. 路书编辑器
  - [x] 9.1 创建路书编辑器页面
    - 实现路书基本信息表单
    - 实现封面图片上传
    - _Requirements: 3.1_
  - [x] 9.2 创建途经点列表组件
    - 实现途经点拖拽排序
    - 实现途经点添加、编辑、删除
    - _Requirements: 2.4, 3.5_
  - [x] 9.3 创建标签选择器组件
    - 实现预设标签选择
    - 实现自定义标签创建
    - _Requirements: 3.7_
  - [x] 9.4 创建日程时间线组件
    - 实现多日行程展示
    - 实现按天分组显示
    - _Requirements: 4.3_

- [x] 10. 路书详情页面
  - [x] 10.1 创建路书详情页面
    - 实现路书信息展示
    - 实现地图路线展示
    - _Requirements: 4.1, 4.2_
  - [x] 10.2 实现收藏和点赞功能
    - 实现收藏按钮和状态
    - 实现点赞按钮和状态
    - _Requirements: 4.4_

- [x] 11. Checkpoint - 确保前端路书功能完成
  - 确保所有测试通过，如有问题请询问用户

## Phase 4: 社交功能

- [ ] 12. 评论功能后端
  - [ ] 12.1 实现评论Repository和Service
    - 实现评论CRUD操作
    - 实现评论树形结构查询
    - _Requirements: 5.1, 5.2, 5.3_
  - [ ] 12.2 编写属性测试：评论Round-Trip
    - **Property 7: 评论数据持久化 Round-Trip**
    - **Validates: Requirements 5.1**
  - [ ] 12.3 编写属性测试：评论排序
    - **Property 8: 评论列表按时间倒序排列**
    - **Validates: Requirements 5.2**
  - [ ] 12.4 编写属性测试：评论父子关系
    - **Property 9: 评论父子关系有效性**
    - **Validates: Requirements 5.3**
  - [ ] 12.5 实现评论API处理器
    - 实现评论相关API接口
    - 实现评论内容验证
    - _Requirements: 5.1, 5.5_
  - [ ] 12.6 编写属性测试：评论内容验证
    - **Property 10: 评论内容验证**
    - **Validates: Requirements 5.5**

- [ ] 13. 评论功能前端
  - [ ] 13.1 创建评论列表组件
    - 实现评论展示和回复
    - 实现评论点赞
    - _Requirements: 5.2, 5.3, 5.6_
  - [ ] 13.2 创建评论输入组件
    - 实现评论输入和字数限制
    - 实现评论提交
    - _Requirements: 5.1, 5.5_

- [ ] 14. 收藏功能
  - [ ] 14.1 实现收藏后端
    - 实现收藏Repository和Service
    - 实现收藏计数更新
    - _Requirements: 4.4_
  - [ ] 14.2 编写属性测试：收藏计数
    - **Property 6: 收藏计数正确性**
    - **Validates: Requirements 4.4**

- [ ] 15. Checkpoint - 确保社交功能完成
  - 确保所有测试通过，如有问题请询问用户

## Phase 5: 分享与导航

- [ ] 16. 分享功能
  - [ ] 16.1 实现分享后端
    - 实现分享码生成
    - 实现分享记录和访问统计
    - _Requirements: 7.1, 7.3_
  - [ ] 16.2 编写属性测试：分享码唯一性
    - **Property 14: 分享码唯一性**
    - **Validates: Requirements 7.1**
  - [ ] 16.3 编写属性测试：访问计数
    - **Property 15: 分享访问计数递增**
    - **Validates: Requirements 7.3**
  - [ ] 16.4 创建分享弹窗组件
    - 实现二维码生成
    - 实现链接复制功能
    - _Requirements: 7.2_

- [ ] 17. 高德导航导出
  - [ ] 17.1 实现导航链接生成服务
    - 实现高德URI Scheme生成
    - 实现分段导航逻辑
    - _Requirements: 6.1, 6.2, 6.3_
  - [ ] 17.2 编写属性测试：导航链接格式
    - **Property 11: 高德导航链接格式正确性**
    - **Validates: Requirements 6.1, 6.4**
  - [ ] 17.3 编写属性测试：途经点顺序
    - **Property 12: 导航链接途经点顺序保持**
    - **Validates: Requirements 6.2**
  - [ ] 17.4 编写属性测试：途经点数量限制
    - **Property 13: 途经点数量超限检测**
    - **Validates: Requirements 6.3**
  - [ ] 17.5 创建导航导出弹窗组件
    - 实现导航模式选择
    - 实现分段导航展示
    - _Requirements: 6.4_

- [ ] 18. Checkpoint - 确保分享与导航功能完成
  - 确保所有测试通过，如有问题请询问用户

## Phase 6: AI智能规划

- [ ] 19. 火山AI集成
  - [ ] 19.1 创建AI服务配置
    - 配置火山AI API密钥
    - 创建AI服务客户端
    - _Requirements: 9.1_
  - [ ] 19.2 实现AI规划Service
    - 实现提示词构建
    - 实现AI响应解析
    - 实现行程导入功能
    - _Requirements: 9.1, 9.2, 9.3, 9.4_
  - [ ] 19.3 实现AI规划API
    - 实现 POST /api/v1/travel/ai/plan 接口
    - 实现 POST /api/v1/travel/ai/chat 接口
    - _Requirements: 9.1, 10.3_

- [ ] 20. AI规划前端
  - [ ] 20.1 创建AI规划页面
    - 实现目的地和天数输入
    - 实现偏好选择
    - 实现预算设置
    - _Requirements: 9.1, 9.3_
  - [ ] 20.2 创建AI结果展示组件
    - 实现行程方案展示
    - 实现一键导入功能
    - _Requirements: 9.2, 9.4_
  - [ ] 20.3 创建AI助手对话组件
    - 实现对话界面
    - 实现问答功能
    - _Requirements: 10.3_

- [ ] 21. Checkpoint - 确保AI功能完成
  - 确保所有测试通过，如有问题请询问用户

## Phase 7: 路书发现与模板

- [ ] 22. 路书广场
  - [ ] 22.1 实现搜索Service
    - 实现全文搜索
    - 实现热门排序
    - _Requirements: 8.1, 8.2_
  - [ ] 22.2 编写属性测试：搜索筛选
    - **Property 16: 搜索筛选结果正确性**
    - **Validates: Requirements 8.2, 8.3**
  - [ ] 22.3 创建路书广场页面
    - 实现热门路书展示
    - 实现搜索和筛选功能
    - _Requirements: 8.1, 8.2, 8.3_

- [ ] 23. 路书模板
  - [ ] 23.1 实现模板后端
    - 实现模板Repository和Service
    - 实现模板使用功能
    - _Requirements: 11.1, 11.2, 11.3_
  - [ ] 23.2 创建模板库页面
    - 实现模板列表展示
    - 实现模板使用功能
    - _Requirements: 11.1, 11.2_

- [ ] 24. Checkpoint - 确保发现与模板功能完成
  - 确保所有测试通过，如有问题请询问用户

## Phase 8: 权限与管理

- [ ] 25. 权限控制
  - [ ] 25.1 实现权限中间件
    - 实现路书访问权限验证
    - 实现所有权验证
    - _Requirements: 15.1, 15.2, 15.3_
  - [ ] 25.2 编写属性测试：权限验证
    - **Property 19: 私有路书访问权限验证**
    - **Validates: Requirements 15.1, 15.2**
  - [ ] 25.3 编写属性测试：所有权验证
    - **Property 20: 路书所有权验证**
    - **Validates: Requirements 15.3**

- [ ] 26. 数据统计
  - [ ] 26.1 实现统计Service
    - 实现浏览量、收藏量统计
    - 实现趋势数据查询
    - _Requirements: 12.1, 12.2_
  - [ ] 26.2 创建统计页面
    - 实现统计图表展示
    - 实现个人数据中心
    - _Requirements: 12.1, 12.2_

- [ ] 27. 后台管理
  - [ ] 27.1 实现管理API
    - 实现路书审核API
    - 实现标签管理API
    - _Requirements: 13.1, 13.2, 13.3_
  - [ ] 27.2 创建管理页面
    - 实现路书管理列表
    - 实现标签管理页面
    - _Requirements: 13.1, 13.2, 13.3_

- [ ] 28. Final Checkpoint - 确保所有功能完成
  - 确保所有测试通过，如有问题请询问用户
