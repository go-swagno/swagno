package swagno3

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/parameter"
	"github.com/go-swagno/swagno/v3/components/security"
	"github.com/go-swagno/swagno/v3/example/models"
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
	if _, exists := schemas["swagno3.TestUser"]; !exists {
		t.Error("TestUser schema is missing")
	}
	if _, exists := schemas["swagno3.TestError"]; !exists {
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
	flows := &security.OAuthFlows{
		AuthorizationCode: &security.OAuthFlow{
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

// TestHidePackageNameV3 verifies that enabling Config.HidePackageName strips the
// package qualifier from schema keys and $ref values (e.g. "swagno3.TestUser"
// becomes "TestUser"), while the default keeps the package-qualified names.
func TestHidePackageNameV3(t *testing.T) {
	buildEndpoint := func() *endpoint.EndPoint {
		return endpoint.New(
			endpoint.POST,
			"/users",
			endpoint.WithBody(TestUser{}),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(TestUser{}, "201", "User created"),
			}),
			endpoint.WithErrors([]response.Response{
				response.New([]TestError{}, "400", "Bad Request"),
			}),
		)
	}

	t.Run("hidden", func(t *testing.T) {
		openapi := New(Config{Title: "Test API", Version: "v1.0.0", HidePackageName: true})
		openapi.AddEndpoint(buildEndpoint())
		doc := string(openapi.MustToJson())

		if strings.Contains(doc, "swagno3.") {
			t.Errorf("expected no package qualifier in output, but found \"swagno3.\":\n%s", doc)
		}
		for _, want := range []string{
			`"#/components/schemas/TestUser"`,
			`"#/components/schemas/TestError"`,
		} {
			if !strings.Contains(doc, want) {
				t.Errorf("expected output to contain %s, but it did not:\n%s", want, doc)
			}
		}

		// schema keys must also be stripped so the refs resolve
		var parsed struct {
			Components struct {
				Schemas map[string]json.RawMessage `json:"schemas"`
			} `json:"components"`
		}
		if err := json.Unmarshal([]byte(doc), &parsed); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"TestUser", "TestError"} {
			if _, ok := parsed.Components.Schemas[key]; !ok {
				keys := make([]string, 0, len(parsed.Components.Schemas))
				for k := range parsed.Components.Schemas {
					keys = append(keys, k)
				}
				t.Errorf("expected schemas to contain key %q, got keys: %v", key, keys)
			}
		}
	})

	t.Run("default keeps package name", func(t *testing.T) {
		openapi := New(Config{Title: "Test API", Version: "v1.0.0"})
		openapi.AddEndpoint(buildEndpoint())
		doc := string(openapi.MustToJson())

		if !strings.Contains(doc, `"#/components/schemas/swagno3.TestUser"`) {
			t.Errorf("expected default output to keep package qualifier \"swagno3.TestUser\":\n%s", doc)
		}
	})
}

// SuccessfulResponse is a deliberate name twin of models.SuccessfulResponse, declared
// in this (swagno3) package, so that enabling HidePackageName makes both strip to the
// same "SuccessfulResponse" name and triggers a collision.
type SuccessfulResponse struct {
	Status string `json:"status"`
}

// TestHidePackageNameCollisionV3 verifies that when HidePackageName strips two distinct
// types from different packages to the same name, generation fails with a
// *NameCollisionError instead of silently dropping one of the schemas.
func TestHidePackageNameCollisionV3(t *testing.T) {
	buildColliding := func() []*endpoint.EndPoint {
		return []*endpoint.EndPoint{
			endpoint.New(
				endpoint.GET,
				"/a",
				endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessfulResponse{}, "200", "OK")}),
			),
			endpoint.New(
				endpoint.GET,
				"/b",
				endpoint.WithSuccessfulReturns([]response.Response{response.New(SuccessfulResponse{}, "200", "OK")}),
			),
		}
	}

	t.Run("ToJson returns NameCollisionError", func(t *testing.T) {
		openapi := New(Config{Title: "Test API", Version: "v1.0.0", HidePackageName: true})
		openapi.AddEndpoints(buildColliding())

		_, err := openapi.ToJson()
		if err == nil {
			t.Fatal("expected a collision error, got nil")
		}
		var collisionErr *NameCollisionError
		if !errors.As(err, &collisionErr) {
			t.Fatalf("expected *NameCollisionError, got %T: %v", err, err)
		}
		if _, ok := collisionErr.Collisions["SuccessfulResponse"]; !ok {
			t.Errorf("expected collision on %q, got %v", "SuccessfulResponse", collisionErr.Collisions)
		}
	})

	t.Run("MustToJson panics", func(t *testing.T) {
		openapi := New(Config{Title: "Test API", Version: "v1.0.0", HidePackageName: true})
		openapi.AddEndpoints(buildColliding())

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected MustToJson to panic on collision, but it did not")
			}
		}()
		openapi.MustToJson()
	})

	t.Run("no collision when HidePackageName is disabled", func(t *testing.T) {
		openapi := New(Config{Title: "Test API", Version: "v1.0.0"})
		openapi.AddEndpoints(buildColliding())

		if _, err := openapi.ToJson(); err != nil {
			t.Fatalf("expected no error with HidePackageName disabled, got %v", err)
		}
	})
}
