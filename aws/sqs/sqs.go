package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// GetSession, get session purpose reuseable
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

// CreateQueue create queue in aws simple queue service
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

// GetQueueURL get url of queue
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

// SendMessage send message to queue
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

// DeleteMessage delete single message in queue
func DeleteMessage(sess *session.Session, queueUrl string, messageHandle *string) error {
	sqsClient := sqs.New(sess)

	_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: messageHandle,
	})

	return err
}

// PurgeQueue delete all message in queue
func PurgeQueue(sess *session.Session, queueUrl string) error {
	sqsClient := sqs.New(sess)
	_, err := sqsClient.PurgeQueue(&sqs.PurgeQueueInput{
		QueueUrl: aws.String(queueUrl),
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

// get list queue
func ListQueues(sess *session.Session) (*sqs.ListQueuesOutput, error) {
	sqsClient := sqs.New(sess)
	list, err := sqsClient.ListQueues(&sqs.ListQueuesInput{
		MaxResults: aws.Int64(10),
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
