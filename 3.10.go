package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma(os.Args[i]))
	}
}

// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	var buf bytes.Buffer
	n := len(s)
	for i, ch := range s {
		if (n-i)%3 == 0 && i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(byte(ch))
	}
	return buf.String()
}
