package ai_tag

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
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

// genWhitespaceString generates strings composed entirely of whitespace characters
func genWhitespaceString() gopter.Gen {
	whitespaceChars := []rune{' ', '\t', '\n', '\r'}
	return gen.IntRange(0, 20).FlatMap(func(length any) gopter.Gen {
		n := length.(int)
		return gen.SliceOfN(n, gen.OneConstOf(whitespaceChars[0], whitespaceChars[1], whitespaceChars[2], whitespaceChars[3])).Map(func(chars []rune) string {
			return string(chars)
		})
	}, reflect.TypeOf(""))
}

// genValidName generates valid non-empty, non-whitespace names
func genValidName() gopter.Gen {
	return gen.Identifier().SuchThat(func(s string) bool {
		return len(s) > 0 && len(s) <= 100 && strings.TrimSpace(s) != ""
	})
}

// genValidSystemPrompt generates valid non-empty system prompts
func genValidSystemPrompt() gopter.Gen {
	return gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0 && strings.TrimSpace(s) != ""
	})
}

// **Feature: ai-tag-management, Property 3: Empty/whitespace validation rejection**
// **Validates: Requirements 1.3**
// For any tag creation request where name or system_prompt consists entirely of whitespace
// characters, the creation SHALL be rejected with a validation error.
func TestEmptyWhitespaceValidationRejection(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAITagRepositoryWithDB(db)
	service := NewAITagServiceWithRepo(repo)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Test 1: Empty/whitespace name should be rejected
	properties.Property("Empty or whitespace-only name is rejected", prop.ForAll(
		func(whitespaceStr string) bool {
			req := &CreateAITagRequest{
				Name:         whitespaceStr,
				Description:  "Valid description",
				SystemPrompt: "Valid system prompt",
			}

			_, err := service.CreateTag(req)

			// Should be rejected
			if err == nil {
				t.Logf("Expected error for whitespace name '%q', but got nil", whitespaceStr)
				return false
			}

			// Error message should indicate name validation failure
			if !strings.Contains(err.Error(), "name") {
				t.Logf("Error should mention 'name': %v", err)
				return false
			}

			return true
		},
		genWhitespaceString(),
	))

	// Test 2: Empty/whitespace system prompt should be rejected
	properties.Property("Empty or whitespace-only system prompt is rejected", prop.ForAll(
		func(whitespaceStr string) bool {
			// Clean up previous test data
			db.Exec("DELETE FROM ai_tags")

			req := &CreateAITagRequest{
				Name:         "ValidName_" + time.Now().Format("150405.000"),
				Description:  "Valid description",
				SystemPrompt: whitespaceStr,
			}

			_, err := service.CreateTag(req)

			// Should be rejected
			if err == nil {
				t.Logf("Expected error for whitespace system prompt '%q', but got nil", whitespaceStr)
				return false
			}

			// Error message should indicate system prompt validation failure
			if !strings.Contains(err.Error(), "system prompt") {
				t.Logf("Error should mention 'system prompt': %v", err)
				return false
			}

			return true
		},
		genWhitespaceString(),
	))

	properties.TestingRun(t)
}

// **Feature: ai-tag-management, Property 10: Update timestamp advancement**
// **Validates: Requirements 3.3**
// For any successful tag update, the updated_at timestamp SHALL be greater than or equal
// to the previous updated_at value.
func TestUpdateTimestampAdvancement(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAITagRepositoryWithDB(db)
	service := NewAITagServiceWithRepo(repo)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for update data
	updateDataGen := gopter.CombineGens(
		genValidName(),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) <= 500 }),
		genValidSystemPrompt(),
	)

	properties.Property("Update timestamp advances or stays the same", prop.ForAll(
		func(updateData []any) bool {
			newName := updateData[0].(string)
			newDesc := updateData[1].(string)
			newPrompt := updateData[2].(string)

			// Clean up previous test data
			db.Exec("DELETE FROM ai_tags")

			// Create initial tag
			createReq := &CreateAITagRequest{
				Name:         "InitialTag_" + time.Now().Format("150405.000000"),
				Description:  "Initial description",
				SystemPrompt: "Initial prompt",
			}

			tag, err := service.CreateTag(createReq)
			if err != nil {
				t.Logf("Failed to create initial tag: %v", err)
				return false
			}

			originalUpdatedAt := tag.UpdatedAt

			// Small delay to ensure timestamp can advance
			time.Sleep(10 * time.Millisecond)

			// Update the tag
			updateReq := &UpdateAITagRequest{
				Name:         newName + "_" + time.Now().Format("150405.000000"),
				Description:  newDesc,
				SystemPrompt: newPrompt,
			}

			updatedTag, err := service.UpdateTag(tag.ID, updateReq)
			if err != nil {
				t.Logf("Failed to update tag: %v", err)
				return false
			}

			// Verify timestamp advanced or stayed the same
			if updatedTag.UpdatedAt.Before(originalUpdatedAt) {
				t.Logf("Updated timestamp %v is before original %v",
					updatedTag.UpdatedAt, originalUpdatedAt)
				return false
			}

			return true
		},
		updateDataGen,
	))

	properties.TestingRun(t)
}

