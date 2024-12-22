# rewinding 2024

- Một năm 2024 sắp kết thúc, là lúc nhìn lại trong 1 năm qua Golang đã có những điểm nào mới.

## Go 1.23:
- Go 1.23 cập nhật một vài điểm mang điến trải nghiệp tốt hơn cho lâp trình viên và hiệu năng tối hợp cho ứng dụng.

### Các Điểm Chính.
1. Build nhanh hơn.
- Nhanh hơn khoảng 15%, tăng công xuất đáng kể trong quá trình CI/CD

2. Sync Package
- `sync.Map` được cải tiến với khả năng xử lý tốt hơn, giúp dễ dàng đáp ứng trong việc xử lý các ứng dụng cần độ trễ thấp.
- code:
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map

	// Lưu trữ các giá trị trong sync.Map
	m.Store("opendev1", "opendev-value1")
	m.Store("opendev2", "opendev-value2")
	m.Store("opendev2", "opendev-value3")

	// Truy xuất giá trị
	if value, ok := m.Load("opendev1"); ok {
		fmt.Println("Found opendev1:", value)
	} else {
		fmt.Println("opendev1 not found")
	}

	// Xóa giá trị
	m.Delete("opendev2")
	if _, ok := m.Load("opendev2"); !ok {
		fmt.Println("opendev2 has been deleted")
	}

	// Duyệt qua các phần tử trong sync.Map
	fmt.Println("Iterating over map:")
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true // return false sẽ dừng việc duyệt
	})
}

```

3. Độc Lập Môi Trường.
- Cho phép lâp trình viên có thể độc lâp từng môi trường, dễ dàng thay đổi qua để chuyển môi trường production.
- Code:
```go
config_dev.go (cho môi trường development)
//go:build dev
package main

import "fmt"

func LoadConfig() {
    fmt.Println("Loaded Development Configuration")
}

config_prod.go (cho môi trường production)
//go:build prod
package main

import "fmt"

func LoadConfig() {
    fmt.Println("Loaded Production Configuration")
}

//
package main

func main() {
    LoadConfig()
}

```

- run: Loaded Development Configuration
```sh
go build -tags dev -o app_dev
./app_dev
```

4. Observability
- Đã tích hợp các tool monitor vào như `Telemetry` `Prometheus` `Grafana`
- Giúp các lâp trình viên dễ debug cũng như tối ưu hiệu năng.
- link: https://go.dev/doc/go1.23#telemetry

5. Xây dựng Opensource Lớn
- Với hiệu năng cao và khả năng biên dịch nhanh, Golang đã trờ thành ngôn ngữ phổ biến để xây dựng các open source lớn: `Docker`, `Kubernetes`, `Terraform`, ...
- k8s: https://github.com/kubernetes/kubernetes
- docker: https://github.com/docker/compose
- terraform: https://github.com/hashicorp/terraform