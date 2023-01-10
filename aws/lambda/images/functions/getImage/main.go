package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	lambda.Start(getImage)
}

type DataGetImageRequest struct {
	PathImage string `json:"pathImage"`
}
type CommonRequest struct {
	RequestId   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
}
type CommonResponse struct {
	RequestId       string `json:"requestId"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

type GetImageRequest struct {
	CommonRequest
	DataInput DataGetImageRequest `json:"data"`
}

func CreateS3Client(req AwsReq) *s3.S3 {
	sess := CreateSession(req)
	client := s3.New(sess)
	return client
}
func PresignUrl(reqAws AwsReq, presign PresignUrlReq) (string, error) {
	var (
		url string
		err error
	)

	clientS3 := CreateS3Client(reqAws)
	req, _ := clientS3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(presign.Bucket),
		Key:    aws.String(presign.KeyPath),
	})
	url, err = req.Presign(presign.UrlTtl)
	if err != nil {
		return url, fmt.Errorf("presign url err %s", err)
	}
	return url, err
}

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func GetenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}
func GetenvInt(key string) (int, error) {
	s, err := GetenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

type PresignedResponse struct {
	Url string `json:"url"`
}

func getImage(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		resp = events.APIGatewayProxyResponse{
			StatusCode: 400,
		}
		err    error
		logger = log.Default()
		reqDTO = GetImageRequest{}
	)
	err = json.Unmarshal([]byte(req.Body), &reqDTO)
	requestId := strings.TrimSpace(reqDTO.RequestId)
	if requestId == "" {
		requestId = req.RequestContext.RequestID
	}
	logger.Printf("%s, \n", requestId)
	if err != nil {
		logger.Println("index face unmarshal body data err:", err)
		resp.Body = "index face unmarshal body data err"
		return resp, nil
	}
	ttl, _ := GetenvInt("IMAGE_URL_TTL")
	if ttl == 0 {
		ttl = 60
	}
	presignReq := PresignUrlReq{
		Bucket:  GetS3Bucket(),
		KeyPath: reqDTO.DataInput.PathImage,
		UrlTtl:  time.Second * time.Duration(ttl),
	}

	var presignUmage = PresignedResponse{}

	presignUmage.Url, err = PresignUrl(AwsReq{}, presignReq)
	if err != nil {
		resp.Body = "get image error"
		resp.StatusCode = 400
		return resp, nil
	}
	resp.Body = Response(err, requestId, presignUmage)
	resp.StatusCode = 200
	return resp, nil
}

func Response(err error, uuid string, data interface{}) string {
	defer func() {
	}()
	respTime := time.Now().Format(time.RFC3339)
	var buf bytes.Buffer
	if err != nil {
		body, errMarshal := json.Marshal(map[string]interface{}{
			"responseId":   uuid,
			"responseTime": respTime,
			// "responseCode":    ParseError(err).Code(),
			// "responseMessage": ParseError(err).Message(),
		})
		if errMarshal != nil {
			log.Default().Println("marshal response err", errMarshal)
		}
		json.HTMLEscape(&buf, body)
		return buf.String()
	}
	mapRes := map[string]interface{}{
		"responseId":      uuid,
		"responseCode":    "00",
		"responseMessage": "successfully",
		"responseTime":    respTime,
	}
	fmt.Println("data response", fmt.Sprintf("%#v", data))
	if data != nil {
		mapRes["data"] = data
	}

	body, errMarshal := json.Marshal(mapRes)
	if errMarshal != nil {
		log.Default().Println("marshal response err", errMarshal)
	}
	json.HTMLEscape(&buf, body)
	return buf.String()
}

type StoreS3Req struct {
	BaseImage      *string // base64 of image
	PathImage      string  // path of image in s3
	CustomEndpoint bool
}

type PresignUrlReq struct {
	Bucket  string
	KeyPath string
	UrlTtl  time.Duration
}

// fn upload image to s3, request iamge is base64
func CheckSizeBase64(baseImage *string) error {
	sizeByte := base64.StdEncoding.DecodedLen(len(*baseImage))
	sizeRe := 5000000
	if sizeByte > sizeRe {
		err := fmt.Errorf("size image large, require size %d and size input %d", sizeRe, sizeByte)
		return err
	}
	return nil
}

type StoreS3Reponse struct {
	PathImage string `json:"pathImage"`
}

func StoreS3(ctx context.Context, reqDTO StoreS3Req) (*s3.PutObjectOutput, StoreS3Reponse, error) {
	var (
		dataRes = StoreS3Reponse{}
		err     error
		putOut  = &s3.PutObjectOutput{}
	)
	if reqDTO.BaseImage == nil {
		return putOut, dataRes, fmt.Errorf("base64 of image is require")
	}
	if len(*reqDTO.BaseImage) == 0 {
		return putOut, dataRes, fmt.Errorf("byte image decode empty")
	}
	splitBase := strings.Split(*reqDTO.BaseImage, "base64,")
	var newBase string
	if len(splitBase) > 1 {
		newBase = splitBase[1]
	} else {
		newBase = splitBase[0]
	}
	byteImage, err := base64.StdEncoding.DecodeString(newBase)
	if err != nil {
		fmt.Println("image decode string err", err)
		return putOut, dataRes, err
	}
	// if len(*baseImage) == 0 {
	// 	return putOut, dataRes, fmt.Errorf("byte image decode empty")
	// }

	mimeType := http.DetectContentType(byteImage)
	fmt.Println("mime type:", mimeType)
	if !strings.Contains(mimeType, "jpeg") && !strings.Contains(mimeType, "png") {
		return putOut, dataRes, fmt.Errorf("mime type error")
	}

	byteReader := bytes.NewReader(byteImage)

	dataRes.PathImage = reqDTO.PathImage
	session := CreateSession(AwsReq{
		CustomEndpoint: reqDTO.CustomEndpoint,
	})
	putOut, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Body:        byteReader,
		Bucket:      aws.String(GetS3Bucket()),
		Key:         aws.String(dataRes.PathImage),
		ContentType: &mimeType,
	})

	return putOut, dataRes, err
}

type AwsReq struct {
	CustomEndpoint bool
	AssumeRole     bool
}

func CreateSession(req AwsReq) *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}
func GetS3Bucket() string {
	bucket := os.Getenv("S3_BUCKET")
	if bucket != "" {
		return bucket
	}
	log.Default().Println("GetS3Bucket default value:", "pkg.BucketDefault")
	return "pkg.BucketDefault"
}
