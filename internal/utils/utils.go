package utils

import (
	"regexp"
	"strings"
)

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

func StripTerminalReturns(s string) string {
	// ansi has a few escape codes that we need to remove because otherwise shiz messed up in the UI
	// \x1b[2K - clear line
	// \x1b[0G - move to beginning of line
	// \x1b[1G - move to beginning of line

	ansiEscape := regexp.MustCompile(`\x1b\[(2K|0G|1G)`) // Regex to remove ANSI escape codes
	controlChars := regexp.MustCompile(`[\r\\b]`)        // Regex to remove control characters

	s = ansiEscape.ReplaceAllString(s, "")
	s = controlChars.ReplaceAllString(s, "")

	return s
}
