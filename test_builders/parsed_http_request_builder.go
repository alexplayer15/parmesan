package test_builder

import (
	"encoding/json"

	"github.com/alexplayer15/parmesan/data"
)

type ParsedHttpRequestBuilder struct {
	parsedHttpRequest data.Request
}

func NewParsedHttpRequest() *ParsedHttpRequestBuilder {
	return &ParsedHttpRequestBuilder{
		parsedHttpRequest: data.Request{
			Headers: make(map[string]string),
		},
	}
}

func (b *ParsedHttpRequestBuilder) WithMethod(method string) *ParsedHttpRequestBuilder {
	b.parsedHttpRequest.Method = method
	return b
}

func (b *ParsedHttpRequestBuilder) WithUrl(url string) *ParsedHttpRequestBuilder {
	b.parsedHttpRequest.Url = url
	return b
}

func (b *ParsedHttpRequestBuilder) WithHeader(key, value string) *ParsedHttpRequestBuilder {
	b.parsedHttpRequest.Headers[key] = value
	return b
}

func (b *ParsedHttpRequestBuilder) WithJSONBody(body any) *ParsedHttpRequestBuilder {
	bytes, _ := json.MarshalIndent(body, "", "  ")
	b.parsedHttpRequest.Body = string(bytes)
	return b
}

func (b *ParsedHttpRequestBuilder) Build() data.Request {
	return b.parsedHttpRequest
}
