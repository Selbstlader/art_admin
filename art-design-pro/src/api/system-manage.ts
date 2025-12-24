import request from '@/utils/http'
import { AppRouteRecord } from '@/types/router'

// 获取用户列表
export function fetchGetUserList(params: Api.SystemManage.UserSearchParams) {
  return request.get<Api.SystemManage.UserList>({
    url: '/api/user/list',
    params
  })
}

// 创建用户
export function fetchCreateUser(data: Api.SystemManage.CreateUserRequest) {
  return request.post({
    url: '/api/user/create',
    data
  })
}

// 更新用户
export function fetchUpdateUser(data: Api.SystemManage.UpdateUserRequest) {
  return request.put({
    url: '/api/user/update',
    data
  })
}

// 删除用户
export function fetchDeleteUser(id: number) {
  return request.del({
    url: `/api/user/delete/${id}`
  })
}

// 重置用户密码
export function fetchResetPassword(data: Api.SystemManage.ResetPasswordRequest) {
  return request.post({
    url: '/api/user/reset-password',
    data
  })
}

// 获取角色列表
export function fetchGetRoleList(params: Api.SystemManage.RoleSearchParams) {
  return request.get<Api.SystemManage.RoleList>({
    url: '/api/role/list',
    params
  })
}

// 获取菜单列表（用户有权限的菜单，用于动态路由）
export function fetchGetMenuList() {
  return request.get<AppRouteRecord[]>({
    url: '/api/system/menus'
  })
}

// 获取所有菜单列表（管理页面用，不过滤权限）
export function fetchGetAllMenus() {
  return request.get<AppRouteRecord[]>({
    url: '/api/menu/all'
  })
}

// 创建菜单
export function fetchCreateMenu(data: any) {
  return request.post({
    url: '/api/menu/create',
    data
  })
}

// 更新菜单
export function fetchUpdateMenu(data: any) {
  return request.put({
    url: '/api/menu/update',
    data
  })
}

// 删除菜单
export function fetchDeleteMenu(id: number) {
  return request.del({
    url: `/api/menu/delete/${id}`
  })
}

// 创建角色
export function fetchCreateRole(data: Api.SystemManage.CreateRoleRequest) {
  return request.post({
    url: '/api/role/create',
    data
  })
}

// 更新角色
export function fetchUpdateRole(data: Api.SystemManage.UpdateRoleRequest) {
  return request.put({
    url: '/api/role/update',
    data
  })
}

// 删除角色
export function fetchDeleteRole(id: number) {
  return request.del({
    url: `/api/role/delete/${id}`
  })
}

// 获取角色权限
export function fetchGetRolePermissions(id: number) {
  return request.get<Api.SystemManage.RolePermissions>({
    url: `/api/role/permissions/${id}`
  })
}

// 更新角色权限
export function fetchUpdateRolePermissions(
  id: number,
  data: Api.SystemManage.UpdateRolePermissionsRequest
) {
  return request.put({
    url: `/api/role/permissions/${id}`,
    data
  })
}

// ========== 字典类型管理 ==========

// 获取字典类型列表
export function fetchGetDictionaryTypeList(params: Api.SystemManage.DictionaryTypeSearchParams) {
  return request.get<Api.SystemManage.DictionaryTypeList>({
    url: '/api/dictionary/type/list',
    params
  })
}

// 根据ID获取字典类型
export function fetchGetDictionaryTypeById(id: number) {
  return request.get<Api.SystemManage.DictionaryTypeItem>({
    url: `/api/dictionary/type/${id}`
  })
}

// 创建字典类型
export function fetchCreateDictionaryType(data: Api.SystemManage.CreateDictionaryTypeRequest) {
  return request.post({
    url: '/api/dictionary/type',
    data
  })
}

// 更新字典类型
export function fetchUpdateDictionaryType(data: Api.SystemManage.UpdateDictionaryTypeRequest) {
  return request.put({
    url: '/api/dictionary/type',
    data
  })
}

// 删除字典类型
export function fetchDeleteDictionaryType(id: number) {
  return request.del({
    url: '/api/dictionary/type',
    data: { id }
  })
}

// ========== 字典数据管理 ==========

// 获取字典数据列表
export function fetchGetDictionaryList(params: Api.SystemManage.DictionarySearchParams) {
  return request.get<Api.SystemManage.DictionaryList>({
    url: '/api/dictionary/list',
    params
  })
}

// 根据ID获取字典数据
export function fetchGetDictionaryById(id: number) {
  return request.get<Api.SystemManage.DictionaryItem>({
    url: `/api/dictionary/${id}`
  })
}

