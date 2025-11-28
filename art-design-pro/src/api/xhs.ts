import request from '@/utils/http'

// ==================== 类型定义 ====================

// 总结风格
export type SummaryStyle = 'concise' | 'detailed' | 'casual'

// 情感倾向
export type Sentiment = 'positive' | 'neutral' | 'negative'

// 专业分析结果
export interface ProfessionalAnalysis {
  category: string
  main_conclusion: string
  detailed_data: Record<string, string>
  risk_warnings: string[]
  recommendations: string[]
}

// 总结请求
export interface SummarizeRequest {
  url?: string
  original_url?: string
  note_title?: string
  content?: string
  image_urls?: string[]
  style?: SummaryStyle
  max_length?: number
  enable_ocr?: boolean
  enable_analysis?: boolean
  save_history?: boolean
}

// 总结响应（不保存历史时）
export interface SummaryResult {
  title: string
  summary: string
  key_points: string[]
  highlights?: string[]
  concerns?: string[]
  tags: string[]
  sentiment: Sentiment
  reading_time?: string
  target_audience?: string
  content_type?: string
  action_items?: string[]
  related_topics?: string[]
  credibility_note?: string
  quick_facts?: Record<string, string>
  ocr_texts?: string[]
  analysis?: ProfessionalAnalysis
  original_url?: string
}

// 总结记录（保存历史时）
export interface XHSSummary {
  id: number
  user_id: number
  original_url: string
  note_title: string
  note_content: string
  image_urls: string[]
  summary_title: string
  summary_content: string
  key_points: string[]
  highlights?: string[]
  concerns?: string[]
  tags: string[]
  sentiment: Sentiment
  reading_time?: string
  target_audience?: string
  content_type?: string
  action_items?: string[]
  related_topics?: string[]
  credibility_note?: string
  quick_facts?: Record<string, string>
  ocr_texts: string[]
  analysis?: ProfessionalAnalysis
  style: SummaryStyle
  is_favorite: boolean
  folder_id: number
  created_at: string
  updated_at: string
}

// OCR请求
export interface OCRRequest {
  image_url?: string
  image_data?: string
}

// 收藏请求
export interface FavoriteRequest {
  is_favorite: boolean
  folder_id?: number
}

// 收藏夹
export interface XHSFolder {
  id: number
  user_id: number
  name: string
  sort_order: number
  created_at: string
  updated_at: string
}

// 分页响应
export interface PageResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 风格配置
export const SummaryStyleConfig: Record<SummaryStyle, { label: string; desc: string }> = {
  concise: { label: '简洁凝练', desc: '仅保留核心信息，适合快速浏览' },
  detailed: { label: '详细全面', desc: '包含完整逻辑链及细节信息' },
  casual: { label: '口语化解读', desc: '生活化语言，易于理解' }
}

// 情感配置
export const SentimentConfig: Record<Sentiment, { label: string; color: string; icon: string }> = {
  positive: { label: '积极', color: '#10B981', icon: '😊' },
  neutral: { label: '中性', color: '#6B7280', icon: '😐' },
  negative: { label: '消极', color: '#EF4444', icon: '😔' }
}

// ==================== API接口 ====================
export const xhsApi = {
  // 生成总结
  summarize: (params: SummarizeRequest) => {
    return request.post<SummaryResult | XHSSummary>({ 
      url: '/api/xhs/summarize', 
      params,
      timeout: 180000 // 3分钟超时，匹配后端配置
    })
  },

  // OCR识别
  ocr: (params: OCRRequest) => {
    return request.post<{ text: string }>({ url: '/api/xhs/ocr', params })
  },

  // 获取历史记录
  getHistory: (params: { page?: number; page_size?: number }) => {
    return request.get<PageResponse<XHSSummary>>({ url: '/api/xhs/history', params })
  },

  // 获取总结详情
  getSummaryDetail: (id: number) => {
    return request.get<XHSSummary>({ url: `/api/xhs/summary/${id}` })
  },

  // 删除总结
  deleteSummary: (id: number) => {
    return request.del({ url: `/api/xhs/summary/${id}` })
  },

  // 更新收藏状态
  updateFavorite: (id: number, params: FavoriteRequest) => {
    return request.put({ url: `/api/xhs/summary/${id}/favorite`, params })
  },

  // 获取收藏列表
  getFavorites: (params: { folder_id?: number; page?: number; page_size?: number }) => {
    return request.get<PageResponse<XHSSummary>>({ url: '/api/xhs/favorites', params })
  },

  // 搜索总结
  search: (params: { keyword: string; page?: number; page_size?: number }) => {
    return request.get<PageResponse<XHSSummary>>({ url: '/api/xhs/search', params })
  },

  // 创建收藏夹
  createFolder: (name: string) => {
    return request.post<XHSFolder>({ url: '/api/xhs/folder', params: { name } })
  },

  // 获取收藏夹列表
  getFolders: () => {
    return request.get<XHSFolder[]>({ url: '/api/xhs/folders' })
  },

  // 更新收藏夹
  updateFolder: (id: number, params: { name: string; sort_order?: number }) => {
    return request.put({ url: `/api/xhs/folder/${id}`, params })
  },

  // 删除收藏夹
  deleteFolder: (id: number) => {
    return request.del({ url: `/api/xhs/folder/${id}` })
  }
}
