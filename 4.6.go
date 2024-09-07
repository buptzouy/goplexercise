package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// squashSpaces squashes each run of adjacent Unicode spaces in a UTF-8-encoded []byte slice into a single ASCII space.
func squashSpaces(data []byte) []byte {
	i := 0 // i is the index for the next character to write
	spaceFound := false

	for j := 0; j < len(data); {
		r, size := utf8.DecodeRune(data[j:]) // Decode the next rune from the byte slice

		if unicode.IsSpace(r) {
			if !spaceFound { // If it's the first space in a run
				data[i] = ' ' // Replace with a single ASCII space
				i++
				spaceFound = true
			}
		} else {
			// Not a space, so copy the rune to the current write position
			copy(data[i:], data[j:j+size])
			i += size
			spaceFound = false
		}
		j += size // Move to the next rune
	}

	return data[:i] // Slice the result up to the last written position
}

func main() {
	// Example byte slice with adjacent Unicode spaces
	data := []byte("Hello,\t世界!  This is a\ntest.  \u2003Another test.")

	// Squash adjacent spaces in-place
	result := squashSpaces(data)

	// Print the resulting slice
	fmt.Printf("%s\n", result)
}
