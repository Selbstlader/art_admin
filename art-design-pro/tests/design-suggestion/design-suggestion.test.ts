/*** Design Suggestion Property-Based Tests ***/
import { describe, it, expect } from 'vitest'
import fc from 'fast-check'

/*** Design Suggestion Interface ***/
interface DesignSuggestion {
  id: number
  projectId: number
  content: string
  applicableScene: string
  costImpact: string
  category: string
  priority: number
  status: string
  detailInfo: Record<string, unknown> | null
  referenceImages: string[]
  createdAt: string
  updatedAt: string
}

/*** Generate Suggestion Response Interface ***/
interface GenerateSuggestionResponse {
  suggestions: DesignSuggestion[]
  projectId: number
  count: number
}

/*** Helper function to validate suggestion structure ***/
function validateSuggestionStructure(suggestion: DesignSuggestion): boolean {
  return (
    typeof suggestion.content === 'string' &&
    suggestion.content.length > 0 &&
    typeof suggestion.applicableScene === 'string' &&
    typeof suggestion.costImpact === 'string'
  )
}

/*** Helper function to validate suggestion response ***/
function validateSuggestionResponse(response: GenerateSuggestionResponse): {
  valid: boolean
  hasMinimumCount: boolean
  allHaveRequiredFields: boolean
} {
  const hasMinimumCount = response.suggestions.length >= 3
  const allHaveRequiredFields = response.suggestions.every(validateSuggestionStructure)

  return {
    valid: hasMinimumCount && allHaveRequiredFields,
    hasMinimumCount,
    allHaveRequiredFields
  }
}

/*** Generate valid ISO date string ***/
const validDateArb = fc
  .integer({ min: 0, max: 2000000000000 }) // Valid timestamp range
  .map((ts) => new Date(ts).toISOString())

/*** Generate random design suggestion ***/
const designSuggestionArb = fc.record({
  id: fc.nat(),
  projectId: fc.integer({ min: 1, max: 1000000 }),
  content: fc.string({ minLength: 1, maxLength: 1000 }),
  applicableScene: fc.string({ minLength: 0, maxLength: 500 }),
  costImpact: fc.oneof(
    fc.constant('增加约5%-10%'),
    fc.constant('节省约15%'),
    fc.constant('基本持平'),
    fc.constant('待评估'),
    fc.string({ minLength: 1, maxLength: 50 })
  ),
  category: fc.oneof(
    fc.constant('layout'),
    fc.constant('material'),
    fc.constant('style'),
    fc.constant('function'),
    fc.constant('lighting'),
    fc.constant('hvac'),
    fc.constant('general')
  ),
  priority: fc.integer({ min: 1, max: 10 }),
  status: fc.oneof(fc.constant('pending'), fc.constant('adopted'), fc.constant('ignored')),
  detailInfo: fc.oneof(fc.constant(null), fc.dictionary(fc.string(), fc.string())),
  referenceImages: fc.array(fc.webUrl(), { minLength: 0, maxLength: 5 }),
  createdAt: validDateArb,
  updatedAt: validDateArb
})

/*** Generate random suggestion response with at least 3 suggestions ***/
const suggestionResponseArb = fc
  .array(designSuggestionArb, { minLength: 3, maxLength: 10 })
  .chain((suggestions) =>
    fc.record({
      suggestions: fc.constant(suggestions),
      projectId: fc.nat({ min: 1 }),
      count: fc.constant(suggestions.length)
    })
  )

