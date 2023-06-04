# keep-alive in http
* `keep-alive` cơ bản là giữ kết nối giữa client và server.
-----
**Lợi ích**
* Giảm kết nối TCP giữa client và server, mỗi client chỉ có duy nhất một kết nối.
* Giảm độ trễ – Việc giảm số lần bắt tay ba bước có thể giúp cải thiện độ trễ của trang web. Điều này đặc biệt đúng với các kết nối SSL/TLS, vốn yêu cầu các chuyến đi khứ hồi bổ sung để mã hóa và xác minh các kết nối.
-----
## Code mẫu client và server
### Ví dụ 1
**Server**
```
package main
import (
	"fmt"
	"net/http"
)
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("receive a request from:", r.RemoteAddr, r.Header)
	w.Write([]byte("ok"))
}
func main() {
	fmt.Println("server keep alive")
	var s = http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(Index),
	}
	s.ListenAndServe()
}
```
**Client**
```
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("client keep alive")

	c := &http.Client{}
	req, err := http.NewRequest("Get", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", *req)

	for i := 0; i < 3; i++ {
		resp, err := c.Do(req)
		if err != nil {
			fmt.Println("http get error:", err)
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read body error:", err)
			return
		}
		fmt.Println("response body:", string(b))
		time.Sleep(time.Second)
	}
}
```
**Kết Quả**
* Server:
```
server keep alive
receive a request from: [::1]:51328 map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/1.1]]
receive a request from: [::1]:51328 map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/1.1]]
receive a request from: [::1]:51328 map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/1.1]]
```
* Client:
```
client keep alive
http.Request{Method:"Get", URL:(*url.URL)(0x1400017a000), Proto:"HTTP/1.1", ProtoMajor:1, ProtoMinor:1, Header:http.Header{}, Body:io.ReadCloser(nil), GetBody:(func() (io.ReadCloser, error))(nil), ContentLength:0, TransferEncoding:[]string(nil), Close:false, Host:"localhost:8080", Form:url.Values(nil), PostForm:url.Values(nil), MultipartForm:(*multipart.Form)(nil), Trailer:http.Header(nil), RemoteAddr:"", RequestURI:"", TLS:(*tls.ConnectionState)(nil), Cancel:(<-chan struct {})(nil), Response:(*http.Response)(nil), ctx:(*context.emptyCtx)(0x1400012e008)}
response body: ok
response body: ok
response body: ok
```
* Có thể thấy, Client đã gọi 3 request và điều chung 1 id `51328`, đồng nghĩa là đã được reuse connection

### Ví dụ 2:
**Client**
* Trong phần tạo ra client, cấu hình disable keep alive:
```
c := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
```
**Kết Quả**
* Server:
```
server keep alive
receive a request from: [::1]:51869 map[Accept-Encoding:[gzip] Connection:[close] User-Agent:[Go-http-client/1.1]]
receive a request from: [::1]:51872 map[Accept-Encoding:[gzip] Connection:[close] User-Agent:[Go-http-client/1.1]]
receive a request from: [::1]:51881 map[Accept-Encoding:[gzip] Connection:[close] User-Agent:[Go-http-client/1.1]]
```
* Client: kêt quả không có gì thay đổi.
* Bây giờ có thể thấy log ở server, mỗi request sẽ có 1 id (51869,51872,51881) khác nhau, đồng nghĩa là connect không được dùng lại.

### Ví dụ 3
* Bây giờ việc config keep-alive sẽ do server quyết định
**Server**
```
func main() {
	fmt.Println("server keep alive")
	var s = http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(Index),
	}
	s.SetKeepAlivesEnabled(false)
	s.ListenAndServe()
}
```
**Client** : không có gì thay đổi.

**Kết quả**
* Server:
```
server keep alive
receive a request from: [::1]:52987 map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/1.1]]
receive a request from: [::1]:52991 map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/1.1]]
receive a request from: [::1]:52994 map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/1.1]]
```
* Client: không có gì thay đổi
* Kết luận là cũng không reuse connection.