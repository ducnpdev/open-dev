# Lets Go

## Chapter 2.4: Routing Requests
- như ví dụ chương 2.3, chúng ta chỉ có 1 route nhưng thực tế application sẽ không vậy, ta cần có nhiều route để làm những chức năng khác nhau:

| URL Pattern |  Handler  | Action |
|:-----|:--------:|------:|
| /   | home| màn hình home |
| /user  |  list user  |   danh sách user |
| /user/create   | create user |    tạo user |

- trong code ta sẽ thêm như thế này:
```go
// list user handler function.
func listuser(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("danh sach user OpenDev..."))
}
// create user handler function.
func createUser(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Create a new user OpenDev..."))
}
---trong func main

mux.HandleFunc("/user", listuser)
mux.HandleFunc("/user/create", createUser)
```

- run source code:
```console
go run main.go
```

- lúc này sẽ mở 2 đường dẫn
```console
http://localhost:4000/user
http://localhost:4000/user/create
```

### Fixed Path and Subtree Patterns
- hiện tại có 2 route đã được tạo và chạy.
- nhưng trong `golang` hoặc trong `servemux` có hỗ trợ 2 loại URL pattern: 
  - subtree paths: là route khớp với pattern ví du
    - /
    - /user/
      - có nghĩa là khi open web tất cả các route kiểu như `/user/*` sẽ match với URL này
  - fixed paths: là 2 route cố định ví dụ:
    - /user
    - /user/create
  
### Quản Lý URL
- trong một vài trường hợp khi `/` chúng ta không muốn hiển thị tất cả, lúc này ta cần xử lý, quản lý url một chút.
- và trả về `page not found` chẳng hạn
- code:
```go
func home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
		return
	}
    w.Write([]byte("Hello from Snippetbox"))
}
```