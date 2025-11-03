package receipt

import (
	"testing"
)

func TestCleanJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain JSON without markdown",
			input:    `{"shop":"Test Store","isValid":true}`,
			expected: `{"shop":"Test Store","isValid":true}`,
		},
		{
			name: "JSON wrapped in markdown code blocks with json tag",
			input: "```json\n{\"shop\":\"Test Store\",\"isValid\":true}\n```",
			expected: `{"shop":"Test Store","isValid":true}`,
		},
		{
			name: "JSON wrapped in markdown code blocks without tag",
			input: "```\n{\"shop\":\"Test Store\",\"isValid\":true}\n```",
			expected: `{"shop":"Test Store","isValid":true}`,
		},
		{
			name: "JSON with leading and trailing whitespace",
			input: "  \n  {\"shop\":\"Test Store\",\"isValid\":true}  \n  ",
			expected: `{"shop":"Test Store","isValid":true}`,
		},
		{
			name: "multiline JSON in code blocks",
			input: "```json\n{\n  \"shop\": \"Test Store\",\n  \"isValid\": true\n}\n```",
			expected: "{\n  \"shop\": \"Test Store\",\n  \"isValid\": true\n}",
		},
		{
			name: "code block with extra whitespace",
			input: "  ```json  \n{\"shop\":\"Test Store\"}\n```  ",
			expected: `{"shop":"Test Store"}`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   \n\t  ",
			expected: "",
		},
		{
			name: "code block with no closing marker",
			input: "```json\n{\"shop\":\"Test Store\"}",
			expected: `{"shop":"Test Store"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanJSONResponse(tt.input)
			if result != tt.expected {
				t.Errorf("cleanJSONResponse() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanJSONResponse_RealWorldExamples(t *testing.T) {
	t.Run("typical OpenAI response with markdown", func(t *testing.T) {
		input := "```json\n{\n    \"shop\": \"Walmart\",\n    \"transactionDate\": \"2024-01-15 14:30\",\n    \"items\": [\n        {\n            \"name\": \"Milk\",\n            \"quantity\": 2,\n            \"price\": 350,\n            \"tax\": 0,\n            \"totalPrice\": 700\n        }\n    ],\n    \"tax\": 0,\n    \"total\": 700,\n    \"isValid\": true\n}\n```"

		result := cleanJSONResponse(input)

		// Should not contain markdown markers
		if containsBackticks(result) {
			t.Errorf("cleanJSONResponse() still contains backticks: %q", result)
		}

		// Should start with {
		if len(result) == 0 || result[0] != '{' {
			t.Errorf("cleanJSONResponse() does not start with '{': %q", result)
		}

		// Should end with }
		if len(result) == 0 || result[len(result)-1] != '}' {
			t.Errorf("cleanJSONResponse() does not end with '}': %q", result)
		}
	})

	t.Run("OpenAI response without markdown", func(t *testing.T) {
		input := `{
    "shop": "Target",
    "isValid": true
}`

		result := cleanJSONResponse(input)

		// Should preserve the JSON structure
		if result != input {
			t.Errorf("cleanJSONResponse() modified plain JSON: got %q, want %q", result, input)
		}
	})
}

func containsBackticks(s string) bool {
	for _, ch := range s {
		if ch == '`' {
			return true
		}
	}
	return false
}
