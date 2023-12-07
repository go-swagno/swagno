package swagno

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/generator"
	"github.com/go-swagno/swagno/http/response"
)

func appendResponses(sourceResponses map[string]endpoint.JsonResponse, additionalResponses []response.Info) map[string]endpoint.JsonResponse {
	responseGenerator := generator.NewResponseGenerator()

	for _, resp := range additionalResponses {
		var responseSchema *parameter.JsonResponseSchema

		switch _resp := resp.(type) {
		case response.CustomResponseType:
			responseSchema = responseGenerator.GenerateJsonResponseScheme(_resp.Model)
		case response.Info:
			responseSchema = responseGenerator.GenerateJsonResponseScheme(_resp)
		}

		sourceResponses[resp.GetReturnCode()] = endpoint.JsonResponse{
			Description: resp.GetDescription(),
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
		path := e.GetPath()

		if s.Paths[path] == nil {
			s.Paths[path] = make(map[string]endpoint.JsonEndPoint)
		}

		method := strings.ToLower(string(e.GetMethod()))

		for _, param := range e.GetParams() {
			if param.GetLocation() == parameter.Form {
				endpoint.WithConsume([]mime.MIME{mime.MULTIFORM})(e)
				break
			}
		}

		parameters := make([]parameter.JsonParameter, 0)
		for _, param := range e.GetParams() {
			pj := param.AsJson()
			if pj.In != parameter.Query.String() {
				pj.Type = ""
			}
			parameters = append(parameters, param.AsJson())
		}

		if bjp := e.GetBodyJsonParameter(); bjp != nil {
			parameters = append(parameters, *bjp)
		}

		// Creates the schema defintion for all successful return and error objects, and then links them in the responses section
		responses := map[string]endpoint.JsonResponse{}
		responses = appendResponses(responses, e.GetSuccessfulReturns())
		responses = appendResponses(responses, e.GetErrors())

		// add each endpoint to paths field of swagger
		je := e.AsJson()
		je.OperationId = method + "-" + path
		je.Parameters = parameters
		je.Responses = responses
		s.Paths[path][method] = je
	}
}

// Generate swagger v2 documentation as json string
func (s Swagger) GenerateDocs() (jsonDocs []byte) {
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
		s.createDefinitions(endpoint.GetSuccessfulReturns())
		s.createDefinitions(endpoint.GetErrors())
	}
}

func (s *Swagger) createDefinitions(r []response.Info) {
	for _, obj := range r {
		s.createDefinition(obj)
	}
}

func (s *Swagger) createDefinition(t interface{}) {
	generator := generator.NewDefinitionGenerator((*s).Definitions)
	generator.CreateDefinition(t)
}
