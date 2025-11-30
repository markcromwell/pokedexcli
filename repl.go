package main

import "strings"

// cleanInput is a stub for input cleaning logic.
func cleanInput(text string) []string {
	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)
	// Lowercase the whole thing
	text = strings.ToLower(text)
	// Split on any whitespace (spaces, tabs, etc.) and ignore empty fields
	words := strings.Fields(text)
	return words
}
