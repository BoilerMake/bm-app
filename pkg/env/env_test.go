package env

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "empty",
			input:    "",
			expected: map[string]string{},
		},
		{
			name:     "single commented line",
			input:    "# This is a comment",
			expected: map[string]string{},
		},
		{
			name:  "single valid line",
			input: "VALUE=key",
			expected: map[string]string{
				"VALUE": "key",
			},
		},
		{
			name:  "two valid lines",
			input: "VALUE=key\nVALUE2=key2",
			expected: map[string]string{
				"VALUE":  "key",
				"VALUE2": "key2",
			},
		},
		{
			name:  "valid comment empty",
			input: "VALUE=key\n# It's ya boy\nVALUE2=key2",
			expected: map[string]string{
				"VALUE":  "key",
				"VALUE2": "key2",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			result, err := parseFile(r)
			if err != nil {
				t.Errorf("Parsing failed: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}
