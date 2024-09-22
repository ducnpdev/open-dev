# Concurrency Pattern

## worker pool
- worker pool là 1 khái niệm khá phổ biến trong khi lập trình, còn trong `golang` mục đích để hạn chế số lượng goroutines được tạo ra để phục vụ một mục đích nhất định.
- pattern này cực kì hữu ích trong quá trình sử dụng, nhắm giới hạn số lượng process cần thực hiện, và quản lý tài nguyên.
  - để handle request từ client
  - để process các job backgrount
  - để ...
### Code:
```golang
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numberJobConst  = 4
	numWorkersConst = 2
)

func doWorker(id int, jobs <-chan int, kqs chan<- int, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Printf("Do Worker %d started job %d\n", id, job)
		time.Sleep(time.Second)
		fmt.Printf("Do Worker %d finished job %d\n", id, job)
		kqs <- job * 2
	}
	wg.Done()
}

func main() {
	jobs := make(chan int, numberJobConst)
	kqs := make(chan int, numberJobConst)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkersConst; i++ {
		wg.Add(1)
		go doWorker(i, jobs, kqs, &wg)
	}

	for j := 1; j <= numberJobConst; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(kqs)

	for kq := range kqs {
		fmt.Println("Worker Result:", kq)
	}
}
```
### kết quả
```sonsole
Do Worker 2 started job 1
Do Worker 1 started job 2
Do Worker 1 finished job 2
Do Worker 1 started job 3
Do Worker 2 finished job 1
Do Worker 2 started job 4
Do Worker 1 finished job 3
Do Worker 2 finished job 4
Worker Result: 4
Worker Result: 2
Worker Result: 6
Worker Result: 8
```

## FanIn FanOut

### FanOUt
Fan-out là một thuật ngữ phổ biến trong các hệ thống phân tán và lập trình hệ thống. Nó mô tả một mô hình mà một tác vụ hoặc yêu cầu được phân chia và gửi đến nhiều thực thể (ví dụ: services, workers, hoặc microservices) để xử lý song song.

Fan-out thường được sử dụng để tăng khả năng mở rộng, hiệu suất, và tốc độ xử lý khi có khối lượng công việc lớn. Dưới đây là một giải thích chi tiết về fan-out:

1. Fan-out trong lập trình hệ thống
- Khi một yêu cầu hoặc một tác vụ từ người dùng hoặc một thành phần của hệ thống cần được xử lý, hệ thống có thể phân tán yêu cầu đó tới nhiều đơn vị xử lý khác nhau. Điều này gọi là fan-out.
- Mỗi đơn vị xử lý sẽ nhận một phần công việc và xử lý nó song song với các đơn vị khác.
- Sau khi hoàn thành, các kết quả có thể được gom lại (fan-in) hoặc được xử lý độc lập tùy vào trường hợp.
2. Fan-out trong hệ thống messaging
- Trong hệ thống message queue (hàng đợi thông điệp), fan-out có thể xuất hiện khi một message được gửi đến nhiều queue hoặc workers khác nhau.
- Ví dụ: Khi một thông báo mới xuất hiện trên hệ thống, nó có thể được đẩy đến nhiều workers để xử lý đồng thời, chẳng hạn như gửi email, ghi log, phân tích dữ liệu, v.v.
3. Ứng dụng của Fan-out
- Hệ thống microservices: Trong các hệ thống microservices, fan-out giúp phân phối công việc giữa nhiều dịch vụ khác nhau. Ví dụ, một yêu cầu tạo đơn hàng có thể cần phải gọi nhiều dịch vụ khác nhau (như xử lý thanh toán, kiểm tra kho, thông báo khách hàng, v.v.) đồng thời.
- Hệ thống messaging và pub/sub: Trong các mô hình publish-subscribe, fan-out được sử dụng để gửi một thông báo đến nhiều subscribers khác nhau. Khi có một sự kiện xảy ra, message sẽ được "fan-out" đến tất cả những dịch vụ quan tâm đến sự kiện đó.
- Xử lý song song: Fan-out được sử dụng rộng rãi trong các hệ thống cần xử lý song song, như việc chia một khối lượng lớn dữ liệu thành nhiều phần nhỏ và phân phối chúng cho nhiều worker hoặc microservice để xử lý cùng lúc.
4. Ví dụ về Fan-out
- Ví dụ 1: Fan-out trong hệ thống xử lý video
  - Một hệ thống xử lý video nhận được một video lớn và muốn chuyển đổi video này sang nhiều định dạng khác nhau (MP4, AVI, MKV).
  - Hệ thống có thể sử dụng fan-out để gửi video này đến nhiều dịch vụ chuyển đổi (mỗi dịch vụ xử lý một định dạng) đồng thời.
- Ví dụ 2: Fan-out trong hệ thống thông báo
  - Khi một sự kiện mới xảy ra trong ứng dụng (ví dụ: một người dùng đăng tải hình ảnh), hệ thống cần thông báo tới nhiều hệ thống khác:
	- Gửi email xác nhận cho người dùng.
	- Lưu lịch sử hoạt động.
	- Cập nhật feed cho người theo dõi.
  - Các hành động này có thể được xử lý song song bằng cách fan-out sự kiện này tới các dịch vụ liên quan.
