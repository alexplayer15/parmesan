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
		httpRequests += generateRequestForPath(path, methods, serverURL)
	}

	return httpRequests, nil
}

func generateRequestForPath(path string, methods map[string]oas_struct.Method, serverURL string) string {
	var httpRequest string
	httpRequest += fmt.Sprintf("### Path: %s\n", path)

	for method, methodData := range methods {
		httpRequest += generateHttpRequestForMethod(method, methodData, path, serverURL)
	}

	return httpRequest
}

func generateHttpRequestForMethod(method string, methodData oas_struct.Method, path, serverURL string) string {
	httpRequest := fmt.Sprintf("#### Summary: %s\n", methodData.Summary)
	httpRequest += fmt.Sprintf("%s %s%s\n", strings.ToUpper(method), serverURL, path)

	httpRequest += handleHeaders(methodData.Parameters)
	httpRequest += "Content-Type: application/json\n\n"
	httpRequest += handleRequestBody(methodData.RequestBody)

	return httpRequest + "\n\n"
}

func handleHeaders(parameters []oas_struct.Parameter) string {
	var headers string
	for _, param := range parameters {
		if param.In == "header" {
			headers += formatHeader(param)
		}
	}
	return headers
}

func formatHeader(param oas_struct.Parameter) string {
	headerValue := param.Example
	if headerValue == "" {
		headerValue = "default-value"
	}
	return fmt.Sprintf("%s: %s\n", param.Name, headerValue)
}

func handleRequestBody(requestBody oas_struct.RequestBody) string {
	var body string
	if content, ok := requestBody.Content["application/json"]; ok {
		body = generateJsonBody(content)
	}
	return body
}

func generateJsonBody(content oas_struct.Content) string {
	var body string
	if content.Schema.Type == "object" {
		body = "{\n"
		for propName, prop := range content.Schema.Properties {
			body += formatJsonProperty(propName, prop)
		}

		if len(content.Schema.Properties) > 0 {
			body = body[:len(body)-2] + "\n"
		}
		body += "}\n"
	}
	return body
}

func formatJsonProperty(propName string, prop oas_struct.Property) string {
	exampleValue := "example value"
	if prop.Example != "" {
		exampleValue = prop.Example
	}
	return fmt.Sprintf("  \"%s\": \"%s\",\n", propName, exampleValue)
}

