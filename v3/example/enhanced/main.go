package main

import (
	"fmt"

	swagno3 "github.com/go-swagno/swagno/v3"
)

func main() {
	// External docs örneği
	externalDocs := swagno3.NewExternalDocs(
		"https://swagger.io",
		"Find more info here",
	)

	// OpenAPI config
	config := swagno3.Config{
		Title:        "Enhanced API",
		Version:      "1.0.0",
		Description:  "OpenAPI 3.0.3 uyumlu API",
		ExternalDocs: externalDocs,
	}

	// OpenAPI instance oluştur
	openapi := swagno3.New(config)

	// JSON çıktı
	json := openapi.ExportOpenAPIDocs("enhanced_api.json")
	fmt.Println("OpenAPI 3.0.3 compatible API created!")
	fmt.Println("ExternalDocs added:", json)
}
