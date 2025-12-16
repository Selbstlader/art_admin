package repository

import (
	"art_admin_backend/internal/model"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupAITagTestDB creates an in-memory SQLite database for testing
func setupAITagTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&model.AITag{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// genAITagForDB generates random AITag instances suitable for database insertion
func genAITagForDB() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }), // Name
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) <= 500 }),              // Description
		gen.Identifier(), // KnowledgeBaseID
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) <= 200 }), // KnowledgeBaseName
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),    // SystemPrompt
		gen.Identifier(),   // ChatAPIKey
		gen.IntRange(0, 1), // Status
	).Map(func(values []interface{}) model.AITag {
		return model.AITag{
			Name:              values[0].(string),
			Description:       values[1].(string),
			KnowledgeBaseID:   values[2].(string),
			KnowledgeBaseName: values[3].(string),
			SystemPrompt:      values[4].(string),
			ChatAPIKey:        values[5].(string),
			Status:            values[6].(int),
		}
	})
}

// **Feature: ai-tag-management, Property 2: Duplicate name rejection on create**
// **Validates: Requirements 1.2**
// For any existing tag name in the database, attempting to create a new tag with the same name
// SHALL result in a rejection error.
func TestAITagDuplicateNameRejection(t *testing.T) {
	db := setupAITagTestDB(t)
	repo := NewAITagRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for AITag data
	aiTagGen := genAITagForDB()

	properties.Property("Duplicate name creation is rejected", prop.ForAll(
		func(tag model.AITag) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM ai_tags")

			// Create the first tag
			firstTag := tag
			if err := repo.Create(&firstTag); err != nil {
				t.Logf("Failed to create first tag: %v", err)
				return false
			}

			// Verify first tag was created
			found, err := repo.FindByID(firstTag.ID)
			if err != nil || found == nil {
				t.Logf("First tag should exist after creation")
				return false
			}

			// Attempt to create a second tag with the same name
			secondTag := model.AITag{
				Name:              tag.Name, // Same name as first tag
				Description:       "Different description",
				KnowledgeBaseID:   "different-kb-id",
				KnowledgeBaseName: "Different KB",
				SystemPrompt:      "Different prompt",
				ChatAPIKey:        "different-api-key",
				Status:            1,
			}

			err = repo.Create(&secondTag)

			// The creation should fail due to unique constraint on name
			if err == nil {
				t.Logf("Second tag with duplicate name should have been rejected")
				return false
			}

			// Verify only one tag exists with this name
			existingTag, err := repo.FindByName(tag.Name)
			if err != nil {
				t.Logf("Should be able to find the original tag: %v", err)
				return false
			}
			if existingTag.ID != firstTag.ID {
				t.Logf("Found tag should be the first one created")
				return false
			}

			return true
		},
		aiTagGen,
	))

	properties.TestingRun(t)
}

// **Feature: ai-tag-management, Property 5: Pagination correctness**
// **Validates: Requirements 2.1**
// For any list request with page P and pageSize S on a database with N tags,
// the returned list SHALL contain at most S items, and the total count SHALL equal N.
func TestAITagPaginationCorrectness(t *testing.T) {
	db := setupAITagTestDB(t)
	repo := NewAITagRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for test parameters
	paramsGen := gopter.CombineGens(
		gen.IntRange(1, 50), // Number of tags to create (N)
		gen.IntRange(1, 10), // Page number (P)
		gen.IntRange(1, 20), // Page size (S)
	)

	properties.Property("Pagination returns correct count and respects page size", prop.ForAll(
		func(params []interface{}) bool {
			numTags := params[0].(int)
			page := params[1].(int)
			pageSize := params[2].(int)

			// Clean up previous test data
			db.Exec("DELETE FROM ai_tags")

			// Create N tags with unique names
			for i := 0; i < numTags; i++ {
				tag := model.AITag{
					Name:         generateUniqueName(i),
					Description:  "Test description",
					SystemPrompt: "Test prompt",
					Status:       1,
				}
				if err := repo.Create(&tag); err != nil {
					t.Logf("Failed to create tag %d: %v", i, err)
					return false
				}
			}

			// Query with pagination
			tags, total, err := repo.List(page, pageSize, "")
			if err != nil {
				t.Logf("Failed to list tags: %v", err)
				return false
			}

			// Verify total count equals N
			if total != int64(numTags) {
				t.Logf("Total count mismatch: expected %d, got %d", numTags, total)
				return false
			}

			// Calculate expected number of items on this page
			offset := (page - 1) * pageSize
			expectedItems := numTags - offset
			if expectedItems < 0 {
				expectedItems = 0
			}
			if expectedItems > pageSize {
				expectedItems = pageSize
			}

			// Verify returned items count
			if len(tags) != expectedItems {
				t.Logf("Items count mismatch: expected %d, got %d (page=%d, pageSize=%d, total=%d)",
					expectedItems, len(tags), page, pageSize, numTags)
				return false
			}

			// Verify returned items don't exceed page size
			if len(tags) > pageSize {
				t.Logf("Returned %d items, exceeds page size %d", len(tags), pageSize)
				return false
			}

			return true
		},
		paramsGen,
	))

	properties.TestingRun(t)
}

