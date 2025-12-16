package model

import (
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: ai-tag-management, Property 16: Serialization round-trip consistency**
// **Validates: Requirements 7.1, 7.2**
// For any AITag object, serializing to JSON and deserializing back SHALL produce an object with identical field values.
func TestAITagJSONRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generator for AITag
	aiTagGen := genAITag()

	properties.Property("AITag JSON round-trip preserves data", prop.ForAll(
		func(original AITag) bool {
			// Serialize to JSON
			jsonData, err := original.ToJSON()
			if err != nil {
				t.Logf("Serialization error: %v", err)
				return false
			}

			// Deserialize from JSON
			var restored AITag
			err = restored.FromJSON(jsonData)
			if err != nil {
				t.Logf("Deserialization error: %v", err)
				return false
			}

			// Compare essential fields
			return compareAITags(original, restored)
		},
		aiTagGen,
	))

	properties.TestingRun(t)
}

// genAITag generates random AITag instances
func genAITag() gopter.Gen {
	return gopter.CombineGens(
		gen.Int64Range(1, 1000000), // ID
		gen.Identifier(),           // Name (always non-empty)
		gen.AlphaString(),          // Description
		gen.Identifier(),           // KnowledgeBaseID
		gen.AlphaString(),          // KnowledgeBaseName
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 500 }), // SystemPrompt (non-empty)
		gen.Identifier(),   // ChatAPIKey
		gen.IntRange(0, 1), // Status
	).Map(func(values []interface{}) AITag {
		now := time.Now().Truncate(time.Second)
		return AITag{
			ID:                values[0].(int64),
			Name:              values[1].(string),
			Description:       values[2].(string),
			KnowledgeBaseID:   values[3].(string),
			KnowledgeBaseName: values[4].(string),
			SystemPrompt:      values[5].(string),
			ChatAPIKey:        values[6].(string),
			Status:            values[7].(int),
			CreatedAt:         now,
			UpdatedAt:         now,
		}
	})
}

// compareAITags compares two AITag instances for equality
func compareAITags(a, b AITag) bool {
	if a.ID != b.ID {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Description != b.Description {
		return false
	}
	if a.KnowledgeBaseID != b.KnowledgeBaseID {
		return false
	}
	if a.KnowledgeBaseName != b.KnowledgeBaseName {
		return false
	}
	if a.SystemPrompt != b.SystemPrompt {
		return false
	}
	if a.ChatAPIKey != b.ChatAPIKey {
		return false
	}
	if a.Status != b.Status {
		return false
	}
	// Compare timestamps with second precision (JSON may lose nanoseconds)
	if !a.CreatedAt.Truncate(time.Second).Equal(b.CreatedAt.Truncate(time.Second)) {
		return false
	}
	if !a.UpdatedAt.Truncate(time.Second).Equal(b.UpdatedAt.Truncate(time.Second)) {
		return false
	}
	return true
}
