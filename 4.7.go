package main

import (
	"fmt"
	"unicode/utf8"
)

// reverseUTF8 反转 UTF-8 编码的 []byte 切片中的字符。
func reverseUTF8(data []byte) {
	// 第一步：反转整个切片，使所有字节的顺序反转。
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	// 第二步：在反转后的切片中，反转每个 rune 内部的字节顺序，以校正它们的内部顺序。
	for i := 0; i < len(data); {
		if data[i] < utf8.RuneSelf { // ASCII 字符（1 字节 rune）
			i++
			continue
		}

		// 解码 rune 以确定其长度
		_, size := utf8.DecodeRune(data[i:])
		// 反转多字节 rune 的字节
		reverseBytes(data[i : i+size])
		// 移动到下一个 rune
		i += size
	}
}

// reverseBytes 反转字节切片的内容。
func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func main() {
	// 示例 UTF-8 编码字符串（字节切片形式）
	data := []byte("Hello, 世界!") // 包含英文和中文字符的字符串

	// 原地反转 UTF-8 编码的字节切片
	reverseUTF8(data)

	// 输出反转后的结果
	fmt.Println(string(data))
}
