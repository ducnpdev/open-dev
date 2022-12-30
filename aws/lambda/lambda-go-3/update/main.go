package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"open-dev/aws/lambda/lambda-go-3/pkg"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

type RequestBodyAPIGW struct {
	RequestID string  `json:"requestId"`
	Data      UserDTO `json:"data"`
}

type UserDTO struct {
	ID    int    `json:"userId"`
	Name  string `json:"name"`
	User  string `json:"userName"`
	Phone string `json:"phone"`
}

func Update(ctx context.Context, eventReq events.APIGatewayProxyRequest) (Response, error) {
	var (
		req  = RequestBodyAPIGW{}
		resp = Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)
	err := json.Unmarshal([]byte(eventReq.Body), &req)
	if err != nil {
		resp.Body = ParseResponse(HttpResponse{
			Uuid: req.RequestID,
			Err:  err,
		})
		// return resp, nil
		req.Data.Name = "update_duc"
		req.Data.ID = 2
	}
	db, err := pkg.InitPostgres()

	if err != nil {
		resp.Body = ParseResponse(HttpResponse{
			Uuid: req.RequestID,
			Err:  err,
		})
		resp.StatusCode = 500
		return resp, nil
	}
	// set http-code 200
	resp.StatusCode = 200
	// save new user
	err = db.Exec(`update game_tet.users set name = ? where id = ?`, req.Data.Name, req.Data.ID).Error
	if err != nil {
		resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Err: err})
		return resp, nil
	}
	resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Data: nil})
	return resp, nil
}

func main() {
	// lambda.Start(Update)
	Update(context.Background(), events.APIGatewayProxyRequest{})
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