// 根据类型编码获取字典数据
export function fetchGetDictionaryByTypeCode(typeCode: string) {
  return request.get<Api.SystemManage.DictionaryItem[]>({
    url: '/api/dictionary/by-type',
    params: { typeCode }
  })
}

// 创建字典数据
export function fetchCreateDictionary(data: Api.SystemManage.CreateDictionaryRequest) {
  return request.post({
    url: '/api/dictionary',
    data
  })
}

// 更新字典数据
export function fetchUpdateDictionary(data: Api.SystemManage.UpdateDictionaryRequest) {
  return request.put({
    url: '/api/dictionary',
    data
  })
}

// 删除字典数据
export function fetchDeleteDictionary(id: number) {
  return request.del({
    url: '/api/dictionary',
    data: { id }
  })
}

// ========== 部门管理 ==========

// 获取部门列表（树形结构）
export function fetchGetDepartmentList(params?: Api.SystemManage.DepartmentSearchParams) {
  return request.get<Api.SystemManage.DepartmentList>({
    url: '/api/department/list',
    params
  })
}

// 根据ID获取部门
export function fetchGetDepartmentById(id: number) {
  return request.get<Api.SystemManage.DepartmentItem>({
    url: `/api/department/${id}`
  })
}

// 创建部门
export function fetchCreateDepartment(data: Api.SystemManage.CreateDepartmentRequest) {
  return request.post({
    url: '/api/department',
    data
  })
}

// 更新部门
export function fetchUpdateDepartment(data: Api.SystemManage.UpdateDepartmentRequest) {
  return request.put({
    url: '/api/department',
    data
  })
}

// 删除部门
export function fetchDeleteDepartment(id: number) {
  return request.del({
    url: `/api/department/${id}`
  })
}

// ========== 操作日志管理 ==========

// 获取操作日志列表
export function fetchGetOperationLogList(params: Api.SystemManage.OperationLogSearchParams) {
  return request.get<Api.SystemManage.OperationLogList>({
    url: '/api/operation-log/list',
    params
  })
}

// 根据ID获取操作日志详情
export function fetchGetOperationLogDetail(id: number) {
  return request.get<Api.SystemManage.OperationLogItem>({
    url: `/api/operation-log/${id}`
  })
}

// 删除操作日志
export function fetchDeleteOperationLog(id: number) {
  return request.del({
    url: `/api/operation-log/${id}`
  })
}

// 批量删除操作日志
export function fetchBatchDeleteOperationLog(ids: number[]) {
  return request.del({
    url: '/api/operation-log/batch-delete',
    data: { ids }
  })
}

// 清理操作日志
export function fetchCleanOperationLog(days: number) {
  return request.post({
    url: '/api/operation-log/clean',
    data: { days }
  })
}

// ========== APP用户管理 ==========

// 获取APP用户列表
export function fetchGetAppUserList(params: Api.SystemManage.AppUserSearchParams) {
  return request.get<Api.SystemManage.AppUserList>({
    url: '/api/app-user/list',
    params
  })
}

// 创建APP用户
export function fetchCreateAppUser(data: Api.SystemManage.CreateAppUserRequest) {
  return request.post({
    url: '/api/app-user/create',
    data
  })
}

// 更新APP用户
export function fetchUpdateAppUser(data: Api.SystemManage.UpdateAppUserRequest) {
  return request.put({
    url: '/api/app-user/update',
    data
  })
}

// 删除APP用户
export function fetchDeleteAppUser(id: number) {
  return request.del({
    url: `/api/app-user/delete/${id}`
  })
}

// 重置APP用户密码
export function fetchResetAppUserPassword(data: Api.SystemManage.ResetAppUserPasswordRequest) {
  return request.post({
    url: '/api/app-user/reset-password',
    data
  })
}

// ========== 用户效果图配额管理 ==========

// 用户配额信息类型
export interface UserRenderQuotaInfo {
  userId: number
  dailyLimit: number
  usedToday: number
  remainingUse: number
}

// 获取用户效果图配额
export function fetchGetUserRenderQuota(userId: number) {
  return request.get<Http.BaseResponse<UserRenderQuotaInfo>>({
    url: '/api/system/user/render-quota',
    params: { userId },
    _fullResponse: true
  })
}

// 更新用户效果图配额
export function fetchUpdateUserRenderQuota(data: { userId: number; dailyLimit: number }) {
  return request.put({
    url: '/api/system/user/render-quota',
    data
  })
}

// 重置用户今日使用次数
export function fetchResetUserRenderQuota(userId: number) {
  return request.post({
    url: '/api/system/user/render-quota/reset',
    params: { userId }
  })
}
