/*** AI Tag API Module ***/
/*** Requirements: 1.1, 2.1, 3.1, 4.1, 6.1 ***/

import request from '@/utils/http'

/*** TypeScript Interfaces for Request/Response ***/

/** AI Tag entity interface */
export interface AITag {
  id: number
  name: string
  description: string
  knowledge_base_id: string
  knowledge_base_name: string
  system_prompt: string
  chat_api_key: string
  status: number
  created_at: string
  updated_at: string
}

/** Request interface for listing AI tags with pagination and search */
export interface AITagListParams {
  current: number
  size: number
  keyword?: string
}

/** Response interface for paginated AI tag list */
export interface AITagListResponse {
  list: AITag[]
  records: AITag[]
  total: number
  page: number
  size: number
}

/** Request interface for creating an AI tag */
export interface CreateAITagRequest {
  name: string
  description?: string
  knowledge_base_id?: string
  knowledge_base_name?: string
  system_prompt: string
  chat_api_key?: string
}

/** Request interface for updating an AI tag */
export interface UpdateAITagRequest {
  name?: string
  description?: string
  knowledge_base_id?: string
  knowledge_base_name?: string
  system_prompt?: string
  chat_api_key?: string
  status?: number
}

/** Response interface for AI tag test result */
export interface AITagTestResponse {
  success: boolean
  response: string
  knowledge_accessed: boolean
  error_message?: string
}

/*** AI Tag API Functions ***/

export const aiTagApi = {
  /**
   * Create a new AI tag
   * Requirements: 1.1 - Create new AI_Tag record and return created tag details
   * @param data - Tag creation data
   * @returns Created AI tag
   */
  create: (data: CreateAITagRequest) => {
    return request.post<AITag>({
      url: '/api/ai-tags',
      params: data,
      showSuccessMessage: true
    })
  },

  /**
   * Get paginated list of AI tags with optional keyword search
   * Requirements: 2.1 - Return all active AI_Tag records with pagination
   * Requirements: 2.2 - Search tags by keyword
   * @param params - Pagination and search parameters
   * @returns Paginated list of AI tags
   */
  list: (params: AITagListParams) => {
    return request.get<AITagListResponse>({
      url: '/api/ai-tags',
      params
    })
  },

  /**
   * Get AI tag details by ID
   * Requirements: 2.3 - Include tag name, description, knowledge base name, etc.
   * @param id - Tag ID
   * @returns AI tag details
   */
  get: (id: number) => {
    return request.get<AITag>({
      url: `/api/ai-tags/${id}`
    })
  },

  /**
   * Update an existing AI tag
   * Requirements: 3.1 - Update AI_Tag record and return updated tag details
   * @param id - Tag ID
   * @param data - Tag update data
   * @returns Updated AI tag
   */
  update: (id: number, data: UpdateAITagRequest) => {
    return request.put<AITag>({
      url: `/api/ai-tags/${id}`,
      params: data,
      showSuccessMessage: true
    })
  },

  /**
   * Soft delete an AI tag
   * Requirements: 4.1 - Perform soft delete by setting status to inactive
   * @param id - Tag ID
   * @returns Deletion result
   */
  delete: (id: number) => {
    return request.del<void>({
      url: `/api/ai-tags/${id}`,
      showSuccessMessage: true
    })
  },

  /**
   * Test AI tag configuration
   * Requirements: 6.1 - Initiate test chat session with predefined test query
   * Requirements: 6.2 - Display AI response and knowledge base access status
   * @param id - Tag ID
   * @returns Test result with success status and response
   */
  test: (id: number) => {
    return request.post<AITagTestResponse>({
      url: `/api/ai-tags/${id}/test`
    })
  }
}

export default aiTagApi
