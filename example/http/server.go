package main

import (
	"fmt"
	"net/http"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/example/models"

	"github.com/go-swagno/swagno-http/swagger"
)

func main() {
	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})
	desc := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id malesuada lorem, et fermentum sapien. Vivamus non pharetra risus, in efficitur leo. Suspendisse sed metus sit amet mi laoreet imperdiet. Donec aliquam eros eu blandit feugiat. Quisque scelerisque justo ac vehicula bibendum. Fusce suscipit arcu nisl, eu maximus odio consequat quis. Curabitur fermentum eleifend tellus, lobortis hendrerit velit varius vitae."

	endpoints := []*endpoint.EndPoint{
		endpoint.New(
			endpoint.GET,
			"/product",
			endpoint.WithTags("product"),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.Product{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")}),
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
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccessfulResponse{}, "200", "Request Accepted")}),
			endpoint.WithErrors([]response.Response{response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")}),
			endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		),
	}

	sw.AddEndpoints(endpoints)
	http.HandleFunc("/swagger/", swagger.SwaggerHandler(sw.MustToJson()))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
