# API 接口文档

## 基础信息
- 基础URL: `http://localhost:48080`
- 认证方式: Bearer Token (JWT)

## 用户管理 API

### 1. 获取用户列表
```http
GET /api/user/list
Authorization: Bearer {token}
```

**Query 参数:**
- `current` (required): 当前页码
- `size` (required): 每页条数
- `userName` (optional): 用户名（模糊搜索）
- `status` (optional): 状态 (0/1)

### 2. 创建用户
```http
POST /api/user/create
Authorization: Bearer {token}
Content-Type: application/json

{
  "userName": "testuser",
  "nickName": "测试用户",
  "password": "123456",
  "email": "test@example.com",
  "userPhone": "13800138000",
  "userGender": "male",
  "status": "1",
  "roleIds": [2]
}
```

### 3. 更新用户
```http
PUT /api/user/update
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 2,
  "nickName": "更新昵称",
  "email": "update@example.com",
  "userPhone": "13800138001",
  "userGender": "female",
  "status": "1",
  "roleIds": [2, 3]
}
```

### 4. 删除用户
```http
DELETE /api/user/delete/{id}
Authorization: Bearer {token}
```

### 5. 重置密码
```http
POST /api/user/reset-password
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 2,
  "newPassword": "newpass123"
}
```

---

## 角色管理 API

### 1. 获取角色列表
```http
GET /api/role/list
Authorization: Bearer {token}
```

**Query 参数:**
- `current` (required): 当前页码
- `size` (required): 每页条数
- `roleName` (optional): 角色名（模糊搜索）
- `enabled` (optional): 是否启用

### 2. 创建角色
```http
POST /api/role/create
Authorization: Bearer {token}
Content-Type: application/json

{
  "roleName": "测试角色",
  "roleCode": "test",
  "description": "测试角色描述",
  "enabled": true,
  "menuIds": [1, 11, 12],
  "buttonIds": [1, 2]
}
```

### 3. 更新角色
```http
PUT /api/role/update
Authorization: Bearer {token}
Content-Type: application/json

{
  "roleId": 2,
  "roleName": "更新角色名",
  "roleCode": "test",
  "description": "更新描述",
  "enabled": true,
  "menuIds": [1, 11],
  "buttonIds": [1]
}
```

### 4. 删除角色
```http
DELETE /api/role/delete/{id}
Authorization: Bearer {token}
```

### 5. 获取角色权限
```http
GET /api/role/permissions/{id}
Authorization: Bearer {token}
```

---

## 菜单管理 API

### 1. 获取用户菜单树
```http
GET /api/system/menus
Authorization: Bearer {token}
```

### 2. 获取所有菜单
```http
GET /api/menu/all
Authorization: Bearer {token}
```

### 3. 创建菜单
```http
POST /api/menu/create
Authorization: Bearer {token}
Content-Type: application/json

{
  "parentId": 0,
  "name": "TestMenu",
  "path": "/test",
  "component": "/index/index",
  "title": "测试菜单",
  "icon": "test",
  "sort": 100,
  "enabled": true
}
```

### 4. 更新菜单
```http
PUT /api/menu/update
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 10,
  "parentId": 0,
  "name": "TestMenu",
  "path": "/test",
  "component": "/index/index",
  "title": "更新菜单",
  "icon": "test",
  "sort": 100,
  "enabled": true
}
```

### 5. 删除菜单
```http
DELETE /api/menu/delete/{id}
Authorization: Bearer {token}
```

---

## 测试示例

### 登录获取 Token
```bash
curl -X POST http://localhost:48080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"userName":"admin","password":"123456"}'
```

### 创建用户示例
```bash
TOKEN="your_token_here"

curl -X POST http://localhost:48080/api/user/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userName": "testuser",
    "nickName": "测试用户",
    "password": "123456",
    "email": "test@example.com",
    "userPhone": "13900139000",
    "userGender": "male",
    "status": "1",
    "roleIds": [2]
  }'
```

### 查询用户列表示例
```bash
curl -X GET "http://localhost:48080/api/user/list?current=1&size=10" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 统一响应格式

### 成功响应
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    // 响应数据
  }
}
```

### 分页响应
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [],
    "total": 100,
    "current": 1,
    "size": 10
  }
}
```

### 错误响应
```json
{
  "code": 400,
  "msg": "错误信息",
  "data": null
}
```

---

## 状态码说明
- `200`: 成功
- `400`: 参数错误
- `401`: 未授权
- `403`: 无权限
- `500`: 服务器错误

---

## Swagger 文档
访问: http://localhost:48080/swagger/index.html

