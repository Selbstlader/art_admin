/*** AI Generate API Module ***/
/*** Provides AI-powered content generation using DeepSeek ***/

import request from '@/utils/http'

/*** TypeScript Interfaces ***/

/** Request interface for generating/polishing prompt */
export interface GeneratePromptRequest {
  keywords: string
  current_prompt?: string
  mode?: 'generate' | 'polish'
}

/** Response interface for generated prompt */
export interface GeneratePromptResponse {
  prompt: string
}

/*** AI Generate API Functions ***/

export const aiGenerateApi = {
  /**
   * Generate or polish a system prompt using AI
   * @param data - Generation request with keywords and optional current prompt
   * @returns Generated or polished prompt
   */
  generatePrompt: (data: GeneratePromptRequest) => {
    return request.post<GeneratePromptResponse>({
      url: '/api/ai-generate/prompt',
      params: data,
      timeout: 120000 // AI generation may take longer, set 120s timeout
    })
  }
}

export default aiGenerateApi
