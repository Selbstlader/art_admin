import request from '@/utils/http'

// 学科相关接口
export const subjectApi = {
  getSubjectList: () => request.get({ url: '/api/learning/subject/list' })
}

// 教材相关接口
export const learningMaterialApi = {
  generateMaterial: (params: any) =>
    request.post({ url: '/api/learning/material/generate', params }), // 异步生成,立即返回
  getMaterialList: (params: any) => request.get({ url: '/api/learning/material/list', params }),
  getMaterialDetail: (id: number) => request.get({ url: `/api/learning/material/${id}` }),
  deleteMaterial: (id: number) => request.del({ url: `/api/learning/material/${id}` }),
  // 任务相关
  getTaskList: (params?: any) => request.get({ url: '/api/learning/material/tasks', params }),
  getTaskDetail: (id: number) => request.get({ url: `/api/learning/material/task/${id}` })
}

// 年级相关接口
export const gradeApi = {
  getGradesBySubject: (subjectId?: number) =>
    request.get({ url: '/api/learning/grades', params: { subject_id: subjectId } })
}

// 类型定义
export interface Subject {
  id: number
  name: string
  icon: string
  description: string
  grade_levels: string[]
  sort: number
}

export interface LearningMaterial {
  id: number
  user_id: number
  subject_id: number
  title: string
  topic: string
  grade: string
  difficulty: number
  summary: string
  content: any
  total_time: number
  view_count: number
  favorite_count: number
  created_at: string
}

export interface GenerateTask {
  id: number
  user_id: number
  subject_id: number
  grade: string
  topic: string
  difficulty: number
  status: number // 0-待处理 1-生成中 2-已完成 3-失败
  material_id: number
  error_msg: string
  created_at: string
  updated_at: string
}
