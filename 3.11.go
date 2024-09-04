package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma3(os.Args[i]))
	}
}

// comma inserts commas in a non-negative decimal integer string.
func comma3(s string) string {
	var sign string
	var decimalPart string
	if strings.HasPrefix(s, "-") {
		sign = "-"
		s = s[1:]
	}
	k := strings.Index(s, ".")
	if k >= 0 {
		decimalPart = s[k:]
		s = s[:k]
	}
	var buf bytes.Buffer
	n := len(s)

	for i, ch := range s {
		if (n-i)%3 == 0 && i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(byte(ch))
	}
	return sign + buf.String() + decimalPart
}
