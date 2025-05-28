package hooks_logic

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexplayer15/parmesan/errors"
	"github.com/alexplayer15/parmesan/request_sender"
	"gopkg.in/yaml.v3"
)

type HooksFile []HookEntry
type HookEntry struct {
	Path   string         `yaml:"path"`
	Method string         `yaml:"method"`
	Body   map[string]any `yaml:"body"`
}

func (h HookEntry) IsEmpty() bool {
	return h.Path == "" && h.Method == "" && len(h.Body) == 0
}

func UnmarshalHooksFile(hooks string) (HooksFile, error) {
	hooksContent, err := os.ReadFile(hooks)
	if err != nil {
		return HooksFile{}, fmt.Errorf("failed to read hooks file %s", hooks)
	}

	ext := strings.TrimPrefix(filepath.Ext(hooks), ".")

	if ext != "yml" && ext != "yaml" {
		return HooksFile{}, fmt.Errorf("hooks file must be YAML, you entered a %s file", ext)
	}

	var hooksFile HooksFile

	if err := yaml.Unmarshal(hooksContent, &hooksFile); err != nil {
		return HooksFile{}, fmt.Errorf("invalid YAML: %w", err)
	}

	return hooksFile, nil
}

func TryAndFindHookForThisRequest(hooks HooksFile, req request_sender.Request) HookEntry {

	//URL has already been validated so no need to return an error
	parsedURL, _ := url.Parse(req.Url)

	for _, hook := range hooks {
		if hook.Path == parsedURL.Path && hook.Method == req.Method {
			return hook
		}
	}

	return HookEntry{}
}

func ModifyRequestBodyUsingHook(matchingHook HookEntry, requestBody string) (string, error) {
	var bodyMap map[string]any
	if err := json.Unmarshal([]byte(requestBody), &bodyMap); err != nil {
		return "", fmt.Errorf("failed to parse request body: %w", err)
	}

	for key, val := range matchingHook.Body {
		// Support nested keys separated by dots
		keys := strings.Split(key, ".")
		if err := updateField(bodyMap, keys, val, []string{}); err != nil {
			return "", fmt.Errorf("failed to apply hook for field '%s': %w", key, err)
		}
	}

	updatedBody, err := json.MarshalIndent(bodyMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to re-encode modified body: %w", err)
	}

	return string(updatedBody), nil
}

func updateField(bodyMap map[string]any, keys []string, newVal any, pathSoFar []string) error {
	if len(keys) == 0 {
		//warn user if hooks they have entered for this path are empty. to do: add in log level flag.
		return nil
	}

	currentKey := keys[0]
	pathSoFar = append(pathSoFar, currentKey)

	if len(keys) == 1 {
		existingVal, exists := bodyMap[currentKey]
		if !exists {
			return errors.NewMissingHookFieldError(strings.Join(pathSoFar, "."))
		}

		if err := validateHookTypesAgainstRequestSchema(existingVal, newVal); err != nil {
			return err
		}

		bodyMap[currentKey] = newVal
		return nil
	}

	next, exists := bodyMap[currentKey]
	if !exists {
		return errors.NewMissingHookFieldError(strings.Join(pathSoFar, "."))
	}

	switch typed := next.(type) {
	case map[string]any:
		return updateField(typed, keys[1:], newVal, pathSoFar)

	case []any:
		for i, item := range typed {
			if itemMap, ok := item.(map[string]any); ok {
				err := updateField(itemMap, keys[1:], newVal, append(pathSoFar, fmt.Sprintf("[%d]", i)))
				if err != nil {
					return err
				}
				typed[i] = itemMap
			}
		}
		bodyMap[currentKey] = typed
		return nil

	default:
		return errors.NewMissingHookFieldError(strings.Join(pathSoFar, "."))
	}
}

func validateHookTypesAgainstRequestSchema(original any, newVal any) error {
	switch original.(type) {
	case string:
		_, ok := newVal.(string)
		if !ok {
			return fmt.Errorf("type mismatch: expected string, got %T", newVal)
		}
	case int:
		_, ok := newVal.(float64)
		if !ok {
			return fmt.Errorf("type mismatch: expected number, got %T", newVal)
		}
	case bool:
		_, ok := newVal.(bool)
		if !ok {
			return fmt.Errorf("type mismatch: expected boolean, got %T", newVal)
		}
	case []any, map[string]any:
		return fmt.Errorf("modifying arrays or objects is not supported")
	default:
		return fmt.Errorf("unsupported target type %T", original)
	}

	return nil
}
