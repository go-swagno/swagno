package swagno3

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/go-swagno/swagno/v3/components/definition"
	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/mime"
	"github.com/go-swagno/swagno/v3/components/parameter"
)

func appendResponses(sourceResponses map[string]endpoint.JsonResponse, additionalResponses []response.Response, hidePackageName bool) map[string]endpoint.JsonResponse {
	responseGenerator := response.NewResponseGenerator(hidePackageName)

	for _, resp := range additionalResponses {
		var responseSchema *parameter.JsonResponseSchema
		var example interface{}
		var examples map[string]interface{}

		switch respType := resp.(type) {
		case response.CustomResponse:
			responseSchema = responseGenerator.Generate(respType.Model)
			example = respType.Example()
			examples = respType.Examples()
		case response.Response:
			responseSchema = responseGenerator.Generate(respType)
		}

		// Support multiple content types, defaulting to application/json
		jsonMediaType := endpoint.MediaType{
			Schema: responseSchema,
		}

		// Add example if available
		if example != nil {
			jsonMediaType.Example = example
		}

		// Add examples if available
		if examples != nil {
			// Convert raw examples to proper ComponentExample objects
			exampleObjects := make(map[string]parameter.ComponentExample)
			for key, value := range examples {
				exampleObjects[key] = parameter.ComponentExample{
					Value: value,
				}
			}
			jsonMediaType.Examples = exampleObjects
		}

		content := map[string]endpoint.MediaType{
			string(mime.JSON): jsonMediaType,
		}

		// TODO: Add support for other content types based on endpoint configuration
		// This could be extended to support:
		// - application/xml
		// - text/plain
		// - multipart/form-data
		// - etc.

		sourceResponses[resp.ReturnCode()] = endpoint.JsonResponse{
			Description: resp.Description(),
			Content:     content,
		}
	}

	return sourceResponses
}

func (o *OpenAPI) generateOpenAPIJson() error {
	if len(o.endpoints) == 0 {
		log.Println("No endpoints found")
		return nil
	}

	// generate schemas component of OpenAPI json: https://spec.openapis.org/oas/v3.0.3#components-object
	if err := o.generateOpenAPIDefinition(); err != nil {
		return err
	}

	// convert all user EndPoint models to 'paths' fields of OpenAPI json
	// https://spec.openapis.org/oas/v3.0.3#paths-object
	for _, e := range o.endpoints {
		path := e.Path()

		// Initialize PathItem if it doesn't exist
		pathItem, exists := o.Paths[path]
		if !exists {
			pathItem = endpoint.PathItem{}
		}

		method := strings.ToLower(string(e.Method()))

		for _, param := range e.Params() {
			if param.Location() == parameter.Form {
				endpoint.WithConsume([]mime.MIME{mime.MULTIFORM})(e)
				break
			}
		}

		parameters := make([]parameter.JsonParameter, 0)
		for _, param := range e.Params() {
			parameters = append(parameters, param.AsJson())
		}

		// Creates the schema definition for all successful return and error objects, and then links them in the responses section
		responses := map[string]endpoint.JsonResponse{}
		responses = appendResponses(responses, e.SuccessfulReturns(), o.hidePackageName)
		responses = appendResponses(responses, e.Errors(), o.hidePackageName)

		// add each endpoint to paths field of OpenAPI
		je := e.AsJson()
		je.OperationId = o.sanitizeOperationID(method + "-" + path)
		je.Parameters = parameters
		je.Responses = responses

		for _, res := range je.Responses {
			content := res.Content[string(mime.JSON)]
			for _, contentType := range je.Produce {
				if contentType != mime.JSON {
					res.Content[string(contentType)] = content
				}
			}
		}

		// Handle request body for OpenAPI 3.0
		if bjp := e.BodyJsonParameter(o.hidePackageName); bjp != nil {
			// Support multiple content types for request body
			content := map[string]endpoint.MediaType{}

			for _, m := range je.Consume {
				mediaType := endpoint.MediaType{
					Schema: bjp.Schema,
				}

				// Add example if available
				if bjp.Example != nil {
					mediaType.Example = bjp.Example
				}

				// Add examples if available
				if bjp.Examples != nil {
					mediaType.Examples = bjp.Examples
				}

				content[string(m)] = mediaType
			}

			// Add form-data support if needed
			for _, param := range e.Params() {
				if param.Location() == parameter.Form {
					mediaType := endpoint.MediaType{
						Schema: bjp.Schema,
						// TODO: Add encoding configuration for form fields
					}

					// Add example if available
					if bjp.Example != nil {
						mediaType.Example = bjp.Example
					}

					// Add examples if available
					if bjp.Examples != nil {
						mediaType.Examples = bjp.Examples
					}

					content["multipart/form-data"] = mediaType
					break
				}
			}

			requestBody := endpoint.RequestBody{
				Description: "Request body",
				Required:    bjp.Required,
				Content:     content,
			}
			je.RequestBody = &requestBody
		}

		// Add operation to PathItem using helper method
		methodType := e.Method()
		pathItem.AddOperation(methodType, &je)

		// Update the PathItem in the map
		o.Paths[path] = pathItem
	}

	return nil
}

// ToJson converts the OpenAPI object into its JSON representation formatted as bytes.
// It returns a slice of bytes containing the OpenAPI documentation in JSON format.
func (o *OpenAPI) ToJson() (jsonDocs []byte, err error) {
	if err := o.generateOpenAPIJson(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(o, "", "  ")
}

// MustToJson same thing as ToJson except for it doesn't return an error.
// It panics if a name collision is detected while generating the document.
func (o OpenAPI) MustToJson() (jsonDocs []byte) {
	if err := o.generateOpenAPIJson(); err != nil {
		panic(err)
	}

	json, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Printf("Error while generating OpenAPI json: %s", err)
	}

	return json
}

// generate "schemas" keys from endpoints: https://spec.openapis.org/oas/v3.0.3#schema-object
// It returns a *NameCollisionError when HidePackageName causes two distinct types
// to map to the same stripped schema name.
func (o *OpenAPI) generateOpenAPIDefinition() error {
	// shared across all createDefinition calls so collisions are detected document-wide
	definitionTypeNames := map[string]map[string]struct{}{}
	for _, endpoint := range o.endpoints {
		if endpoint.Body.Content != nil {
			o.createDefinition(endpoint.Body.Content, definitionTypeNames)
		}
		o.createDefinitions(endpoint.SuccessfulReturns(), definitionTypeNames)
		o.createDefinitions(endpoint.Errors(), definitionTypeNames)
	}
	return collisionError(definitionTypeNames)
}

func (o *OpenAPI) createDefinitions(r []response.Response, definitionTypeNames map[string]map[string]struct{}) {
	for _, obj := range r {
		o.createDefinition(obj, definitionTypeNames)
	}
}

func (o *OpenAPI) createDefinition(t interface{}, definitionTypeNames map[string]map[string]struct{}) {
	if o.Components == nil {
		o.Components = &Components{
			Schemas: make(map[string]definition.Schema),
		}
	}
	if o.Components.Schemas == nil {
		o.Components.Schemas = make(map[string]definition.Schema)
	}

	generator := definition.NewDefinitionGenerator(o.Components.Schemas, o.hidePackageName, definitionTypeNames)
	generator.CreateDefinition(t)
}

func (s *OpenAPI) sanitizeOperationID(operationID string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(operationID, "/", "_"), "{", ""), "}", "")
}
