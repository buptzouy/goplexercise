package main

import "fmt"

// eliminateAdjDuplicates removes adjacent duplicates from a slice of strings in-place.
func eliminateAdjDuplicates(strings []string) []string {
	if len(strings) == 0 {
		return strings
	}

	// i is the index to place the next unique element.
	i := 0
	for j := 1; j < len(strings); j++ {
		if strings[i] != strings[j] {
			i++
			// Move the unique element to the next position.
			strings[i] = strings[j]
		}
	}
	// Slice up to i+1 to get the result without the adjacent duplicates.
	return strings[:i+1]
}

func main() {
	// Example input slice with adjacent duplicates
	strings := []string{"a", "a", "b", "b", "c", "a", "a", "d", "d", "d", "e"}

	// Remove adjacent duplicates in-place
	result := eliminateAdjDuplicates(strings)

	// Print the resulting slice
	fmt.Println(result)
}
