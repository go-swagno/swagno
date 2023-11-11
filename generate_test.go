package swagno

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/go-swagno/swagno/example/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TODO add tests so that it works the way it did before
// with one error and successful return responses
// and multiple errors

func TestSwagger_GenerateDocs(t *testing.T) {
	type fields struct {
		Swagger             string
		Info                swaggerInfo
		Paths               map[string]map[string]swaggerEndpoint
		BasePath            string
		Host                string
		Definitions         map[string]swaggerDefinition
		Schemes             []string
		Tags                []SwaggerTag
		SecurityDefinitions map[string]swaggerSecurityDefinition
	}
	tests := []struct {
		name         string
		fields       fields
		wantJsonDocs []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			swagger := Swagger{
				Swagger:             tt.fields.Swagger,
				Info:                tt.fields.Info,
				Paths:               tt.fields.Paths,
				BasePath:            tt.fields.BasePath,
				Host:                tt.fields.Host,
				Definitions:         tt.fields.Definitions,
				Schemes:             tt.fields.Schemes,
				Tags:                tt.fields.Tags,
				SecurityDefinitions: tt.fields.SecurityDefinitions,
			}
			if gotJsonDocs := swagger.GenerateDocs(); !reflect.DeepEqual(gotJsonDocs, tt.wantJsonDocs) {
				t.Errorf("Swagger.GenerateDocs() = %v, want %v", gotJsonDocs, tt.wantJsonDocs)
			}
		})
	}
}

type TestErrorResponses struct{}

func (e TestErrorResponses) GetErrors() []ResponseInfo {
	return []ResponseInfo{models.UnsuccessfulResponse{}, models.PageNotFound{}}
}

func TestPlayground(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name      string
		endpoints []Endpoint
		file      string
	}{
		{
			name: "test 1",
			endpoints: []Endpoint{
				EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil),
				EndPoint(GET, "/product/{id}", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
				EndPoint(GET, "/product/{id}/detail", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
				EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
			},
			file: "testdata/test-1.json",
		},
		{
			name: "test 2",
			endpoints: []Endpoint{
				EndPoint(GET, "/product", "product", Params(), nil, models.SuccessfulResponse{}, TestErrorResponses{}, "", nil),
			},
			file: "testdata/test-2.json",
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatalf("Error reading file: %v", err)
			}

			testSwagger := &Swagger{}
			if err := json.Unmarshal(jsonData, testSwagger); err != nil {
				t.Fatal(err)
			}

			c := Config{"Swagger API", "1.0", "localhost", "/v2", nil, nil}
			sw := CreateNewSwagger(c)
			sw.AddEndpoints(tc.endpoints)
			sw.generateSwaggerObject()

			if diff := cmp.Diff(testSwagger, sw, cmpopts.IgnoreFields(Swagger{}, "SecurityDefinitions", "endpoints")); diff != "" {
				t.Errorf("Swagger() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
