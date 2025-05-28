package test_builder

import (
	"encoding/json"
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

func (b *HTTPRequestBuilder) WithJSONBody(body any) *HTTPRequestBuilder {
	bytes, _ := json.MarshalIndent(body, "", "  ")
	b.body = string(bytes)
	return b
}

func (b *HTTPRequestBuilder) Build() *HTTPRequestBuilder {
	return b
}
