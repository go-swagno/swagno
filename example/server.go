package main

import (
	"fmt"
	"net/http"

	"github.com/domhoward14/swagno"
	"github.com/domhoward14/swagno/components/endpoint"
	"github.com/domhoward14/swagno/components/parameter"
	"github.com/domhoward14/swagno/example/models"
	"github.com/domhoward14/swagno/http/response"

	"github.com/go-swagno/swagno-http/swagger"
)

func main() {
	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})
	desc := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id malesuada lorem, et fermentum sapien. Vivamus non pharetra risus, in efficitur leo. Suspendisse sed metus sit amet mi laoreet imperdiet. Donec aliquam eros eu blandit feugiat. Quisque scelerisque justo ac vehicula bibendum. Fusce suscipit arcu nisl, eu maximus odio consequat quis. Curabitur fermentum eleifend tellus, lobortis hendrerit velit varius vitae."

	endpoints := []*endpoint.EndPoint{
		endpoint.New(
			endpoint.WithMethod(endpoint.GET),
			endpoint.WithPath("/product/page"),
			endpoint.WithTags("product"),
			endpoint.WithSuccessfulReturns([]response.Info{models.UnsuccessfulResponse{}}),
			endpoint.WithErrors([]response.Info{models.EmptySuccessfulResponse{}}),
			endpoint.WithDescription(desc),
			endpoint.WithProduce([]string{"application/json", "application/xml"}),
			endpoint.WithConsume([]string{"application/json"}),
		),
		endpoint.New(
			endpoint.WithMethod(endpoint.GET),
			endpoint.WithPath("/product"),
			endpoint.WithTags("product"),
			endpoint.WithParams([]parameter.Parameter{parameter.IntParam("id", parameter.WithRequired(true))}),
			endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
			endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
		),
		endpoint.New(
			endpoint.WithMethod(endpoint.GET),
			endpoint.WithPath("/product/{id}/detail"),
			endpoint.WithTags("product"),
			endpoint.WithParams([]parameter.Parameter{parameter.IntParam("id", parameter.WithRequired(true))}),
			endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
			endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
		),
		endpoint.New(
			endpoint.WithMethod(endpoint.POST),
			endpoint.WithPath("/product"),
			endpoint.WithTags("product"),
			endpoint.WithParams([]parameter.Parameter{}),
			endpoint.WithBody(models.ProductPost{}),
			endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
			endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
		),
	}

	// TODO make support for popular http server libs so that the creation of endpoints and handlers can happen in one funciton
	sw.AddEndpoints(endpoints)
	http.HandleFunc("/swagger/", swagger.SwaggerHandler(sw.GenerateDocs()))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
