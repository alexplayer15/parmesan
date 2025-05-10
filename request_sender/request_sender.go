package request_sender

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/alexplayer15/parmesan/errors"
)

type Request struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    string
}

func ParseHttpRequestFile(httpRequestFile string) ([]Request, error) {
	var requests []Request

	if len(httpRequestFile) == 0 {
		return []Request{}, errors.ErrEmptyHTTPFile
	}

	lines := strings.Split(httpRequestFile, "\n")

	var blocks [][]string
	var currentBlock []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		//check if its the start of a new request
		if strings.HasPrefix(line, "####") {
			//to avoid appending empty block on first request
			if len(currentBlock) > 0 {
				blocks = append(blocks, currentBlock)
			}
			currentBlock = []string{}
			continue
		}

		//append empty lines to be able to differentiate between headers and body later
		currentBlock = append(currentBlock, line)

	}

	if len(currentBlock) > 0 {
		blocks = append(blocks, currentBlock)
	}

	for _, block := range blocks {
		if len(block) == 0 {
			return []Request{}, fmt.Errorf("empty request block")
		}
		req, err := getRequestInfo(block)
		if err != nil {
			return []Request{}, fmt.Errorf("failed to extract req from block: %w", err)
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func getRequestInfo(block []string) (Request, error) {
	var request Request
	method, err := getMethod(block)
	if err != nil {
		return Request{}, fmt.Errorf("failed to extract method: %w", err)
	}
	url, err := getURL(block)
	if err != nil {
		return Request{}, fmt.Errorf("failed to extract URL: %w", err)
	}
	headers, bodyStartingIndex, err := getHeaders(block)
	if err != nil {
		return Request{}, fmt.Errorf("failed to extract headers: %w", err)
	}
	body, err := getBody(block, bodyStartingIndex)
	if err != nil {
		return Request{}, fmt.Errorf("failed to extract body: %w", err)
	}
	request.Method = method
	request.Url = url
	request.Headers = headers
	request.Body = body

	return request, nil
}

func getMethod(block []string) (string, error) {
	firstLine := block[0]
	parts := strings.SplitN(firstLine, " ", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid request line: %q", firstLine)
	}

	method := parts[0]
	if err := ValidateHTTPMethod(method); err != nil {
		return "", err
	}

	return method, nil
}

func getURL(block []string) (string, error) {
	firstLine := block[0]
	parts := strings.SplitN(firstLine, " ", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid request line: %q", firstLine)
	}

	url := parts[1]
	if err := validateURL(url); err != nil {
		return "", err
	}

	return url, nil
}

func getHeaders(block []string) (map[string]string, int, error) {
	headers := make(map[string]string)

	for i := 1; i < len(block); i++ {
		line := block[i]

		if line == "" {
			return headers, i, nil
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			return map[string]string{}, 0, fmt.Errorf("header is in invalid format %q", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		headers[key] = value
	}

	return headers, len(block), nil
}

func getBody(block []string, bodyStartingIndex int) (string, error) {
	if bodyStartingIndex >= len(block) {
		return "", nil
	}

	var bodyLines []string

	for i := bodyStartingIndex; i < len(block); i++ {
		line := block[i]
		bodyLines = append(bodyLines, line)
	}

	body := strings.Join(bodyLines, "\n")

	return body, nil
}

func ValidateHTTPMethod(httpMethod string) error {
	upperHTTPMethod := strings.ToUpper(httpMethod)
	allowedMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "*"}

	if slices.Contains(allowedMethods, upperHTTPMethod) {
		return nil
	}

	return errors.NewValidationError(httpMethod)
}

func validateURL(rawURL string) error {
	allowedURLPrefixes := []string{"http://", "https://"}

	hasValidPrefix := false
	for _, prefix := range allowedURLPrefixes {
		if strings.HasPrefix(rawURL, prefix) {
			hasValidPrefix = true
			break
		}
	}

	if !hasValidPrefix {
		return errors.NewValidationError(rawURL)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return errors.NewValidationError(rawURL)
	}

	if parsedURL.Host == "" {
		return errors.NewValidationError(rawURL)
	}

	if parsedURL.Path == "" {
		return errors.NewValidationError(rawURL)
	}

	return nil
}

func SendHTTPRequest(req Request) (string, int, error) {
	request, err := http.NewRequest(req.Method, req.Url, bytes.NewBuffer([]byte(req.Body)))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", response.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return string(responseBody), response.StatusCode, nil
}
