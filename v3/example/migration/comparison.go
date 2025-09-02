// Package migration demonstrates the differences between Swagger 2.0 and OpenAPI 3.0 usage
package migration

import (
	// Swagger 2.0 imports
	"github.com/go-swagno/swagno"
	swaggerEndpoint "github.com/go-swagno/swagno/components/endpoint"
	swaggerResponse "github.com/go-swagno/swagno/components/http/response"
	swaggerParameter "github.com/go-swagno/swagno/components/parameter"

	// OpenAPI 3.0 imports
	v3 "github.com/go-swagno/swagno/v3"
	openapiEndpoint "github.com/go-swagno/swagno/v3/components/endpoint"
	openapiResponse "github.com/go-swagno/swagno/v3/components/http/response"
	openapiParameter "github.com/go-swagno/swagno/v3/components/parameter"
)

type User struct {
	ID       uint64  `json:"id" example:"1"`
	Name     string  `json:"name" example:"John Doe"`
	Email    *string `json:"email,omitempty" example:"john@example.com"`
	IsActive bool    `json:"is_active" example:"true"`
}

type CreateUserRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type ErrorResponse struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Something went wrong"`
}

// ExampleSwagger20 demonstrates Swagger 2.0 usage
func ExampleSwagger20() []byte {
	// Create Swagger 2.0 instance
	sw := swagno.New(swagno.Config{
		Title:       "User API",
		Version:     "v1.0.0",
		Description: "User management API using Swagger 2.0",
		Host:        "api.example.com",
		Path:        "/v1",
		Contact: &swagno.Contact{
			Name:  "API Support",
			Email: "support@example.com",
		},
	})

	// Security (limited options in 2.0)
	sw.SetBasicAuth("Basic authentication required")
	sw.SetApiKeyAuth("X-API-Key", "header", "API key authentication")

	// Define endpoints
	endpoints := []*swaggerEndpoint.EndPoint{
		swaggerEndpoint.New(
			swaggerEndpoint.GET,
			"/users/{id}",
			swaggerEndpoint.WithParams(
				swaggerParameter.IntParam("id", swaggerParameter.Path,
					swaggerParameter.WithRequired(),
					swaggerParameter.WithDescription("User ID"),
				),
			),
			swaggerEndpoint.WithSuccessfulReturns([]swaggerResponse.Response{
				swaggerResponse.New(User{}, "200", "User found"),
			}),
			swaggerEndpoint.WithErrors([]swaggerResponse.Response{
				swaggerResponse.New(ErrorResponse{}, "404", "User not found"),
			}),
		),
		swaggerEndpoint.New(
			swaggerEndpoint.POST,
			"/users",
			swaggerEndpoint.WithBody(CreateUserRequest{}),
			swaggerEndpoint.WithSuccessfulReturns([]swaggerResponse.Response{
				swaggerResponse.New(User{}, "201", "User created"),
			}),
		),
	}

	sw.AddEndpoints(endpoints)
	return sw.MustToJson()
}