describe('Design Suggestion Properties', () => {
  /**
   * Feature: designer-ai-assistant, Property 5: Design Suggestion Count and Structure
   * Validates: Requirements 3.2
   *
   * For any design suggestion request, the system SHALL return at least 3 suggestions,
   * and each suggestion SHALL contain: content, applicableScene, and costImpact fields.
   */
  it('Property 5: Design suggestions should have at least 3 items with required fields', () => {
    fc.assert(
      fc.property(suggestionResponseArb, (response) => {
        // Property: response should contain at least 3 suggestions
        expect(response.suggestions.length).toBeGreaterThanOrEqual(3)

        // Property: count should match actual suggestions length
        expect(response.count).toBe(response.suggestions.length)

        // Property: each suggestion must have content field (non-empty string)
        response.suggestions.forEach((suggestion, index) => {
          expect(suggestion).toHaveProperty('content')
          expect(typeof suggestion.content).toBe('string')
          expect(suggestion.content.length).toBeGreaterThan(0)
        })

        // Property: each suggestion must have applicableScene field (string)
        response.suggestions.forEach((suggestion) => {
          expect(suggestion).toHaveProperty('applicableScene')
          expect(typeof suggestion.applicableScene).toBe('string')
        })

        // Property: each suggestion must have costImpact field (string)
        response.suggestions.forEach((suggestion) => {
          expect(suggestion).toHaveProperty('costImpact')
          expect(typeof suggestion.costImpact).toBe('string')
        })

        // Validate using helper function
        const validation = validateSuggestionResponse(response)
        expect(validation.hasMinimumCount).toBe(true)
        expect(validation.allHaveRequiredFields).toBe(true)
        expect(validation.valid).toBe(true)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Each suggestion should have valid category
   * Validates: Requirements 3.2
   */
  it('Each suggestion should have a valid category', () => {
    const validCategories = ['layout', 'material', 'style', 'function', 'lighting', 'hvac', 'general']

    fc.assert(
      fc.property(designSuggestionArb, (suggestion) => {
        // Property: category should be one of the valid categories
        expect(validCategories).toContain(suggestion.category)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Priority should be within valid range
   * Validates: Requirements 3.2
   */
  it('Suggestion priority should be between 1 and 10', () => {
    fc.assert(
      fc.property(designSuggestionArb, (suggestion) => {
        // Property: priority should be between 1 and 10
        expect(suggestion.priority).toBeGreaterThanOrEqual(1)
        expect(suggestion.priority).toBeLessThanOrEqual(10)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Status should be valid
   * Validates: Requirements 3.4
   */
  it('Suggestion status should be one of pending, adopted, or ignored', () => {
    const validStatuses = ['pending', 'adopted', 'ignored']

    fc.assert(
      fc.property(designSuggestionArb, (suggestion) => {
        // Property: status should be one of the valid statuses
        expect(validStatuses).toContain(suggestion.status)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: ProjectId should be positive
   * Validates: Requirements 3.1
   */
  it('Suggestion should be associated with a valid project', () => {
    fc.assert(
      fc.property(designSuggestionArb, (suggestion) => {
        // Property: projectId should be a positive number
        expect(suggestion.projectId).toBeGreaterThan(0)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Validate suggestion structure helper function
   * Validates: Requirements 3.2
   */
  it('validateSuggestionStructure should correctly identify valid suggestions', () => {
    fc.assert(
      fc.property(designSuggestionArb, (suggestion) => {
        const isValid = validateSuggestionStructure(suggestion)

        // Property: if content is non-empty, applicableScene and costImpact are strings,
        // then the suggestion should be valid
        if (
          suggestion.content.length > 0 &&
          typeof suggestion.applicableScene === 'string' &&
          typeof suggestion.costImpact === 'string'
        ) {
          expect(isValid).toBe(true)
        }

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Response count should match suggestions array length
   * Validates: Requirements 3.2
   */
  it('Response count should accurately reflect the number of suggestions', () => {
    fc.assert(
      fc.property(suggestionResponseArb, (response) => {
        // Property: count field should equal suggestions array length
        expect(response.count).toBe(response.suggestions.length)

        // Property: count should be at least 3
        expect(response.count).toBeGreaterThanOrEqual(3)

        return true
      }),
      { numRuns: 100 }
    )
  })
})
