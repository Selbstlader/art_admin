# WebRTC 视频通话功能 - 本地测试指南

## ✅ 已完成的工作

### 1. 前端实现
- ✅ WebRTC 工具类 (`/src/utils/webrtc/`)
  - `types.ts` - 类型定义和常量
  - `peer-manager.ts` - P2P 连接管理
  - `webrtc-manager.ts` - WebRTC 主管理器
  - `index.ts` - 统一导出

- ✅ 独立视频通话页面 (`/src/views/video-call/index.vue`)
  - 房间列表和创建功能
  - 支持 1-4 人视频网格布局
  - 音视频控制按钮
  - 实时状态显示
  - 完整的通话流程

- ✅ 路由配置 (`/src/router/modules/video-call.ts`)
  - 路径: `/video-call`
  - 独立页面,不依赖聊天室

### 2. 后端实现
- ✅ WebSocket 信令服务器扩展 (`/internal/pkg/websocket/client.go`)
  - 支持 WebRTC Offer/Answer/ICE Candidate 转发
  - 支持通话控制消息
  - 支持媒体状态同步

## 🚀 本地测试步骤

### 前置条件
1. **浏览器要求**: Chrome/Edge/Firefox (推荐 Chrome)
2. **HTTPS 或 localhost**: WebRTC 需要安全上下文
3. **摄像头和麦克风**: 确保设备正常工作

### 步骤 1: 启动后端服务

```bash
cd art_admin_backend
go run cmd/server/main.go
```

后端将在 `http://localhost:48080` 启动

### 步骤 2: 启动前端服务

```bash
cd art-design-pro
pnpm dev
```

前端将在 `http://localhost:5173` 启动

### 步骤 3: 登录系统

使用提供的 Token 或登录:
- 用户名: `admin`
- 密码: (根据你的系统配置)

### 步骤 4: 访问视频通话页面

1. 在浏览器中访问: `http://localhost:5173/video-call`
2. 或在系统菜单中找到"视频通话"入口

### 步骤 5: 创建或加入房间

**方式 1: 创建新房间**
1. 在页面顶部输入房间名称
2. 点击"创建房间"按钮
3. 自动进入视频通话界面

**方式 2: 加入现有房间**
1. 在房间列表中选择一个房间
2. 点击"加入"按钮
3. 进入视频通话界面

### 步骤 6: 多设备测试

#### 方法 1: 同一台电脑多个浏览器窗口
1. 打开 Chrome 浏览器,登录用户 A
2. 打开 Chrome 隐身窗口,登录用户 B (或使用不同浏览器)
3. 两个用户都进入同一个聊天室
4. 任一用户点击"开始视频通话"
5. 观察视频连接建立过程

#### 方法 2: 不同设备测试
1. 电脑 A: 登录用户 A,进入聊天室
2. 电脑 B: 登录用户 B,进入同一聊天室
3. 任一用户发起视频通话

## 🔍 测试检查点

### 1. 媒体设备权限
- [ ] 浏览器弹出摄像头/麦克风权限请求
- [ ] 允许权限后能看到本地视频预览
- [ ] 本地视频镜像翻转显示

### 2. WebSocket 连接
- [ ] 打开浏览器开发者工具 -> Network -> WS
- [ ] 确认 WebSocket 连接成功 (`ws://localhost:48080/api/chat/ws?roomId=xxx`)
- [ ] 查看信令消息收发 (offer, answer, ice_candidate)

### 3. P2P 连接
- [ ] 打开 Chrome DevTools -> Console
- [ ] 查看 WebRTC 连接日志
- [ ] 确认 ICE 连接状态变为 "connected"

### 4. 视频通话功能
- [ ] 能看到对方的视频画面
- [ ] 能听到对方的声音
- [ ] 点击静音按钮,对方听不到声音
- [ ] 点击关闭摄像头,对方看不到画面
- [ ] 点击挂断,通话正常结束

### 5. 多人测试 (3-4人)
- [ ] 3 个用户同时加入聊天室
- [ ] 视频网格布局正确显示
- [ ] 每个用户都能看到其他人的视频
- [ ] 任一用户离开,其他人画面自动调整

## 🐛 常见问题排查

### 问题 1: 无法获取摄像头/麦克风

**原因**: 
- 浏览器权限被拒绝
- 设备被其他应用占用
- 非 HTTPS 环境 (localhost 除外)

**解决**:
```bash
# 检查浏览器权限设置
chrome://settings/content/camera
chrome://settings/content/microphone

# 确保使用 localhost 或 HTTPS
```

### 问题 2: WebSocket 连接失败

**检查**:
```bash
# 后端是否运行
curl http://localhost:48080/api/user/info \
  -H 'Authorization: Bearer YOUR_TOKEN'

# WebSocket 端点是否正确
ws://localhost:48080/api/chat/ws?roomId=1
```

### 问题 3: 看不到对方视频

**调试步骤**:
1. 打开 Chrome DevTools -> Console
2. 查找错误信息
3. 检查 WebRTC 连接状态:
```javascript
// 在 Console 中执行
console.log('ICE Connection State:', peerConnection.iceConnectionState)
console.log('Connection State:', peerConnection.connectionState)
```

### 问题 4: ICE 连接失败

**原因**: NAT 穿透失败

**解决**: 配置 TURN 服务器
```typescript
// 在 types.ts 中添加 TURN 服务器
export const DEFAULT_ICE_SERVERS: RTCIceServer[] = [
  { urls: 'stun:stun.l.google.com:19302' },
  {
    urls: 'turn:your-turn-server.com:3478',
    username: 'username',
    credential: 'password'
  }
]
```

### 问题 5: 视频卡顿或延迟

**优化**:
1. 降低视频分辨率
```typescript
// 修改 DEFAULT_MEDIA_CONSTRAINTS
video: {
  width: { ideal: 640 },
  height: { ideal: 480 },
  frameRate: { ideal: 15 }
}
```

2. 检查网络带宽
3. 使用有线网络代替 WiFi

## 📊 调试工具

### Chrome WebRTC Internals
访问 `chrome://webrtc-internals/` 查看详细的 WebRTC 统计信息:
- ICE 候选信息
- 连接状态
- 音视频码率
- 丢包率

### 后端日志
查看后端控制台输出的 WebSocket 消息:
```
Client registered: UserID=1, Username=admin, RoomID=1
Received signal: webrtc_offer from: 1
Received signal: webrtc_answer from: 2
```

## 🎯 下一步优化建议

### 短期 (1-2 天)
1. [ ] 添加通话邀请/接听 UI
2. [ ] 显示网络质量指示器
3. [ ] 添加屏幕共享功能

### 中期 (1 周)
1. [ ] 通话录制功能
2. [ ] 聊天记录与视频通话集成
3. [ ] 美颜滤镜

### 长期 (2-4 周)
1. [ ] SFU 架构支持更多人 (10+ 人)
2. [ ] 虚拟背景
3. [ ] 云端录制和回放

## 📞 技术支持

如遇到问题:
1. 查看浏览器 Console 错误信息
2. 查看后端日志
3. 检查 `chrome://webrtc-internals/`
4. 参考本文档的常见问题部分

## 🔗 相关资源

- [WebRTC API 文档](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API)
- [STUN/TURN 服务器列表](https://gist.github.com/sagivo/3a4b2f2c7ac6e1b5267c2f1f59ac6c6b)
- [WebRTC 最佳实践](https://webrtc.org/getting-started/overview)

---

**测试愉快! 🎉**
