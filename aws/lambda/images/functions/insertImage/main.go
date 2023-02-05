package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	lambda.Start(insertImage)
}

type DataSaveImageS3 struct {
	Image     string `json:"image"` // base64
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	FileName  string `json:"file_name"`
	PartnerID string `json:"partner_id"`
	UserID    string `json:"user_id"`
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

type SaveImageS3 struct {
	CommonRequest
	DataInput DataSaveImageS3 `json:"data"`
}

func insertImage(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		resp = events.APIGatewayProxyResponse{
			StatusCode: 400,
		}
		err    error
		logger = log.Default()
		reqDTO = SaveImageS3{}
	)
	fmt.Println("insert image")
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
	imageBase64 := strings.TrimSpace(reqDTO.DataInput.Image)
	if imageBase64 == "" {
		resp.Body = "base64 required"
		return resp, nil
	}
	pathImage := os.Getenv("IMAGE_BLOG") + reqDTO.DataInput.Phone + reqDTO.DataInput.FileName + reqDTO.DataInput.Name // image.FileNameS3(os.Getenv("FOLDER_INDEX"), reqDTO.DataInput.CCCD, fmt.Sprintf("%s", reqDTO.DataInput.FileName))
	logger.Printf("%s, start store image at pathImage %s \n", requestId, pathImage)
	_, dataStore, err := StoreS3(ctx, StoreS3Req{
		BaseImage:      &imageBase64,
		PathImage:      pathImage,
		CustomEndpoint: false,
	})
	if err != nil {
		logger.Printf("%s, store image err %s \n", requestId, err.Error())
		resp.Body = "store image err"
		resp.StatusCode = 400
		return resp, nil
	}

	fmt.Println(pathImage, dataStore)

	resp.Body = "store image successfully"
	resp.StatusCode = 200
	return resp, nil
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
