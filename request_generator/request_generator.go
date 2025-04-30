package request_generator

import (
	"fmt"

	"github.com/alexplayer15/parmesan/data"
)

func GenerateHttpRequest(oas oas_struct.OAS, path string) (string, error) {
	if len(oas.Servers) == 0 {
		return "", fmt.Errorf("no server URL found in OAS")
	}

	serverURL := oas.Servers[0].URL
	methodData := oas.Paths[path]["get"]
	requestBody := methodData.RequestBody

	httpRequest := fmt.Sprintf("### %s\n", methodData.Summary)
	httpRequest += fmt.Sprintf("GET %s%s\n", serverURL, path)
	httpRequest += "Content-Type: application/json\n\n"

	if content, ok := requestBody.Content["application/json"]; ok {
		if content.Schema.Type == "object" {
			httpRequest += "{\n"
			for propName, prop := range content.Schema.Properties {
				exampleValue := "example value"
				if prop.Example != "" {
					exampleValue = prop.Example
				}
				httpRequest += fmt.Sprintf("  \"%s\": \"%s\",\n", propName, exampleValue)
			}
			// Remove last comma
			if len(content.Schema.Properties) > 0 {
				httpRequest = httpRequest[:len(httpRequest)-2] + "\n"
			}
			httpRequest += "}\n"
		}
	}
	return httpRequest, nil
}