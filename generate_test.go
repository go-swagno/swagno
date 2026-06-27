package swagno

import (
	"encoding/json"
	"os"
	"strings"
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

			if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(Swagger{}, "endpoints", "hidePackageName"), cmpopts.IgnoreFields(definition.DefinitionProperties{}, "Example", "IsRequired")); diff != "" {
				t.Errorf("JsonSwagger() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

// TestHidePackageName verifies that enabling Config.HidePackageName strips the
// package qualifier from definition keys and $ref values (e.g. "models.ProductPost"
// becomes "ProductPost"), while the default keeps the package-qualified names.
func TestHidePackageName(t *testing.T) {
	endpoints := []*endpoint.EndPoint{
		endpoint.New(
			endpoint.POST,
			"/product",
			endpoint.WithBody(models.ProductPost{}),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessfulResponse{}, "201", "Request Accepted")}),
			endpoint.WithErrors([]response.Response{response.New([]models.UnsuccessfulResponse{}, "400", "Bad Request")}),
		),
	}

	t.Run("hidden", func(t *testing.T) {
		sw := New(Config{Title: "Testing API", Version: "v1.0.0", HidePackageName: true})
		sw.AddEndpoints(endpoints)
		doc := string(sw.MustToJson())

		if strings.Contains(doc, "models.") {
			t.Errorf("expected no package qualifier in output, but found \"models.\":\n%s", doc)
		}
		for _, want := range []string{
			`"#/definitions/ProductPost"`,
			`"#/definitions/SuccessfulResponse"`,
			`"#/definitions/UnsuccessfulResponse"`,
		} {
			if !strings.Contains(doc, want) {
				t.Errorf("expected output to contain %s, but it did not:\n%s", want, doc)
			}
		}

		// definition keys must also be stripped so the refs resolve
		var parsed struct {
			Definitions map[string]json.RawMessage `json:"definitions"`
		}
		if err := json.Unmarshal([]byte(doc), &parsed); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"ProductPost", "SuccessfulResponse", "UnsuccessfulResponse"} {
			if _, ok := parsed.Definitions[key]; !ok {
				t.Errorf("expected definitions to contain key %q, got keys: %v", key, keysOf(parsed.Definitions))
			}
		}
	})

	t.Run("default keeps package name", func(t *testing.T) {
		sw := New(Config{Title: "Testing API", Version: "v1.0.0"})
		sw.AddEndpoints(endpoints)
		doc := string(sw.MustToJson())

		if !strings.Contains(doc, `"#/definitions/models.ProductPost"`) {
			t.Errorf("expected default output to keep package qualifier \"models.ProductPost\":\n%s", doc)
		}
	})
}

func keysOf(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
