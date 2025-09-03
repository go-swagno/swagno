package main

import (
	"fmt"

	v3 "github.com/go-swagno/swagno/v3"
)

func main() {
	// External docs örneği
	externalDocs := v3.NewExternalDocs(
		"https://swagger.io",
		"Find more info here",
	)

	// OpenAPI config
	config := v3.Config{
		Title:        "Enhanced API",
		Version:      "1.0.0",
		Description:  "OpenAPI 3.0.3 uyumlu API",
		ExternalDocs: externalDocs,
	}

	// OpenAPI instance oluştur
	openapi := v3.New(config)

	// JSON çıktı
	json := openapi.ExportOpenAPIDocs("enhanced_api.json")
	fmt.Println("OpenAPI 3.0.3 uyumlu API oluşturuldu!")
	fmt.Println("ExternalDocs eklendi:", json)
}
