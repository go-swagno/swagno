package v3

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/parameter"
)

type TestUser struct {
	ID       uint64  `json:"id" example:"1"`
	Name     string  `json:"name" example:"John Doe"`
	Email    *string `json:"email,omitempty" example:"john@example.com"`
	IsActive bool    `json:"is_active" example:"true"`
}

type TestError struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Something went wrong"`
}

func TestOpenAPIGeneration(t *testing.T) {
	// Create OpenAPI instance
	openapi := New(Config{
		Title:       "Test API",
		Version:     "v1.0.0",
		Description: "Test API for OpenAPI 3.0",
	})

	// Add server
	openapi.AddServer("https://api.test.com/v1", "Test server")

	// Add security
	openapi.SetBearerAuth("JWT", "JWT Bearer authentication")

	// Create test endpoint
	testEndpoint := endpoint.New(
		endpoint.GET,
		"/users/{id}",
		endpoint.WithParams(
			parameter.IntParam("id", parameter.Path,
				parameter.WithRequired(),
				parameter.WithDescription("User ID"),
			),
		),
		endpoint.WithSuccessfulReturns([]response.Response{
			response.New(TestUser{}, "200", "User found"),
		}),
		endpoint.WithErrors([]response.Response{
			response.New(TestError{}, "404", "User not found"),
		}),
	)

	openapi.AddEndpoint(testEndpoint)

	// Generate JSON
	jsonBytes, err := openapi.ToJson()
	if err != nil {
		t.Fatalf("Failed to generate JSON: %v", err)
	}

	// Parse JSON to validate structure
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	// Validate OpenAPI version
	if result["openapi"] != "3.0.3" {
		t.Errorf("Expected openapi version 3.0.3, got %v", result["openapi"])
	}

	// Validate info section
	info, ok := result["info"].(map[string]interface{})
	if !ok {
		t.Fatal("Info section is missing or invalid")
	}
	if info["title"] != "Test API" {
		t.Errorf("Expected title 'Test API', got %v", info["title"])
	}
	if info["version"] != "v1.0.0" {
		t.Errorf("Expected version 'v1.0.0', got %v", info["version"])
	}

	// Validate servers
	servers, ok := result["servers"].([]interface{})
	if !ok || len(servers) == 0 {
		t.Fatal("Servers section is missing or empty")
	}
	firstServer := servers[0].(map[string]interface{})
	if firstServer["url"] != "https://api.test.com/v1" {
		t.Errorf("Expected server URL 'https://api.test.com/v1', got %v", firstServer["url"])
	}

	// Validate components
	components, ok := result["components"].(map[string]interface{})
	if !ok {
		t.Fatal("Components section is missing")
	}

	// Validate schemas (should contain TestUser and TestError)
	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		t.Fatal("Schemas section is missing")
	}
	if _, exists := schemas["v3.TestUser"]; !exists {
		t.Error("TestUser schema is missing")
	}
	if _, exists := schemas["v3.TestError"]; !exists {
		t.Error("TestError schema is missing")
	}

	// Validate security schemes
	securitySchemes, ok := components["securitySchemes"].(map[string]interface{})
	if !ok {
		t.Fatal("Security schemes section is missing")
	}
	bearerAuth, exists := securitySchemes["bearerAuth"].(map[string]interface{})
	if !exists {
		t.Fatal("Bearer auth security scheme is missing")
	}
	if bearerAuth["type"] != "http" || bearerAuth["scheme"] != "bearer" {
		t.Error("Bearer auth configuration is incorrect")
	}

	// Validate paths
	paths, ok := result["paths"].(map[string]interface{})
	if !ok || len(paths) == 0 {
		t.Fatal("Paths section is missing or empty")
	}
	userPath, exists := paths["/users/{id}"].(map[string]interface{})
	if !exists {
		t.Fatal("User path is missing")
	}
	getMethod, exists := userPath["get"].(map[string]interface{})
	if !exists {
		t.Fatal("GET method is missing")
	}

	// Validate responses structure
	responses, ok := getMethod["responses"].(map[string]interface{})
	if !ok {
		t.Fatal("Responses section is missing")
	}
	response200, exists := responses["200"].(map[string]interface{})
	if !exists {
		t.Fatal("200 response is missing")
	}

	// Validate response content structure (OpenAPI 3.0)
	content, ok := response200["content"].(map[string]interface{})
	if !ok {
		t.Fatal("Response content section is missing")
	}
	jsonContent, exists := content["application/json"].(map[string]interface{})
	if !exists {
		t.Fatal("JSON content type is missing")
	}
	schema, exists := jsonContent["schema"].(map[string]interface{})
	if !exists {
		t.Fatal("Response schema is missing")
	}
	schemaRef, ok := schema["$ref"].(string)
	if !ok || !strings.Contains(schemaRef, "#/components/schemas/") {
		t.Error("Schema reference should point to components/schemas")
	}
}

