import request from '@/utils/http'

// ==================== 类型定义 ====================

// 情绪类型
export type MoodType =
  | 'happy'
  | 'sad'
  | 'anxious'
  | 'calm'
  | 'angry'
  | 'excited'
  | 'tired'
  | 'stressed'

// 情绪类型配置
export const MoodTypeConfig: Record<
  MoodType,
  { label: string; icon: string; color: string; tendency: string }
> = {
  happy: { label: '开心', icon: '😊', color: '#FFD700', tendency: '积极' },
  sad: { label: '难过', icon: '😢', color: '#4169E1', tendency: '消极' },
  anxious: { label: '焦虑', icon: '😰', color: '#FF6347', tendency: '消极' },
  calm: { label: '平静', icon: '😌', color: '#90EE90', tendency: '积极' },
  angry: { label: '愤怒', icon: '😠', color: '#DC143C', tendency: '消极' },
  excited: { label: '兴奋', icon: '🤩', color: '#FF69B4', tendency: '积极' },
  tired: { label: '疲惫', icon: '😴', color: '#808080', tendency: '中性' },
  stressed: { label: '压力', icon: '😫', color: '#8B4513', tendency: '消极' }
}

// 情绪记录
export interface MoodRecord {
  id: number
  user_id: number
  mood_type: MoodType
  intensity: number
  note: string
  triggers: string[]
  activities: string[]
  location: string
  created_at: string
  updated_at: string
}

// 创建情绪记录请求
export interface CreateMoodRecordRequest {
  mood_type: MoodType
  intensity: number
  note?: string
  triggers?: string[]
  activities?: string[]
  location?: string
}

// 日历数据
export interface MoodCalendarData {
  date: string
  mood_type: MoodType
  avg_intensity: number
  count: number
}

// 统计数据
export interface MoodStatistics {
  total_records: number
  avg_intensity: number
  mood_distribution: Record<MoodType, number>
  daily_avg: number
}

// 趋势数据
export interface MoodTrend {
  date: string
  avg_intensity: number
  mood_type: MoodType
  count: number
}

// ==================== 情绪管理接口 ====================
export const moodRecordApi = {
  // 创建情绪记录
  create: (params: CreateMoodRecordRequest) => {
    return request.post({ url: '/api/mood-record/create', params })
  },
  // 获取情绪记录列表
  getList: (params: { page?: number; pageSize?: number }) => {
    return request.get({ url: '/api/mood-record/list', params })
  },
  // 获取情绪记录详情
  getDetail: (id: number) => {
    return request.get({ url: `/api/mood-record/${id}` })
  },
  // 按日期查询
  getByDate: (date: string) => {
    return request.get({ url: `/api/mood-record/date/${date}` })
  },
  // 获取日历数据
  getCalendar: (params: { year: number; month: number }) => {
    return request.get({ url: '/api/mood-record/calendar', params })
  },
  // 获取统计数据
  getStatistics: (params?: { days?: number }) => {
    return request.get({ url: '/api/mood-record/statistics', params })
  },
  // 获取趋势数据
  getTrends: (params?: { days?: number }) => {
    return request.get({ url: '/api/mood-record/trends', params })
  },
  // 更新情绪记录
  update: (params: Partial<MoodRecord> & { id: number }) => {
    return request.put({ url: '/api/mood-record/update', params })
  },
  // 删除情绪记录
  delete: (id: number) => {
    return request.del({ url: `/api/mood-record/delete/${id}` })
  }
}

// ==================== 日记管理接口 ====================
export interface JournalEntry {
  id: number
  user_id: number
  title: string
  content: string
  mood: string
  tags: string[]
  images: string[]
  is_private: boolean
  created_at: string
  updated_at: string
}

export interface CreateJournalRequest {
  title: string
  content: string
  mood?: string
  tags?: string[]
  is_private?: boolean
}

export const journalApi = {
  // 创建日记
  create: (params: CreateJournalRequest) => {
    return request.post({ url: '/api/journal-entry/create', params })
  },
  // 获取日记列表
  getList: (params?: { page?: number; pageSize?: number }) => {
    return request.get({ url: '/api/journal-entry/list', params })
  },
  // 获取日记详情
  getDetail: (id: number) => {
    return request.get({ url: `/api/journal-entry/${id}` })
  },
  // 获取所有标签
  getTags: () => {
    return request.get({ url: '/api/journal-entry/tags' })
  },
  // 按标签查询
  getByTag: (tag: string) => {
    return request.get({ url: '/api/journal-entry/by-tag', params: { tag } })
  },
  // 更新标签
  updateTags: (id: number, tags: string[]) => {
    return request.put({ url: `/api/journal-entry/${id}/tags`, params: { tags } })
  },
  // 添加图片
  addImages: (id: number, images: string[]) => {
    return request.post({ url: `/api/journal-entry/${id}/images`, params: { images } })
  },
  // 删除图片
  deleteImages: (id: number, images: string[]) => {
    return request.del({ url: `/api/journal-entry/${id}/images`, params: { images } })
  },
  // 搜索日记
  search: (keyword: string) => {
    return request.get({ url: '/api/journal-entry/search', params: { keyword } })
  },
  // 更新日记
  update: (id: number, params: Partial<CreateJournalRequest>) => {
    return request.put({ url: `/api/journal-entry/${id}`, params })
  },
  // 删除日记
  delete: (id: number) => {
    return request.del({ url: `/api/journal-entry/${id}` })
  }
}

