package chain_logic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexplayer15/parmesan/data"
	"github.com/alexplayer15/parmesan/errors"
	"github.com/stretchr/testify/assert/yaml"
)

func OrderRequests(requests []data.Request, rules data.RuleSet) ([]data.Request, error) {

	var orderedRequests []data.Request

	for _, rule := range rules {
		request, err := findRequestAssociatedWithRule(requests, rule)
		if err != nil {
			return []data.Request{}, err
		}

		orderedRequests = append(orderedRequests, request)

	}

	return orderedRequests, nil
}

func ApplyInjectionRules(request data.Request, rules data.RuleSet, extractedValues map[string]any) (data.Request, error) {
	rule, err := findRuleAssociatedWithRequest(request, rules)
	if err != nil {
		return data.Request{}, err
	}

	if rule.Inject == nil {
		return data.Request{}, errors.MissingInjectionRuleError(request.Url)
	}

	if rule.Inject.Headers != nil {
		for _, header := range rule.Inject.Headers {
			val, ok := extractedValues[parseFromKey(header.From)]
			if !ok {
				return data.Request{}, errors.MissingHeaderValue(header.From)
			}
			request.Headers[header.Name] = fmt.Sprintf("%v", val)
		}
	}

	if rule.Inject.Body != nil {
		var bodyMap map[string]any
		if err := json.Unmarshal([]byte(request.Body), &bodyMap); err != nil || bodyMap == nil {
			bodyMap = make(map[string]any)
		}

		for _, field := range rule.Inject.Body {
			val, ok := extractedValues[parseFromKey(field.From)]
			if !ok {
				return data.Request{}, errors.MissingBodyValue(field.From)
			}

			switch field.Type {
			case "array":
				switch v := val.(type) {
				case []any:
				case []string:
					converted := make([]any, len(v))
					for i, s := range v {
						converted[i] = s
					}
					val = converted
				default:
					val = []any{val}
				}
			case "string":
				val = fmt.Sprintf("%v", val)
			}

			if err := setJSONPathValue(bodyMap, field.Path, val); err != nil {
				return data.Request{}, fmt.Errorf("injection failed at path %s: %w", field.Path, err)
			}
		}

		bodyBytes, err := json.Marshal(bodyMap)
		if err != nil {
			return data.Request{}, fmt.Errorf("failed to encode modified request body: %w", err)
		}
		request.Body = string(bodyBytes)
	}

	return request, nil
}

func ApplyExtractionRules(responseBody any, headers http.Header, rules data.RuleSet, request data.Request) (map[string]any, error) {
	rule, err := findRuleAssociatedWithRequest(request, rules)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)

	if rule.Extract == nil {
		return nil, fmt.Errorf("you have not defined any extraction rules for %T", request)
	}

	if rule.Extract.Body != nil {
		bodyValues, err := extractBody(responseBody, rule)
		if err != nil {
			return nil, err
		}
		for k, v := range bodyValues {
			result[k] = v
		}
	}

	if rule.Extract.Headers != nil {
		headerValues, err := extractHeaders(headers, rule)
		if err != nil {
			return nil, err
		}
		for k, v := range headerValues {
			result[k] = v
		}
	}

	return result, nil
}

func UnmarshalRulesFile(rulesFile string) (data.RuleSet, error) {

	rulesFileContent, err := os.ReadFile(rulesFile)
	if err != nil {
		return data.RuleSet{}, fmt.Errorf("error reading rules file %s", rulesFile)
	}

	ext := strings.TrimPrefix(filepath.Ext(rulesFile), ".")

	if ext != "yml" && ext != "yaml" {
		return data.RuleSet{}, fmt.Errorf("rules file must be YAML, you entered a %s file", ext)
	}

	var rules data.RuleSet

	if err := yaml.Unmarshal(rulesFileContent, &rules); err != nil {
		return data.RuleSet{}, fmt.Errorf("invalid YAML: %w", err)
	}

	return rules, nil
}

func findRequestAssociatedWithRule(requests []data.Request, rule data.Rule) (data.Request, error) {

	for _, request := range requests {
		parsedUrl, _ := url.Parse(request.Url)
		if rule.Path == parsedUrl.Path && rule.Method == request.Method {
			return request, nil
		}
	}

	return data.Request{}, errors.ErrNoMatchingRequestsInRulesFile
}

