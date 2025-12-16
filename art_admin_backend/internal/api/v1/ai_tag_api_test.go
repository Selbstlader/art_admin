package v1

import (
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	aiTagSvc "art_admin_backend/internal/service/ai_tag"
	"reflect"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

/*** Property-Based Tests for AI Tag API ***/

// genValidAITagData generates valid AI tag data for testing
func genValidAITagData() gopter.Gen {
	return gopter.CombineGens(
		gen.Int64Range(1, 1000000), // ID
		gen.Identifier().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }), // Name
		gen.AlphaString().Map(func(s string) string {
			if len(s) > 500 {
				return s[:500]
			}
			return s
		}), // Description
		gen.Identifier().Map(func(s string) string {
			if len(s) > 100 {
				return s[:100]
			}
			return s
		}), // KnowledgeBaseID
		gen.AlphaString().Map(func(s string) string {
			if len(s) > 200 {
				return s[:200]
			}
			return s
		}), // KnowledgeBaseName
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }), // SystemPrompt
		gen.Identifier().Map(func(s string) string {
			if len(s) > 200 {
				return s[:200]
			}
			return s
		}), // ChatAPIKey
		gen.IntRange(0, 1), // Status
	)
}

// **Feature: ai-tag-management, Property 7: Response structure completeness**
// *For any* tag in a list response, the tag object SHALL contain non-null values for:
// id, name, status, created_at, and system_prompt.
// **Validates: Requirements 2.3**
func TestProperty7_ResponseStructureCompleteness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(time.Now().UnixNano())

	properties := gopter.NewProperties(parameters)

	properties.Property("Response structure contains all required non-null fields", prop.ForAll(
		func(data []any) bool {
			id := data[0].(int64)
			name := data[1].(string)
			desc := data[2].(string)
			kbID := data[3].(string)
			kbName := data[4].(string)
			prompt := data[5].(string)
			apiKey := data[6].(string)
			status := data[7].(int)

			// Create model
			tag := &model.AITag{
				ID:                id,
				Name:              name,
				Description:       desc,
				KnowledgeBaseID:   kbID,
				KnowledgeBaseName: kbName,
				SystemPrompt:      prompt,
				ChatAPIKey:        apiKey,
				Status:            status,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}

			// Convert model to response using the helper function
			resp := convertToAITagResponse(tag)

			// Verify all required fields are present and non-null/non-zero
			// id must be non-zero
			if resp.ID == 0 {
				t.Logf("ID is zero")
				return false
			}

			// name must be non-empty
			if resp.Name == "" {
				t.Logf("Name is empty")
				return false
			}

			// system_prompt must be non-empty
			if resp.SystemPrompt == "" {
				t.Logf("SystemPrompt is empty")
				return false
			}

			// status must be valid (0 or 1)
			if resp.Status != 0 && resp.Status != 1 {
				t.Logf("Status is invalid: %d", resp.Status)
				return false
			}

			// created_at must be non-zero
			if resp.CreatedAt.IsZero() {
				t.Logf("CreatedAt is zero")
				return false
			}

			return true
		},
		genValidAITagData(),
	))

	properties.TestingRun(t)
}

