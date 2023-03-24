package main

import (
	"fmt"
	"open-dev/aws/apiwg"
	"time"
)

func main() {
	// kafka.MainKafka() /
	// Example2()
	// MulChannel()
	// concurrency.Example1Once()
	// concurrency.Example2Once()
	// apiwg.CreateAPIKey()
	apiwg.DeleteAPIKey()

}

func SelectDefault() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}
func Example2() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c) // (1)
	}()
	fmt.Println("Blocking on read...")
	select {
	case <-c: // (2)
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

// multiple channel read
func MulChannel() {
	ch1 := make(chan interface{})
	close(ch1)
	ch2 := make(chan interface{})
	close(ch2)
	var ch1Count, ch2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-ch1:
			ch1Count++
		case <-ch2:
			ch2Count++
		}
	}
	fmt.Printf("ch1Count: %d\n ch2Count: %d\n", ch1Count, ch2Count)
}

//

/*
---
*/

// package main

// import (
// 	"bytes"
// 	"context"
// 	"crypto/x509"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/google/uuid"
// )

// type MerchantInfo struct {
// 	MerchantName string `json:"merchantName"`
// 	MerchantCode string `json:"merchantCode"`
// }
// type ClientHttpRequest struct {
// 	Body        interface{}
// 	Method      string
// 	Url         string
// 	Header      map[string]string
// 	ContentType string
// 	Query       string
// }

// func main() {
// 	caCert, err := ioutil.ReadFile("rootCA.crt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)

// 	cl := NewClientHttp(10)
// 	req := ClientHttpRequest{
// 		Body: Test{
// 			Request: uuid.NewString(),
// 			Data: Us{
// 				User: "0774641820",
// 			},
// 		},
// 		Method: "POST",
// 		Url:    "https://openbanking.hdbank.com.vn/vjacargo/partner-service/authen",
// 		Header: map[string]string{},
// 	}
// 	a, b := cl.Post(context.Background(), req)
// 	fmt.Println(a, b)
// 	// convertStrToMapStruct()
// 	// TestBuildRequest()
// 	// utils.GetTranID()
// 	// fmt.Println(time.Now().Format("02012006"))
// 	// fmt.Println(time.Now().Format(time.RFC3339Nano))

// 	// sdf := utils.ValidatorUUID("\\uu 53e3bbfb-493d-4eae-a36c-e336fde30f85")
// 	// fmt.Println(sdf)

// 	// timeNow := time.Now()
// 	// year, m, day := timeNow.Date()
// 	// fmt.Println(year, day, int(m))
// 	// ctx := contexts.Background()
// 	// connectRedis(ctx)
// 	// TransactionId := "EPASS29Y13013"
// 	// fmt.Println(len(TransactionId))
// 	// if len(TransactionId) < 6 || len(TransactionId) > 30 {
// 	// 	fmt.Println("sssssss")
// 	// }
// 	// s := string(`{"responseCode":"01","responseMessage":"signature invalid","responseTime":"2022-12-15T13:38:23.455+07:00"}`)
// 	// data := ResponseData{}
// 	// err := json.Unmarshal([]byte(s), &data)
// 	// if err != nil {
// 	// 	return
// 	// }
// 	// fmt.Println("ResponseCode: ", data.ResponseCode)
// 	// fmt.Println("ResponseMessage: ", data.ResponseMessage)
// 	// fmt.Println("ResponseTime: ", data.ResponseTime)
// 	// qrcode.GenQrCode(context.Background(), qrcode.QRCodeRequest{
// 	// 	MerchantId:  "000000000000000",
// 	// 	TerminalId:  "06800972-indicator=virtual",
// 	// 	AccountNo:   "QR220220230856131KZ",
// 	// 	AcqId:       "970437",
// 	// 	AccountName: "TRAN BICH THUY",
// 	// 	AddInfo:     "",
// 	// 	Amount:      50000,
// 	// })

// 	// fmt.Println(usecase.ParserCreatedTime("2022-11-30T11:40:50.201253Z"))

// }

