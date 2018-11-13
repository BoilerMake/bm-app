package env

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	tt := []struct {
		name      string
		input     string
		expected  map[string]string
		shouldErr bool
	}{
		{
			name:      "empty",
			input:     "",
			expected:  map[string]string{},
			shouldErr: false,
		},
		{
			name:      "single commented line",
			input:     "# This is a comment",
			expected:  map[string]string{},
			shouldErr: false,
		},
		{
			name:  "single valid line",
			input: "VALUE=key",
			expected: map[string]string{
				"VALUE": "key",
			},
			shouldErr: false,
		},
		{
			name:  "two valid lines",
			input: "VALUE=key\nVALUE2=key2",
			expected: map[string]string{
				"VALUE":  "key",
				"VALUE2": "key2",
			},
			shouldErr: false,
		},
		{
			name:  "valid comment empty",
			input: "VALUE=key\n# It's ya boy\nVALUE2=key2",
			expected: map[string]string{
				"VALUE":  "key",
				"VALUE2": "key2",
			},
			shouldErr: false,
		},
		{
			name:      "no assignment",
			input:     "VALUE=",
			expected:  map[string]string{},
			shouldErr: false,
		},
		{
			name:  "equals in value",
			input: "VALUE=URL:==with--characters++=like=equals==",
			expected: map[string]string{
				"VALUE": "URL:==with--characters++=like=equals==",
			},
			shouldErr: false,
		},
		{
			name:      "invalid line",
			input:     "THING",
			expected:  map[string]string{},
			shouldErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			result, err := parseFile(r)
			if err != nil && !tc.shouldErr {
				t.Errorf("Parsing failed: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}
