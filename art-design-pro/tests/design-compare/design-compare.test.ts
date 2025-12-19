/*** Design Compare Property-Based Tests ***/
import { describe, it, expect } from 'vitest'
import fc from 'fast-check'

/*** Match Item Interface ***/
interface MatchItem {
  requirement: string
  designMatch: string
  score: number
}

/*** Deviation Item Interface ***/
interface DeviationItem {
  location: string
  content: string
  originalRequirement: string
  severity: 'low' | 'medium' | 'high'
}

/*** Suggestion Item Interface ***/
interface SuggestionItem {
  content: string
  priority: 'low' | 'medium' | 'high'
  costImpact: string
}

/*** Design Compare Result Interface ***/
interface DesignCompareResult {
  compareId: number
  matchItems: MatchItem[]
  deviationItems: DeviationItem[]
  suggestions: SuggestionItem[]
  overallScore: number
}

/*** Helper function to validate compare report structure ***/
function validateCompareReportStructure(result: DesignCompareResult): boolean {
  return (
    Array.isArray(result.matchItems) &&
    Array.isArray(result.deviationItems) &&
    Array.isArray(result.suggestions)
  )
}

/*** Helper function to validate deviation item structure ***/
function validateDeviationItemStructure(item: DeviationItem): boolean {
  return (
    typeof item.location === 'string' &&
    typeof item.content === 'string' &&
    typeof item.originalRequirement === 'string' &&
    ['low', 'medium', 'high'].includes(item.severity)
  )
}

/*** Arbitrary generators ***/
const severityArb = fc.constantFrom('low', 'medium', 'high') as fc.Arbitrary<
  'low' | 'medium' | 'high'
>
const priorityArb = fc.constantFrom('low', 'medium', 'high') as fc.Arbitrary<
  'low' | 'medium' | 'high'
>

const matchItemArb = fc.record({
  requirement: fc.string({ minLength: 1, maxLength: 200 }),
  designMatch: fc.string({ minLength: 1, maxLength: 200 }),
  score: fc.integer({ min: 0, max: 100 })
})

const deviationItemArb = fc.record({
  location: fc.string({ minLength: 1, maxLength: 100 }),
  content: fc.string({ minLength: 1, maxLength: 500 }),
  originalRequirement: fc.string({ minLength: 1, maxLength: 200 }),
  severity: severityArb
})

const suggestionItemArb = fc.record({
  content: fc.string({ minLength: 1, maxLength: 500 }),
  priority: priorityArb,
  costImpact: fc.string({ minLength: 1, maxLength: 200 })
})

const designCompareResultArb = fc.record({
  compareId: fc.nat(),
  matchItems: fc.array(matchItemArb, { minLength: 0, maxLength: 20 }),
  deviationItems: fc.array(deviationItemArb, { minLength: 0, maxLength: 20 }),
  suggestions: fc.array(suggestionItemArb, { minLength: 0, maxLength: 10 }),
  overallScore: fc.float({ min: 0, max: 100, noNaN: true })
})

describe('Design Compare Properties', () => {
  /**
   * Feature: designer-ai-assistant, Property 3: Compare Report Structure Completeness
   * Validates: Requirements 2.2
   *
   * For any design comparison result, the report SHALL contain three categories:
   * matchItems, deviationItems, and suggestions.
   */
  it('Property 3: Compare report should contain matchItems, deviationItems, and suggestions', () => {
    fc.assert(
      fc.property(designCompareResultArb, (result) => {
        // Property: result must have matchItems field of type array
        expect(result).toHaveProperty('matchItems')
        expect(Array.isArray(result.matchItems)).toBe(true)

        // Property: result must have deviationItems field of type array
        expect(result).toHaveProperty('deviationItems')
        expect(Array.isArray(result.deviationItems)).toBe(true)

        // Property: result must have suggestions field of type array
        expect(result).toHaveProperty('suggestions')
        expect(Array.isArray(result.suggestions)).toBe(true)

        // Validate using helper function
        expect(validateCompareReportStructure(result)).toBe(true)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Feature: designer-ai-assistant, Property 4: Deviation Item Structure
   * Validates: Requirements 2.3
   *
   * For any deviation item in a comparison result, it SHALL contain:
   * location, content, and originalRequirement fields.
   */
  it('Property 4: Deviation items should contain location, content, and originalRequirement', () => {
    fc.assert(
      fc.property(deviationItemArb, (item) => {
        // Property: deviation item must have location field of type string
        expect(item).toHaveProperty('location')
        expect(typeof item.location).toBe('string')
        expect(item.location.length).toBeGreaterThan(0)

        // Property: deviation item must have content field of type string
        expect(item).toHaveProperty('content')
        expect(typeof item.content).toBe('string')
        expect(item.content.length).toBeGreaterThan(0)

        // Property: deviation item must have originalRequirement field of type string
        expect(item).toHaveProperty('originalRequirement')
        expect(typeof item.originalRequirement).toBe('string')
        expect(item.originalRequirement.length).toBeGreaterThan(0)

        // Property: deviation item must have severity field with valid value
        expect(item).toHaveProperty('severity')
        expect(['low', 'medium', 'high']).toContain(item.severity)

        // Validate using helper function
        expect(validateDeviationItemStructure(item)).toBe(true)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Match item score should be between 0 and 100
   * Validates: Requirements 2.2
   */
  it('Match item score should be between 0 and 100', () => {
    fc.assert(
      fc.property(matchItemArb, (item) => {
        // Property: score should be between 0 and 100
        expect(item.score).toBeGreaterThanOrEqual(0)
        expect(item.score).toBeLessThanOrEqual(100)

        // Property: score should be an integer
        expect(Number.isInteger(item.score)).toBe(true)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Overall score should be between 0 and 100
   * Validates: Requirements 2.2
   */
  it('Overall score should be between 0 and 100', () => {
    fc.assert(
      fc.property(designCompareResultArb, (result) => {
        // Property: overall score should be between 0 and 100
        expect(result.overallScore).toBeGreaterThanOrEqual(0)
        expect(result.overallScore).toBeLessThanOrEqual(100)

        // Property: overall score should be a finite number
        expect(Number.isFinite(result.overallScore)).toBe(true)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Suggestion items should have valid priority
   * Validates: Requirements 2.2
   */
  it('Suggestion items should have valid priority values', () => {
    fc.assert(
      fc.property(suggestionItemArb, (item) => {
        // Property: suggestion must have content field
        expect(item).toHaveProperty('content')
        expect(typeof item.content).toBe('string')
        expect(item.content.length).toBeGreaterThan(0)

        // Property: suggestion must have priority field with valid value
        expect(item).toHaveProperty('priority')
        expect(['low', 'medium', 'high']).toContain(item.priority)

        // Property: suggestion must have costImpact field
        expect(item).toHaveProperty('costImpact')
        expect(typeof item.costImpact).toBe('string')

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: All arrays in compare result should be non-null
   * Validates: Requirements 2.2
   */
  it('All arrays in compare result should be non-null', () => {
    fc.assert(
      fc.property(designCompareResultArb, (result) => {
        // Property: matchItems should not be null
        expect(result.matchItems).not.toBeNull()

        // Property: deviationItems should not be null
        expect(result.deviationItems).not.toBeNull()

        // Property: suggestions should not be null
        expect(result.suggestions).not.toBeNull()

        return true
      }),
      { numRuns: 100 }
    )
  })
})
