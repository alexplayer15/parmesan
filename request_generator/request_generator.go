package request_generator

import (
	"fmt"
	"log"
	"strings"
	"time"

	oas_struct "github.com/alexplayer15/parmesan/data"
)

func GenerateHttpRequest(oas oas_struct.OAS) (string, error) {
	serverURL := oas.Servers[0].URL
	if serverURL == "" {
		return "", fmt.Errorf("server URL is empty")
	}

	var httpRequests strings.Builder

	for path, methods := range oas.Paths {
		fullURL := joinURL(serverURL, path)
		err := generateRequestForPath(&httpRequests, fullURL, methods, oas)
		if err != nil {
			return "", fmt.Errorf("failed to generate request for path %s: %w", path, err)
		}
	}

	return httpRequests.String(), nil
}

func generateRequestForPath(builder *strings.Builder, fullURL string, methods map[string]oas_struct.Method, oas oas_struct.OAS) error {
	for method, methodData := range methods {
		err := generateHttpRequestForMethod(builder, method, methodData, fullURL, oas)
		if err != nil {
			return fmt.Errorf("failed to generate HTTP request for method %s: %w", method, err)
		}
	}
	return nil
}

func generateHttpRequestForMethod(builder *strings.Builder, method string, methodData oas_struct.Method, fullURL string, oas oas_struct.OAS) error {
	body, err := handleRequestBody(methodData.RequestBody, oas, fullURL)
	if err != nil {
		return fmt.Errorf("failed to handle request body: %w", err)
	}

	builder.WriteString(fmt.Sprintf("#### Summary: %s\n", methodData.Summary))
	builder.WriteString(fmt.Sprintf("%s %s\n", strings.ToUpper(method), fullURL))
	builder.WriteString(handleHeaders(methodData.Parameters))
	builder.WriteString("Content-Type: application/json\n\n")
	builder.WriteString(body)
	builder.WriteString("\n\n")

	return nil
}

func handleHeaders(parameters []oas_struct.Parameter) string {
	var builder strings.Builder
	for _, param := range parameters {
		if param.In != "header" {
			continue
		}
		headerValue := param.Example
		if headerValue == "" {
			headerValue = "default-value"
		}
		fmt.Fprintf(&builder, "%s: %s\n", param.Name, headerValue)
	}
	return builder.String()
}

func handleRequestBody(requestBody oas_struct.RequestBody, oas oas_struct.OAS, path string) (string, error) {
	content, ok := requestBody.Content["application/json"]
	if !ok {
		log.Printf("[WARNING] Path '%s' has no 'application/json' request body. Skipping body generation.", path)
		return "", nil
	}

	schema, err := resolveSchema(content.Schema, oas)
	if err != nil {
		return "", fmt.Errorf("failed to resolve schema: %w", err)
	}

	body, err := generateJsonFromSchema(schema, oas)
	if err != nil {
		return "", err
	}
	return body, nil
}

