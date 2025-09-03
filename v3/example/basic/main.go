package main

import (
	v3 "github.com/go-swagno/swagno/v3"
	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/parameter"
	"github.com/go-swagno/swagno/v3/components/tag"
)

// User represents a user model
type User struct {
	ID       uint64 `json:"id" example:"1"`
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	IsActive *bool  `json:"is_active,omitempty" example:"true"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"User not found"`
}

func main() {
	// Create OpenAPI 3.0 instance
	openapi := v3.New(v3.Config{
		Title:       "User Management API",
		Version:     "v1.0.0",
		Description: "A simple API for managing users",
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

	// Add servers
	openapi.AddServer("https://api.example.com/v1", "Production server")
	openapi.AddServer("https://staging-api.example.com/v1", "Staging server")

	// Add security schemes
	openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")
	openapi.SetApiKeyAuth("X-API-Key", "header", "API key authentication")

	// Add tags
	openapi.AddTags(
		tag.New("users", "User management operations",
			tag.WithExternalDocs("https://docs.example.com/users", "User API Documentation"),
		),
	)

	// Define endpoints
	endpoints := []*endpoint.EndPoint{
		// GET /users - List all users
		endpoint.New(
			endpoint.GET,
			"/users",
			endpoint.WithTags("users"),
			endpoint.WithSummary("List users"),
			endpoint.WithDescription("Retrieve a list of all users"),
			endpoint.WithParams(
				parameter.IntParam("page", parameter.Query,
					parameter.WithDescription("Page number"),
					parameter.WithDefault(1),
					parameter.WithMin(1),
					parameter.WithExample(1),
				),
				parameter.IntParam("limit", parameter.Query,
					parameter.WithDescription("Number of users per page"),
					parameter.WithDefault(10),
					parameter.WithMin(1),
					parameter.WithMax(100),
					parameter.WithExample(10),
				),
			),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New([]User{}, "200", "List of users"),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(ErrorResponse{}, "400", "Bad Request"),
				response.New(ErrorResponse{}, "500", "Internal Server Error"),
			}),
			endpoint.WithSecurity([]map[string][]string{
				{"bearerAuth": {}},
			}),
		),

		// POST /users - Create a new user
		endpoint.New(
			endpoint.POST,
			"/users",
			endpoint.WithTags("users"),
			endpoint.WithSummary("Create user"),
			endpoint.WithDescription("Create a new user"),
			endpoint.WithBody(CreateUserRequest{}),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(User{}, "201", "User created successfully"),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(ErrorResponse{}, "400", "Bad Request"),
				response.New(ErrorResponse{}, "409", "User already exists"),
				response.New(ErrorResponse{}, "500", "Internal Server Error"),
			}),
			endpoint.WithSecurity([]map[string][]string{
				{"bearerAuth": {}},
			}),
		),

		// GET /users/{id} - Get a specific user
		endpoint.New(
			endpoint.GET,
			"/users/{id}",
			endpoint.WithTags("users"),
			endpoint.WithSummary("Get user by ID"),
			endpoint.WithDescription("Retrieve a specific user by their ID"),
			endpoint.WithParams(
				parameter.IntParam("id", parameter.Path,
					parameter.WithRequired(),
					parameter.WithDescription("User ID"),
					parameter.WithMin(1),
					parameter.WithExample(1),
				),
			),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(User{}, "200", "User found"),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(ErrorResponse{}, "400", "Bad Request"),
				response.New(ErrorResponse{}, "404", "User not found"),
				response.New(ErrorResponse{}, "500", "Internal Server Error"),
			}),
			endpoint.WithSecurity([]map[string][]string{
				{"bearerAuth": {}},
				{"apiKeyAuth": {}},
			}),
		),

		// PUT /users/{id} - Update a user
		endpoint.New(
			endpoint.PUT,
			"/users/{id}",
			endpoint.WithTags("users"),
			endpoint.WithSummary("Update user"),
			endpoint.WithDescription("Update an existing user"),
			endpoint.WithParams(
				parameter.IntParam("id", parameter.Path,
					parameter.WithRequired(),
					parameter.WithDescription("User ID"),
					parameter.WithMin(1),
					parameter.WithExample(1),
				),
			),
			endpoint.WithBody(CreateUserRequest{}),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(User{}, "200", "User updated successfully"),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(ErrorResponse{}, "400", "Bad Request"),
				response.New(ErrorResponse{}, "404", "User not found"),
				response.New(ErrorResponse{}, "500", "Internal Server Error"),
			}),
			endpoint.WithSecurity([]map[string][]string{
				{"bearerAuth": {}},
			}),
		),

		// DELETE /users/{id} - Delete a user
		endpoint.New(
			endpoint.DELETE,
			"/users/{id}",
			endpoint.WithTags("users"),
			endpoint.WithSummary("Delete user"),
			endpoint.WithDescription("Delete an existing user"),
			endpoint.WithParams(
				parameter.IntParam("id", parameter.Path,
					parameter.WithRequired(),
					parameter.WithDescription("User ID"),
					parameter.WithMin(1),
					parameter.WithExample(1),
				),
			),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(struct{}{}, "204", "User deleted successfully"),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(ErrorResponse{}, "400", "Bad Request"),
				response.New(ErrorResponse{}, "404", "User not found"),
				response.New(ErrorResponse{}, "500", "Internal Server Error"),
			}),
			endpoint.WithSecurity([]map[string][]string{
				{"bearerAuth": {}},
			}),
		),
	}

	// Add endpoints to OpenAPI
	openapi.AddEndpoints(endpoints)

	// Optionally save to file
	openapi.ExportOpenAPIDocs("openapi.json")
}
