package utils

import "strings"

// SmartTruncate truncates the input string to the specified limit,
// preferentially removing intermediate vowels before cutting the end.
// It adds the specified ellipsis if truncation occurs.
func SmartTruncate(s string, limit int, ellipsis string) string {
	if len(s) <= limit || limit <= len(ellipsis) {
		return s
	}

	// Preserve at least 3 characters at the start and
	if limit < 6+len(ellipsis) {
		return s[:limit-len(ellipsis)] + ellipsis
	}

	vowels := "aeiouAEIOU"
	result := []rune(s)
	vowelPositions := make([]int, 0)

	// Find positions of vowels, excluding first and last two characters
	for i := 2; i < len(result)-2; i++ {
		if strings.ContainsRune(vowels, result[i]) {
			vowelPositions = append(vowelPositions, i)
		}
	}

	// Remove vowels from the middle outwards
	for i := len(vowelPositions)/2 - 1; i >= 0; i-- {
		if len(result)-len(ellipsis) <= limit {
			break
		}
		result = append(result[:vowelPositions[i]], result[vowelPositions[i]+1:]...)
		for j := i + 1; j < len(vowelPositions); j++ {
			vowelPositions[j]--
		}
	}
	for i := len(vowelPositions) / 2; i < len(vowelPositions); i++ {
		if len(result)-len(ellipsis) <= limit {
			break
		}
		result = append(result[:vowelPositions[i]], result[vowelPositions[i]+1:]...)
		for j := i + 1; j < len(vowelPositions); j++ {
			vowelPositions[j]--
		}
	}

	// If still too long, truncate from the end
	if len(result) > limit {
		result = result[:limit-len(ellipsis)]
	}

	return string(result) + ellipsis
}