// generateUniqueName generates a unique name for testing
func generateUniqueName(index int) string {
	return "TestTag_" + strings.Repeat("x", index%50) + "_" + string(rune('A'+index%26))
}

// **Feature: ai-tag-management, Property 6: Keyword search completeness**
// **Validates: Requirements 2.2**
// For any keyword K and set of tags, the search result SHALL include all and only tags
// where name OR description contains K as a substring (case-insensitive).
func TestAITagKeywordSearchCompleteness(t *testing.T) {
	db := setupAITagTestDB(t)
	repo := NewAITagRepositoryWithDB(db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for keyword - use a simple fixed-length approach to avoid discards
	keywordGen := gen.SliceOfN(4, gen.AlphaLowerChar()).Map(func(chars []rune) string {
		return string(chars)
	})

	// Generator for test parameters
	paramsGen := gopter.CombineGens(
		gen.IntRange(5, 20), // Number of tags to create
		keywordGen,          // Keyword to search (4 lowercase letters)
	)

	properties.Property("Keyword search returns all and only matching tags", prop.ForAll(
		func(params []interface{}) bool {
			numTags := params[0].(int)
			keyword := params[1].(string)

			// Clean up previous test data
			db.Exec("DELETE FROM ai_tags")

			// Create tags - some with keyword in name, some in description, some without
			var expectedMatches []int64
			keywordLower := strings.ToLower(keyword)

			for i := 0; i < numTags; i++ {
				var name, description string

				switch i % 4 {
				case 0:
					// Keyword in name
					name = "Prefix" + keyword + "Suffix"
					description = "No match here"
				case 1:
					// Keyword in description
					name = "NoMatchName" + string(rune('A'+i%26))
					description = "Contains " + keyword + " in description"
				case 2:
					// Keyword in both
					name = keyword + "InName"
					description = "Also " + keyword + " here"
				default:
					// No keyword
					name = "UnrelatedName" + string(rune('A'+i%26))
					description = "Unrelated description"
				}

				tag := model.AITag{
					Name:         name + "_" + string(rune('0'+i%10)),
					Description:  description,
					SystemPrompt: "Test prompt",
					Status:       1,
				}

				if err := repo.Create(&tag); err != nil {
					t.Logf("Failed to create tag %d: %v", i, err)
					return false
				}

				// Track expected matches (case-insensitive)
				if strings.Contains(strings.ToLower(tag.Name), keywordLower) ||
					strings.Contains(strings.ToLower(tag.Description), keywordLower) {
					expectedMatches = append(expectedMatches, tag.ID)
				}
			}

			// Search with keyword
			tags, _, err := repo.List(1, 100, keyword)
			if err != nil {
				t.Logf("Failed to search tags: %v", err)
				return false
			}

			// Verify all returned tags contain the keyword (case-insensitive)
			for _, tag := range tags {
				nameLower := strings.ToLower(tag.Name)
				descLower := strings.ToLower(tag.Description)
				if !strings.Contains(nameLower, keywordLower) && !strings.Contains(descLower, keywordLower) {
					t.Logf("Tag %d (%s) doesn't contain keyword '%s' in name or description",
						tag.ID, tag.Name, keyword)
					return false
				}
			}

			// Verify count matches expected
			if len(tags) != len(expectedMatches) {
				t.Logf("Expected %d matches, got %d", len(expectedMatches), len(tags))
				return false
			}

			// Verify all expected matches are in the result
			resultIDs := make(map[int64]bool)
			for _, tag := range tags {
				resultIDs[tag.ID] = true
			}
			for _, expectedID := range expectedMatches {
				if !resultIDs[expectedID] {
					t.Logf("Expected tag %d not found in results", expectedID)
					return false
				}
			}

			return true
		},
		paramsGen,
	))

	properties.TestingRun(t)
}
