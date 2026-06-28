package swagno

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/go-swagno/swagno/components/definition"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
)

func appendResponses(sourceResponses map[string]endpoint.JsonResponse, additionalResponses []response.Response, hidePackageName bool) map[string]endpoint.JsonResponse {
	responseGenerator := response.NewResponseGenerator(hidePackageName)

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
			Schema:      responseSchema,
		}
	}

	return sourceResponses
}

func (s *Swagger) generateSwaggerJson() error {
	if len(s.endpoints) == 0 {
		log.Println("No endpoints found")
		return nil
	}

	// generate definition object of swagger json: https://swagger.io/specification/v2/#definitions-object
	if err := s.generateSwaggerDefinition(); err != nil {
		return err
	}

	// convert all user EndPoint models to 'path' fields of swagger json
	// https://swagger.io/specification/v2/#paths-object
	for _, e := range s.endpoints {
		path := e.Path()

		if s.Paths[path] == nil {
			s.Paths[path] = make(map[string]endpoint.JsonEndPoint)
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

		if bjp := e.BodyJsonParameter(s.hidePackageName); bjp != nil {
			parameters = append(parameters, *bjp)
		}

		// Creates the schema defintion for all successful return and error objects, and then links them in the responses section
		responses := map[string]endpoint.JsonResponse{}
		responses = appendResponses(responses, e.SuccessfulReturns(), s.hidePackageName)
		responses = appendResponses(responses, e.Errors(), s.hidePackageName)

		// add each endpoint to paths field of swagger
		je := e.AsJson()
		je.OperationId = method + "-" + path
		je.Parameters = parameters
		je.Responses = responses
		s.Paths[path][method] = je
	}

	return nil
}

// ToJSON converts the Swagger object into its JSON representation formatted as bytes.
// It returns a slice of bytes containing the Swagger documentation in JSON format.
func (s *Swagger) ToJson() (jsonDocs []byte, err error) {
	if err := s.generateSwaggerJson(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(s, "", "  ")
}

// MustToJson same thing as ToJson except for it doesn't return an error.
// It panics if a name collision is detected while generating the document.
func (s Swagger) MustToJson() (jsonDocs []byte) {
	if err := s.generateSwaggerJson(); err != nil {
		panic(err)
	}

	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("Error while generating swagger json: %s", err)
	}

	return json
}

// generate "definitions" keys from endpoints: https://swagger.io/specification/v2/#definitions-object
// It returns a *NameCollisionError when HidePackageName causes two distinct types
// to map to the same stripped definition name.
func (s *Swagger) generateSwaggerDefinition() error {
	// shared across all createDefinition calls so collisions are detected document-wide
	definitionTypeNames := map[string]map[string]struct{}{}
	for _, endpoint := range s.endpoints {
		if endpoint.Body.Content != nil {
			s.createDefinition(endpoint.Body.Content, definitionTypeNames)
		}
		s.createDefinitions(endpoint.SuccessfulReturns(), definitionTypeNames)
		s.createDefinitions(endpoint.Errors(), definitionTypeNames)
	}
	return collisionError(definitionTypeNames)
}

func (s *Swagger) createDefinitions(r []response.Response, definitionTypeNames map[string]map[string]struct{}) {
	for _, obj := range r {
		s.createDefinition(obj, definitionTypeNames)
	}
}

func (s *Swagger) createDefinition(t interface{}, definitionTypeNames map[string]map[string]struct{}) {
	generator := definition.NewDefinitionGenerator((*s).Definitions, s.hidePackageName, definitionTypeNames)
	generator.CreateDefinition(t)
}
