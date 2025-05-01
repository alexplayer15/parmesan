package request_generator

import (
	"fmt"
	"strings"

	oas_struct "github.com/alexplayer15/parmesan/data"
)

func GenerateHttpRequest(oas oas_struct.OAS) (string, error) {
	if len(oas.Servers) == 0 {
		return "", fmt.Errorf("no server URL found in OAS")
	}

	serverURL := oas.Servers[0].URL
	httpRequests := ""

	for path, methods := range oas.Paths {
		httpRequests += generateRequestForPath(path, methods, serverURL, oas)
	}

	return httpRequests, nil
}

func generateRequestForPath(path string,
	methods map[string]oas_struct.Method,
	serverURL string,
	oas oas_struct.OAS) string {

	var httpRequest string
	httpRequest += fmt.Sprintf("### Path: %s\n", path)

	for method, methodData := range methods {
		httpRequest += generateHttpRequestForMethod(method, methodData, path, serverURL, oas)
	}

	return httpRequest
}

func generateHttpRequestForMethod(method string,
	methodData oas_struct.Method,
	path,
	serverURL string,
	oas oas_struct.OAS) string {

	httpRequest := fmt.Sprintf("#### Summary: %s\n", methodData.Summary)
	httpRequest += fmt.Sprintf("%s %s%s\n", strings.ToUpper(method), serverURL, path)

	httpRequest += handleHeaders(methodData.Parameters)
	httpRequest += "Content-Type: application/json\n\n"
	httpRequest += handleRequestBody(methodData.RequestBody, oas)

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

func handleRequestBody(requestBody oas_struct.RequestBody, oas oas_struct.OAS) string {
	var body string
	if content, ok := requestBody.Content["application/json"]; ok {
		schema := content.Schema

		// If schema is a $ref, resolve it
		if schema.Ref != "" {
			resolvedSchema, err := resolveRef(schema.Ref, oas)
			if err != nil {

				return ""
			}
			schema = resolvedSchema
		}

		body = generateJsonBody(schema, oas)
	}
	return body
}

func resolveRef(ref string, oas oas_struct.OAS) (oas_struct.Schema, error) {
	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		return oas_struct.Schema{}, fmt.Errorf("unsupported ref format: %s", ref)
	}

	name := strings.TrimPrefix(ref, prefix)
	schema, ok := oas.Components.Schemas[name]
	if !ok {
		return oas_struct.Schema{}, fmt.Errorf("schema not found: %s", name)
	}
	return schema, nil
}

func generateJsonBody(schema oas_struct.Schema, oas oas_struct.OAS) string {
	var body string
	if schema.Type == "object" {
		body = "{\n"
		for propName, prop := range schema.Properties {
			body += formatJsonProperty(propName, prop, oas)
		}

		if len(schema.Properties) > 0 {
			body = body[:len(body)-2] + "\n"
		}
		body += "}\n"
	}
	return body
}

func formatJsonProperty(propName string, prop oas_struct.Property, oas oas_struct.OAS) string {
	if prop.Example != nil {
		switch v := prop.Example.(type) {
		case string:
			return fmt.Sprintf("  \"%s\": \"%s\",\n", propName, v)
		case int:
			return fmt.Sprintf("  \"%s\": %d,\n", propName, v)
		case []interface{}:
			formattedItems := []string{}
			for _, item := range v {
				switch val := item.(type) {
				case string:
					formattedItems = append(formattedItems, fmt.Sprintf("\"%s\"", val))
				case int, int64, float64:
					formattedItems = append(formattedItems, fmt.Sprintf("%v", val))
				case map[string]interface{}:
					objFields := []string{}
					for key, value := range val {
						objFields = append(objFields, fmt.Sprintf("\"%s\": \"%v\"", key, value))
					}
					formattedItems = append(formattedItems, fmt.Sprintf("{%s}", strings.Join(objFields, ", ")))
				default:
					formattedItems = append(formattedItems, fmt.Sprintf("\"%v\"", val))
				}
			}
			return fmt.Sprintf("  \"%s\": [%s],\n", propName, strings.Join(formattedItems, ", "))
		default:
			return fmt.Sprintf("  \"%s\": \"%v\",\n", propName, v)
		}
	}

	if prop.Type == "array" && prop.Items != nil {
		if prop.Items.Ref != "" {

			resolvedSchema, err := resolveRef(prop.Items.Ref, oas)
			if err == nil {
				objectBody := generateJsonBody(resolvedSchema, oas)
				objectBody = strings.TrimSpace(objectBody)

				indentedObjectBody := indentJson(objectBody, 4)
				return fmt.Sprintf("  \"%s\": [\n%s\n  ],\n", propName, indentedObjectBody)
			}
		} else if prop.Items.Type == "object" {

			objectBody := generateJsonBody(*prop.Items, oas)
			objectBody = strings.TrimSpace(objectBody)

			indentedObjectBody := indentJson(objectBody, 4)
			return fmt.Sprintf("  \"%s\": [\n%s\n  ],\n", propName, indentedObjectBody)
		}
		// Just empty array if no further details
		return fmt.Sprintf("  \"%s\": [],\n", propName)
	}

	// Final fallback
	return fmt.Sprintf("  \"%s\": \"example value\",\n", propName)
}

func indentJson(json string, spaces int) string {
	lines := strings.Split(json, "\n")
	indent := strings.Repeat(" ", spaces)
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}