// //	func timeZone() {
// //		createdDate := "2022-10-22T08:40:55+07:00"
// //		created, b := utils.TimeParse("", createdDate)
// //		fmt.Println(b)
// //		sec := time.Since(created).Seconds()
// //		secInt := int(sec)
// //		fmt.Println("timeExpire:", secInt)
// //	}
// type Test struct {
// 	Request string `json:"requestId"`
// 	Data    Us
// }
// type Us struct {
// 	User string `json:"userId"`
// }

// func TestBuildRequest() {
// 	ctx := context.Background()
// 	req := ClientHttpRequest{
// 		Body: Test{
// 			Request: uuid.NewString(),
// 			Data: Us{
// 				User: "0774641820",
// 			},
// 		},
// 		Method: "POST",
// 		Url:    "https://openbanking.hdbank.com.vn/vjacargo/partner-service/authen",
// 		Header: map[string]string{},
// 	}
// 	reqhttp, err := buildRequestHttp(ctx, req)
// 	if err != nil {
// 		fmt.Println(reqhttp, err)
// 	}
// 	// httpResp, err = c.Do(reqhttp)
// 	// return httpResp, err
// }

// const (
// 	// contentType
// 	MimeJSON  = "application/json"
// 	URLEncode = "application/x-www-form-urlencoded"
// 	// timeDuration
// 	TimeoutHttp = 30 * time.Second
// )

// type ClientHttp interface {
// 	Post(ctx context.Context, req ClientHttpRequest) (*http.Response, error)
// }

// type clientHttp struct {
// 	client *http.Client
// }

// func timeoutHttp(timeout int) time.Duration {
// 	if timeout == 0 {
// 		return time.Duration(TimeoutHttp)
// 	}
// 	return time.Duration(timeout) * time.Second
// }

// func NewClientHttp(timeout int) ClientHttp {
// 	client := &http.Client{
// 		Timeout:   timeoutHttp(timeout),
// 		Transport: getTransport(),
// 	}

// 	return &clientHttp{
// 		client: client,
// 	}
// }

// func getTransport() *http.Transport {
// 	tr := &http.Transport{
// 		// MaxIdleConns:       10,
// 		// IdleConnTimeout:    30 * time.Second,
// 		// DisableCompression: true,
// 	}
// 	return tr
// }

// // build common header
// func buildHeader(mapHeader map[string]string) (header http.Header) {
// 	header = make(http.Header)
// 	for key, value := range mapHeader {
// 		header.Set(key, value)
// 	}
// 	return header
// }

// // build body of api
// func buildBody(ctx context.Context, contentType string, bodyReq interface{}) (*bytes.Reader, error) {
// 	var (
// 		body     *bytes.Reader
// 		err      error
// 		bodyByte []byte
// 	)
// 	switch contentType {
// 	default:
// 		bodyByte, err = json.Marshal(bodyReq)
// 	}
// 	if err != nil {
// 		return body, err
// 	}
// 	body = bytes.NewReader(bodyByte)
// 	return body, err
// }

// // build request data of http
// func buildRequestHttp(ctx context.Context, req ClientHttpRequest) (*http.Request, error) {
// 	var (
// 		httpReq *http.Request
// 		err     error
// 	)
// 	body, err := buildBody(ctx, req.ContentType, req.Body)
// 	if err != nil {
// 		return httpReq, err
// 	}
// 	httpReq, err = http.NewRequestWithContext(ctx, req.Method, req.Url, body)
// 	if err != nil {
// 		return httpReq, err
// 	}
// 	httpReq.Header = buildHeader(req.Header)
// 	if req.Query != "" {
// 		httpReq.URL.RawQuery = req.Query
// 	}
// 	return httpReq, err
// }

// // post api
// func (h *clientHttp) Post(ctx context.Context,
// 	req ClientHttpRequest) (httpResp *http.Response, err error) {
// 	req.Method = http.MethodPost
// 	reqhttp, err := buildRequestHttp(ctx, req)
// 	if err != nil {
// 		return httpResp, err
// 	}
// 	httpResp, err = h.client.Do(reqhttp)
// 	fmt.Println(httpResp, err)
// 	return httpResp, err
// }
