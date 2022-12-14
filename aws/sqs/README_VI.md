# simple queue service

## chuẩn bị
- cài đặt sdk aws, run
```go
go get -u github.com/aws/aws-sdk-go/
```
- cấu hình access-key and secret-key, nếu bạn không biết thì đọc tham khảo link: https://viblo.asia/p/cau-hinh-aws-credential-zOQJwYPxVMP 

## tạo một session
- tạo một session để những dùng trong những function khác
- profile được dùng là default, để biết profile, run trên linux/mac
```
cat ~/.aws/credentials 
```
- region thì lấy region muốn tạo 1 queue
```go
func GetSession() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-west-1"),
		},
	})
	if err != nil {
		panic(err)
	}
	return sess
}

```
## cách tạo một queue?
- sử dụng function CreateQueue trong aws sdk để tạo ra 1 queue cho việc test.
- có một số parameter quan trọng như
    - queueName: name của queue mà bạn muốn tạo
    - DelaySeconds: message bạn muốn giữ lại trước khi gửi đi, 
    - VisibilityTimeout: thời gian hiện thị của một message

```go
func CreateQueue(sess *session.Session, queueName string) (*sqs.CreateQueueOutput, error) {
	sqsClient := sqs.New(sess)
	result, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: map[string]*string{
			"DelaySeconds":      aws.String("0"),
			"VisibilityTimeout": aws.String("60"),
		},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
```

## get url của một queue