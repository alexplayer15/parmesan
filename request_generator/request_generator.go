package request_generator

import (
	"fmt"
	"strings"

	"github.com/alexplayer15/parmesan/data"
)

func GenerateHttpRequest(oas oas_struct.OAS) (string, error) {
	if len(oas.Servers) == 0 {
		return "", fmt.Errorf("no server URL found in OAS")
	}

	serverURL := oas.Servers[0].URL
	httpRequests := ""

	for path, methods := range oas.Paths {
		for method, methodData := range methods {
			httpRequest := fmt.Sprintf("### %s\n", methodData.Summary)
			httpRequest += fmt.Sprintf("%s %s%s\n", strings.ToUpper(method), serverURL, path)

			httpRequest += handleHeaders(methodData.Parameters)

			httpRequest += "Content-Type: application/json\n\n"

			httpRequest += handleRequestBody(methodData.RequestBody)

			httpRequests += httpRequest + "\n\n"
		}
	}

	return httpRequests, nil
}

func handleHeaders(parameters []oas_struct.Parameter) string {
	var headers string

	for _, param := range parameters {
		if param.In == "header" {
			// Default to "default-value" if no example is provided
			headerValue := param.Example
			if headerValue == "" {
				headerValue = "default-value"
			}
			headers += fmt.Sprintf("%s: %s\n", param.Name, headerValue)
		}
	}

	return headers
}

func handleRequestBody(requestBody oas_struct.RequestBody) string {
	var body string
	if content, ok := requestBody.Content["application/json"]; ok {

		if content.Schema.Type == "object" {
			body = "{\n"
			for propName, prop := range content.Schema.Properties {
				// Use the example value from the OAS spec
				exampleValue := "example value"
				if prop.Example != "" {
					exampleValue = prop.Example
				}
				body += fmt.Sprintf("  \"%s\": \"%s\",\n", propName, exampleValue)
			}

			if len(content.Schema.Properties) > 0 {
				body = body[:len(body)-2] + "\n"
			}
			body += "}\n"
		}
	}

	return body
}

