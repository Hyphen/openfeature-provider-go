package toggle

import "testing"

func TestGreeting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty name returns default greeting",
			input:    "",
			expected: "Hello, World!",
		},
		{
			name:     "non-empty name returns personalized greeting",
			input:    "Alice",
			expected: "Hello, Alice!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Greeting(tt.input)
			if result != tt.expected {
				t.Errorf("Greeting(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