func resolveSchema(schema oas_struct.Schema, oas oas_struct.OAS) (oas_struct.Schema, error) {
	if schema.Ref == "" {
		return schema, nil
	}
	return resolveRef(schema.Ref, oas)
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

func generateJsonFromSchema(schema oas_struct.Schema, oas oas_struct.OAS) (string, error) {
	if len(schema.OneOf) > 0 {
		schema = schema.OneOf[0]
	}

	if len(schema.AllOf) > 0 {
		expanded, err := expandAllOfSchema(schema, oas)
		if err == nil {
			schema = expanded
		}
	}

	if schema.Type != "object" && len(schema.Properties) == 0 {
		return "", nil //come back to handle errors properly
	}

	var builder strings.Builder
	builder.WriteString("{\n")
	for propName, prop := range schema.Properties {
		formattedJsonProperty, err := formatJsonProperty(propName, prop, oas)
		if err != nil {
			return "", err
		}
		builder.WriteString(formattedJsonProperty)
	}

	if len(schema.Properties) > 0 {
		objectStr := builder.String()
		objectStr = objectStr[:len(objectStr)-2] + "\n"
		builder.Reset()
		builder.WriteString(objectStr)
	}

	builder.WriteString("}")
	return builder.String(), nil
}

func expandAllOfSchema(schema oas_struct.Schema, oas oas_struct.OAS) (oas_struct.Schema, error) {
	combined := oas_struct.Schema{
		Type:       "object",
		Properties: make(map[string]oas_struct.Property),
	}

	for name, prop := range schema.Properties {
		combined.Properties[name] = prop
	}

	for _, item := range schema.AllOf {
		if item.Ref != "" {
			resolved, err := resolveRef(item.Ref, oas)
			if err != nil {
				return combined, err
			}

			expanded, err := expandAllOfSchema(resolved, oas)
			if err != nil {
				return combined, err
			}
			for name, prop := range expanded.Properties {
				combined.Properties[name] = prop
			}
		} else {

			for name, prop := range item.Properties {
				combined.Properties[name] = prop
			}
		}
	}

	return combined, nil
}

func formatJsonProperty(propName string, prop oas_struct.Property, oas oas_struct.OAS) (string, error) {
	if prop.Example != nil {
		return formatPropertyExample(propName, prop.Example), nil
	}

	if len(prop.AllOf) > 0 {
		expandedSchema, err := expandAllOfProperty(prop, oas)
		if err != nil {
			return "", err
		}
		objectBody, err := generateJsonFromSchema(expandedSchema, oas)
		if err != nil {
			return "", err
		}
		objectBody = strings.TrimSpace(objectBody)
		indentedObjectBody := indentJson(objectBody, 2)
		return fmt.Sprintf("  \"%s\": %s,\n", propName, indentedObjectBody), nil
	}

	if len(prop.OneOf) > 0 {
		selected := prop.OneOf[0]

		if selected.Ref != "" {
			if resolved, err := resolveRef(selected.Ref, oas); err == nil {
				selected = resolved
			}
		}

		objectBody, err := generateJsonFromSchema(selected, oas)
		if err != nil {
			return "", err
		}
		indented := indentJson(strings.TrimSpace(objectBody), 2)

		return fmt.Sprintf("\"%s\": %s,\n", propName, indented), nil
	}

	if prop.Ref != "" {
		referredSchema, err := resolveRef(prop.Ref, oas)
		if err != nil {
			return "", err
		}

		if referredSchema.Example != nil {
			return formatPropertyExample(propName, referredSchema.Example), nil
		}

		switch referredSchema.Type {
		case "array":
			propFromSchema := oas_struct.Property{
				Type:  referredSchema.Type,
				Items: referredSchema.Items,
			}
			formattedArrayProperty, err := formatArrayProperty(propName, propFromSchema, oas)
			if err != nil {
				return "", err
			}

			return formattedArrayProperty, nil
		case "object":
			objectBody, err := generateJsonFromSchema(referredSchema, oas)
			if err != nil {
				return "", err
			}
			objectBody = strings.TrimSpace(objectBody)
			indentedObjectBody := indentJson(objectBody, 2)
			return fmt.Sprintf("  \"%s\": %s,\n", propName, indentedObjectBody), nil

		default:
			if referredSchema.Default != nil {
				return fmt.Sprintf("  \"%s\": \"%v\",\n", propName, referredSchema.Default), nil
			}
			return getFallbackValue(propName, oas_struct.Property{
				Type: referredSchema.Type,
			}), nil
		}
	}

	if prop.Type == "array" && prop.Items != nil {
		formattedArrayProperty, err := formatArrayProperty(propName, prop, oas)
		if err != nil {
			return "", err
		}
		return formattedArrayProperty, nil
	}

	if prop.Type == "object" && len(prop.Properties) > 0 {
		nestedBody, err := generateObjectFromProperties(prop.Properties, oas)
		if err != nil {
			return "", err
		}
		nestedBody = strings.TrimSpace(nestedBody)
		indented := indentJson(nestedBody, 2)
		return fmt.Sprintf("  \"%s\":%s,\n", propName, indented), nil
	}

	return getFallbackValue(propName, prop), nil
}

func expandAllOfProperty(prop oas_struct.Property, oas oas_struct.OAS) (oas_struct.Schema, error) {
	combined := oas_struct.Schema{
		Type:       "object",
		Properties: make(map[string]oas_struct.Property),
	}

	for name, p := range prop.Properties {
		combined.Properties[name] = p
	}

	for _, item := range prop.AllOf {
		if item.Ref != "" {
			resolved, err := resolveRef(item.Ref, oas)
			if err != nil {
				return combined, err
			}
			expanded, err := expandAllOfSchema(resolved, oas)
			if err != nil {
				return combined, err
			}
			for name, p := range expanded.Properties {
				combined.Properties[name] = p
			}
		} else {
			for name, p := range item.Properties {
				combined.Properties[name] = p
			}
		}
	}

	return combined, nil
}

func generateObjectFromProperties(properties map[string]oas_struct.Property, oas oas_struct.OAS) (string, error) {
	var builder strings.Builder
	builder.WriteString("{\n")

	for name, prop := range properties {
		formattedJsonProperty, err := formatJsonProperty(name, prop, oas)
		if err != nil {
			return "", err
		}
		builder.WriteString(formattedJsonProperty)
	}
	if len(properties) > 0 {
		str := builder.String()
		str = str[:len(str)-2] + "\n"
		builder.Reset()
		builder.WriteString(str)
	}
	builder.WriteString("}")
	return builder.String(), nil
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

func formatPropertyExample(propName string, example any) string {
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

func formatArrayProperty(propName string, prop oas_struct.Property, oas oas_struct.OAS) (string, error) {

	if prop.Items.Ref != "" {
		resolvedSchema, err := resolveRef(prop.Items.Ref, oas)
		if err != nil {
			return "", err
		}
		if resolvedSchema.Example != nil {
			if exampleStr, ok := resolvedSchema.Example.(string); ok {
				return fmt.Sprintf("  \"%s\": [\"%s\"],\n", propName, exampleStr), nil
			}
		}
		if len(resolvedSchema.AllOf) > 0 {
			expanded, err := expandAllOfSchema(resolvedSchema, oas)
			if err == nil {
				resolvedSchema = expanded
			}
		}
		objectBody, err := generateJsonFromSchema(resolvedSchema, oas)
		if err != nil {
			return "", err
		}
		indentedObjectBody := indentJson(objectBody, 4)
		return fmt.Sprintf("  \"%s\": [\n%s\n  ],\n", propName, indentedObjectBody), nil
	} else if prop.Items.Type == "object" {
		objectBody, err := generateJsonFromSchema(*prop.Items, oas)
		if err != nil {
			return "", err
		}
		indentedObjectBody := indentJson(objectBody, 4)
		return fmt.Sprintf("  \"%s\": [\n%s\n  ],\n", propName, indentedObjectBody), nil
	}

	return fmt.Sprintf("  \"%s\": [],\n", propName), nil
}

func indentJson(json string, spaces int) string {
	lines := strings.Split(json, "\n")
	indent := strings.Repeat(" ", spaces)
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}

func joinURL(baseURL, path string) string {

	baseURL = strings.TrimSuffix(baseURL, "/")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return baseURL + path
}
