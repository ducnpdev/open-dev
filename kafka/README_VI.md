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
