# Consumer Advanced
## Chuẩn Bị:
1. create topic:
- create topic with 1 partition
```sh
kafka-topics --create --topic consumer-rebalance --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1
```
- describe topic
```sh
kafka-topics --describe --topic consumer-rebalance --bootstrap-server localhost:9092
```
        ```
        Topic: consumer-rebalance       TopicId: cLCaoJPETCKu19544xdL0A PartitionCount: 1       ReplicationFactor: 1    Configs: segment.bytes=1073741824
        Topic: consumer-rebalance       Partition: 0    Leader: 0       Replicas: 0     Isr: 0
        ```
2. Test Push and Consume Message Topic: `consumer-rebalance `
* push message:
  - edit const `TopicName` là topic name thành `consumer-rebalance` in file `producer/main.go`
```go
go run producer/main.go
```
* receive message:
  - edit const `TopicName` là topic name thành `consumer-rebalance` in file `consumer/main.go`
```go
go run consumer/main.go
```

## Consumer Rebalance
- thực hiện việc test consumer khi topic có 1 partition và nhiều hơn 1 partition
### Topic Có 1 partition.
<!-- 1. join new consumer to group -->
- start producer và consumer 1 cách bình thường:
```go
go run consumer/main.go
---
go run producer/main.go
```

- tiếp đến thì start 1 consumer nữa, để xem có thể consumer được không.
```go
go run consumer/main.go
```
  - không thể nào consume message vì topic `consumer-rebalance` chỉ có 1 partition.

#### Topic có 3 partition
1. Thực hiện việc tạo 1 `topic` mới có 3 partition.
- create topic with 3 partition: `--partitions 3`
```sh
kafka-topics --create --topic consumer-rebalance-3p --bootstrap-server localhost:9092 --partitions 3
```
- describe topic
```sh
kafka-topics --describe --topic consumer-rebalance-3p --bootstrap-server localhost:9092
```
        
        Topic: consumer-rebalance-3p    TopicId: bFdZwzJeTMupOzssawol6g PartitionCount: 3      ReplicationFactor: 1    Configs: segment.bytes=1073741824
        Topic: consumer-rebalance-3p    Partition: 0    Leader: 0     Replicas: 0      Isr: 0
        Topic: consumer-rebalance-3p    Partition: 1    Leader: 0     Replicas: 0      Isr: 0
        Topic: consumer-rebalance-3p    Partition: 2    Leader: 0     Replicas: 0      Isr: 0
        
1. Test
- push message vào kafka
  - đổi filed `TopicName` thành `consumer-rebalance-3p`
  - run
    ```go
    go run producer/main.go
    ```
- start 1 consumer:
  - đổi filed `TopicName` thành `consumer-rebalance-3p`
  - run:
  ```go
  go run consumer/main.go
  ```
  - Kết quả:
    - video: 
    - lúc này consumer này sẽ balance consume message điều trong cả 3 partition:
      ```
      message at topic/partition/offset consumer-rebalance-3p/0/0: 1 = timeNow:2024-06-09T16:14:56.438471+07:00 uuid:62aea132-8568-40ea-a83d-4e9f4d038053
      kafka reader: committed offsets for group rebalance-1: 
              topic: consumer-rebalance-3p
                      partition 0: 1
      message at topic/partition/offset consumer-rebalance-3p/0/1: 4 = timeNow:2024-06-09T16:14:59.453256+07:00 uuid:54ae3550-7a92-45ca-963f-108bbaefc94d
      kafka reader: committed offsets for group rebalance-1: 
              topic: consumer-rebalance-3p
                      partition 0: 2
      ```

- start consumer thứ 2:

  - run:
    ```go
    go run consumer/main.go
    ```
  - vẫn consumer message, nhưng chỉ consume đúng 1 partion mà thôi, như trong trường hợp này thì 1 consumer sẽ consume partition 0 và partition 1,2 sẽ consume bởi consumer kia.
  - video:

- start consumer thứ 3:

  - run:
    ```go
    go run consumer/main.go
    ```
  - cũng tương tự như khi start consumer thứ 2, nhưng chỉ consume đúng 1 partion mà thôi.
  - video:

### consumer out group
bây giờ ta sẽ test trường hợp 1 consumer out group
- theo dõi 1 thời gian sẽ thấy tất cả các consumer sẽ re-assign all:
  ```log
  kafka reader: no messages received from kafka within the allocated time for partition 2 of consumer-rebalance-3p at offset 722
  kafka reader: stopped heartbeat for group rebalance-1
  kafka reader: stopped commit for group rebalance-1
  kafka reader: joined group rebalance-1 as member main@nguyens-MacBook-Pro-4.local (github.com/segmentio/kafka-go)-f08c577f-aea1-4244-ae00-9f7da09ed3b6 in generation 10
  kafka reader: selected as leader for group, rebalance-1
  kafka reader: using 'range' balancer to assign group, rebalance-1
  kafka reader: found member: main@nguyens-MacBook-Pro-4.local (github.com/segmentio/kafka-go)-f08c577f-aea1-4244-ae00-9f7da09ed3b6/[]byte(nil)
  kafka reader: found member: main@nguyens-MacBook-Pro-4.local (github.com/segmentio/kafka-go)-2769d806-2eef-4e01-a5a7-4438de66b1fb/[]byte(nil)
  kafka reader: found topic/partition: consumer-rebalance-3p/0
  kafka reader: found topic/partition: consumer-rebalance-3p/1
  kafka reader: found topic/partition: consumer-rebalance-3p/2
  kafka reader: assigned member/topic/partitions main@nguyens-MacBook-Pro-4.local (github.com/segmentio/kafka-go)-2769d806-2eef-4e01-a5a7-4438de66b1fb/consumer-rebalance-3p/[0]
  kafka reader: assigned member/topic/partitions main@nguyens-MacBook-Pro-4.local (github.com/segmentio/kafka-go)-f08c577f-aea1-4244-ae00-9f7da09ed3b6/consumer-rebalance-3p/[1 2]
  kafka reader: joinGroup succeeded for response, rebalance-1.  generationID=10, memberID=main@nguyens-MacBook-Pro-4.local (github.com/segmentio/kafka-go)-f08c577f-aea1-4244-ae00-9f7da09ed3b6
  kafka reader: Joined group rebalance-1 as member main@nguyens-MacBook-Pro-4.local (github.com/segmentio/kafka-go)-f08c577f-aea1-4244-ae00-9f7da09ed3b6 in generation 10
  kafka reader: Syncing 2 assignments for generation 10 as member main@nguyens-MacBook-Pro-4.local (github.com/segmentio/kafka-go)-f08c577f-aea1-4244-ae00-9f7da09ed3b6
  kafka reader: sync group finished for group, rebalance-1
  kafka reader: started heartbeat for group, rebalance-1 [3s]
  kafka reader: subscribed to topics and partitions: map[{topic:consumer-rebalance-3p partition:1}:591 {topic:consumer-rebalance-3p partition:2}:621]
  kafka reader: initializing kafka reader for partition 2 of consumer-rebalance-3p starting at offset 621
  kafka reader: started commit for group rebalance-1
  kafka reader: initializing kafka reader for partition 1 of consumer-rebalance-3p starting at offset 591
  kafka reader: the kafka reader for partition 1 of consumer-rebalance-3p is seeking to offset 591
  kafka reader: the kafka reader for partition 2 of consumer-rebalance-3p is seeking to offset 621
  ```
- bây giờ sẽ chỉ còn 2 consumer, và 1 consumer sẽ consume 2 partition, 1 consumer sẽ consumer 1 partition.
 
### consumer multi group
1. run group 1:
- edit group name:
  - in file consumer/main.go, line 16
    ```go
    Group         = "rebalance-1"
    ```
- run cmd:
```go
go run consumer/main.go
```

2. run group 2:
- edit group name:
  - in file consumer/main.go, line 16
    ```go
    Group         = "rebalance-2"
    ```
- run cmd:
```go
go run consumer/main.go
```

3. run group 3:
- edit group name:
  - in file consumer/main.go, line 16
    ```go
    Group         = "rebalance-3"
    ```
- run cmd:
```go
go run consumer/main.go
```
#### tổng kết:
- Với topic `consumer-rebalance` có `1 partition` hiện tại có 3 group consume là `rebalance-1`,`rebalance-2`,`rebalance-3`
- khi 1 group nào xảy ra vấn đề thì các group còn lại vẫn hoạt động bình thường.
- Thử đặt ra vấn đề là group `rebalance-1` gặp vấn đề và down-time service, thì khi rollout lại dịch vụ vẫn hoạt động bình thường không có gì xảy ra cả. Consume tiếp offset đã consume trước đó.