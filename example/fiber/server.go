package main

import (
	"fmt"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/example/models"
	"github.com/go-swagno/swagno/http/response"
	"github.com/gofiber/fiber/v2"
)

func main() {
	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})
	desc := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id malesuada lorem, et fermentum sapien. Vivamus non pharetra risus, in efficitur leo. Suspendisse sed metus sit amet mi laoreet imperdiet. Donec aliquam eros eu blandit feugiat. Quisque scelerisque justo ac vehicula bibendum. Fusce suscipit arcu nisl, eu maximus odio consequat quis. Curabitur fermentum eleifend tellus, lobortis hendrerit velit varius vitae."

	endpoints := []*endpoint.EndPoint{
		endpoint.New(
			endpoint.WithMethod(endpoint.GET),
			endpoint.WithPath("/product"),
			endpoint.WithTags("product"),
			endpoint.WithSuccessfulReturns([]response.Info{response.New([]models.Product{}, "200", "Product List")}),
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
			endpoint.WithSuccessfulReturns([]response.Info{response.New(models.MapTest{
				"data": models.Product{},
			}, "200", "")}),
			endpoint.WithErrors([]response.Info{response.New(map[string]interface{}{
				"error":     "Not Authorized",
				"errorCode": 401,
			}, "401", "Not Authorized")}),
			endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		),
		endpoint.New(
			endpoint.WithMethod(endpoint.POST),
			endpoint.WithPath("/product"),
			endpoint.WithTags("product"),
			endpoint.WithBody(models.ProductPost{}),
			endpoint.WithSuccessfulReturns([]response.Info{response.New(models.Product{}, "201", "Created Product")}),
			endpoint.WithErrors([]response.Info{response.New([]models.ErrorResponse{}, "400", "")}),
			endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		),
	}

	sw.AddEndpoints(endpoints)

	app := fiber.New()
	swagger.SwaggerHandler(app, sw.GenerateDocs(), swagger.WithPrefix("/swagger"))

	fmt.Println(app.Listen(":8080"))
}