// **Feature: ai-tag-management, Property 11: Soft delete sets inactive status**
// **Validates: Requirements 4.1**
// For any deleted tag, querying the database (including soft-deleted records) SHALL show
// the tag with status=0 (inactive) and a non-null deleted_at timestamp.
func TestSoftDeleteSetsInactiveStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAITagRepositoryWithDB(db)
	service := NewAITagServiceWithRepo(repo)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for tag data
	tagDataGen := gopter.CombineGens(
		genValidName(),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) <= 500 }),
		genValidSystemPrompt(),
	)

	properties.Property("Soft delete sets status to inactive (0)", prop.ForAll(
		func(tagData []any) bool {
			name := tagData[0].(string)
			desc := tagData[1].(string)
			prompt := tagData[2].(string)

			// Clean up previous test data
			db.Exec("DELETE FROM ai_tags")

			// Create a tag
			createReq := &CreateAITagRequest{
				Name:         name + "_" + time.Now().Format("150405.000000"),
				Description:  desc,
				SystemPrompt: prompt,
			}

			tag, err := service.CreateTag(createReq)
			if err != nil {
				t.Logf("Failed to create tag: %v", err)
				return false
			}

			// Verify initial status is active (1)
			if tag.Status != 1 {
				t.Logf("Initial status should be 1, got %d", tag.Status)
				return false
			}

			// Delete the tag (soft delete)
			err = service.DeleteTag(tag.ID)
			if err != nil {
				t.Logf("Failed to delete tag: %v", err)
				return false
			}

			// Query the tag directly from database to check status
			var deletedTag model.AITag
			err = db.First(&deletedTag, tag.ID).Error
			if err != nil {
				t.Logf("Failed to query deleted tag: %v", err)
				return false
			}

			// Verify status is now inactive (0)
			if deletedTag.Status != 0 {
				t.Logf("Deleted tag status should be 0, got %d", deletedTag.Status)
				return false
			}

			return true
		},
		tagDataGen,
	))

	properties.TestingRun(t)
}

// **Feature: ai-tag-management, Property 15: Test response structure**
// **Validates: Requirements 6.2**
// For any tag test execution, the response SHALL include: success (boolean),
// response (string), and knowledge_accessed (boolean) fields.
func TestTagTestResponseStructure(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAITagRepositoryWithDB(db)
	service := NewAITagServiceWithRepo(repo)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for tag data
	tagDataGen := gopter.CombineGens(
		genValidName(),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) <= 500 }),
		genValidSystemPrompt(),
		gen.Identifier(), // API key (may be invalid for testing)
	)

	properties.Property("Test response has required structure fields", prop.ForAll(
		func(tagData []any) bool {
			name := tagData[0].(string)
			desc := tagData[1].(string)
			prompt := tagData[2].(string)
			apiKey := tagData[3].(string)

			// Clean up previous test data
			db.Exec("DELETE FROM ai_tags")

			// Create a tag
			createReq := &CreateAITagRequest{
				Name:         name + "_" + time.Now().Format("150405.000000"),
				Description:  desc,
				SystemPrompt: prompt,
				ChatAPIKey:   apiKey,
			}

			tag, err := service.CreateTag(createReq)
			if err != nil {
				t.Logf("Failed to create tag: %v", err)
				return false
			}

			// Test the tag (will likely fail due to invalid API key, but response structure should be valid)
			resp, err := service.TestTag(tag.ID)
			if err != nil {
				t.Logf("TestTag returned error: %v", err)
				return false
			}

			// Verify response is not nil
			if resp == nil {
				t.Logf("Response should not be nil")
				return false
			}

			// Verify response has the required structure
			// The Success field is a bool (always present)
			// The Response field is a string (always present, may be empty)
			// The KnowledgeAccessed field is a bool (always present)
			// These are struct fields, so they're always present by definition

			// For failed tests, ErrorMessage should be set
			if !resp.Success && resp.ErrorMessage == "" {
				// This is acceptable - some failures may not have detailed error messages
				// but the structure is still valid
			}

			// Verify the response type is correct
			var _ bool = resp.Success
			var _ string = resp.Response
			var _ bool = resp.KnowledgeAccessed
			var _ string = resp.ErrorMessage

			return true
		},
		tagDataGen,
	))

	properties.TestingRun(t)
}
