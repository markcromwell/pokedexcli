// repl_test.go
package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "   single   ",
			expected: []string{"single"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "  multiple   spaces   here  ",
			expected: []string{"multiple", "spaces", "here"},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) { // Subtest for each case (nice for debugging)
			actual := cleanInput(c.input)

			// Check slice lengths match
			if len(actual) != len(c.expected) {
				t.Errorf("Expected %d words, got %d for input '%s'. Actual: %v", len(c.expected), len(actual), c.input, actual)
				return
			}

			// Check each word matches (lowercased, trimmed)
			for i := range actual {
				if actual[i] != c.expected[i] {
					t.Errorf("Word %d mismatch for input '%s': expected '%s', got '%s'", i, c.input, c.expected[i], actual[i])
				}
			}
		})
	}
}
