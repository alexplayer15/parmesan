package test_helpers

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExtractBody(jsonWithHeaders string) (string, error) {
	start := strings.Index(jsonWithHeaders, "{")
	if start == -1 {
		return "", fmt.Errorf("no JSON body found in the input")
	}
	return jsonWithHeaders[start:], nil
}

func ExtractHeaders(t *testing.T, result string) (map[string]string, error) {
	t.Helper()
	parts := strings.SplitN(result, "\n\n", 2)
	require.Len(t, parts, 2, "expected result to have headers and body separated by \\n\\n")

	headers := parts[0]

	headerLines := strings.Split(headers, "\n")
	headerMap := make(map[string]string)
	for _, line := range headerLines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headerMap[key] = value
		}
	}

	return headerMap, nil
}

func getFieldFromJSON(jsonString string, fieldName string) (any, error) {
	var resultBody map[string]any
	err := json.Unmarshal([]byte(jsonString), &resultBody)
	fmt.Println("Career Object:", resultBody)

	fieldValue := resultBody[fieldName]

	return fieldValue, err
}

func AssertJSONHasArrayWithObject(t *testing.T, jsonString string, fieldName string, requiredKeys []string) {
	t.Helper()

	fieldValue, err := getFieldFromJSON(jsonString, fieldName)
	assert.NoError(t, err, "Should unmarshal successfully")

	arr, ok := fieldValue.([]any)
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

	fieldValue, err := getFieldFromJSON(jsonString, fieldName)
	assert.NoError(t, err, "Should unmarshal successfully")

	arr := fieldValue.([]any)

	firstItem := arr[0].(map[string]any)

	for key, expectedValue := range expectedExamples {
		actualValue, exists := firstItem[key]
		assert.True(t, exists, "expected key %s to exist in object inside %s", key, fieldName)
		assert.Equal(t, expectedValue, actualValue, "expected value for key %s in object inside %s to be %v, got %v", key, fieldName, expectedValue, actualValue)
	}
}

func AssertJSONHasObject(t *testing.T, jsonString string, fieldName string, requiredKeys []string) {
	t.Helper()

	fieldValue, err := getFieldFromJSON(jsonString, fieldName)
	assert.NoError(t, err, "Should unmarshal successfully")

	obj, ok := fieldValue.(map[string]any)
	assert.True(t, ok, "field %s should be an object", fieldName)
	assert.GreaterOrEqual(t, len(obj), 1, "object %s should have at least one item", fieldName)

	for _, key := range requiredKeys {
		_, exists := obj[key]
		assert.True(t, exists, "object in array %s should contain key %s", fieldName, key)
	}
}

func AssertJSONExamplesForObject(t *testing.T, jsonString, fieldName string, expectedExamples map[string]any) {
	t.Helper()

	fieldValue, err := getFieldFromJSON(jsonString, fieldName)
	assert.NoError(t, err, "Should unmarshal successfully")

	obj := fieldValue.(map[string]any)

	for key, expectedValue := range expectedExamples {
		actualValue, exists := obj[key]
		assert.True(t, exists, "expected key %s to exist in object inside %s", key, fieldName)
		assert.Equal(t, expectedValue, actualValue, "expected value for key %s in object inside %s to be %v, got %v", key, fieldName, expectedValue, actualValue)
	}
}

func AssertJSONHasXAmountOfArrays(t *testing.T, jsonString string, expectedArrayAmount int) {
	t.Helper()

	var resultBody map[string]any
	err := json.Unmarshal([]byte(jsonString), &resultBody)
	assert.NoError(t, err, "Should unmarshal successfully")

	arrayCount := 0

	for _, v := range resultBody {
		if _, ok := v.([]any); ok {
			arrayCount++
		}
	}

	assert.Equal(t, expectedArrayAmount, arrayCount, "There should be exactly %d array(s) in the response body", expectedArrayAmount)
}

func AssertJSONHasXAmountOfObjects(t *testing.T, jsonString string, expectedObjAmount int) {
	t.Helper()

	var resultBody map[string]any
	err := json.Unmarshal([]byte(jsonString), &resultBody)
	assert.NoError(t, err, "Should unmarshal successfully")

	objCount := 0

	for _, v := range resultBody {
		if _, ok := v.(map[string]any); ok {
			objCount++
		}
	}

	assert.Equal(t, expectedObjAmount, objCount, "There should be exactly %d object(s) in the response body", expectedObjAmount)
}
