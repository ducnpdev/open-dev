package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Response events.APIGatewayProxyResponse

// func SendMessage(myMessage map[string]string) (*sqs.SendMessageOutput, error) {
// 	sess, err := session.NewSession()
// 	if err != nil {
// 		return nil, err
// 	}
// 	svc := sqs.New(sess)
// 	body, err := json.Marshal(myMessage)
// 	if err != nil {
// 		return nil, err
// 	}
// 	messageGroupId := "groupid"
// 	res, err := svc.SendMessage(&sqs.SendMessageInput{
// 		MessageBody:    aws.String(string(body)),
// 		MessageGroupId: &messageGroupId,
// 		QueueUrl:
// 	})
// 	return res, err
// }

func SendMessage(messageBody string) (*sqs.SendMessageOutput, error) {
	log.Default().Println("start-SendMessage:")
	ses := GetSession()
	// if err != nil {
	// 	log.Default().Println("get-session:", err)
	// 	return nil, err
	// }
	sqsClient := sqs.New(ses)
	for i := 0; i < 100; i++ {

		_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
			QueueUrl:    aws.String("https://sqs.us-east-1.amazonaws.com/064038607558/sqs-demo"),
			MessageBody: aws.String(messageBody + time.Now().Format("2006-01-02T15:04:05.000-07:00")),
			// DelaySeconds: aws.Int64(0),

		})
		if err != nil {
			log.Default().Println("SendMessage:", err)
			return nil, err
		}

	}

	return nil, nil
}

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

func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer
	msgResp := "push message to queue"
	msgSend := "this is message"
	_, err := SendMessage(msgSend)
	if err != nil {
		msgResp = fmt.Sprintf("send msg to queue err %s", err)
	}
	log.Default().Println("resp:", msgResp)
	body, err := json.Marshal(map[string]interface{}{
		"message": msgResp,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "push-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
	// Handler(context.Background())
}
