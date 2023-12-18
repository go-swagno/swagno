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
			Schema:      responseSchema,
		}
	}

	return sourceResponses
}

func (s *Swagger) generateSwaggerJson() {
	if len(s.endpoints) == 0 {
		log.Println("No endpoints found")
		return
	}

	// generate definition object of swagger json: https://swagger.io/specification/v2/#definitions-object
	s.generateSwaggerDefinition()

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

		if bjp := e.BodyJsonParameter(); bjp != nil {
			parameters = append(parameters, *bjp)
		}

		// Creates the schema defintion for all successful return and error objects, and then links them in the responses section
		responses := map[string]endpoint.JsonResponse{}
		responses = appendResponses(responses, e.SuccessfulReturns())
		responses = appendResponses(responses, e.Errors())

		// add each endpoint to paths field of swagger
		je := e.AsJson()
		je.OperationId = method + "-" + path
		je.Parameters = parameters
		je.Responses = responses
		s.Paths[path][method] = je
	}
}

// ToJSON converts the Swagger object into its JSON representation formatted as bytes.
// It returns a slice of bytes containing the Swagger documentation in JSON format.
func (s *Swagger) ToJson() (jsonDocs []byte, err error) {
	s.generateSwaggerJson()
	return json.MarshalIndent(s, "", "  ")
}

// MustToJson same thing as ToJson except for it doesn't return an error.
func (s Swagger) MustToJson() (jsonDocs []byte) {
	s.generateSwaggerJson()

	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("Error while generating swagger json: %s", err)
	}

	return json
}

// generate "definitions" keys from endpoints: https://swagger.io/specification/v2/#definitions-object
func (s *Swagger) generateSwaggerDefinition() {
	for _, endpoint := range s.endpoints {
		if endpoint.Body != nil {
			s.createDefinition(endpoint.Body)
		}
		s.createDefinitions(endpoint.SuccessfulReturns())
		s.createDefinitions(endpoint.Errors())
	}
}

func (s *Swagger) createDefinitions(r []response.Response) {
	for _, obj := range r {
		s.createDefinition(obj)
	}
}

func (s *Swagger) createDefinition(t interface{}) {
	generator := definition.NewDefinitionGenerator((*s).Definitions)
	generator.CreateDefinition(t)
}
