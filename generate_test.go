package swagno

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/go-swagno/swagno/components/definition"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/example/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var desc = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id malesuada lorem, et fermentum sapien. Vivamus non pharetra risus, in efficitur leo. Suspendisse sed metus sit amet mi laoreet imperdiet. Donec aliquam eros eu blandit feugiat. Quisque scelerisque justo ac vehicula bibendum. Fusce suscipit arcu nisl, eu maximus odio consequat quis. Curabitur fermentum eleifend tellus, lobortis hendrerit velit varius vitae."

func TestSwaggerGeneration(t *testing.T) {
	var id float64 = 1234
	testCases := []struct {
		name      string
		endpoints []*endpoint.EndPoint
		file      string
		debug     bool
	}{
		{
			name: "Basic Functionality Test",
			// TODO i don't want to ruin the current design by adding in validation checks for user input, but a nice compromise
			// would be to have an optional function that iterates through the endpoints and validates all the endpoints for syntactical
			// errors. This way the Swagno isn't more restrictive then the actual OpenAPI parser is when rendering
			// but the client still has an option to check for errors if they wish to do so.
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
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
					endpoint.WithConsume([]mime.MIME{mime.JSON}),
					endpoint.WithSummary("this is a test summary"),
				),
				endpoint.New(
					endpoint.GET,
					"/arraydeeplynested",
					endpoint.WithSuccessfulReturns([]response.Response{response.New([]models.ComplexSuccessfulResponse{}, "200", "OK")}),
					endpoint.WithDescription(desc),
					endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
					endpoint.WithConsume([]mime.MIME{mime.JSON}),
					endpoint.WithSummary("this is a test summary"),
				),
			},
			file: "testdata/expected_output/dnmt.json",
		},
		{
			name: "test map[string]any",
			endpoints: []*endpoint.EndPoint{endpoint.New(
				endpoint.PUT,
				"/test-map",
				endpoint.WithTags("product"),
				endpoint.WithSuccessfulReturns([]response.Response{response.New(
					map[string]any{
						"code": float64(200),
						"data": map[string]any{
							"id":  &id,
							"id2": "asd",
							"id3": 123.23,
							"id4": []int{12, 34},
							"id5": []string{"asd", "asd"},
						},
					},
					"201",
					"Request Accepted",
				)},
				),
			)},
			file: "testdata/expected_output/map.json",
		},
		{
			name: "test composable type with primitive type",
			endpoints: []*endpoint.EndPoint{
				endpoint.New(
					endpoint.PUT,
					"/composable-primitive",
					endpoint.WithTags("product"),
					endpoint.WithSuccessfulReturns([]response.Response{
						response.New(
							models.Response{
								Status: "success",
								Data:   "asd",
							},
							"200",
							"Request Accepted",
						),
						response.New(
							models.Response{
								Status: "success",
								Data:   1,
							},
							"201",
							"Request Accepted",
						),
						response.New(
							models.Response{
								Status: "success",
								Data:   true,
							},
							"203",
							"Request Accepted",
						),
						response.New(
							models.Response{
								Status: "success",
								Data:   &id,
							},
							"204",
							"Request Accepted",
						),
					}),
				),
			},
			file:  "testdata/expected_output/composable-primitive.json",
			debug: true,
		},
		{
			name: "test composable type with primitive array type",
			endpoints: []*endpoint.EndPoint{
				endpoint.New(
					endpoint.PUT,
					"/composable-array-primitive",
					endpoint.WithTags("product"),
					endpoint.WithSuccessfulReturns([]response.Response{

						response.New(
							models.Response{
								Status: "success",
								Data:   []*float64{&id, &id},
								Errors: []float64{1, 2},
							},
							"200",
							"Request Accepted",
						),
						response.New(
							models.Response{
								Status: "success",
								Data:   []string{"asd", "asd"},
								Errors: []int{1, 2},
							},
							"201",
							"Request Accepted",
						),
						response.New(
							models.Response{
								Status: "success",
								Data:   []bool{true, false},
								Errors: []string{"asd", "asf"},
							},
							"202",
							"Request Accepted",
						),
					},
					),
				),
			},
			file: "testdata/expected_output/composable-array-of-primitive.json",
		},
		{
			name: "test composable type with object type",
			endpoints: []*endpoint.EndPoint{
				endpoint.New(
					endpoint.PUT,
					"/composable-object",
					endpoint.WithTags("product"),
					endpoint.WithSuccessfulReturns([]response.Response{response.New(
						models.Response{
							Status: "success",
							Data:   models.PostBody{},
							Errors: models.UnsuccessfulResponse{},
						},
						"200",
						"Request Accepted",
					)},
					),
				),
			},
			file: "testdata/expected_output/composable-object.json",
		},

		{
			name: "test composable type with array of object type",
			endpoints: []*endpoint.EndPoint{
				endpoint.New(
					endpoint.PUT,
					"/composable-array-of-object",
					endpoint.WithTags("product"),
					endpoint.WithSuccessfulReturns([]response.Response{response.New(
						models.Response{
							Status: "success",
							Data:   []models.PostBody{},
							Errors: []models.UnsuccessfulResponse{},
						},
						"200",
						"Request Accepted",
					)},
					),
				),
			},
			file: "testdata/expected_output/composable-array-of-object.json",
		},
	}

	// Iterate through test cases./
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

			if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(Swagger{}, "endpoints"),
				cmpopts.IgnoreFields(definition.DefinitionProperties{}, "Example", "IsRequired"),
				cmpopts.SortSlices(func(a, b string) bool { return a < b }),
			); diff != "" {
				t.Errorf("JsonSwagger() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}
