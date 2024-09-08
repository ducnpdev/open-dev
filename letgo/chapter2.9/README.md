# Lets Go

## Chapter 2.9: Serving Static Files
- đôi lúc sẽ có 1 vài yêu cầu đó là load file static trong project để hiện lên web.
- các file dạng như html, image, css, ...

### File Server
- trong golang package net/http sẽ có thể load file trực tiếp thông qua `http.FileServer`
- bây giờ ta sẽ load file với route bắt đầu bắt /static/*
- structure folder sẽ như này
```
root/
│
├── css/
│   ├── main.css
│
├── img/
│   └── code.png
│
├── js/
│   ├── main.js
```
- router:

| URL Pattern |  Handler  | Action |
|:-----|:--------:|------:|
| ANY /static/   | http.FileServer |   load file static |

- để tạo route load static file:
```go
fileServer := http.FileServer(http.Dir("./ui/static"))
```

### Code

```go
package main

import (
	"log"
	"net/http"
)

func fileStatis(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("danh sach user OpenDev..."))
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
```