// ==================== 冥想内容接口 ====================
export type MeditationType = 'breathing' | 'sleep' | 'stress' | 'focus' | 'anxiety' | 'body_scan'
export type DifficultyLevel = 'beginner' | 'intermediate' | 'advanced'

export const MeditationTypeConfig: Record<
  MeditationType,
  { label: string; description: string; duration: string }
> = {
  breathing: { label: '呼吸练习', description: '快速放松、专注力训练', duration: '5-15分钟' },
  sleep: { label: '睡眠冥想', description: '入睡困难、改善睡眠', duration: '15-30分钟' },
  stress: { label: '减压放松', description: '工作压力大、情绪紧张', duration: '10-20分钟' },
  focus: { label: '专注训练', description: '提高注意力', duration: '10-15分钟' },
  anxiety: { label: '焦虑缓解', description: '焦虑情绪、紧张不安', duration: '10-20分钟' },
  body_scan: { label: '身体扫描', description: '身体觉察、深度放松', duration: '15-25分钟' }
}

export const DifficultyConfig: Record<DifficultyLevel, { label: string; color: string }> = {
  beginner: { label: '初级', color: '#67C23A' },
  intermediate: { label: '中级', color: '#E6A23C' },
  advanced: { label: '高级', color: '#F56C6C' }
}

export interface MeditationContent {
  id: number
  title: string
  description: string
  category: MeditationType
  difficulty: DifficultyLevel
  duration: number // 秒
  audio_url: string
  cover_image: string
  play_count: number
  favorite_count: number
  created_at: string
}

export interface PlayRecord {
  id: number
  user_id: number
  content_id: number
  duration: number
  progress: number
  completed: boolean
  last_played_at: string
  play_count: number
}

export interface UpdatePlayRecordRequest {
  content_id: number
  duration: number
  progress: number
  completed: boolean
}

export const meditationApi = {
  // 获取内容列表
  getList: (params?: {
    page?: number
    pageSize?: number
    category?: string
    difficulty?: string
  }) => {
    return request.get({ url: '/api/meditation-content/list', params })
  },
  // 获取内容详情
  getDetail: (id: number) => {
    return request.get({ url: `/api/meditation-content/${id}` })
  },
  // 按分类查询
  getByCategory: (category: MeditationType) => {
    return request.get({ url: `/api/meditation-content/category/${category}` })
  },
  // 按难度查询
  getByDifficulty: (level: DifficultyLevel) => {
    return request.get({ url: `/api/meditation-content/difficulty/${level}` })
  },
  // 获取热门内容
  getPopular: () => {
    return request.get({ url: '/api/meditation-content/popular' })
  },
  // 添加收藏
  addFavorite: (id: number) => {
    return request.post({ url: `/api/meditation-content/${id}/favorite` })
  },
  // 取消收藏
  removeFavorite: (id: number) => {
    return request.del({ url: `/api/meditation-content/${id}/favorite` })
  },
  // 检查是否收藏
  checkFavorite: (id: number) => {
    return request.get({ url: `/api/meditation-content/${id}/favorite/check` })
  },
  // 获取收藏列表
  getFavorites: () => {
    return request.get({ url: '/api/meditation-content/favorites' })
  },
  // 更新播放记录
  updatePlayRecord: (params: UpdatePlayRecordRequest) => {
    return request.post({ url: '/api/meditation-content/play-record', params })
  },
  // 获取播放记录
  getPlayRecord: (contentId: number) => {
    return request.get({ url: `/api/meditation-content/${contentId}/play-record` })
  },
  // 获取播放历史
  getPlayHistory: () => {
    return request.get({ url: '/api/meditation-content/play-history' })
  },
  // 获取最近播放
  getRecentlyPlayed: () => {
    return request.get({ url: '/api/meditation-content/recently-played' })
  },
  // 获取继续播放列表
  getContinuePlaying: () => {
    return request.get({ url: '/api/meditation-content/continue-playing' })
  },
  // 获取播放统计
  getPlayStats: () => {
    return request.get({ url: '/api/meditation-content/play-stats' })
  }
}

