# kafka golang example
- install kafka-go
```go
go get https://github.com/segmentio/kafka-go
```
## producer
- create file main.go
- copy and paste:
```go
package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var (
	BrokerAddress = "localhost:9092"
	TopicLogging  = "cmak"
)

func main() {
	Produce(context.Background())
}
func Produce(ctx context.Context) {
	i := 0

	l := log.New(os.Stdout, "kafka writer: ", 0)
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{BrokerAddress},
		Topic:   TopicLogging,
		// assign the logger to the writer
		Logger: l,
		Async:  false,
	})

	for {
		writeMsg(ctx, w, i)
		time.Sleep(time.Millisecond)
	}
}

func writeMsg(ctx context.Context, w *kafka.Writer, i int) {
	msgs := kafka.Message{
		Key:   []byte(strconv.Itoa(i)),
		Value: []byte("timeNow:" + time.Now().Format(time.RFC3339Nano) + " " + "uuid:" + uuid.New().String()),
	}
	err := w.WriteMessages(ctx, msgs)
	if err != nil {
		panic("could not write message " + err.Error())
	}

}
```
## consumer
- create file main.go
- copy and paste:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "topic"
	TopicLogging  = "cmak"
	Topic1        = "topic1"
	Topic2        = "topic2"
	BrokerAddress = "localhost:9092"
	Group         = "cmak-consumer-group-1"
)

func MainKafka() {
	s := make(chan bool)
	go Consume(context.Background(), TopicLogging, 0)
	<-s
}

func Consume(ctx context.Context, topic string, pa int) {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	l := log.New(os.Stdout, "kafka reader: ", 0)
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{BrokerAddress},
		Topic:   topic,
		GroupID: Group,
		// assign the logger to the reader
		Logger: l,
		Dialer: dialer,
	})
	for {

		m, err := r.FetchMessage(ctx)

		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}
```

## Kafka Advanced Concepts
### Topic Advanced
#### Change configuration a kafka topic
trước khi run cli, chắc chắn là kafka đã start thành công.
- Đầu tiên, tạo một chủ đề có tên là chủ đề được định cấu hình với 3 phân vùng và hệ số sao chép là 1, sử dụng chủ đề Kafka CLI, kafka-topics
```sh
kafka-topics --bootstrap-server localhost:9092 --create --topic configured-topic --partitions 3 --replication-factor 1
``` 
- Mô tả chủ đề để kiểm tra xem có bất kỳ cài đặt ghi đè cấu hình nào cho chủ đề này không.
```sh
kafka-topics --bootstrap-server localhost:9092 --describe --topic configured-topic
```

- Không có tập ghi đè cấu hình nào.
  - Đặt giá trị min.insync.replicas cho chủ đề được định cấu hình chủ đề thành 2
```sh
kafka-configs --bootstrap-server localhost:9092 --alter --entity-type topics --entity-name configured-topic --add-config min.insync.replicas=2
```
```sh
kafka-topics --bootstrap-server localhost:9092 --describe --topic configured-topic
Topic: configured-topic	TopicId: CDU7SBxBQ1mzJGnuH68-cQ	PartitionCount: 3	ReplicationFactor: 1	Configs: min.insync.replicas=2
	Topic: configured-topic	Partition: 0	Leader: 2	Replicas: 2	Isr: 2
	Topic: configured-topic	Partition: 1	Leader: 3	Replicas: 3	Isr: 3
	Topic: configured-topic	Partition: 2	Leader: 1	Replicas: 1	Isr: 1
