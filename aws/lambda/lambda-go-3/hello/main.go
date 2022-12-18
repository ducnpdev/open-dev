package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

type RequestBodyAPIGW struct {
	RequestID string      `json:"requestId"`
	Data      interface{} `json:"data"`
}

// type ResponseBodyAPIGW struct {
// 	RequestID string `json:"request_id"`
// 	Message   string `json:"message"`
// }

func Handler(ctx context.Context, eventReq events.APIGatewayProxyRequest) (Response, error) {
	var (
		req  = RequestBodyAPIGW{}
		resp = Response{
			StatusCode:      404,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type":           "application/json",
				"X-MyCompany-Func-Reply": "hello-handler",
			},
		}
	)
	err := json.Unmarshal([]byte(eventReq.Body), &req)
	if err != nil {
		resp.Body = ParseResponse(HttpResponse{
			Uuid: req.RequestID,
			Err:  err})
		return resp, nil
	}
	resp.StatusCode = 200
	resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Data: req.Data})
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

type HttpResponse struct {
	Uuid string // uuid, indicator per api
	Err  error
	Time string // time tracing
	Data interface{}
}

func ParseResponse(respBody HttpResponse) string {
	respBody.Time = time.Now().Format("2006-01-02T15:04:05.000-07:00")
	if respBody.Err != nil {
		return responseErr(respBody)
	}
	return responseOk(respBody)
}

func responseOk(respBody HttpResponse) string {
	var buf bytes.Buffer
	mapRes := map[string]interface{}{
		"responseId":      respBody.Uuid,
		"responseMessage": "successfully",
		"responseTime":    respBody.Time,
	}
	if respBody.Data != nil {
		mapRes["data"] = respBody.Data
	}
	body, errMarshal := json.Marshal(mapRes)
	if errMarshal != nil {
		log.Default().Println("marshal response err", errMarshal)
	}
	json.HTMLEscape(&buf, body)
	return buf.String()
}

func responseErr(respBody HttpResponse) string {
	var buf bytes.Buffer
	mapRes := map[string]interface{}{
		"responseId":      respBody.Uuid,
		"responseMessage": respBody.Err.Error(),
		"responseTime":    respBody.Time,
	}

	body, errMarshal := json.Marshal(mapRes)
	if errMarshal != nil {
		log.Default().Println("marshal response err", errMarshal)
	}
	json.HTMLEscape(&buf, body)
	return buf.String()
}