// ExampleOpenAPI30 demonstrates OpenAPI 3.0 usage
func ExampleOpenAPI30() []byte {
	// Create OpenAPI 3.0 instance
	openapi := v3.New(v3.Config{
		Title:       "User API",
		Version:     "v1.0.0",
		Description: "User management API using OpenAPI 3.0",
		Contact: &v3.Contact{
			Name:  "API Support",
			Email: "support@example.com",
			URL:   "https://example.com/support",
		},
		License: &v3.License{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
	})

	// Multiple servers support
	openapi.AddServer("https://api.example.com/v1", "Production server")
	openapi.AddServer("https://staging-api.example.com/v1", "Staging server")
	openapi.AddServer("http://localhost:8080/v1", "Development server")

	// Enhanced security schemes
	openapi.SetBasicAuth("Basic authentication required")
	openapi.SetBearerAuth("JWT", "JWT Bearer token authentication")
	openapi.SetApiKeyAuth("X-API-Key", "header", "API key authentication")

	// OAuth2 with proper flow definition
	oauthFlows := &v3.OAuthFlows{
		AuthorizationCode: &v3.OAuthFlow{
			AuthorizationUrl: "https://example.com/oauth/authorize",
			TokenUrl:         "https://example.com/oauth/token",
			Scopes: map[string]string{
				"read":  "Read access to user data",
				"write": "Write access to user data",
			},
		},
	}
	openapi.SetOAuth2Auth(oauthFlows, "OAuth2 authentication")

	// Define endpoints with enhanced features
	endpoints := []*openapiEndpoint.EndPoint{
		openapiEndpoint.New(
			openapiEndpoint.GET,
			"/users/{id}",
			openapiEndpoint.WithSummary("Get user by ID"),
			openapiEndpoint.WithDescription("Retrieve a specific user by their unique identifier"),
			openapiEndpoint.WithParams(
				openapiParameter.IntParam("id", openapiParameter.Path,
					openapiParameter.WithRequired(),
					openapiParameter.WithDescription("Unique identifier for the user"),
					openapiParameter.WithExample(123),
					openapiParameter.WithMin(1),
				),
			),
			openapiEndpoint.WithSuccessfulReturns([]openapiResponse.Response{
				openapiResponse.New(User{}, "200", "User found successfully"),
			}),
			openapiEndpoint.WithErrors([]openapiResponse.Response{
				openapiResponse.New(ErrorResponse{}, "400", "Invalid user ID"),
				openapiResponse.New(ErrorResponse{}, "404", "User not found"),
				openapiResponse.New(ErrorResponse{}, "500", "Internal server error"),
			}),
			openapiEndpoint.WithSecurity([]map[string][]string{
				{"bearerAuth": {}},
				{"apiKeyAuth": {}},
			}),
		),
		openapiEndpoint.New(
			openapiEndpoint.POST,
			"/users",
			openapiEndpoint.WithSummary("Create new user"),
			openapiEndpoint.WithDescription("Create a new user account"),
			openapiEndpoint.WithBody(CreateUserRequest{}),
			openapiEndpoint.WithSuccessfulReturns([]openapiResponse.Response{
				openapiResponse.New(User{}, "201", "User created successfully"),
			}),
			openapiEndpoint.WithErrors([]openapiResponse.Response{
				openapiResponse.New(ErrorResponse{}, "400", "Invalid request body"),
				openapiResponse.New(ErrorResponse{}, "409", "User already exists"),
				openapiResponse.New(ErrorResponse{}, "500", "Internal server error"),
			}),
			openapiEndpoint.WithSecurity([]map[string][]string{
				{"bearerAuth": {}},
			}),
		),
	}

	openapi.AddEndpoints(endpoints)
	return openapi.MustToJson()
}

// Key differences between Swagger 2.0 and OpenAPI 3.0:
//
// 1. **Package Import**:
//    - v2: github.com/go-swagno/swagno
//    - v3: github.com/go-swagno/swagno/v3
//
// 2. **Constructor**:
//    - v2: swagno.New()
//    - v3: v3.New()
//
// 3. **Server Configuration**:
//    - v2: Host and Path fields in config
//    - v3: Multiple servers with AddServer()
//
// 4. **Schema References**:
//    - v2: #/definitions/ModelName
//    - v3: #/components/schemas/ModelName
//
// 5. **Security Schemes**:
//    - v2: Limited to basic, apiKey, oauth2
//    - v3: Adds bearer tokens, OpenID Connect, better OAuth2 flows
//
// 6. **Request Bodies**:
//    - v2: Body parameter with schema
//    - v3: Dedicated requestBody object with content types
//
// 7. **Response Structure**:
//    - v2: Simple schema reference
//    - v3: Content types with media type objects
//
// 8. **Parameter Features**:
//    - v2: Basic parameter definition
//    - v3: Enhanced with examples, styles, explode, deprecated
//
// 9. **Schema Features**:
//    - v2: Basic type definitions
//    - v3: Nullable, readOnly, writeOnly, deprecated, better composition
//
// 10. **Specification Version**:
//     - v2: "swagger": "2.0"
//     - v3: "openapi": "3.0.3"
