/**
 * namespace: Api
 *
 * 所有接口相关类型定义
 * 在.vue文件使用会报错，需要在 eslint.config.mjs 中配置 globals: { Api: 'readonly' }
 */

declare namespace Api {
  /** 通用类型 */
  namespace Common {
    /** 分页参数 */
    interface PaginationParams {
      /** 当前页码 */
      current: number
      /** 每页条数 */
      size: number
      /** 总条数 */
      total: number
    }

    /** 通用搜索参数 */
    type CommonSearchParams = Pick<PaginationParams, 'current' | 'size'>

    /** 分页响应基础结构 */
    interface PaginatedResponse<T = any> {
      records: T[]
      current: number
      size: number
      total: number
    }

    /** 启用状态 */
    type EnableStatus = '1' | '2'
  }

  /** 认证类型 */
  namespace Auth {
    /** 登录参数 */
    interface LoginParams {
      userName: string
      password: string
    }

    /** 登录响应 */
    interface LoginResponse {
      token: string
      refreshToken: string
    }

    /** 用户信息 */
    interface UserInfo {
      buttons: string[]
      roles: string[]
      userId: number
      userName: string
      email: string
      avatar?: string
    }
  }

  /** 系统管理类型 */
  namespace SystemManage {
    /** 用户列表 */
    type UserList = Api.Common.PaginatedResponse<UserListItem>

    /** 用户列表项 */
    interface UserListItem {
      id: number
      avatar: string
      status: string
      userName: string
      userGender: string
      nickName: string
      userPhone: string
      userEmail: string
      userRoles: string[]
      createBy: string
      createTime: string
      updateBy: string
      updateTime: string
    }

    /** 用户搜索参数 */
    type UserSearchParams = Partial<
      Pick<UserListItem, 'id' | 'userName' | 'userGender' | 'userPhone' | 'userEmail' | 'status'> &
        Api.Common.CommonSearchParams
    >

    /** 创建用户请求 */
    interface CreateUserRequest {
      userName: string
      nickName: string
      password: string
      email: string
      userPhone: string
      userGender: string
      avatar?: string
      status: string
      roleIds: number[]
    }

    /** 更新用户请求 */
    interface UpdateUserRequest {
      id: number
      nickName: string
      email: string
      userPhone: string
      userGender: string
      avatar?: string
      status: string
      roleIds: number[]
    }

    /** 重置密码请求 */
    interface ResetPasswordRequest {
      id: number
      newPassword: string
    }

    /** 角色列表 */
    type RoleList = Api.Common.PaginatedResponse<RoleListItem>

    /** 角色列表项 */
    interface RoleListItem {
      roleId: number
      roleName: string
      roleCode: string
      description: string
      enabled: boolean
      createTime: string
    }

    /** 角色搜索参数 */
    type RoleSearchParams = Partial<
      Pick<RoleListItem, 'roleId' | 'roleName' | 'roleCode' | 'description' | 'enabled'> &
        Api.Common.CommonSearchParams
    >

    /** 创建角色请求 */
    interface CreateRoleRequest {
      roleName: string
      roleCode: string
      description: string
      enabled: boolean
      menuIds?: number[]
      buttonIds?: number[]
    }

    /** 更新角色请求 */
    interface UpdateRoleRequest {
      roleId: number
      roleName: string
      roleCode: string
      description: string
      enabled: boolean
      menuIds?: number[]
      buttonIds?: number[]
    }

    /** 角色权限响应 */
    interface RolePermissions {
      menuIds: number[]
      buttonIds: number[]
    }

    /** 更新角色权限请求 */
    interface UpdateRolePermissionsRequest {
      menuIds: number[]
      buttonIds: number[]
    }

    // ========== 字典类型管理 ==========

    /** 字典类型列表 */
    type DictionaryTypeList = Api.Common.PaginatedResponse<DictionaryTypeItem>

