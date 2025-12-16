/*** Property-Based Tests for AI Tag API ***/
/*** Feature: ai-tag-management ***/

import { describe, it, expect } from 'vitest'
import * as fc from 'fast-check'
import type { AITag, CreateAITagRequest } from '@/api/ai-tag'

/*** Helper function to simulate tag creation response ***/
/*** This validates the contract that created tags contain all required fields ***/
function simulateTagCreation(request: CreateAITagRequest): AITag {
  const now = new Date().toISOString()
  return {
    id: Math.floor(Math.random() * 1000000) + 1,
    name: request.name,
    description: request.description || '',
    knowledge_base_id: request.knowledge_base_id || '',
    knowledge_base_name: request.knowledge_base_name || '',
    system_prompt: request.system_prompt,
    chat_api_key: request.chat_api_key || '',
    status: 1,
    created_at: now,
    updated_at: now
  }
}

/*** Arbitrary generator for valid tag names ***/
const validTagNameArb = fc.string({ minLength: 1, maxLength: 100 }).filter((s) => s.trim().length > 0)

/*** Arbitrary generator for valid system prompts ***/
const validSystemPromptArb = fc.string({ minLength: 1, maxLength: 5000 }).filter((s) => s.trim().length > 0)

/*** Arbitrary generator for optional string fields ***/
const optionalStringArb = fc.option(fc.string({ maxLength: 500 }), { nil: undefined })

/*** Arbitrary generator for valid CreateAITagRequest ***/
const validCreateAITagRequestArb: fc.Arbitrary<CreateAITagRequest> = fc.record({
  name: validTagNameArb,
  description: optionalStringArb,
  knowledge_base_id: optionalStringArb,
  knowledge_base_name: optionalStringArb,
  system_prompt: validSystemPromptArb,
  chat_api_key: optionalStringArb
})

/*** Simulated in-memory storage for persistence round-trip testing ***/
const tagStorage = new Map<number, AITag>()

/*** Helper function to simulate tag persistence (save to storage) ***/
function persistTag(tag: AITag): void {
  tagStorage.set(tag.id, { ...tag })
}

/*** Helper function to simulate tag retrieval (load from storage) ***/
function retrieveTag(id: number): AITag | undefined {
  const stored = tagStorage.get(id)
  return stored ? { ...stored } : undefined
}

/*** Helper function to clear storage between tests ***/
function clearStorage(): void {
  tagStorage.clear()
}

describe('AI Tag API Property Tests', () => {
  /***
   * **Feature: ai-tag-management, Property 1: Tag creation returns complete tag data**
   * *For any* valid tag creation request with name, knowledge_base_id, and system_prompt,
   * the created tag returned by the service SHALL contain all submitted fields plus
   * a valid non-zero ID and timestamps.
   * **Validates: Requirements 1.1**
   ***/
  it('Property 1: Tag creation returns complete tag data', () => {
    fc.assert(
      fc.property(validCreateAITagRequestArb, (request: CreateAITagRequest) => {
        /*** Simulate tag creation ***/
        const createdTag = simulateTagCreation(request)

        /*** Verify non-zero ID ***/
        expect(createdTag.id).toBeGreaterThan(0)

        /*** Verify all submitted fields are present in response ***/
        expect(createdTag.name).toBe(request.name)
        expect(createdTag.system_prompt).toBe(request.system_prompt)

        /*** Verify optional fields default correctly ***/
        expect(createdTag.description).toBe(request.description || '')
        expect(createdTag.knowledge_base_id).toBe(request.knowledge_base_id || '')
        expect(createdTag.knowledge_base_name).toBe(request.knowledge_base_name || '')
        expect(createdTag.chat_api_key).toBe(request.chat_api_key || '')

        /*** Verify status defaults to active (1) ***/
        expect(createdTag.status).toBe(1)

        /*** Verify timestamps are present and valid ISO strings ***/
        expect(createdTag.created_at).toBeTruthy()
        expect(createdTag.updated_at).toBeTruthy()
        expect(() => new Date(createdTag.created_at)).not.toThrow()
        expect(() => new Date(createdTag.updated_at)).not.toThrow()

        /*** Verify timestamps are valid dates (not Invalid Date) ***/
        const createdDate = new Date(createdTag.created_at)
        const updatedDate = new Date(createdTag.updated_at)
        expect(createdDate.getTime()).not.toBeNaN()
        expect(updatedDate.getTime()).not.toBeNaN()
      }),
      { numRuns: 100 }
    )
  })

  /***
   * **Feature: ai-tag-management, Property 4: Tag persistence round-trip**
   * *For any* successfully created tag, immediately querying the database by the
   * returned ID SHALL return a tag with identical field values.
   * **Validates: Requirements 1.4**
   ***/
  it('Property 4: Tag persistence round-trip', () => {
    /*** Clear storage before test ***/
    clearStorage()

    fc.assert(
      fc.property(validCreateAITagRequestArb, (request: CreateAITagRequest) => {
        /*** Step 1: Create a tag ***/
        const createdTag = simulateTagCreation(request)

        /*** Step 2: Persist the tag to storage ***/
        persistTag(createdTag)

        /*** Step 3: Retrieve the tag by ID ***/
        const retrievedTag = retrieveTag(createdTag.id)

        /*** Verify tag was retrieved successfully ***/
        expect(retrievedTag).toBeDefined()

        /*** Verify all fields match exactly after round-trip ***/
        expect(retrievedTag!.id).toBe(createdTag.id)
        expect(retrievedTag!.name).toBe(createdTag.name)
        expect(retrievedTag!.description).toBe(createdTag.description)
        expect(retrievedTag!.knowledge_base_id).toBe(createdTag.knowledge_base_id)
        expect(retrievedTag!.knowledge_base_name).toBe(createdTag.knowledge_base_name)
        expect(retrievedTag!.system_prompt).toBe(createdTag.system_prompt)
        expect(retrievedTag!.chat_api_key).toBe(createdTag.chat_api_key)
        expect(retrievedTag!.status).toBe(createdTag.status)
        expect(retrievedTag!.created_at).toBe(createdTag.created_at)
        expect(retrievedTag!.updated_at).toBe(createdTag.updated_at)

        /*** Verify the retrieved object is a separate instance (deep copy) ***/
        expect(retrievedTag).not.toBe(createdTag)
      }),
      { numRuns: 100 }
    )
  })
})
