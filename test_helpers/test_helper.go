package test_helpers

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExtractBody(jsonWithHeaders string) (string, error) {
	start := strings.Index(jsonWithHeaders, "{")
	if start == -1 {
		return "", fmt.Errorf("no JSON body found in the input")
	}
	return jsonWithHeaders[start:], nil
}

func AssertJSONHasArrayWithObject(t *testing.T, jsonString string, fieldName string, requiredKeys []string) {
	t.Helper()

	var resultBody map[string]any
	err := json.Unmarshal([]byte(jsonString), &resultBody)
	assert.NoError(t, err, "should unmarshal JSON successfully")

	value, ok := resultBody[fieldName]
	assert.True(t, ok, "field %s should exist", fieldName)

	arr, ok := value.([]any)
	assert.True(t, ok, "field %s should be an array", fieldName)
	assert.GreaterOrEqual(t, len(arr), 1, "array %s should have at least one item", fieldName)

	firstItem, ok := arr[0].(map[string]any)
	assert.True(t, ok, "first item of array %s should be an object", fieldName)

	for _, key := range requiredKeys {
		_, exists := firstItem[key]
		assert.True(t, exists, "object in array %s should contain key %s", fieldName, key)
	}
}