```

- Bây giờ, bạn có thể thấy có một bộ ghi đè cấu hình chủ đề (ở bên phải của đầu ra) - min.insync.replicas=2.
  - Bạn có thể xóa ghi đè cấu hình bằng cách chuyển --delete-config thay cho cờ --add-config.
```sh
kafka-configs --bootstrap-server localhost:9092 --alter --entity-type topics --entity-name configured-topic --delete-config min.insync.replicas
```
- Describe the topic to make sure the configuration override has been removed.

### topic internal: Segment và Indexes
- Đơn vị lưu trữ cơ bản của Kafka là bản sao phân vùng. Khi bạn tạo một chủ đề, trước tiên Kafka sẽ quyết định cách phân bổ các phân vùng giữa các nhà môi giới. Nó phân phối bản sao đồng đều giữa các nhà môi giới.

- Kafka brokers chia mỗi phân vùng thành các phân đoạn. Mỗi phân đoạn được lưu trữ trong một tệp dữ liệu duy nhất trên đĩa gắn liền với nhà môi giới. Theo mặc định, mỗi phân đoạn chứa 1 GB dữ liệu hoặc một tuần dữ liệu, tùy theo giới hạn nào đạt được trước. Khi nhà môi giới Kafka nhận được dữ liệu cho một phân vùng, khi đạt đến giới hạn phân đoạn, nó sẽ đóng tệp và bắt đầu một phân vùng mới:
![Logo của dự án](https://github.com/ducnpdev/open-dev/blob/master/kafka/images/segment.png)

- Chỉ có một phân đoạn HOẠT ĐỘNG tại bất kỳ thời điểm nào - một phân đoạn đang được ghi vào. Một phân đoạn chỉ có thể bị xóa nếu nó đã được đóng trước đó. Kích thước của một phân khúc được kiểm soát bởi hai cấu hình Nhà môi giới (cũng có thể được sửa đổi ở cấp chủ đề)
  - log.segment.bytes: kích thước tối đa của một phân đoạn tính bằng byte (mặc định 1 GB)
  - log.segment.ms: thời gian Kafka sẽ đợi trước khi thực hiện phân đoạn nếu chưa đầy (mặc định 1 tuần)
- Kafka brokers giữ một tệp xử lý mở cho mọi phân đoạn trong mọi phân vùng - ngay cả các phân đoạn không hoạt động. Điều này dẫn đến số lượng xử lý tệp đang mở thường cao và hệ điều hành phải được điều chỉnh cho phù hợp.

#### Kafka Topic Segments and Indexes
- Kafka allows consumers to start fetching messages from any available offset. In order to help brokers quickly locate the message for a given offset, Kafka maintains two indexes for each segment:
  - An offset to position index - It helps Kafka know what part of a segment to read to find a message
  - A timestamp to offset index - It allows Kafka to find messages with a specific timestamp
![Logo của dự án](https://github.com/ducnpdev/open-dev/blob/master/kafka/images/segment1.png)

#### Inspecting the Kafka Directory Structure
- Kafka lưu trữ tất cả dữ liệu của nó trong một thư mục trên đĩa môi giới. Thư mục này được chỉ định bằng thuộc tính log.dirs trong tệp cấu hình của nhà môi giới. Ví dụ,
```code
# A comma separated list of directories under which to store log files
log.dirs=/tmp/kafka-logs
```
- Khám phá thư mục và nhận thấy rằng có một thư mục cho mỗi phân vùng chủ đề. Tất cả các phân đoạn của phân vùng đều nằm bên trong thư mục phân vùng. Ở đây, chủ đề có tên là configure-topic có ba phân vùng, mỗi phân vùng có một thư mục - configure-topic-0, configure-topic-1 và configure-topic-2.
![Logo của dự án](https://github.com/ducnpdev/open-dev/blob/master/kafka/images/segment2.png)
- Đi xuống một thư mục cho một phân vùng chủ đề. Lưu ý các chỉ mục - thời gian và độ lệch cho phân đoạn và chính tệp phân đoạn nơi các thông báo được lưu trữ.
![Logo của dự án](https://github.com/ducnpdev/open-dev/blob/master/kafka/images/segment3.png)

#### Considerations for Segment Configurations
Chúng ta hãy xem lại cấu hình cho các phân đoạn và tìm hiểu tầm quan trọng của chúng.
- log.segment.bytes Khi các tin nhắn được tạo cho nhà môi giới Kafka, chúng sẽ được thêm vào phân đoạn hiện tại cho phân vùng. Sau khi phân đoạn đã đạt đến kích thước được chỉ định bởi tham số log.segment.bytes (mặc định là 1 GB), phân đoạn sẽ bị đóng và một phân đoạn mới sẽ được mở.
  - Kích thước phân đoạn nhỏ hơn có nghĩa là các tệp phải được đóng và phân bổ thường xuyên hơn, điều này làm giảm hiệu quả chung của việc ghi đĩa.
  - Khi một phân khúc đã bị đóng, nó có thể được coi là hết hạn. Việc điều chỉnh kích thước của các phân khúc có thể quan trọng nếu các chủ đề có tỷ lệ sản xuất thấp. Có kích thước phân đoạn nhỏ có nghĩa là Kafka phải mở nhiều tệp, điều này có thể dẫn đến lỗi Quá nhiều tệp đang mở.

- log.segment.ms Một cách khác để kiểm soát thời điểm đóng phân đoạn là sử dụng tham số log.segment.ms, tham số này chỉ định khoảng thời gian sau đó phân đoạn sẽ được đóng. Mặc định là 1 tuần. Kafka sẽ đóng một phân đoạn khi đạt đến giới hạn kích thước hoặc khi đạt đến giới hạn thời gian, tùy điều kiện nào đến trước.
  - Khi sử dụng giới hạn phân đoạn dựa trên thời gian, điều quan trọng là phải xem xét tác động lên hiệu suất đĩa khi nhiều phân đoạn được đóng đồng thời.
  - Quyết định xem bạn có muốn nén hàng ngày thay vì hàng tuần hay không
  

### Consumer Advanced
#### Các cách để đọc message từ consumer
  - khi xử lý 1 message từ `kafka`, cần phải lựa chọn thời điểm để commit offsets của nó, và các cách đọc message khác nhau sẽ ảnh hưởng đến các thiết kế và ứng dụng của bạn.

##### At Most Once Delivery
- trong trường hợp này, khi thực hiện đọc message từ `kafka` thì message đó sẽ commit ngay lập tức.
- nếu quá trình xử lý message xảy ra lỗi, thì không thể nào khôi phục được message đó -> mất message.
- phù hợp cho những ứng dụng có thể cho phép data bị mất.
![Logo của dự án](https://github.com/ducnpdev/open-dev/blob/master/kafka/images/at_most_once_delivery.png)
1. Create topic and config.
- 

##### At Least Once Delivery (thường được xử dụng)
- đọc message từ `kafka` ít nhất 1 lần,
- việc đọc message được nhiều lần từ `kafka` sẽ dẫn đến trường hợp trùng message. khi 1 message gặp vấn đề trong lúc xử lý, chúng ta có thể đọc lại.
- phù hợp cho những hệ thống không được mất message.
```txt
Idempotent Processing: Make sure your processing is idempotent (i.e. processing again the messages won’t impact your systems)
```
![Logo của dự án](https://github.com/ducnpdev/open-dev/blob/master/kafka/images/at_least_once_delivery.png)

#### Exactly Once Delivery
- một số ứng dụng ngoài việc đọc message ít nhất 1 lần( không mất dữ liệu ) mà còn yêu cầu là chính xác đúng 1 lần, 1 message là được xử lý đúng 1 lần.
- điều này giúp cho `kafka` đáp ứng được 1 số vấn đề yêu cầu hệ thống xử lý đúng 1 lần như api payment
- để config chỗ này chúng ta cần chỉnh: `processing.guarantee=exactly.once`

#### Summary
- At most once: `offsets` sẽ được `commit` sau khi message được nhận, nếu xảy ra lỗi xử lý, message sẽ bị mất.
- At least once: `offsets` sẽ được `commit` sau khi message được xử lý xong, nếu gặp vấn đề thì hoàn toàn có thể đọc lại message -> message nhiều lúc sẽ bị duplicated, sử dụng `idempotent-key` để xử lý.
- Exactly Once Delivery: sẽ phù hợp với việc xử lý trong api transaction hoặc cơ chế `kafka-streams api`

```txt
Cuối Cùng:
Đối với hầu hết các ứng dụng, bạn nên sử dụng quy trình xử lý 'Ít nhất một lần' và đảm bảo các phép biến đổi/xử lý của bạn là bình thường.
```