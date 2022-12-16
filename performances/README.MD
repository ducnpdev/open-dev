# performance
## Standard
- Test simple of function return error: https://opendev.hashnode.dev/golang-test-performance-function-standard-1
  - comparing 2 function:
  ```go
    func newErr() error {
        return errors.new("this is error")
    }
    func fmtErr() error {
        return fmt.Errorf("this is error")
    }
  ```
- Test convert string to int of 3 function: https://opendev.hashnode.dev/golang-test-performance-function-standard-1
    - comparing 3 function:
    - comparing 2 function:
    ```go
    func MethodInt(i int) string {
        return strconv.FormatInt(int64(i), 10)
    }
    func MethodItoa(i int) string {
        return strconv.Itoa(i)
    }
    func MethodFmt(i int) string {
        return fmt.Sprintf("%d", i)
    }
    ```