1. Fan-out và Fan-in
- Fan-out: Một yêu cầu được phân phối tới nhiều nơi để xử lý đồng thời.
- Fan-in: Sau khi các tác vụ song song hoàn thành, kết quả được tập hợp lại (gộp) thành một kết quả cuối cùng.
1. Lợi ích của Fan-out
- Tăng hiệu suất: Việc xử lý song song cho phép tận dụng nhiều tài nguyên hơn, giảm thời gian xử lý tổng thể.
- Khả năng mở rộng: Fan-out cho phép hệ thống dễ dàng mở rộng bằng cách thêm nhiều worker hoặc dịch vụ xử lý hơn khi khối lượng công việc tăng.
- Phân chia trách nhiệm: Mỗi thành phần của hệ thống chỉ cần xử lý một phần công việc nhất định, giúp hệ thống dễ bảo trì và phát triển.
1. Thách thức của Fan-out
- Đồng bộ hóa: Sau khi phân chia công việc, cần có cách để đồng bộ kết quả lại (nếu cần).
- Quản lý lỗi: Xử lý lỗi phức tạp hơn, vì một trong nhiều tác vụ song song có thể thất bại, và cần có cơ chế retry hoặc recovery.
- Quá tải hệ thống: Fan-out không giới hạn có thể dẫn đến việc quá tải tài nguyên nếu không được quản lý cẩn thận.

### Fan In
Fan-in là một khái niệm đối lập với fan-out trong các hệ thống phân tán và xử lý song song. Nó mô tả quá trình gom kết quả từ nhiều nguồn (hoặc tác vụ song song) và hợp nhất lại để tạo ra một kết quả cuối cùng. Fan-in thường được sử dụng sau một quá trình fan-out để đồng bộ hoặc kết hợp các kết quả riêng lẻ từ nhiều tác vụ khác nhau.

1. Fan-in trong hệ thống phân tán
- Khi một yêu cầu hoặc tác vụ được chia nhỏ và xử lý bởi nhiều dịch vụ hoặc worker (fan-out), fan-in sẽ thực hiện việc thu thập các kết quả từ những tác vụ này và kết hợp chúng thành một kết quả tổng hợp.
- Quá trình fan-in giúp hệ thống gom lại các phần dữ liệu hoặc kết quả nhỏ lẻ đã xử lý để phục vụ các mục đích như hiển thị thông tin tổng hợp hoặc đưa ra quyết định dựa trên toàn bộ kết quả.
2. Ứng dụng của Fan-in
Fan-in rất hữu ích trong các hệ thống yêu cầu tổng hợp hoặc kết hợp dữ liệu từ nhiều nguồn hoặc các luồng xử lý song song. Một số ứng dụng phổ biến của fan-in gồm:
- Xử lý song song: Khi một tác vụ lớn được chia nhỏ và xử lý song song, fan-in sẽ tập hợp kết quả của tất cả các tác vụ nhỏ đó lại với nhau để đưa ra kết quả cuối cùng.
- Microservices: Trong hệ thống microservices, nhiều dịch vụ có thể xử lý các phần của một yêu cầu (ví dụ: lấy dữ liệu từ nhiều nguồn khác nhau), và fan-in giúp gom tất cả dữ liệu này lại để trả về cho người dùng.
- Dịch vụ message queue: Sau khi fan-out nhiều message đến các worker để xử lý song song, fan-in giúp thu thập kết quả hoặc sự phản hồi từ các worker này.
3. Ví dụ về Fan-in
* Ví dụ 1: Hệ thống tìm kiếm phân tán
- Trong một hệ thống tìm kiếm lớn, dữ liệu có thể được chia thành nhiều phần và xử lý bởi nhiều server khác nhau.
- Khi người dùng thực hiện tìm kiếm, yêu cầu tìm kiếm sẽ được fan-out đến nhiều server để xử lý đồng thời.
- Sau đó, fan-in sẽ gom kết quả tìm kiếm từ các server này và trả về cho người dùng một kết quả tổng hợp.

