package swagno3

import (
	"encoding/json"
	"testing"

	"github.com/go-swagno/swagno/v3/components/definition"
	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/extensions"
	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/mime"
	"github.com/go-swagno/swagno/v3/components/parameter"
	"github.com/go-swagno/swagno/v3/components/security"
	"github.com/go-swagno/swagno/v3/components/tag"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type TestExtPayload struct {
	ID   uint64 `json:"id" example:"1"`
	Name string `json:"name" example:"Alice"`
}

func TestExtensionsCombined(t *testing.T) {
	cmpOpts := []cmp.Option{
		cmpopts.IgnoreUnexported(OpenAPI{}, endpoint.EndPoint{}, endpoint.JsonEndPoint{}),
		cmpopts.EquateEmpty(),
		cmpopts.IgnoreFields(definition.SchemaProperty{}, "IsRequired"),
		cmpopts.SortSlices(func(a, b string) bool { return a < b }),
	}

	t.Run("Operations", func(t *testing.T) {
		testEndpoint := endpoint.New(
			endpoint.GET,
			"/product",
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(TestExtPayload{}, "200", "OK"),
			}),
			endpoint.WithExtension("x-rate-limit", 100),
			endpoint.WithExtensions(extensions.Extensions{
				"x-internal-id": "product.list",
			}),
		)

		openapi := New(Config{Title: "Ext API", Version: "1.0.0"})
		openapi.AddEndpoint(testEndpoint)
		openapi.generateOpenAPIJson()

		expected := &OpenAPI{
			OpenAPI: "3.0.3",
			Info:    Info{Title: "Ext API", Version: "1.0.0"},
			Servers: []Server{{URL: "/"}},
			Paths: map[string]endpoint.PathItem{
				"/product": {
					Get: &endpoint.JsonEndPoint{
						OperationId: "get-_product",
						Consume:     []mime.MIME{mime.JSON},
						Produce:     []mime.MIME{mime.JSON},
						Responses: map[string]endpoint.JsonResponse{
							"200": {
								Description: "OK",
								Content: map[string]endpoint.MediaType{
									"application/json": {
										Schema: &parameter.JsonResponseSchema{
											Ref: "#/components/schemas/swagno3.TestExtPayload",
										},
									},
								},
							},
						},
						Extensions: extensions.Extensions{
							"x-rate-limit":  100,
							"x-internal-id": "product.list",
						},
					},
				},
			},
			Components: &Components{
				Schemas: map[string]definition.Schema{
					"swagno3.TestExtPayload": {
						Type: "object",
						Properties: map[string]definition.SchemaProperty{
							"id":   {Type: "integer", Example: float64(1)},
							"name": {Type: "string", Example: "Alice"},
						},
						Required: []string{"id", "name"},
					},
				},
				SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme),
			},
			Tags: []tag.Tag{},
		}

		if diff := cmp.Diff(expected, openapi, cmpOpts...); diff != "" {
			t.Errorf("OpenAPI mismatch (-expected +got):\n%s", diff)
		}
	})

	t.Run("DocumentAndInfo", func(t *testing.T) {
		openapi := New(Config{
			Title:          "Ext API",
			Version:        "1.0.0",
			Extensions:     extensions.Extensions{"x-audience": "internal"},
			InfoExtensions: extensions.Extensions{"x-logo": map[string]string{"url": "https://example.com/logo.png"}},
		})
		openapi.generateOpenAPIJson()

		expected := &OpenAPI{
			OpenAPI: "3.0.3",
			Info: Info{
				Title:      "Ext API",
				Version:    "1.0.0",
				Extensions: extensions.Extensions{"x-logo": map[string]string{"url": "https://example.com/logo.png"}},
			},
			Servers:    []Server{{URL: "/"}},
			Paths:      map[string]endpoint.PathItem{},
			Extensions: extensions.Extensions{"x-audience": "internal"},
			Components: &Components{
				Schemas:         map[string]definition.Schema{},
				SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme),
			},
			Tags: []tag.Tag{},
		}

		if diff := cmp.Diff(expected, openapi, cmpOpts...); diff != "" {
			t.Errorf("OpenAPI mismatch (-expected +got):\n%s", diff)
		}
	})

	t.Run("TagsAndServers", func(t *testing.T) {
		openapi := New(Config{Title: "Ext API", Version: "1.0.0"})
		openapi.AddServer("https://api.example.com/v1", "Production")
		openapi.Servers[0].Extensions = extensions.Extensions{"x-env": "prod"}

		openapi.AddTags(tag.New("users", "User endpoints",
			tag.WithExternalDocs("https://docs.example.com/users", "User docs"),
		))
		openapi.Tags[0].Extensions = extensions.Extensions{"x-display-name": "Users"}
		openapi.Tags[0].ExternalDocs.Extensions = extensions.Extensions{"x-last-reviewed": "2026-04-01"}

		openapi.generateOpenAPIJson()

		expected := &OpenAPI{
			OpenAPI: "3.0.3",
			Info:    Info{Title: "Ext API", Version: "1.0.0"},
			Servers: []Server{{
				URL:         "https://api.example.com/v1",
				Description: "Production",
				Extensions:  extensions.Extensions{"x-env": "prod"},
			}},
			Paths: map[string]endpoint.PathItem{},
			Components: &Components{
				Schemas:         map[string]definition.Schema{},
				SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme),
			},
			Tags: []tag.Tag{{
				Name:        "users",
				Description: "User endpoints",
				ExternalDocs: &tag.ExternalDocs{
					URL:         "https://docs.example.com/users",
					Description: "User docs",
					Extensions:  extensions.Extensions{"x-last-reviewed": "2026-04-01"},
				},
				Extensions: extensions.Extensions{"x-display-name": "Users"},
			}},
		}

		if diff := cmp.Diff(expected, openapi, cmpOpts...); diff != "" {
			t.Errorf("OpenAPI mismatch (-expected +got):\n%s", diff)
		}
	})

	t.Run("SchemasAndParameters", func(t *testing.T) {
		testEndpoint := endpoint.New(
			endpoint.GET,
			"/users/{id}",
			endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(TestExtPayload{}, "200", "OK"),
			}),
		)

		openapi := New(Config{Title: "Ext API", Version: "1.0.0"})
		openapi.AddEndpoint(testEndpoint)
		openapi.generateOpenAPIJson()

		// Post-generation mutations: users set extensions on auto-generated
		// schemas and parameters via the exposed structs.
		schemaKey := "swagno3.TestExtPayload"
		s := openapi.Components.Schemas[schemaKey]
		s.Extensions = extensions.Extensions{"x-internal-id": "payload-v1"}
		openapi.Components.Schemas[schemaKey] = s

		pi := openapi.Paths["/users/{id}"]
		pi.Get.Parameters[0].Extensions = extensions.Extensions{"x-source": "auth"}
		openapi.Paths["/users/{id}"] = pi

		expected := &OpenAPI{
			OpenAPI: "3.0.3",
			Info:    Info{Title: "Ext API", Version: "1.0.0"},
			Servers: []Server{{URL: "/"}},
			Paths: map[string]endpoint.PathItem{
				"/users/{id}": {
					Get: &endpoint.JsonEndPoint{
						OperationId: "get-_users_id",
						Consume:     []mime.MIME{mime.JSON},
						Produce:     []mime.MIME{mime.JSON},
						Parameters: []parameter.JsonParameter{{
							Name:     "id",
							In:       "path",
							Required: true,
							Schema:   &parameter.JsonResponseSchema{Type: "integer"},
							Extensions: extensions.Extensions{"x-source": "auth"},
						}},
						Responses: map[string]endpoint.JsonResponse{
							"200": {
								Description: "OK",
								Content: map[string]endpoint.MediaType{
									"application/json": {
										Schema: &parameter.JsonResponseSchema{
											Ref: "#/components/schemas/swagno3.TestExtPayload",
										},
									},
								},
							},
						},
					},
				},
			},
			Components: &Components{
				Schemas: map[string]definition.Schema{
					"swagno3.TestExtPayload": {
						Type: "object",
						Properties: map[string]definition.SchemaProperty{
							"id":   {Type: "integer", Example: float64(1)},
							"name": {Type: "string", Example: "Alice"},
						},
						Required:   []string{"id", "name"},
						Extensions: extensions.Extensions{"x-internal-id": "payload-v1"},
					},
				},
				SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme),
			},
			Tags: []tag.Tag{},
		}

		if diff := cmp.Diff(expected, openapi, cmpOpts...); diff != "" {
			t.Errorf("OpenAPI mismatch (-expected +got):\n%s", diff)
		}
	})
}

// TestMergeRejectsNonExtensionKeys is a unit test for the shared helper that
// filters keys not prefixed with "x-" at serialization time.
func TestMergeRejectsNonExtensionKeys(t *testing.T) {
	out, err := extensions.Merge(struct {
		A string `json:"a"`
	}{A: "1"}, extensions.Extensions{
		"x-ok": "yes",
		"bad":  "nope",
	})
	if err != nil {
		t.Fatalf("merge failed: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if m["x-ok"] != "yes" {
		t.Errorf("x-ok missing")
	}
	if _, ok := m["bad"]; ok {
		t.Errorf("non-extension key leaked: %v", m)
	}
}
