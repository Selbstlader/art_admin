import request from '@/utils/http'

// 项目模板相关接口
export const projectTemplateApi = {
  // 获取模板列表
  getTemplateList: (params: any) => {
    return request.get({ url: '/api/project/template/list', params })
  },
  // 获取模板详情
  getTemplateDetail: (id: number) => {
    return request.get({ url: `/api/project/template/${id}` })
  },
  // 创建模板
  createTemplate: (params: any) => {
    return request.post({ url: '/api/project/template', params })
  },
  // 更新模板
  updateTemplate: (params: any) => {
    return request.put({ url: '/api/project/template', params })
  },
  // 删除模板
  deleteTemplate: (id: number) => {
    return request.del({ url: `/api/project/template/${id}` })
  }
}

// 项目相关接口
export const projectApi = {
  // 获取项目列表
  getProjectList: (params: any) => {
    return request.get({ url: '/api/project/list', params })
  },
  // 获取项目详情
  getProjectDetail: (id: number) => {
    return request.get({ url: `/api/project/${id}` })
  },
  // 创建项目
  createProject: (params: any) => {
    return request.post({ url: '/api/project', params })
  },
  // 从模板创建项目
  createProjectFromTemplate: (params: any) => {
    return request.post({ url: '/api/project/from-template', params })
  },
  // 更新项目
  updateProject: (params: any) => {
    return request.put({ url: '/api/project', params })
  },
  // 删除项目
  deleteProject: (id: number) => {
    return request.del({ url: `/api/project/${id}` })
  },
  // 获取项目统计
  getProjectStatistics: () => {
    return request.get({ url: '/api/project/statistics' })
  }
}

// 任务相关接口
export const taskApi = {
  // 获取任务列表
  getTaskList: (params: any) => {
    return request.get({ url: '/api/project/task/list', params })
  },
  // 获取任务详情
  getTaskDetail: (id: number) => {
    return request.get({ url: `/api/project/task/${id}` })
  },
  // 创建任务
  createTask: (params: any) => {
    return request.post({ url: '/api/project/task', params })
  },
  // 更新任务
  updateTask: (params: any) => {
    return request.put({ url: '/api/project/task', params })
  },
  // 批量更新任务
  batchUpdateTask: (params: any) => {
    return request.put({ url: '/api/project/task/batch', params })
  },
  // 快速更新任务(甘特图实时更新)
  quickUpdateTask: (params: any) => {
    return request.put({ url: '/api/project/task/quick', params })
  },
  // 删除任务
  deleteTask: (id: number) => {
    return request.del({ url: `/api/project/task/${id}` })
  },
  // 获取甘特图数据
  getGanttData: (projectId: number) => {
    return request.get({ url: '/api/project/task/gantt', params: { projectId } })
  },
  // 创建任务依赖
  createDependency: (params: any) => {
    return request.post({ url: '/api/project/task/dependency', params })
  },
  // 删除任务依赖
  deleteDependency: (id: number) => {
    return request.del({ url: `/api/project/task/dependency/${id}` })
  },
  // 创建任务评论
  createComment: (params: any) => {
    return request.post({ url: '/api/project/task/comment', params })
  }
}
