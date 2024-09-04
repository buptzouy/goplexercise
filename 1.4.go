// Dup2 prints the count and text of lines that appear more than once
// in the input. It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	fileMap := make(map[string][]string)
	files := os.Args[1:]

	if len(files) == 0 {
		countLines(os.Stdin, counts, fileMap)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, fileMap)
			f.Close()
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
			fmt.Printf("Files: %v\n", fileMap[line])
		}
	}
}

// countLines reads lines from a file and records the counts of each line
// and the files in which they occur.
func countLines(f *os.File, counts map[string]int, fileMap map[string][]string) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		counts[line]++
		// 记录文件名到 fileMap 中，避免重复添加相同文件名
		if !contains(fileMap[line], f.Name()) {
			fileMap[line] = append(fileMap[line], f.Name())
		}
	}
	// NOTE: ignoring potential errors from input.Err()
}

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
