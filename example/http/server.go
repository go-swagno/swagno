package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/example/models"
	"github.com/go-swagno/swagno/http/response"

	"github.com/go-swagno/swagno-http/swagger"
)

type Response1 struct {
	Id            uint64                `json:"id"`
	Amount        float32               `json:"amount"`
	Amounts       []float32             `json:"amounts"`
	Token         string                `json:"-"`
	Email         string                `json:"email"`
	Phones        []Response1           `json:"phones"`
	CreatedBy     *uint64               `json:"created_by"`
	CreatedAt     *time.Time            `json:"created_at"`
	UpdatedAt     *time.Time            `json:"updated_at"`
	SomeBoolField *bool                 `json:"some_bool_field"`
	Test          []Response2           `json:"response2"`
	TestMap       map[string]*Response2 `json:"test_map"`
}

func (s Response1) GetDescription() string {
	return "Deneme"
}
func (s Response1) GetReturnCode() string {
	return "200"
}

type Response2 struct {
	Id int `json:"id" example:"123"`
}

func (s Response2) GetDescription() string {
	return "Test"
}
func (s Response2) GetReturnCode() string {
	return "201"
}

type Response3 []int

func (s Response3) GetDescription() string {
	return "Test"
}
func (s Response3) GetReturnCode() string {
	return "500"
}

type Response4 map[string]Response2

func (s Response4) GetDescription() string {
	return "Test"
}
func (s Response4) GetReturnCode() string {
	return "204"
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func main() {
	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})
	desc := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id malesuada lorem, et fermentum sapien. Vivamus non pharetra risus, in efficitur leo. Suspendisse sed metus sit amet mi laoreet imperdiet. Donec aliquam eros eu blandit feugiat. Quisque scelerisque justo ac vehicula bibendum. Fusce suscipit arcu nisl, eu maximus odio consequat quis. Curabitur fermentum eleifend tellus, lobortis hendrerit velit varius vitae."

	endpoints := []*endpoint.EndPoint{
		endpoint.New(
			endpoint.WithMethod(endpoint.GET),
			endpoint.WithPath("/product"),
			endpoint.WithTags("product"),
			endpoint.WithSuccessfulReturns([]response.Info{response.New(Response1{}, "200", "OKEY")}),
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
			endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}, response.New(Response4{
				"test": Response2{},
			}, "200", "")}),
			endpoint.WithErrors([]response.Info{response.New(map[string]interface{}{
				"error":     "Not Authorized",
				"errorCode": 401,
			}, "401", "")}),
			endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		),
		endpoint.New(
			endpoint.WithMethod(endpoint.POST),
			endpoint.WithPath("/product"),
			endpoint.WithTags("product"),
			endpoint.WithBody(models.ProductPost{}),
			endpoint.WithSuccessfulReturns([]response.Info{response.New(SuccessResponse{Data: Response1{}}, "200", "")}),
			endpoint.WithErrors([]response.Info{response.New([]Response2{}, "400", "")}),
			endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		),
	}

	sw.AddEndpoints(endpoints)
	http.HandleFunc("/swagger/", swagger.SwaggerHandler(sw.GenerateDocs()))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
