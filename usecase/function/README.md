# function
## variadic function
### khái niệm:
- variadic function là một function mà cho phép truyền vào bất kì một số nào
- trong golang, để có thể truyền variadic function ta có thể dùng `...`
- xem ví dụ:
```go
func total(numbers ...int) {
	fmt.Print(numbers, " ")
	total := 0
	for _, num := range numbers {
		total += num
	}
	fmt.Println(total)
}
```

### calling:
- để call func total ta dùng:
```
	total(1, 2, 3) => output: 6
	total(1) => output: 1
```
- hoặc nếu bạn muốn truyền vào func một slice,
```
numbers := []int{1, 2, 3, 4}
total(numbers...)
=> output: 10
```
- call variadic function với một tham số khác.
```
func prNumber(sep string, numbers ...int) {
	for i, number := range numbers {
		if i > 0 {
			fmt.Print(sep)
		}
		fmt.Print(number)
	}
}
prNumber(",", 11, 22, 33)
=> output: 11,22,33
```