    /** 字典类型列表项 */
    interface DictionaryTypeItem {
      id: number
      typeName: string
      typeCode: string
      description: string
      enabled: boolean
      remark: string
      createBy: string
      createTime: string
      updateBy: string
      updateTime: string
    }

    /** 字典类型搜索参数 */
    type DictionaryTypeSearchParams = Partial<
      Pick<DictionaryTypeItem, 'id' | 'typeName' | 'typeCode' | 'description' | 'enabled'> &
        Api.Common.CommonSearchParams
    >

    /** 创建字典类型请求 */
    interface CreateDictionaryTypeRequest {
      typeName: string
      typeCode: string
      description: string
      enabled: boolean
      remark: string
    }

    /** 更新字典类型请求 */
    interface UpdateDictionaryTypeRequest {
      id: number
      typeName: string
      typeCode: string
      description: string
      enabled: boolean
      remark: string
    }

    // ========== 字典数据管理 ==========

    /** 字典数据列表 */
    type DictionaryList = Api.Common.PaginatedResponse<DictionaryItem>

    /** 字典数据列表项 */
    interface DictionaryItem {
      id: number
      typeCode: string
      label: string
      value: string
      orderNum: number
      enabled: boolean
      remark: string
      createBy: string
      createTime: string
      updateBy: string
      updateTime: string
      dictionaryType?: DictionaryTypeItem
    }

    /** 字典数据搜索参数 */
    type DictionarySearchParams = Partial<
      Pick<DictionaryItem, 'id' | 'typeCode' | 'label' | 'value' | 'enabled'> &
        Api.Common.CommonSearchParams
    >

    /** 创建字典数据请求 */
    interface CreateDictionaryRequest {
      typeCode: string
      label: string
      value: string
      orderNum: number
      enabled: boolean
      remark: string
    }

    /** 更新字典数据请求 */
    interface UpdateDictionaryRequest {
      id: number
      typeCode: string
      label: string
      value: string
      orderNum: number
      enabled: boolean
      remark: string
    }

    // ========== 部门管理 ==========

    /** 部门列表（树形结构） */
    type DepartmentList = DepartmentListItem[]

    /** 部门列表项 */
    interface DepartmentListItem {
      deptId: number
      parentId?: number
      parentName?: string
      deptName: string
      deptCode: string
      orderNum: number
      leader?: string
      phone?: string
      email?: string
      status: number // 1-正常 0-停用
      createBy?: string
      createTime: string
      updateBy?: string
      updateTime?: string
      children?: DepartmentListItem[]
      hasChildren?: boolean
    }

    /** 部门项（用于详情） */
    type DepartmentItem = DepartmentListItem

    /** 部门搜索参数 */
    interface DepartmentSearchParams {
      deptName?: string
      deptCode?: string
      status?: number
    }

    /** 创建部门请求 */
    interface CreateDepartmentRequest {
      parentId?: number | null
      deptName: string
      deptCode: string
      orderNum: number
      leader?: string
      phone?: string
      email?: string
      status: number
    }

    /** 更新部门请求 */
    interface UpdateDepartmentRequest {
      deptId: number
      parentId?: number | null
      deptName: string
      deptCode: string
      orderNum: number
      leader?: string
      phone?: string
      email?: string
      status: number
    }

    // ========== 操作日志管理 ==========

    /** 操作日志列表 */
    type OperationLogList = Api.Common.PaginatedResponse<OperationLogItem>

    /** 操作日志列表项 */
    interface OperationLogItem {
      id: number
      module: string
      businessType: string
      requestMethod: string
      requestUrl: string
      operatorName: string
      operatorIp: string
      operatorAddr: string
      requestParam: string
      responseData: string
      status: number // 1-成功 0-失败
      errorMsg: string
      costTime: number
      userAgent: string
      operationTime: string
    }

    /** 操作日志搜索参数 */
    type OperationLogSearchParams = Partial<
      Pick<OperationLogItem, 'module' | 'businessType' | 'operatorName' | 'status'> &
        Api.Common.CommonSearchParams & {
          startTime?: string
          endTime?: string
        }
    >
  }
}