func TestMustToJson(t *testing.T) {
	openapi := New(Config{
		Title:   "Test API",
		Version: "v1.0.0",
	})

	// Should not panic
	jsonBytes := openapi.MustToJson()
	if len(jsonBytes) == 0 {
		t.Error("MustToJson returned empty result")
	}

	// Validate it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Errorf("MustToJson generated invalid JSON: %v", err)
	}
}

func TestSecurityMethods(t *testing.T) {
	openapi := New(Config{
		Title:   "Test API",
		Version: "v1.0.0",
	})

	// Test different security methods
	openapi.SetBasicAuth("Basic auth")
	openapi.SetBearerAuth("JWT", "Bearer auth")
	openapi.SetApiKeyAuth("X-API-Key", "header", "API key auth")

	// Test OAuth2
	flows := &OAuthFlows{
		AuthorizationCode: &OAuthFlow{
			AuthorizationUrl: "https://example.com/oauth/authorize",
			TokenUrl:         "https://example.com/oauth/token",
			Scopes: map[string]string{
				"read": "Read access",
			},
		},
	}
	openapi.SetOAuth2Auth(flows, "OAuth2 auth")

	// Test OpenID Connect
	openapi.SetOpenIdConnectAuth("https://example.com/.well-known/openid_configuration", "OpenID Connect")

	jsonBytes, err := openapi.ToJson()
	if err != nil {
		t.Fatalf("Failed to generate JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	components := result["components"].(map[string]interface{})
	securitySchemes := components["securitySchemes"].(map[string]interface{})

	// Validate all security schemes are present
	expectedSchemes := []string{"basicAuth", "bearerAuth", "apiKeyAuth", "oauth2", "openIdConnect"}
	for _, scheme := range expectedSchemes {
		if _, exists := securitySchemes[scheme]; !exists {
			t.Errorf("Security scheme %s is missing", scheme)
		}
	}
}

func TestMultipleServers(t *testing.T) {
	openapi := New(Config{
		Title:   "Test API",
		Version: "v1.0.0",
	})

	openapi.AddServer("https://api.example.com/v1", "Production")
	openapi.AddServer("https://staging-api.example.com/v1", "Staging")
	openapi.AddServer("http://localhost:8080/v1", "Development")

	jsonBytes, err := openapi.ToJson()
	if err != nil {
		t.Fatalf("Failed to generate JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	servers, ok := result["servers"].([]interface{})
	if !ok {
		t.Fatal("Servers section is missing")
	}

	if len(servers) != 3 {
		t.Errorf("Expected 3 servers, got %d", len(servers))
	}

	// Validate first server
	firstServer := servers[0].(map[string]interface{})
	if firstServer["url"] != "https://api.example.com/v1" {
		t.Errorf("First server URL is incorrect")
	}
	if firstServer["description"] != "Production" {
		t.Errorf("First server description is incorrect")
	}
}
