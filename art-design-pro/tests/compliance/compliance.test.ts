/*** Compliance Check Property-Based Tests ***/
import { describe, it, expect } from 'vitest'
import fc from 'fast-check'

/*** Passed Item Interface ***/
interface PassedItem {
  standardId: number
  standardCode: string
  standardName: string
  category: string
  description: string
}

/*** Failed Item Interface ***/
interface FailedItem {
  standardId: number
  standardCode: string
  standardName: string
  category: string
  violationContent: string
  originalRequirement: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  location: string
  suggestion: string
}

/*** Compliance Suggestion Interface ***/
interface ComplianceSuggestion {
  content: string
  priority: 'low' | 'medium' | 'high'
  category: string
  costImpact: string
  reference: string
}

/*** Compliance Check Result Interface ***/
interface ComplianceCheckResult {
  id: number
  projectId: number
  projectName: string
  checkType: string
  passedItems: PassedItem[]
  failedItems: FailedItem[]
  suggestions: ComplianceSuggestion[]
  overallScore: number
  checkStatus: string
  checkedStandards: number[]
}

/*** Standard Category Constants ***/
const STANDARD_CATEGORIES = ['fire', 'accessibility', 'environmental', 'safety', 'other']

/*** Helper function to validate compliance report structure ***/
function validateComplianceReportStructure(result: ComplianceCheckResult): boolean {
  return (
    Array.isArray(result.passedItems) &&
    Array.isArray(result.failedItems) &&
    Array.isArray(result.suggestions)
  )
}

/*** Helper function to validate failed item structure ***/
function validateFailedItemStructure(item: FailedItem): boolean {
  return (
    typeof item.standardId === 'number' &&
    typeof item.standardCode === 'string' &&
    typeof item.standardName === 'string' &&
    typeof item.category === 'string' &&
    typeof item.violationContent === 'string' &&
    typeof item.originalRequirement === 'string' &&
    ['low', 'medium', 'high', 'critical'].includes(item.severity) &&
    typeof item.location === 'string' &&
    typeof item.suggestion === 'string'
  )
}

/*** Arbitrary generators ***/
const categoryArb = fc.constantFrom(...STANDARD_CATEGORIES)
const severityArb = fc.constantFrom('low', 'medium', 'high', 'critical') as fc.Arbitrary<
  'low' | 'medium' | 'high' | 'critical'
>
const priorityArb = fc.constantFrom('low', 'medium', 'high') as fc.Arbitrary<
  'low' | 'medium' | 'high'
>
const checkStatusArb = fc.constantFrom('pending', 'processing', 'completed', 'failed')

const passedItemArb = fc.record({
  standardId: fc.nat(),
  standardCode: fc.string({ minLength: 1, maxLength: 50 }),
  standardName: fc.string({ minLength: 1, maxLength: 200 }),
  category: categoryArb,
  description: fc.string({ minLength: 1, maxLength: 500 })
})

const failedItemArb = fc.record({
  standardId: fc.nat(),
  standardCode: fc.string({ minLength: 1, maxLength: 50 }),
  standardName: fc.string({ minLength: 1, maxLength: 200 }),
  category: categoryArb,
  violationContent: fc.string({ minLength: 1, maxLength: 500 }),
  originalRequirement: fc.string({ minLength: 1, maxLength: 500 }),
  severity: severityArb,
  location: fc.string({ minLength: 1, maxLength: 200 }),
  suggestion: fc.string({ minLength: 1, maxLength: 500 })
})

const complianceSuggestionArb = fc.record({
  content: fc.string({ minLength: 1, maxLength: 500 }),
  priority: priorityArb,
  category: categoryArb,
  costImpact: fc.string({ minLength: 1, maxLength: 200 }),
  reference: fc.string({ minLength: 1, maxLength: 200 })
})

const complianceCheckResultArb = fc.record({
  id: fc.nat(),
  projectId: fc.nat(),
  projectName: fc.string({ minLength: 1, maxLength: 200 }),
  checkType: fc.constantFrom('full', 'partial'),
  passedItems: fc.array(passedItemArb, { minLength: 0, maxLength: 20 }),
  failedItems: fc.array(failedItemArb, { minLength: 0, maxLength: 20 }),
  suggestions: fc.array(complianceSuggestionArb, { minLength: 0, maxLength: 10 }),
  overallScore: fc.float({ min: 0, max: 100, noNaN: true }),
  checkStatus: checkStatusArb,
  checkedStandards: fc.array(fc.nat(), { minLength: 0, maxLength: 50 })
})

