package swagno

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/domhoward14/swagno/components/definition"
	"github.com/domhoward14/swagno/components/endpoint"
	"github.com/domhoward14/swagno/components/mime"
	"github.com/domhoward14/swagno/components/parameter"
	"github.com/domhoward14/swagno/example/models"
	"github.com/domhoward14/swagno/http/response"
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
			// TODO i don't want to ruin the current design by adding in validation checks for user input, but a nice compromise
			// would be to have an optional function that iterates through the endpoints and validates all the endpoints for syntactical
			// errors. This way the Swagno isn't more restrictive then the actual OpenAPI parser is when rendering
			// but the client still has an option to check for errors if they wish to do so.
			endpoints: []*endpoint.EndPoint{
				endpoint.New(
					endpoint.WithMethod(endpoint.GET),
					endpoint.WithPath("/product"),
					endpoint.WithTags("product"),
					endpoint.WithSuccessfulReturns([]response.Info{models.UnsuccessfulResponse{}}),
					endpoint.WithErrors([]response.Info{models.EmptySuccessfulResponse{}}),
					endpoint.WithDescription(desc),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
					endpoint.WithConsume([]mime.MIME{mime.JSON}),
					endpoint.WithSummary("this is a test summary"),
				),
				endpoint.New(
					endpoint.WithMethod(endpoint.GET),
					endpoint.WithPath("/product/{id}"),
					endpoint.WithTags("product"),
					endpoint.WithParams(parameter.IntParam("id", parameter.WithIn(parameter.Path), parameter.WithRequired())),
					endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
					endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
				),
				endpoint.New(
					endpoint.WithMethod(endpoint.GET),
					endpoint.WithPath("/product/{id}/detail"),
					endpoint.WithTags("product"),
					endpoint.WithParams(parameter.IntParam("id", parameter.WithIn(parameter.Path), parameter.WithRequired())),
					endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
					endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
				),
				endpoint.New(
					endpoint.WithMethod(endpoint.POST),
					endpoint.WithPath("/product"),
					endpoint.WithTags("product"),
					endpoint.WithBody(models.ProductPost{}),
					endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
					endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
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

			if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(Swagger{}, "endpoints"), cmpopts.IgnoreFields(definition.DefinitionProperties{}, "Example")); diff != "" {
				t.Errorf("JsonSwagger() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}
