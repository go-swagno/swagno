package swagno

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/domhoward14/swagno/components/endpoint"
	"github.com/domhoward14/swagno/components/parameter"
	"github.com/domhoward14/swagno/components/response"
	"github.com/domhoward14/swagno/example/models"
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
					endpoint.WithMethod(endpoint.GET),
					endpoint.WithPath("/product"),
					endpoint.WithTags("product"),
					endpoint.WithSuccessfulReturns([]response.Info{models.UnsuccessfulResponse{}}),
					endpoint.WithErrors([]response.Info{models.EmptySuccessfulResponse{}}),
					endpoint.WithDescription(desc),
					endpoint.WithProduce([]string{"application/json", "application/xml"}),
					endpoint.WithConsume([]string{"application/json"}),
					endpoint.WithSummary("this is a test summary"),
				),
				endpoint.New(
					endpoint.WithMethod(endpoint.GET),
					endpoint.WithPath("/product/{id}"),
					endpoint.WithTags("product"),
					endpoint.WithParams([]parameter.Parameter{parameter.IntParam("id", parameter.WithRequired(true))}),
					endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
					endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
					endpoint.WithProduce([]string{"application/json", "application/xml"}),
				),
				endpoint.New(
					endpoint.WithMethod(endpoint.GET),
					endpoint.WithPath("/product/{id}/detail"),
					endpoint.WithTags("product"),
					endpoint.WithParams([]parameter.Parameter{parameter.IntParam("id", parameter.WithRequired(true))}),
					endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
					endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
					endpoint.WithProduce([]string{"application/json", "application/xml"}),
				),
				endpoint.New(
					endpoint.WithMethod(endpoint.POST),
					endpoint.WithPath("/product"),
					endpoint.WithTags("product"),
					endpoint.WithParams([]parameter.Parameter{}),
					endpoint.WithBody(models.ProductPost{}),
					endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
					endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
					endpoint.WithProduce([]string{"application/json", "application/xml"}),
				),
			},
			file: "testdata/expected_output/bft.json",
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
			got.generateSwaggerJson()

			if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(JsonSwagger{}, "endpoints"), cmpopts.IgnoreFields(jsonDefinitionProperties{}, "Example")); diff != "" {
				t.Errorf("JsonSwagger() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}
