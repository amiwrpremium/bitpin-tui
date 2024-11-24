package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func StringEndsWith(s string, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

func FormatWithCommas(input string) string {
	// Parse the string into an integer
	num, err := strconv.Atoi(input)
	if err != nil {
		return input
	}

	// Format the number with commas
	formatted := fmt.Sprintf("%d", num)

	// Reverse the string for easier grouping
	reversed := reverseString(formatted)

	// Add commas every three digits
	var result strings.Builder
	for i, char := range reversed {
		if i > 0 && i%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(char)
	}

	// Reverse it back to get the final result
	return reverseString(result.String())
}

// Helper function to reverse a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
