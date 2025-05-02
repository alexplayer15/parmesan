package request_generator

import (
	"fmt"
	"strings"
	"time"

	oas_struct "github.com/alexplayer15/parmesan/data"
)

func GenerateHttpRequest(oas oas_struct.OAS) (string, error) {
	if len(oas.Servers) == 0 {
		return "", fmt.Errorf("no server URL found in OAS")
	}

	serverURL := oas.Servers[0].URL
	if serverURL == "" {
		return "", fmt.Errorf("server URL is empty")
	}

	var httpRequests strings.Builder

	for path, methods := range oas.Paths {
		fullURL := joinURL(serverURL, path)
		generateRequestForPath(&httpRequests, fullURL, methods, path, oas)
	}

	return httpRequests.String(), nil
}

func joinURL(baseURL, path string) string {

	baseURL = strings.TrimSuffix(baseURL, "/")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return baseURL + path
}

func generateRequestForPath(builder *strings.Builder,
	fullURL string,
	methods map[string]oas_struct.Method,
	path string,
	oas oas_struct.OAS) {

	builder.WriteString(fmt.Sprintf("### Path: %s\n", path))

	for method, methodData := range methods {
		generateHttpRequestForMethod(builder, method, methodData, fullURL, oas)
	}
}

func generateHttpRequestForMethod(builder *strings.Builder,
	method string,
	methodData oas_struct.Method,
	fullURL string,
	oas oas_struct.OAS) {

	builder.WriteString(fmt.Sprintf("#### Summary: %s\n", methodData.Summary))
	builder.WriteString(fmt.Sprintf("%s %s\n", strings.ToUpper(method), fullURL))
	builder.WriteString(handleHeaders(methodData.Parameters))
	builder.WriteString("Content-Type: application/json\n\n")
	builder.WriteString(handleRequestBody(methodData.RequestBody, oas))
	builder.WriteString("\n\n")
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
		return formatExampleProperty(propName, prop.Example)
	}

	// Handle oneOf
	if len(prop.OneOf) > 0 {
		firstOption := prop.OneOf[0]
		// Pretend that the property is the first oneOf option
		return formatJsonProperty(propName, firstOption, oas)
	}

	// Handle $ref
	if prop.Ref != "" {
		referredSchema, err := resolveRef(prop.Ref, oas)
		if err == nil {
			if referredSchema.Example != nil {
				return formatExampleProperty(propName, referredSchema.Example)
			}

			switch referredSchema.Type {
			case "array":
				propFromSchema := oas_struct.Property{
					Type:  referredSchema.Type,
					Items: referredSchema.Items,
				}
				return formatArrayProperty(propName, propFromSchema, oas)

			case "object":
				objectBody := generateJsonBody(referredSchema, oas)
				objectBody = strings.TrimSpace(objectBody)
				indentedObjectBody := indentJson(objectBody, 2)
				return fmt.Sprintf("  \"%s\": %s,\n", propName, indentedObjectBody)

			default:
				if referredSchema.Default != nil {
					return fmt.Sprintf("  \"%s\": \"%v\",\n", propName, referredSchema.Default)
				}
				return getFallbackValue(propName, oas_struct.Property{
					Type: referredSchema.Type,
				})
			}
		}
	}

	if prop.Type == "array" && prop.Items != nil {
		return formatArrayProperty(propName, prop, oas)
	}

	return getFallbackValue(propName, prop)
}

func getFallbackValue(propName string, prop oas_struct.Property) string {

	switch prop.Type {
	case "string":
		switch prop.Format {
		case "date":
			return fmt.Sprintf("  \"%s\": \"2022-01-01\",\n", propName)
		case "date-time":
			return fmt.Sprintf("  \"%s\": \"2022-01-01T00:00:00Z\",\n", propName)
		default:
			return fmt.Sprintf("  \"%s\": \"example value\",\n", propName)
		}
	case "integer", "number":
		return fmt.Sprintf("  \"%s\": 0,\n", propName)
	case "boolean":
		return fmt.Sprintf("  \"%s\": false,\n", propName)
	case "object":
		return fmt.Sprintf("  \"%s\": {},\n", propName)
	default:
		// Unknown type fallback
		return fmt.Sprintf("  \"%s\": null,\n", propName)
	}
}

func formatExampleProperty(propName string, example any) string {
	switch v := example.(type) {
	case []interface{}:
		formattedItems := []string{}
		for _, item := range v {
			formattedItems = append(formattedItems, formatExampleValue(item))
		}
		return fmt.Sprintf("  \"%s\": [%s],\n", propName, strings.Join(formattedItems, ", "))
	default:
		return fmt.Sprintf("  \"%s\": %s,\n", propName, formatExampleValue(v))
	}
}

func formatExampleValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", val)
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	case time.Time:
		return fmt.Sprintf("\"%s\"", val.Format("2006-01-02"))
	case map[string]interface{}:
		objFields := []string{}
		for key, value := range val {
			objFields = append(objFields, fmt.Sprintf("\"%s\": \"%v\"", key, value))
		}
		return fmt.Sprintf("{%s}", strings.Join(objFields, ", "))
	default:
		return fmt.Sprintf("\"%v\"", val)
	}
}

func formatArrayProperty(propName string, prop oas_struct.Property, oas oas_struct.OAS) string {
	if prop.Items.Ref != "" {
		resolvedSchema, err := resolveRef(prop.Items.Ref, oas)
		if err == nil {
			// Try to use example from resolved schema if available
			if resolvedSchema.Example != nil {
				if exampleStr, ok := resolvedSchema.Example.(string); ok {
					return fmt.Sprintf("  \"%s\": [\"%s\"],\n", propName, exampleStr)
				}
			}
			// No direct example, generate a full object
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

	return fmt.Sprintf("  \"%s\": [],\n", propName)
}

func indentJson(json string, spaces int) string {
	lines := strings.Split(json, "\n")
	indent := strings.Repeat(" ", spaces)
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}
