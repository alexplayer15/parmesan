package test_builder

import (
	"encoding/json"
	"fmt"
	"strings"
)

type HTTPRequestBuilder struct {
	summary string
	method  string
	url     string
	headers map[string]string
	body    string
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		summary: "No Summary",
		method:  "GET",
		url:     "/",
		headers: map[string]string{},
		body:    "",
	}
}

func (b *HTTPRequestBuilder) WithSummary(summary string) *HTTPRequestBuilder {
	b.summary = summary
	return b
}

func (b *HTTPRequestBuilder) WithMethod(method string) *HTTPRequestBuilder {
	b.method = method
	return b
}

func (b *HTTPRequestBuilder) WithURL(url string) *HTTPRequestBuilder {
	b.url = url
	return b
}

func (b *HTTPRequestBuilder) WithHeader(key, value string) *HTTPRequestBuilder {
	b.headers[key] = value
	return b
}

func (b *HTTPRequestBuilder) WithBody(body string) *HTTPRequestBuilder {
	b.body = body
	return b
}

func (b *HTTPRequestBuilder) WithJSONBody(body interface{}) *HTTPRequestBuilder {
	bytes, _ := json.MarshalIndent(body, "", "  ")
	b.body = string(bytes)
	return b
}

func (b *HTTPRequestBuilder) Build() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("#### Summary: %s\n", b.summary))
	sb.WriteString(fmt.Sprintf("%s %s\n", b.method, b.url))
	for k, v := range b.headers {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	sb.WriteString("\n")
	sb.WriteString(b.body)
	sb.WriteString("\n\n")
	return sb.String()
}
