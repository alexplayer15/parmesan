package test_builder

import (
	"fmt"
	"strings"
)

type HttpFileBuilder struct {
	httpRequests []*HTTPRequestBuilder
}

func NewHTTPFileBuilder() *HttpFileBuilder {
	return &HttpFileBuilder{
		httpRequests: []*HTTPRequestBuilder{},
	}
}

func (b *HttpFileBuilder) WithHTTPRequests(httpRequests []*HTTPRequestBuilder) *HttpFileBuilder {
	for _, req := range httpRequests {
		b.httpRequests = append(httpRequests, req)
	}
	return b
}

func (b *HttpFileBuilder) Build() string {
	var sb strings.Builder
	for _, req := range b.httpRequests {
		sb.WriteString(fmt.Sprintf("#### Summary: %s\n", req.summary))
		sb.WriteString(fmt.Sprintf("%s %s\n", req.method, req.url))
		for k, v := range req.headers {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
		sb.WriteString("\n")
		sb.WriteString(req.body)
		sb.WriteString("\n\n")
	}
	return sb.String()
}
