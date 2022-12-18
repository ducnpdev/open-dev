# simple queue service

## chuẩn bị
- cài đặt sdk aws, run
```go
go get -u github.com/aws/aws-sdk-go/
```
- cấu hình access-key and secret-key, nếu bạn không biết thì đọc tham khảo link: https://viblo.asia/p/cau-hinh-aws-credential-zOQJwYPxVMP 

## tạo một session
- tạo một session để tái sử dụng trong những function khác
- profile được dùng là default, để biết profile run trên linux/mac
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
## cách tạo một queue
- sử dụng function `CreateQueue` trong aws-sdk để tạo 1 queue cho việc test.
- có một số parameter quan trọng cần lưu ý:
    - **queueName**: name của queue mà bạn muốn tạo
    - **DelaySeconds**: thời gian message bạn muốn giữ lại trước khi gửi đi, 
    - **VisibilityTimeout**: thời gian một message trước khi expire.

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

## get URL của một queue
- Tất cả các apis điều yêu cầu url của queue, vì thế sử dụng `GetQueueURL` để truy vấn url của *queue-name*
```go
func main() {
	queueName := "sqs-demo"
	ses := GetSession()
	url, err := GetQueueURL(ses, queueName)
	if err != nil {
		panic(err)
	}
	fmt.Println("url:", *url.QueueUrl)
}
func GetSession() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		panic(err)
	}
	return sess
}

func GetQueueURL(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

```
- output:
```
url: https://sqs.us-east-1.amazonaws.com/********/sqs-demo
```

## Gửi message đến queue
- chúng ta sẽ sử dụng hàm `SendMessage` từ aws-sdk để gửi message đến queue
- một vài paramters quan trọng cần chú ý trong lúc gửi message:
	- **QueueUrl**: url của queue muốn gửi message đến.
	- **MessageBody**: body được gửi đến queue, có thể là json-string, string
```go

func main() {
	queueName := "sqs-demo"
	ses := GetSession()
	url, err := GetQueueURL(ses, queueName)
	if err != nil {
		panic(err)
	}
	msgBody := "this is body of message"
	_, err = SendMessage(ses, *url.QueueUrl, msgBody)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent successfully")
}

func SendMessage(sess *session.Session, queueUrl string, messageBody string) (*sqs.SendMessageOutput, error) {
	sqsClient := sqs.New(sess)

	sendOut, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		return nil, err
	}

	return sendOut, nil
}

func GetSession() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		panic(err)
	}
	return sess
}

func GetQueueURL(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
```
- output:
```
Message sent successfully
```

## Nhận message từ queue
- chúng ta sẽ sử dụng hàm `ReceiveMessage` từ aws-sdk để lấy message từ queue
- một vài paramters quan trọng cần được chú ý trong lúc nhận message:
	- **QueueUrl**: url của queue muốn nhận message.
	- **MaxNumberOfMessages**: tổng số message có thể nhận được.
```go

func main() {
	queueName := "sqs-demo"
	ses := GetSession()
	url, err := GetQueueURL(ses, queueName)
	if err != nil {
		panic(err)
	}
	msgRes, err := GetMessages(ses, *url.QueueUrl, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message Body: " + *msgRes.Messages[0].Body)
	fmt.Println("Message Handle: " + *msgRes.Messages[0].ReceiptHandle)
}
func GetMessages(sess *session.Session, queueUrl string, maxMessages int) (*sqs.ReceiveMessageOutput, error) {
	sqsClient := sqs.New(sess)
	msgResult, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(int64(maxMessages)),
	})

	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func GetSession() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		panic(err)
	}
	return sess
}

func GetQueueURL(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
```

- output:
```
Message Body: this is body of message
Message Handle: AQEBqfjXi***
```

## Xoá message trong queue
- khi nhận message từ queue, message sẽ không tự động xoá ra khỏi queue.
- consumer khác có thể nhận message sau thời gian **VisibilityTimeout** hết hạn.
- để đảm bảo rằng không bị trùng message thì cần phải xoá.
- chúng ta sẽ sử dụng hàm `DeleteMessage` từ aws-sdk để xoá message từ queue
- cần cung cấp **ReceiptHandle** trong thành phần của method
```go

func main() {
	queueName := "sqs-demo"
	ses := GetSession()
	url, err := GetQueueURL(ses, queueName)
	if err != nil {
		panic(err)
	}
	msgRes, err := GetMessages(ses, *url.QueueUrl, 1)
	if err != nil {
		panic(err)
	}
	err = DeleteMessage(ses, *url.QueueUrl, msgRes.Messages[0].ReceiptHandle)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete message successfully")
}
func DeleteMessage(sess *session.Session, queueUrl string, messageHandle *string) error {
	sqsClient := sqs.New(sess)

	_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: messageHandle,
	})

	return err
}
func GetMessages(sess *session.Session, queueUrl string, maxMessages int) (*sqs.ReceiveMessageOutput, error) {
	sqsClient := sqs.New(sess)
	msgResult, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(int64(maxMessages)),
	})

	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func GetSession() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		panic(err)
	}
	return sess
}

func GetQueueURL(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
```
- output:
```
Delete message successfully
```


## Xoá tất cả message trong queue
- chúng ta sẽ sử dụng hàm `PurgeQueue` từ aws-sdk để xoá tất cả message trong queue
- một vài paramters quan trọng cần được chú ý trong lúc xoá tất cả message:
	- **QueueUrl**: url của queue muốn xoá tất cả.
```go

func main() {
	queueName := "sqs-demo"
	ses := GetSession()
	url, err := GetQueueURL(ses, queueName)
	if err != nil {
		panic(err)
	}
	err = PurgeQueue(ses, *url.QueueUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete all message successfully")
}
func PurgeQueue(sess *session.Session, queueUrl string) error {
	sqsClient := sqs.New(sess)
	_, err := sqsClient.PurgeQueue(&sqs.PurgeQueueInput{
		QueueUrl: aws.String(queueUrl),
	})
	return err
}

func GetSession() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		panic(err)
	}
	return sess
}

func GetQueueURL(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
```
- output:
```
Delete all message successfully
```