// Test that list response conversion preserves all fields
func TestProperty7_ListResponseStructureCompleteness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(time.Now().UnixNano())

	properties := gopter.NewProperties(parameters)

	// Generator for list size
	listSizeGen := gen.IntRange(1, 10)

	properties.Property("List response contains all required fields for each tag", prop.ForAll(
		func(listSize int) bool {
			// Create list of tags
			modelList := make([]*model.AITag, listSize)
			for i := 0; i < listSize; i++ {
				modelList[i] = &model.AITag{
					ID:                int64(i + 1),
					Name:              "TestTag_" + time.Now().Format("150405.000000"),
					Description:       "Test description",
					KnowledgeBaseID:   "kb_123",
					KnowledgeBaseName: "Test KB",
					SystemPrompt:      "Test prompt",
					ChatAPIKey:        "api_key_123",
					Status:            1,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				}
			}

			svcResp := &aiTagSvc.AITagListResponse{
				List:  modelList,
				Total: int64(listSize),
				Page:  1,
				Size:  10,
			}

			// Convert to DTO response
			resp := convertToAITagListResponse(svcResp)

			// Verify list length matches
			if len(resp.List) != listSize {
				t.Logf("List length mismatch: expected %d, got %d", listSize, len(resp.List))
				return false
			}

			// Verify pagination fields
			if resp.Total != int64(listSize) {
				t.Logf("Total mismatch: expected %d, got %d", listSize, resp.Total)
				return false
			}
			if resp.Page != 1 {
				t.Logf("Page mismatch: expected 1, got %d", resp.Page)
				return false
			}
			if resp.Size != 10 {
				t.Logf("Size mismatch: expected 10, got %d", resp.Size)
				return false
			}

			// Verify each tag in the list has required fields
			for i, tagResp := range resp.List {
				if tagResp.ID == 0 {
					t.Logf("Tag %d: ID is zero", i)
					return false
				}
				if tagResp.Name == "" {
					t.Logf("Tag %d: Name is empty", i)
					return false
				}
				if tagResp.SystemPrompt == "" {
					t.Logf("Tag %d: SystemPrompt is empty", i)
					return false
				}
				if tagResp.Status != 0 && tagResp.Status != 1 {
					t.Logf("Tag %d: Status is invalid: %d", i, tagResp.Status)
					return false
				}
				if tagResp.CreatedAt.IsZero() {
					t.Logf("Tag %d: CreatedAt is zero", i)
					return false
				}
			}

			return true
		},
		listSizeGen,
	))

	properties.TestingRun(t)
}

// Test that test response has required structure
// **Feature: ai-tag-management, Property 15: Test response structure**
// (Additional validation at API layer for response DTO)
func TestProperty15_TestResponseStructureAtAPILayer(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(time.Now().UnixNano())

	properties := gopter.NewProperties(parameters)

	// Generator for test response data
	testRespDataGen := gopter.CombineGens(
		gen.Bool(),        // Success
		gen.AlphaString(), // Response
		gen.Bool(),        // KnowledgeAccessed
		gen.AlphaString(), // ErrorMessage
	)

	properties.Property("Test response has all required fields", prop.ForAll(
		func(data []any) bool {
			success := data[0].(bool)
			respStr := data[1].(string)
			kbAccessed := data[2].(bool)
			errMsg := data[3].(string)

			resp := &response.AITagTestResponse{
				Success:           success,
				Response:          respStr,
				KnowledgeAccessed: kbAccessed,
				ErrorMessage:      errMsg,
			}

			// Verify the structure is correct by accessing fields
			// These are compile-time checks essentially
			_ = resp.Success
			_ = resp.Response
			_ = resp.KnowledgeAccessed
			_ = resp.ErrorMessage

			// Verify types using reflection
			respType := reflect.TypeOf(*resp)

			// Check Success field exists and is bool
			successField, ok := respType.FieldByName("Success")
			if !ok || successField.Type.Kind() != reflect.Bool {
				t.Logf("Success field missing or wrong type")
				return false
			}

			// Check Response field exists and is string
			responseField, ok := respType.FieldByName("Response")
			if !ok || responseField.Type.Kind() != reflect.String {
				t.Logf("Response field missing or wrong type")
				return false
			}

			// Check KnowledgeAccessed field exists and is bool
			kaField, ok := respType.FieldByName("KnowledgeAccessed")
			if !ok || kaField.Type.Kind() != reflect.Bool {
				t.Logf("KnowledgeAccessed field missing or wrong type")
				return false
			}

			return true
		},
		testRespDataGen,
	))

	properties.TestingRun(t)
}
