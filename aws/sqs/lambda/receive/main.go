package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func GetSession() *session.Session {
	log.Default().Println("start-GetSession:")
	return session.Must(session.NewSessionWithOptions(session.Options{
		// Profile: "default",
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	}))
}
func Handler(ctx context.Context, e events.SQSEvent) error {
	ses := GetSession()

	for i, record := range e.Records {
		log.Default().Println("i:", i)
		log.Default().Println("record:", fmt.Sprintf("%#v", record))
		DeleteMessage(ses, &record.ReceiptHandle)
	}
	return nil
}
func DeleteMessage(sess *session.Session, messageHandle *string) error {
	sqsClient := sqs.New(sess)

	_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/064038607558/sqs-demo"),
		ReceiptHandle: messageHandle,
	})

	return err
}

func main() {
	fmt.Println("func-receive:", time.Now().Format("2006-01-02T15:04:05.000-07:00"))
	lambda.Start(Handler)
	// Handler(context.Background())
}
