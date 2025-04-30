package oas_struct

type Property struct {
    Type        string `json:"type" yaml:"type"`
    Description string `json:"description" yaml:"description"`
    Example     string `json:"example" yaml:"example"`
}

type Schema struct {
    Type       string              `json:"type" yaml:"type"`
    Properties map[string]Property `json:"properties" yaml:"properties"`
}

type Content struct {
    Schema Schema `json:"schema" yaml:"schema"`
}

type RequestBody struct {
    Content map[string]Content `json:"content" yaml:"content"`
}

type Method struct {
    Summary     string      `json:"summary" yaml:"summary"`
    Description string      `json:"description" yaml:"description"`
    RequestBody RequestBody `json:"requestBody" yaml:"requestBody"`
}

type Server struct {
    URL string `json:"url" yaml:"url"`
}

type Info struct {
	Title string 
	Description string 
	Version string
}

type OAS struct {
    OpenAPI string                         `json:"openapi" yaml:"openapi"`
    Info    Info
    Servers []Server
    Paths   map[string]map[string]Method    `json:"paths" yaml:"paths"`
}
