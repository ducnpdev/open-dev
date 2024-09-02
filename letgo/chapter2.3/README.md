# Lets Go

## Chapter 2.3: Web Application Basics
- Đầu tiên cần 1 func handler, để xử lý logic. Giống như MVC, handler dùng để xử lý logic, nhận request từ client và response data thông qua HTTP.

- Cần có 1 router, trong này mình dùng http.NewServeMux, sàu này có thể dùng Gin, Echo.

- Sau cùng là 1 máy chủ, mục đích để dợi nhận request từ client, sau này có thể dùng Nginx hoặc Apache.
```go
package main
import (
	"log"
	"net/http"
)
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from OpenDev"))
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
```

- Chú ý: khi start `go run main.go` web sẽ lắng nghe trên port 4000, mở web tại đường dẫn: http://localhost:4000