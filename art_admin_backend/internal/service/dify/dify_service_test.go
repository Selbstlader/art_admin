package dify

import (
	"art_admin_backend/internal/dto/request"
	difyPkg "art_admin_backend/internal/pkg/dify"
	"reflect"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func genKnowledgeBaseID() gopter.Gen {
	return gen.Identifier().SuchThat(func(s string) bool {
		return len(s) > 0 && len(s) <= 100
	})
}

func genSystemPrompt() gopter.Gen {
	return gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0 && strings.TrimSpace(s) != ""
	})
}

// **Feature: ai-tag-management, Property 12: Chat uses tag's knowledge base ID**
// **Validates: Requirements 5.1**
func TestChatUsesTagKnowledgeBaseID(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	testDataGen := gopter.CombineGens(
		genKnowledgeBaseID(),
		genSystemPrompt(),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
		gen.Identifier(),
	)

	properties.Property("buildChatRequestData includes knowledge_base_id in inputs", prop.ForAll(
		func(testData []any) bool {
			knowledgeBaseID := testData[0].(string)
			systemPrompt := testData[1].(string)
			query := testData[2].(string)
			user := testData[3].(string)

			service := &DifyService{
				baseURL: "http://test.example.com",
				timeout: 30,
			}

			req := &request.DifyChatRequest{
				Query: query,
				User:  user,
			}

			data := service.buildChatRequestData(req, systemPrompt, knowledgeBaseID)

			inputs, ok := data["inputs"].(map[string]interface{})
			if !ok {
				t.Logf("inputs should be a map, got %T", data["inputs"])
				return false
			}

			kbID, ok := inputs["knowledge_base_id"].(string)
			if !ok {
				t.Logf("knowledge_base_id should be a string in inputs")
				return false
			}

			if kbID != knowledgeBaseID {
				t.Logf("knowledge_base_id mismatch: expected %s, got %s", knowledgeBaseID, kbID)
				return false
			}

			return true
		},
		testDataGen,
	))

	properties.TestingRun(t)
}

// **Feature: ai-tag-management, Property 13: Chat prepends system prompt**
// **Validates: Requirements 5.2**
func TestChatPrependsSystemPrompt(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	testDataGen := gopter.CombineGens(
		genSystemPrompt(),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 }),
		gen.Identifier(),
	)

	properties.Property("buildChatRequestData includes system_prompt in inputs", prop.ForAll(
		func(testData []any) bool {
			systemPrompt := testData[0].(string)
			query := testData[1].(string)
			user := testData[2].(string)

			service := &DifyService{
				baseURL: "http://test.example.com",
				timeout: 30,
			}

			req := &request.DifyChatRequest{
				Query: query,
				User:  user,
			}

			data := service.buildChatRequestData(req, systemPrompt, "")

			inputs, ok := data["inputs"].(map[string]interface{})
			if !ok {
				t.Logf("inputs should be a map, got %T", data["inputs"])
				return false
			}

			prompt, ok := inputs["system_prompt"].(string)
			if !ok {
				t.Logf("system_prompt should be a string in inputs")
				return false
			}

			if prompt != systemPrompt {
				t.Logf("system_prompt mismatch: expected %s, got %s", systemPrompt, prompt)
				return false
			}

			return true
		},
		testDataGen,
	))

	properties.TestingRun(t)
}

// **Feature: ai-tag-management, Property 14: Chat API key selection**
// **Validates: Requirements 5.3**
func TestChatAPIKeySelection(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	apiKeyGen := gen.OneGenOf(
		gen.Const(""),
		gen.Identifier().SuchThat(func(s string) bool { return len(s) > 0 }),
	)

	properties.Property("CreateClientWithAPIKey returns correct client based on API key", prop.ForAll(
		func(apiKey any) bool {
			apiKeyStr := apiKey.(string)

			defaultClient := difyPkg.NewClient("default-api-key", "http://test.example.com", 30)
			service := &DifyService{
				baseURL:    "http://test.example.com",
				timeout:    30,
				chatClient: defaultClient,
			}

			client := service.CreateClientWithAPIKey(apiKeyStr)

			if apiKeyStr == "" {
				if client != service.chatClient {
					t.Logf("Empty API key should return default chat client")
					return false
				}
			} else {
				if client == service.chatClient {
					t.Logf("Non-empty API key should return a new client, not the default")
					return false
				}
			}

			return true
		},
		apiKeyGen.WithLabel("apiKey"),
	))

	properties.TestingRun(t)
}

func init() {
	_ = reflect.TypeOf("")
}
