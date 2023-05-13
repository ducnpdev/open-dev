# context package
### khá niệm
1. Các khái niệm cơ bản.
- hiện tại trong một chương trình thì sẽ đối mặt với nhiều các vấn đề như timeouts, cancelltion hoặc lỗi từ một hệ thống khác. Thì để hỗ trợ và giải quyết vấn đề này, đa số những developer sẽ dùng `channel` cách này hoạt động tốt nhưng về lâu dài sẽ có một số hạn chế.
- Để giải quyết những vấn đề này, trong version `Go 1.7` , context package đã được đưa vào `standard library`
- thì trong `context` sẽ cố một số func sau:
```go
func Background() Context //  (1)
func TODO() Context  // (2)
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) func WithValue(parent Context, key, val interface{}) Context
```
  - (1): Background đơn giản chỉ trả về một context rỗng
  - (2): TODO thường ít dùng trong chương trình, nhưng nếu trong một func cần truyền vào context mà không có context cha truyền vào thì dùng TODO, nó cũng trả về một context rỗng.
  => Thường thì Background và TODO sẽ nằm tại top của một luồng api hoặc một follow
2. Trong `context` interface sẽ có:
```go
type Context interface {
// Deadline returns the time when work done on behalf of this
// context should be canceled. Deadline returns ok==false when no
// deadline is set. Successive calls to Deadline return the same // results.
Deadline() (deadline time.Time, ok bool)
// Done returns a channel that's closed when work done on behalf // of this context should be canceled. Done may return nil if this // context can never be canceled. Successive calls to Done return // the same value.
Done() <-chan struct{}
// Err returns a non-nil error value after Done is closed. Err
// returns Canceled if the context was canceled or
// DeadlineExceeded if the context's deadline passed. No other
// values for Err are defined.  After Done is closed, successive
// calls to Err return the same value.
Err() error
// Value returns the value associated with this context for key, // or nil if no value is associated with key. Successive calls to // Value with the same key returns the same result.
Value(key interface{}) interface{}
}
```

### Một số trường hợp cần dùng
#### cancellation
1. vấn đề
  - Một goroutine cha muốn cancel.
  - Một goroutine có thể muốn cancel func children. 
  - Có một số vấn đề như bị blocking cũng cần cancel.
2. Hướng xử lý:
  - Để giải quyết vấn đề thì context package cung cấp 3 func:
```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) // (1)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) // (2)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) // (3)
```
  - (1): sẽ return một context mới, nó sẽ được close khi func `CancelFunc` được gọi .
  - (2): sẽ return một context mới, nó sẽ được close khi một đồng hồ thời gian được truyền vào trong deadline.
  - (3): sẽ return một context mới, nó sẽ được close sau một khoảng thời gian truyền vào.
3. Xử lý:
  - Nếu trong follow func phía sau cần cancel với bất kì lý do nào đó thì chỉ cần truyền context vào từng func. Hoặc nếu không sử dụng thì đó chỉ cần truyền vào và không xử lý gì cả.
  - Code ví dụ:
```go

```