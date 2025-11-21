// Dify 知识库相关类型
export interface DifyDataset {
  id: string
  name: string
  description: string
  provider: string
  permission: string
  data_source_type: string
  indexing_technique: string
  app_count: number
  document_count: number
  word_count: number
  created_by: string
  created_at: number
  updated_by: string
  updated_at: number
  embedding_model: string
  embedding_model_provider: string
  embedding_available: boolean
}

export interface DifyDatasetListResponse {
  data: DifyDataset[]
  has_more: boolean
  limit: number
  total: number
  page: number
}

// Dify 会话相关类型
export interface DifyConversation {
  id: string
  name: string
  inputs: Record<string, any>
  status: string
  introduction: string
  created_at: number
  updated_at: number
}

export interface DifyConversationListResponse {
  limit: number
  has_more: boolean
  data: DifyConversation[]
}

// Dify 消息相关类型
export interface DifyMessage {
  id: string
  conversation_id: string
  inputs: Record<string, any>
  query: string
  answer: string
  message_files: DifyMessageFile[]
  feedback?: DifyFeedback
  retriever_resources: DifyRetrieverResource[]
  created_at: number
}

export interface DifyMessageFile {
  id: string
  type: string
  url: string
  belongs_to: string
}

export interface DifyFeedback {
  rating: string
}

export interface DifyRetrieverResource {
  position: number
  dataset_id: string
  dataset_name: string
  document_id: string
  document_name: string
  segment_id: string
  score: number
  content: string
}

export interface DifyMessageListResponse {
  limit: number
  has_more: boolean
  data: DifyMessage[]
}

// Dify 对话响应类型
export interface DifyChatResponse {
  event: string
  task_id?: string
  message_id: string
  conversation_id: string
  mode: string
  answer: string
  metadata: {
    annotation_reply: any
    retriever_resources: any[]
    usage: {
      prompt_tokens: number
      prompt_unit_price: string
      prompt_price_unit: string
      prompt_price: string
      completion_tokens: number
      completion_unit_price: string
      completion_price_unit: string
      completion_price: string
      total_tokens: number
      total_price: string
      currency: string
      latency: number
    }
  }
  created_at: number
}

// Dify 流式事件类型
export interface DifyChatStreamEvent {
  event: string
  task_id?: string
  message_id?: string
  conversation_id?: string
  answer?: string
  created_at?: number
  metadata?: any
}

// Dify 建议问题响应类型
export interface DifySuggestedQuestionsResponse {
  result: string
  data: string[]
}

// Dify 对话请求参数类型
export interface DifyChatRequest {
  query: string
  user: string
  conversation_id?: string
  inputs?: Record<string, any>
  files?: any[]
  auto_generate_name?: boolean
}
