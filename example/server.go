package main

import (
	"fmt"
	"net/http"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/example/models"
	"github.com/go-swagno/swagno/http/response"

	"github.com/go-swagno/swagno-http/swagger"
)

func main() {
	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})
	desc := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id malesuada lorem, et fermentum sapien. Vivamus non pharetra risus, in efficitur leo. Suspendisse sed metus sit amet mi laoreet imperdiet. Donec aliquam eros eu blandit feugiat. Quisque scelerisque justo ac vehicula bibendum. Fusce suscipit arcu nisl, eu maximus odio consequat quis. Curabitur fermentum eleifend tellus, lobortis hendrerit velit varius vitae."

	endpoints := []*endpoint.EndPoint{
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
	}

	// TODO make support for popular http server libs so that the creation of endpoints and handlers can happen in one funciton
	sw.AddEndpoints(endpoints)
	http.HandleFunc("/swagger/", swagger.SwaggerHandler(sw.GenerateDocs()))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