// ==================== 目标管理接口 ====================
export type GoalType = 'mood_record' | 'meditation' | 'journal'
export type GoalPeriod = 'daily' | 'weekly' | 'monthly'

export const GoalTypeConfig: Record<GoalType, { label: string; icon: string; unit: string }> = {
  mood_record: { label: '情绪记录', icon: '😊', unit: '次' },
  meditation: { label: '冥想练习', icon: '🧘', unit: '次' },
  journal: { label: '写日记', icon: '📝', unit: '篇' }
}

export const GoalPeriodConfig: Record<GoalPeriod, { label: string }> = {
  daily: { label: '每日' },
  weekly: { label: '每周' },
  monthly: { label: '每月' }
}

export interface UserGoal {
  id: number
  user_id: number
  goal_type: GoalType
  target_value: number
  current_value: number
  period: GoalPeriod
  is_active: boolean
  start_date: string
  end_date: string
  created_at: string
  updated_at: string
}

// 目标类型（用于goals页面）
export interface Goal {
  id: number
  user_id: number
  title: string
  category: string
  description?: string
  target_value: number
  current_value: number
  unit: string
  status: 'active' | 'completed' | 'paused' | 'abandoned'
  start_date: string
  end_date: string
  reminder_time?: string
  streak_days?: number
  total_checkins?: number
  created_at: string
  updated_at: string
}

export interface CreateGoalRequest {
  goal_type: GoalType
  target_value: number
  period: GoalPeriod
  title?: string
  category?: string
  description?: string
  unit?: string
  start_date?: string
  end_date?: string
  reminder_time?: string
}

export interface GoalProgress {
  id: number
  goal_id: number
  progress_date: string
  current_value: number
  is_completed: boolean
}

export const goalApi = {
  // 创建目标
  create: (params: Partial<Goal> & { goal_type?: GoalType; period?: GoalPeriod }) => {
    return request.post({ url: '/api/goals', params })
  },
  // 获取目标列表
  getList: (params?: { status?: string }) => {
    return request.get({ url: '/api/goals', params })
  },
  // 更新目标
  update: (params: { id: number } & Partial<Goal>) => {
    return request.put({ url: `/api/goals/${params.id}`, params })
  },
  // 删除目标
  delete: (id: number) => {
    return request.del({ url: `/api/goals/${id}` })
  },
  // 更新进度
  updateProgress: (id: number, params: { current_value: number }) => {
    return request.put({ url: `/api/goals/${id}/progress`, params })
  },
  // 获取进度历史
  getProgressHistory: (id: number) => {
    return request.get({ url: `/api/goals/${id}/progress` })
  }
}

// ==================== 打卡管理接口 ====================
export type CheckinType = 'mood' | 'meditation' | 'journal' | 'general'

export const CheckinTypeConfig: Record<CheckinType, { label: string; icon: string }> = {
  mood: { label: '情绪', icon: '😊' },
  meditation: { label: '冥想', icon: '🧘' },
  journal: { label: '日记', icon: '📝' },
  general: { label: '通用', icon: '✓' }
}

export interface UserCheckin {
  id: number
  user_id: number
  checkin_date: string
  checkin_type: CheckinType
  note: string
  created_at: string
}

export interface CheckinRequest {
  checkin_type: CheckinType
  note?: string
}

export interface TodayCheckinStatus {
  mood: boolean
  meditation: boolean
  journal: boolean
  general: boolean
  total: number
}

export interface CheckinStats {
  mood_streak: number
  meditation_streak: number
  journal_streak: number
  total_streak: number
  last_30_days: number
  type_count: Record<CheckinType, number>
}

export const checkinApi = {
  // 打卡
  checkin: (params: CheckinRequest) => {
    return request.post({ url: '/api/checkin', params })
  },
  // 创建打卡（目标打卡）
  create: (params: { goal_id: number; value: number; note?: string }) => {
    return request.post({ url: '/api/checkin/goal', params })
  },
  // 获取今日打卡状态
  getTodayStatus: () => {
    return request.get({ url: '/api/checkin/today' })
  },
  // 获取打卡历史
  getHistory: (params?: { days?: number }) => {
    return request.get({ url: '/api/checkin/history', params })
  },
  // 获取连续天数
  getStreak: (params?: { type?: CheckinType }) => {
    return request.get({ url: '/api/checkin/streak', params })
  },
  // 获取打卡日历
  getCalendar: (params: { year: number; month: number }) => {
    return request.get({ url: '/api/checkin/calendar', params })
  },
  // 获取打卡统计
  getStats: () => {
    return request.get({ url: '/api/checkin/stats' })
  }
}