* Ví dụ 2: Xử lý đơn hàng phức tạp
- Một yêu cầu đặt hàng có thể yêu cầu nhiều bước xử lý, chẳng hạn như kiểm tra tồn kho, tính toán thuế, và xử lý thanh toán.
- Hệ thống có thể fan-out các tác vụ này cho các dịch vụ khác nhau và sử dụng fan-in để gom lại các kết quả: xác nhận hàng có sẵn, thuế đã được tính, và thanh toán thành công, rồi trả về kết quả cuối cùng là đơn hàng hoàn tất.
4. Lợi ích của Fan-in
- Tối ưu hóa tài nguyên: Fan-in cho phép kết hợp các kết quả từ nhiều luồng xử lý song song, giúp tận dụng tài nguyên hệ thống một cách hiệu quả hơn.
- Đảm bảo tính nhất quán: Khi nhiều dịch vụ hoặc worker hoàn thành tác vụ của mình, fan-in giúp kết hợp lại và đưa ra một kết quả cuối cùng nhất quán.
- Quản lý kết quả phức tạp: Trong các hệ thống xử lý dữ liệu lớn hoặc phân tán, việc thu thập và xử lý kết quả từ nhiều nơi có thể phức tạp, và fan-in giúp gom chúng lại một cách có tổ chức.
5. Thách thức của Fan-in
- ồng bộ hóa: Đảm bảo rằng tất cả các tác vụ đều hoàn thành và kết quả từ các nguồn khác nhau có thể kết hợp một cách chính xác. Nếu một trong các tác vụ song song gặp lỗi, điều này có thể ảnh hưởng đến quá trình fan-in.
- Xử lý lỗi: Nếu một trong những nguồn dữ liệu bị lỗi hoặc không phản hồi, hệ thống fan-in cần có cơ chế để xử lý tình huống này mà không làm gián đoạn toàn bộ quy trình.
- Tính nhất quán: Đảm bảo rằng dữ liệu hoặc kết quả từ các nguồn khác nhau có thể kết hợp và đồng bộ một cách đúng đắn.
6. Fan-in và Fan-out cùng nhau
- Fan-out và fan-in thường được sử dụng kết hợp trong các hệ thống phân tán để tối ưu hóa quá trình xử lý dữ liệu.
  - Fan-out chia một công việc lớn thành nhiều phần nhỏ và xử lý song song.
  - Fan-in gom kết quả từ các phần nhỏ này lại để tạo thành kết quả tổng hợp.

### Code
```go
package main

import (
	"fmt"
	"sync"
)

const (
	numProducers = 2
	numConsumers = 2
)

func producer(id int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 5; i++ {
		ch <- i
		fmt.Printf("Fanout Producer %d produced %d\n", id, i)
	}
}

func consumer(id int, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range in {
		out <- v * 2
		fmt.Printf("FanIn Consumer %d processed %d\n", id, v)
	}
}

func main() {

	input := make(chan int, 10)
	output := make(chan int, 10)
	var wg sync.WaitGroup
	for i := 1; i <= numProducers; i++ {
		wg.Add(1)
		go producer(i, input, &wg)
	}
	wg.Wait()
	close(input)
	for i := 1; i <= numConsumers; i++ {
		wg.Add(1)
		go consumer(i, input, output, &wg)
	}
	wg.Wait()
	close(output)
	for result := range output {
		fmt.Println("Ketqua:", result)
	}
}

```

### KQ:
```console
Fanout Producer 2 produced 1
Fanout Producer 2 produced 2
Fanout Producer 2 produced 3
Fanout Producer 2 produced 4
Fanout Producer 2 produced 5
Fanout Producer 1 produced 1
Fanout Producer 1 produced 2
Fanout Producer 1 produced 3
Fanout Producer 1 produced 4
Fanout Producer 1 produced 5
FanIn Consumer 2 processed 1
FanIn Consumer 2 processed 1
FanIn Consumer 2 processed 2
FanIn Consumer 2 processed 3
FanIn Consumer 2 processed 4
FanIn Consumer 2 processed 2
FanIn Consumer 2 processed 3
FanIn Consumer 2 processed 4
FanIn Consumer 2 processed 5
FanIn Consumer 1 processed 5
Ketqua: 2
Ketqua: 2
Ketqua: 4
Ketqua: 6
Ketqua: 8
Ketqua: 10
Ketqua: 4
Ketqua: 6
Ketqua: 8
Ketqua: 10
```

## Pipeline Pattern
Pipeline pattern là một mô hình thiết kế (design pattern) trong lập trình, đặc biệt hữu ích khi xử lý chuỗi các tác vụ (tasks) mà mỗi tác vụ có thể xử lý dữ liệu một cách tuần tự. Mỗi bước trong pipeline có thể được coi là một công đoạn (stage), và dữ liệu sẽ đi qua từng công đoạn để được xử lý dần dần.
### Đặc điểm của Pipeline Pattern:
1. Chia nhỏ tác vụ: Bài toán lớn được chia nhỏ thành nhiều bước hoặc nhiều tác vụ nhỏ hơn. Mỗi bước trong pipeline chịu trách nhiệm xử lý một phần công việc.
2. Tuần tự hóa xử lý: Các bước thực thi tuần tự, bước sau sử dụng kết quả của bước trước. Dữ liệu được "chảy" qua từng bước giống như nước qua các ống trong một hệ thống dẫn nước.
3. Tăng khả năng tái sử dụng: Mỗi bước có thể là một phần độc lập, dễ tái sử dụng trong các pipeline khác.
4. Tính mở rộng: Pipeline pattern dễ mở rộng khi bạn muốn thêm các bước mới hoặc loại bỏ các bước không cần thiết.
### Cách hoạt động của Pipeline Pattern:
1. Input: Dữ liệu đầu vào được chuyển đến bước đầu tiên của pipeline.
2. Processing: Mỗi bước trong pipeline thực hiện một công việc cụ thể, sau đó chuyển kết quả đến bước tiếp theo.
3. Output: Kết quả cuối cùng sẽ được xử lý qua tất cả các bước trong pipeline và được trả về.