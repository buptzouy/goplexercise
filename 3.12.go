package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// isAnagram checks if two strings are anagrams of each other.
func isAnagram(s1, s2 string) bool {
	// Remove non-letter characters and convert to lower case
	s1 = cleanString(s1)
	s2 = cleanString(s2)

	// Check if lengths are different
	if len(s1) != len(s2) {
		return false
	}

	// Convert strings to slices of runes for sorting
	s1Runes := []rune(s1)
	s2Runes := []rune(s2)

	// Sort runes
	sort.Slice(s1Runes, func(i, j int) bool {
		return s1Runes[i] < s1Runes[j]
	})
	sort.Slice(s2Runes, func(i, j int) bool {
		return s2Runes[i] < s2Runes[j]
	})

	// Convert sorted runes back to strings and compare
	return string(s1Runes) == string(s2Runes)
}

// cleanString removes non-letter characters and converts to lower case
func cleanString(s string) string {
	var cleanedString strings.Builder
	for _, ch := range s {
		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' {
			cleanedString.WriteRune(ch)
		}
	}
	return strings.ToLower(cleanedString.String())
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide two strings to compare.")
		return
	}
	s1 := os.Args[1]
	s2 := os.Args[2]

	if isAnagram(s1, s2) {
		fmt.Println("The strings are anagrams.")
	} else {
		fmt.Println("The strings are not anagrams.")
	}
}
