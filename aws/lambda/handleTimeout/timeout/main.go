package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func process(ctx context.Context) error {
	log.Default().Println("start timeout")
	time.Sleep(time.Duration(time.Second) * 3)
	log.Default().Println("end timeout")
	return nil
}

func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "world-handler",
		},
	}

	// add context timeout
	ctx, cannel := context.WithTimeout(ctx, 1*time.Second)
	defer cannel()

	// process at here
	go process(ctx)

	select {
	case <-ctx.Done():
		log.Default().Println("cancel timeout", ctx.Err())
		resp.StatusCode = http.StatusRequestTimeout
	}

	log.Default().Println("end timeout")

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