func findRuleAssociatedWithRequest(request data.Request, rules data.RuleSet) (data.Rule, error) {

	for _, rule := range rules {
		parsedUrl, _ := url.Parse(request.Url)
		if rule.Path == parsedUrl.Path && rule.Method == request.Method {
			return rule, nil
		}
	}

	return data.Rule{}, fmt.Errorf("no rule associated with the requests defined in the rules")
}

func extractBody(response any, rule data.Rule) (map[string]any, error) {
	result := make(map[string]any)

	var body map[string]any
	body, ok := response.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected JSON object at root, got %T", response)
	}

	for _, item := range rule.Extract.Body {
		path := item.Path
		key := item.As

		value, err := extractJSONPath(body, path)
		if err != nil {
			return nil, fmt.Errorf("failed to extract path %s: %v", path, err)
		}

		result[key] = value
	}

	return result, nil
}

var arrayIndexRegex = regexp.MustCompile(`^(\w+)\[(\d+)\]$`)

var wildcardRegex = regexp.MustCompile(`^(\w+)\[\*\]$`)

func extractJSONPath(data any, path string) (any, error) {
	parts := strings.Split(path, ".")
	var current []any
	current = []any{data}

	for _, part := range parts {
		var next []any

		for _, item := range current {
			// Handle wildcards like "flights[*]"
			if matches := wildcardRegex.FindStringSubmatch(part); len(matches) == 2 {
				field := matches[1]
				obj, ok := item.(map[string]any)
				if !ok {
					continue
				}
				list, ok := obj[field].([]any)
				if !ok {
					continue
				}
				next = append(next, list...)
				continue
			}

			// Handle array indexing like "flights[2]"
			if matches := arrayIndexRegex.FindStringSubmatch(part); len(matches) == 3 {
				field := matches[1]
				index, _ := strconv.Atoi(matches[2])
				obj, ok := item.(map[string]any)
				if !ok {
					continue
				}
				list, ok := obj[field].([]any)
				if !ok || index >= len(list) {
					continue
				}
				next = append(next, list[index])
				continue
			}

			// Normal key lookup
			obj, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if val, exists := obj[part]; exists {
				next = append(next, val)
			}
		}

		if len(next) == 0 {
			return nil, fmt.Errorf("path component '%s' matched no elements", part)
		}
		current = next
	}

	// Return first result if multiple
	return current[0], nil
}
func extractHeaders(headers http.Header, rule data.Rule) (map[string]any, error) {
	result := make(map[string]any)

	for _, headerRule := range rule.Extract.Headers {
		for key, vals := range headers {
			if strings.EqualFold(key, headerRule.Name) && len(vals) > 0 {
				result[headerRule.As] = vals[0] //only supports one value per header
				break
			}
		}
	}

	return result, nil
}

func parseFromKey(from string) string {
	parts := strings.Split(from, ".")
	if len(parts) == 2 {
		return parts[1]
	}
	return from
}

func setJSONPathValue(root map[string]any, path string, value any) error {
	parts := strings.Split(path, ".")
	current := root

	for i, part := range parts {
		// Last part is where we set the value
		isLast := i == len(parts)-1

		// Handle array syntax: e.g. fares[0]
		if matches := arrayIndexRegex.FindStringSubmatch(part); len(matches) == 3 {
			key := matches[1]
			index, _ := strconv.Atoi(matches[2])

			// Ensure key exists
			child, exists := current[key]
			if !exists {
				child = []any{}
			}

			// Ensure it's an array
			array, ok := child.([]any)
			if !ok {
				return fmt.Errorf("path '%s' is not an array", key)
			}

			// Extend array if needed
			for len(array) <= index {
				array = append(array, map[string]any{})
			}

			if isLast {
				array[index] = value
			} else {
				// Prepare next map level
				nextMap, ok := array[index].(map[string]any)
				if !ok {
					nextMap = map[string]any{}
					array[index] = nextMap
				}
				current[key] = array
				current = nextMap
			}
		} else {
			// Plain key
			if isLast {
				current[part] = value
			} else {
				child, exists := current[part]
				if !exists {
					child = map[string]any{}
					current[part] = child
				}

				nextMap, ok := child.(map[string]any)
				if !ok {
					return fmt.Errorf("path '%s' is not a valid object", part)
				}
				current = nextMap
			}
		}
	}

	return nil
}
