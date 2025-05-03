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

func AssertJSONExamplesForObjectsInAnArray(t *testing.T, jsonString, fieldName string, expectedExamples map[string]any) {
	t.Helper()

	var resultBody map[string]any
	err := json.Unmarshal([]byte(jsonString), &resultBody)
	assert.NoError(t, err, "should unmarshal JSON successfully")

	arrayField := resultBody[fieldName]
	arr := arrayField.([]any)

	firstItem := arr[0].(map[string]any)

	for key, expectedValue := range expectedExamples {
		actualValue, exists := firstItem[key]
		assert.True(t, exists, "expected key %s to exist in object inside %s", key, fieldName)
		assert.Equal(t, expectedValue, actualValue, "expected value for key %s in object inside %s to be %v, got %v", key, fieldName, expectedValue, actualValue)
	}
}

func AssertJSONHasObject(t *testing.T, jsonString string, fieldName string, requiredKeys []string) {
	t.Helper()

	var resultBody map[string]any
	err := json.Unmarshal([]byte(jsonString), &resultBody)
	assert.NoError(t, err, "should unmarshal JSON successfully")

	value, ok := resultBody[fieldName]
	assert.True(t, ok, "field %s should exist", fieldName)

	obj, ok := value.(map[string]any)
	assert.True(t, ok, "field %s should be an object", fieldName)
	assert.GreaterOrEqual(t, len(obj), 1, "object %s should have at least one item", fieldName)

	for _, key := range requiredKeys {
		_, exists := obj[key]
		assert.True(t, exists, "object in array %s should contain key %s", fieldName, key)
	}
}

func AssertJSONExamplesForObject(t *testing.T, jsonString, fieldName string, expectedExamples map[string]any) {
	t.Helper()

	var resultBody map[string]any
	err := json.Unmarshal([]byte(jsonString), &resultBody)
	assert.NoError(t, err, "should unmarshal JSON successfully")

	objField := resultBody[fieldName]
	obj := objField.(map[string]any)

	for key, expectedValue := range expectedExamples {
		actualValue, exists := obj[key]
		assert.True(t, exists, "expected key %s to exist in object inside %s", key, fieldName)
		assert.Equal(t, expectedValue, actualValue, "expected value for key %s in object inside %s to be %v, got %v", key, fieldName, expectedValue, actualValue)
	}
}
