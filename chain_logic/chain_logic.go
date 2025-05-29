package chain_logic

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexplayer15/parmesan/data"
	"github.com/alexplayer15/parmesan/request_sender"
	"github.com/stretchr/testify/assert/yaml"
)

func OrderRequests(requests []request_sender.Request, rules data.RuleSet) ([]request_sender.Request, error) {

	var orderedRequests []request_sender.Request

	for _, rule := range rules {
		request, err := findRequestAssociatedWithRule(requests, rule)
		if err != nil {
			return []request_sender.Request{}, err
		}

		orderedRequests = append(orderedRequests, request)

	}

	return orderedRequests, nil
}

func ApplyInjectionRules(request request_sender.Request, rules data.RuleSet) (request_sender.Request, error) {

	rule, err := findRuleAssociatedWithRequest(request, rules)
	if err != nil {
		return request_sender.Request{}, err
	}

	if rule.Inject == nil {
		return request_sender.Request{}, fmt.Errorf("you have not defined any injection rules for %T", request)
	}

	// if rule.Inject.Body != nil {
	// 	injectRequestBody(request, rule, rules)
	// }

	return request, nil

}

func ApplyExtractionRules(response any, rules data.RuleSet, request request_sender.Request) (map[string]any, error) {
	rule, err := findRuleAssociatedWithRequest(request, rules)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)

	if rule.Extract != nil {
		if rule.Extract.Body != nil {
			bodyValues, err := extractBody(response, rule)
			if err != nil {
				return nil, err
			}
			for k, v := range bodyValues {
				result[k] = v
			}
		}

		// You can implement extractHeaders similarly
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

func findRequestAssociatedWithRule(requests []request_sender.Request, rule data.Rule) (request_sender.Request, error) {

	for _, request := range requests {
		parsedUrl, _ := url.Parse(request.Url)
		if rule.Path == parsedUrl.Path && rule.Method == request.Method {
			return request, nil
		}
	}

	return request_sender.Request{}, fmt.Errorf("no request associated with the requests defined in the rules")
}

func findRuleAssociatedWithRequest(request request_sender.Request, rules data.RuleSet) (data.Rule, error) {

	for _, rule := range rules {
		parsedUrl, _ := url.Parse(request.Url)
		if rule.Path == parsedUrl.Path && rule.Method == request.Method {
			return rule, nil
		}
	}

	return data.Rule{}, fmt.Errorf("no rule associated with the requests defined in the rules")
}

// func applyBodyInjection(request request_sender.Request, rule data.Rule, rules data.RuleSet) (request_sender.Request, error) {
// 	rule.Inject.Body
// }

func extractBody(response any, rule data.Rule) (map[string]any, error) {
	result := make(map[string]any)

	body, ok := response.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("response is not a JSON object")
	}

	for _, item := range rule.Extract.Body {
		path := item.Path
		key := item.As

		// Support nested paths like "data.details.uri"
		value, err := extractJSONPath(body, path)
		if err != nil {
			return nil, fmt.Errorf("failed to extract path %s: %v", path, err)
		}

		result[key] = value
	}

	return result, nil
}

func extractJSONPath(data map[string]any, path string) (any, error) {
	parts := strings.Split(path, ".")
	var current any = data

	for _, part := range parts {
		obj, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("path '%s' does not point to a valid object", path)
		}
		val, exists := obj[part]
		if !exists {
			return nil, fmt.Errorf("path '%s' not found", path)
		}
		current = val
	}

	return current, nil
}
