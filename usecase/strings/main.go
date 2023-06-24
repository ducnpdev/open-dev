package main

import (
	"fmt"
	"strings"
)

func main() {
	strs := []string{"b111231sdfsdfsaf21312sdafsda", "b111231sdfsdfsaf21312sdafsda", "b111231sdfsdfsaf21312sdafsda", "b111231sdfsdfsaf21312sdafsda"}
	fmt.Println(Concat(strs))
	fmt.Println(Concatv2(strs))
	fmt.Println(Concatv3(strs))
	fmt.Println(Concatv4(strs))
}

func Concatv4(values []string) string {
	return values[0] + values[1] + values[2] + values[3]
}

func Concat(values []string) string {
	s := ""
	for _, value := range values {
		s += value
	}
	return s
}

func Concatv2(values []string) string {
	sb := strings.Builder{}
	for _, value := range values {
		_, _ = sb.WriteString(value)
	}
	return sb.String()
}

func Concatv3(values []string) string {
	total := 0
	for i := 0; i < len(values); i++ {
		total += len(values[i])
	}
	sb := strings.Builder{}
	sb.Grow(total)
	for _, value := range values {
		_, _ = sb.WriteString(value)
	}
	return sb.String()
}

// sdf
// sdf
func Trim() {
	// TrimRight/TrimLeft
	fmt.Println(strings.TrimRight("123xxoxo", "xo")) // 123
	fmt.Println(strings.TrimLeft("xxoxo123", "xo"))  // 123

	// TrimSuffix/TrimPrefix
	fmt.Println(strings.TrimSuffix("123xxoxoxo", "xo")) // 123xxoxo
	fmt.Println(strings.TrimPrefix("xoxoxo123", "xo"))  // xoxo123
}
