package v3

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

func appendResponses(sourceResponses map[string]endpoint.JsonResponse, additionalResponses []response.Response) map[string]endpoint.JsonResponse {
	responseGenerator := response.NewResponseGenerator()

	for _, resp := range additionalResponses {
		var responseSchema *parameter.JsonResponseSchema

		switch respType := resp.(type) {
		case response.CustomResponse:
			responseSchema = responseGenerator.Generate(respType.Model)
		case response.Response:
			responseSchema = responseGenerator.Generate(respType)
		}

		sourceResponses[resp.ReturnCode()] = endpoint.JsonResponse{
			Description: resp.Description(),
			Content: map[string]endpoint.MediaType{
				"application/json": {
					Schema: responseSchema,
				},
			},
		}
	}

	return sourceResponses
}

func (o *OpenAPI) generateOpenAPIJson() {
	if len(o.endpoints) == 0 {
		log.Println("No endpoints found")
		return
	}

	// generate schemas component of OpenAPI json: https://spec.openapis.org/oas/v3.0.3#components-object
	o.generateOpenAPIDefinition()

	// convert all user EndPoint models to 'paths' fields of OpenAPI json
	// https://spec.openapis.org/oas/v3.0.3#paths-object
	for _, e := range o.endpoints {
		path := e.Path()

		if o.Paths[path] == nil {
			o.Paths[path] = make(map[string]endpoint.JsonEndPoint)
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
			pj := param.AsJson()
			if pj.In != parameter.Query.String() {
				pj.Type = ""
			}
			parameters = append(parameters, param.AsJson())
		}

		// Creates the schema definition for all successful return and error objects, and then links them in the responses section
		responses := map[string]endpoint.JsonResponse{}
		responses = appendResponses(responses, e.SuccessfulReturns())
		responses = appendResponses(responses, e.Errors())

		// add each endpoint to paths field of OpenAPI
		je := e.AsJson()
		je.OperationId = method + "-" + path
		je.Parameters = parameters
		je.Responses = responses

		// Handle request body for OpenAPI 3.0
		if bjp := e.BodyJsonParameter(); bjp != nil {
			requestBody := endpoint.RequestBody{
				Description: "Request body",
				Required:    true,
				Content: map[string]endpoint.MediaType{
					"application/json": {
						Schema: bjp.Schema,
					},
				},
			}
			je.RequestBody = &requestBody
		}

		o.Paths[path][method] = je
	}
}

// ToJson converts the OpenAPI object into its JSON representation formatted as bytes.
// It returns a slice of bytes containing the OpenAPI documentation in JSON format.
func (o *OpenAPI) ToJson() (jsonDocs []byte, err error) {
	o.generateOpenAPIJson()
	return json.MarshalIndent(o, "", "  ")
}

// MustToJson same thing as ToJson except for it doesn't return an error.
func (o OpenAPI) MustToJson() (jsonDocs []byte) {
	o.generateOpenAPIJson()

	json, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Printf("Error while generating OpenAPI json: %s", err)
	}

	return json
}

// generate "schemas" keys from endpoints: https://spec.openapis.org/oas/v3.0.3#schema-object
func (o *OpenAPI) generateOpenAPIDefinition() {
	for _, endpoint := range o.endpoints {
		if endpoint.Body != nil {
			o.createDefinition(endpoint.Body)
		}
		o.createDefinitions(endpoint.SuccessfulReturns())
		o.createDefinitions(endpoint.Errors())
	}
}

func (o *OpenAPI) createDefinitions(r []response.Response) {
	for _, obj := range r {
		o.createDefinition(obj)
	}
}

func (o *OpenAPI) createDefinition(t interface{}) {
	if o.Components == nil {
		o.Components = &Components{
			Schemas: make(map[string]definition.Schema),
		}
	}
	if o.Components.Schemas == nil {
		o.Components.Schemas = make(map[string]definition.Schema)
	}

	generator := definition.NewDefinitionGenerator(o.Components.Schemas)
	generator.CreateDefinition(t)
}