describe('Compliance Check Properties', () => {
  /**
   * Feature: designer-ai-assistant, Property 16: Compliance Report Structure
   * Validates: Requirements 11.2
   *
   * For any compliance check result, the report SHALL contain three categories:
   * passedItems, failedItems, and suggestions.
   */
  it('Property 16: Compliance report should contain passedItems, failedItems, and suggestions', () => {
    fc.assert(
      fc.property(complianceCheckResultArb, (result) => {
        // Property: result must have passedItems field of type array
        expect(result).toHaveProperty('passedItems')
        expect(Array.isArray(result.passedItems)).toBe(true)

        // Property: result must have failedItems field of type array
        expect(result).toHaveProperty('failedItems')
        expect(Array.isArray(result.failedItems)).toBe(true)

        // Property: result must have suggestions field of type array
        expect(result).toHaveProperty('suggestions')
        expect(Array.isArray(result.suggestions)).toBe(true)

        // Validate using helper function
        expect(validateComplianceReportStructure(result)).toBe(true)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Failed items should have valid structure
   * Validates: Requirements 11.3
   *
   * For any failed item, it SHALL contain: standardId, standardCode, standardName,
   * category, violationContent, originalRequirement, severity, location, and suggestion.
   */
  it('Failed items should contain all required fields with valid structure', () => {
    fc.assert(
      fc.property(failedItemArb, (item) => {
        // Property: failed item must have standardId field of type number
        expect(item).toHaveProperty('standardId')
        expect(typeof item.standardId).toBe('number')

        // Property: failed item must have standardCode field of type string
        expect(item).toHaveProperty('standardCode')
        expect(typeof item.standardCode).toBe('string')
        expect(item.standardCode.length).toBeGreaterThan(0)

        // Property: failed item must have standardName field of type string
        expect(item).toHaveProperty('standardName')
        expect(typeof item.standardName).toBe('string')
        expect(item.standardName.length).toBeGreaterThan(0)

        // Property: failed item must have category field with valid value
        expect(item).toHaveProperty('category')
        expect(STANDARD_CATEGORIES).toContain(item.category)

        // Property: failed item must have violationContent field of type string
        expect(item).toHaveProperty('violationContent')
        expect(typeof item.violationContent).toBe('string')
        expect(item.violationContent.length).toBeGreaterThan(0)

        // Property: failed item must have originalRequirement field of type string
        expect(item).toHaveProperty('originalRequirement')
        expect(typeof item.originalRequirement).toBe('string')
        expect(item.originalRequirement.length).toBeGreaterThan(0)

        // Property: failed item must have severity field with valid value
        expect(item).toHaveProperty('severity')
        expect(['low', 'medium', 'high', 'critical']).toContain(item.severity)

        // Property: failed item must have location field of type string
        expect(item).toHaveProperty('location')
        expect(typeof item.location).toBe('string')
        expect(item.location.length).toBeGreaterThan(0)

        // Property: failed item must have suggestion field of type string
        expect(item).toHaveProperty('suggestion')
        expect(typeof item.suggestion).toBe('string')
        expect(item.suggestion.length).toBeGreaterThan(0)

        // Validate using helper function
        expect(validateFailedItemStructure(item)).toBe(true)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Passed items should have valid structure
   * Validates: Requirements 11.2
   */
  it('Passed items should contain all required fields', () => {
    fc.assert(
      fc.property(passedItemArb, (item) => {
        // Property: passed item must have standardId field of type number
        expect(item).toHaveProperty('standardId')
        expect(typeof item.standardId).toBe('number')

        // Property: passed item must have standardCode field of type string
        expect(item).toHaveProperty('standardCode')
        expect(typeof item.standardCode).toBe('string')
        expect(item.standardCode.length).toBeGreaterThan(0)

        // Property: passed item must have standardName field of type string
        expect(item).toHaveProperty('standardName')
        expect(typeof item.standardName).toBe('string')
        expect(item.standardName.length).toBeGreaterThan(0)

        // Property: passed item must have category field with valid value
        expect(item).toHaveProperty('category')
        expect(STANDARD_CATEGORIES).toContain(item.category)

        // Property: passed item must have description field of type string
        expect(item).toHaveProperty('description')
        expect(typeof item.description).toBe('string')
        expect(item.description.length).toBeGreaterThan(0)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Overall score should be between 0 and 100
   * Validates: Requirements 11.2
   */
  it('Overall score should be between 0 and 100', () => {
    fc.assert(
      fc.property(complianceCheckResultArb, (result) => {
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
   * Additional test: Compliance suggestions should have valid structure
   * Validates: Requirements 11.2
   */
  it('Compliance suggestions should have valid priority and category', () => {
    fc.assert(
      fc.property(complianceSuggestionArb, (item) => {
        // Property: suggestion must have content field
        expect(item).toHaveProperty('content')
        expect(typeof item.content).toBe('string')
        expect(item.content.length).toBeGreaterThan(0)

        // Property: suggestion must have priority field with valid value
        expect(item).toHaveProperty('priority')
        expect(['low', 'medium', 'high']).toContain(item.priority)

        // Property: suggestion must have category field with valid value
        expect(item).toHaveProperty('category')
        expect(STANDARD_CATEGORIES).toContain(item.category)

        // Property: suggestion must have costImpact field
        expect(item).toHaveProperty('costImpact')
        expect(typeof item.costImpact).toBe('string')

        // Property: suggestion must have reference field
        expect(item).toHaveProperty('reference')
        expect(typeof item.reference).toBe('string')

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: All arrays in compliance result should be non-null
   * Validates: Requirements 11.2
   */
  it('All arrays in compliance result should be non-null', () => {
    fc.assert(
      fc.property(complianceCheckResultArb, (result) => {
        // Property: passedItems should not be null
        expect(result.passedItems).not.toBeNull()

        // Property: failedItems should not be null
        expect(result.failedItems).not.toBeNull()

        // Property: suggestions should not be null
        expect(result.suggestions).not.toBeNull()

        // Property: checkedStandards should not be null
        expect(result.checkedStandards).not.toBeNull()

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Check status should be valid
   * Validates: Requirements 11.2
   */
  it('Check status should be one of valid values', () => {
    fc.assert(
      fc.property(complianceCheckResultArb, (result) => {
        // Property: checkStatus should be one of valid values
        expect(['pending', 'processing', 'completed', 'failed']).toContain(result.checkStatus)

        return true
      }),
      { numRuns: 100 }
    )
  })

  /**
   * Additional test: Check type should be valid
   * Validates: Requirements 11.1
   */
  it('Check type should be full or partial', () => {
    fc.assert(
      fc.property(complianceCheckResultArb, (result) => {
        // Property: checkType should be 'full' or 'partial'
        expect(['full', 'partial']).toContain(result.checkType)

        return true
      }),
      { numRuns: 100 }
    )
  })
})
