package oas_struct

type Property struct {
	Type        string              `json:"type" yaml:"type"`
	Format      string              `json:"format,omitempty" yaml:"format,omitempty"`
	Description string              `json:"description" yaml:"description"`
	Example     any                 `json:"example" yaml:"example"`
	Items       *Schema             `json:"items,omitempty" yaml:"items,omitempty"`
	Ref         string              `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	OneOf       []Property          `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AllOf       []Property          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type Schema struct {
	Ref        string              `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Type       string              `json:"type,omitempty" yaml:"type,omitempty"`
	Properties map[string]Property `json:"properties,omitempty" yaml:"properties,omitempty"`
	Example    any                 `json:"example" yaml:"example"`
	Default    any                 `json:"default" yaml:"default"`
	Items      *Schema             `json:"items,omitempty" yaml:"items,omitempty"`
	AllOf      []Property          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
}

type Content struct {
	Schema Schema `json:"schema" yaml:"schema"`
}

type RequestBody struct {
	Content map[string]Content `json:"content" yaml:"content"`
}

type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
	Schema      Schema `json:"schema" yaml:"schema"`
	Example     string `json:"example" yaml:"example"`
}

type Method struct {
	Summary     string      `json:"summary" yaml:"summary"`
	Description string      `json:"description" yaml:"description"`
	Parameters  []Parameter `json:"parameters" yaml:"parameters"`
	RequestBody RequestBody `json:"requestBody" yaml:"requestBody"`
}

type Server struct {
	URL string `json:"url" yaml:"url"`
}

type Info struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version" yaml:"version"`
}

type Components struct {
	Schemas map[string]Schema `json:"schemas" yaml:"schemas"`
}

type OAS struct {
	OpenAPI    string                       `json:"openapi" yaml:"openapi"`
	Info       Info                         `json:"info" yaml:"info"`
	Servers    []Server                     `json:"servers" yaml:"servers"`
	Paths      map[string]map[string]Method `json:"paths" yaml:"paths"`
	Components Components                   `json:"components" yaml:"components"`
}
