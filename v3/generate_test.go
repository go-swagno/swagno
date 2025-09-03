package v3

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/go-swagno/swagno/v3/components/definition"
	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/mime"
	"github.com/go-swagno/swagno/v3/components/parameter"
	"github.com/go-swagno/swagno/v3/example/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var desc = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id malesuada lorem, et fermentum sapien. Vivamus non pharetra risus, in efficitur leo. Suspendisse sed metus sit amet mi laoreet imperdiet. Donec aliquam eros eu blandit feugiat. Quisque scelerisque justo ac vehicula bibendum. Fusce suscipit arcu nisl, eu maximus odio consequat quis. Curabitur fermentum eleifend tellus, lobortis hendrerit velit varius vitae."

func TestSwaggerGeneration(t *testing.T) {
	testCases := []struct {
		name      string
		endpoints []*endpoint.EndPoint
		file      string
	}{
		{
			name: "Basic Functionality Test",
			endpoints: []*endpoint.EndPoint{
				endpoint.New(
					endpoint.GET,
					"/product",
					endpoint.WithTags("product"),
					endpoint.WithErrors([]response.Response{response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")}),
					endpoint.WithSuccessfulReturns([]response.Response{response.New(models.EmptySuccessfulResponse{}, "200", "OK")}),
					endpoint.WithDescription(desc),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
					endpoint.WithConsume([]mime.MIME{mime.JSON}),
					endpoint.WithSummary("this is a test summary"),
				),
				endpoint.New(
					endpoint.GET,
					"/product/{id}",
					endpoint.WithTags("product"),
					endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
					endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessfulResponse{}, "201", "Request Accepted")}),
					endpoint.WithErrors([]response.Response{response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")}),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
				),
				endpoint.New(
					endpoint.GET,
					"/product/{id}/detail",
					endpoint.WithTags("product"),
					endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
					endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessfulResponse{}, "201", "Request Accepted")}),
					endpoint.WithErrors([]response.Response{response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")}),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
				),
				endpoint.New(
					endpoint.POST,
					"/product",
					endpoint.WithTags("product"),
					endpoint.WithBody(models.ProductPost{}),
					endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessfulResponse{}, "201", "Request Accepted")}),
					endpoint.WithErrors([]response.Response{response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")}),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
				),
			},
			file: "testdata/expected_output/bft.json",
		},
		{
			name: "Deeply Nested Model Test",
			endpoints: []*endpoint.EndPoint{
				endpoint.New(
					endpoint.GET,
					"/deeplynested",
					endpoint.WithSuccessfulReturns([]response.Response{response.New(models.ComplexSuccessfulResponse{}, "200", "OK")}),
					endpoint.WithDescription(desc),
					endpoint.WithProduce([]mime.MIME{mime.JSON}),
					endpoint.WithConsume([]mime.MIME{mime.JSON}),
					endpoint.WithSummary("this is a test summary"),
				),
				endpoint.New(
					endpoint.GET,
					"/arraydeeplynested",
					endpoint.WithSuccessfulReturns([]response.Response{response.New([]models.ComplexSuccessfulResponse{}, "200", "OK")}),
					endpoint.WithDescription(desc),
					endpoint.WithProduce([]mime.MIME{mime.JSON}),
					endpoint.WithConsume([]mime.MIME{mime.JSON}),
					endpoint.WithSummary("this is a test summary"),
				),
			},
			file: "testdata/expected_output/dnmt.json",
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			expectedJsonData, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatalf("Error reading file: %v", err)
			}

			want := New(Config{Title: "Testing API", Version: "v1.0.0"})
			if err := json.Unmarshal(expectedJsonData, want); err != nil {
				t.Fatal(err)
			}

			got := New(Config{Title: "Testing API", Version: "v1.0.0"})
			got.AddEndpoints(tc.endpoints)
			got.generateOpenAPIJson()

			if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(OpenAPI{}, "endpoints"), cmpopts.IgnoreFields(definition.SchemaProperty{}, "Example", "IsRequired"), cmpopts.IgnoreFields(endpoint.JsonEndPoint{}, "Consume", "Produce")); diff != "" {
				t.Errorf("OpenAPIJson() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}
