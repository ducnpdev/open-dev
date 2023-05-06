package main

import "fmt"

func total(numbers ...int) {
	fmt.Print(numbers, " ")
	total := 0
	for _, num := range numbers {
		total += num
	}
	fmt.Println(total)
}

func main() {
	numbers := []int{1, 2, 3, 4}
	total(numbers...)
	total(1, 2, 3)

	prNumber(",", 11, 22, 33)
}

func prNumber(sep string, numbers ...int) {
	for i, number := range numbers {
		if i > 0 {
			fmt.Print(sep)
		}
		fmt.Print(number)
	}
}
