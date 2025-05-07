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
	body, err := handleRequestBody(methodData.RequestBody, oas, fullURL, method)
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

func handleRequestBody(requestBody oas_struct.RequestBody, oas oas_struct.OAS, path string, method string) (string, error) {
	content, ok := requestBody.Content["application/json"]
	if !ok {
		log.Printf("[WARNING] %s %s has no 'application/json' request body. Skipping body generation.", method, path)
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
		formattedJsonProperty, err := generateJsonFromProperty(propName, prop, oas)
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
		var props map[string]oas_struct.Property
		if item.Ref != "" {
			resolved, err := resolveRef(item.Ref, oas)
			if err != nil {
				return combined, err
			}

			expanded, err := expandAllOfSchema(resolved, oas)
			if err != nil {
				return combined, err
			}
			props = expanded.Properties
		} else {
			props = item.Properties
		}

		for name, prop := range props {
			combined.Properties[name] = prop
		}
	}

	return combined, nil
}

func generateJsonFromProperty(propName string, prop oas_struct.Property, oas oas_struct.OAS) (string, error) {
	if prop.Example != nil {
		return formatProperty(propName, prop.Example), nil
	}

	resolvedSchema, err := resolveProperty(prop, oas)
	if err != nil {
		return "", err
	}

	if resolvedSchema.Example != nil {
		return formatProperty(propName, resolvedSchema.Example), nil
	}

	switch resolvedSchema.Type {
	case "object":
		body, err := generateJsonFromSchema(resolvedSchema, oas)
		if err != nil {
			return "", err
		}
		return formatObjectProperty(propName, body), nil

	case "array":
		arrayProp := oas_struct.Property{
			Type:  resolvedSchema.Type,
			Items: resolvedSchema.Items,
		}
		return generateJsonFromArray(propName, arrayProp, oas)
	default:
		if resolvedSchema.Default != nil {
			return formatProperty(propName, resolvedSchema.Default), nil
		}
		return getFallbackValue(propName, prop), nil
	}
}

func resolveProperty(prop oas_struct.Property, oas oas_struct.OAS) (oas_struct.Schema, error) {
	if prop.Ref != "" {
		return resolveRef(prop.Ref, oas)
	}

	if len(prop.OneOf) > 0 {
		selected := prop.OneOf[0]
		if selected.Ref != "" {
			return resolveRef(selected.Ref, oas)
		}
		return selected, nil
	}

	if len(prop.AllOf) > 0 {
		expanded, err := expandAllOfProperty(prop, oas)
		if err != nil {
			return oas_struct.Schema{}, err
		}
		return expanded, nil
	}

	return oas_struct.Schema{
		Type:       prop.Type,
		Properties: prop.Properties,
		Items:      prop.Items,
		Default:    prop.Default,
	}, nil
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

func formatProperty(propName string, value any) string {
	return fmt.Sprintf("  \"%s\": %s,\n", propName, formatPropertyValue(value))
}

func formatPropertyValue(v any) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", val)
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	case time.Time:
		return fmt.Sprintf("\"%s\"", val.Format("2006-01-02"))
	case map[string]any:
		objFields := []string{}
		for key, value := range val {
			objFields = append(objFields, fmt.Sprintf("\"%s\": \"%v\"", key, value))
		}
		return fmt.Sprintf("{%s}", strings.Join(objFields, ", "))
	case []any:
		formattedItems := []string{}
		for _, item := range val {
			formattedItems = append(formattedItems, formatPropertyValue(item))
		}
		return fmt.Sprintf("[%s]", strings.Join(formattedItems, ", "))
	default:
		return fmt.Sprintf("\"%v\"", val)
	}
}

func generateJsonFromArray(propName string, prop oas_struct.Property, oas oas_struct.OAS) (string, error) {
	itemSchema := prop.Items
	if itemSchema == nil {
		return fmt.Sprintf("  \"%s\": [],\n", propName), nil
	}

	resolvedItem, err := resolveSchema(*itemSchema, oas)
	if err != nil {
		return "", err
	}

	if resolvedItem.Example != nil {
		return fmt.Sprintf("  \"%s\": [%s],\n", propName, formatPropertyValue(resolvedItem.Example)), nil
	}

	body, err := generateJsonFromSchema(resolvedItem, oas)
	if err != nil {
		return "", err
	}

	indentedBody := indentJson(body, 4, false)
	return fmt.Sprintf("  \"%s\": [\n%s\n  ],\n", propName, indentedBody), nil
}

func formatObjectProperty(propName string, body string) string {
	return fmt.Sprintf("  \"%s\": %s,\n", propName, indentJson(body, 2, true))
}

func indentJson(json string, spaces int, skipFirstLineIndent bool) string {
	json = strings.TrimSpace(json)

	lines := strings.Split(json, "\n")
	if len(lines) == 0 {
		return json
	}

	indent := strings.Repeat(" ", spaces)
	var builder strings.Builder

	for i, line := range lines {
		if skipFirstLineIndent && i == 0 {
			builder.WriteString(line + "\n")
			continue
		}
		builder.WriteString(indent + line + "\n")
	}

	return strings.TrimSuffix(builder.String(), "\n")
}

func joinURL(baseURL, path string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	path = "/" + strings.TrimLeft(path, "/")
	return baseURL + path
}
