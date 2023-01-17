package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	"github.com/jinzhu/copier"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv(envP string) error {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Default().Println("Getwd", dir)
	err = godotenv.Load(envP)
	if err != nil {
		panic(err)
	}
	return nil
}

func main() {
	lambda.Start(detailArticle)
	// LoadEnv(".env")
	// detailArticle(context.Background(), events.APIGatewayProxyRequest{})
}

type DataGetListArticle struct {
	// PathImage string `json:"pathImage"`
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

type GetListArticle struct {
	CommonRequest
	DataInput []ArticleDTO `json:"data"`
}
type ArticleDTO struct {
	Title       string `json:"title,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty" `
	Content     string `json:"content,omitempty"`
	Image       string `json:"image,omitempty"`
	Status      string `json:"status,omitempty"`
	Author      string `json:"author,omitempty"`
	Uuid        string `json:"uuid,omitempty"`
	IsTop       int    `json:"is_top"`
	Position    int    `json:"position"`
	CountView   int    `json:"count_view,omitempty"`
	UserID      uint   `json:"user_id,omitempty"`
	AppID       uint   `json:"app_id,omitempty"`
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

func detailArticle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		resp = events.APIGatewayProxyResponse{
			StatusCode: 400,
		}
		err    error
		uuid   string
		logger = log.Default()
	)
	fmt.Println("lambdaRequest:", fmt.Sprintf("%#v", req))
	requestId := req.RequestContext.RequestID
	fmt.Println("req.PathParameters:", fmt.Sprintf("%#v", req.PathParameters))
	fmt.Println("req.PathParameters:", fmt.Sprintf("%#v", req.QueryStringParameters))

	uuid = req.QueryStringParameters["uuid"]
	if uuid == "" {
		log.Println("uuid from query emtpy")
		uuid = "1673279461112223000"
	}
	lisst, err := getArticleDetailDB(ctx, uuid)
	if err != nil {
		logger.Println("get detail-article unmarshal body data err:", err)
		resp.Body = "get detail-article from database err"
		resp.StatusCode = 400
		return resp, nil
	}

	resp.Body = Response(err, requestId, lisst)
	resp.StatusCode = 200
	return resp, nil
}

type Article struct {
	gorm.Model
	Title       string `json:"title,omitempty" gorm:"type:varchar(1000)"`
	Slug        string `json:"slug,omitempty" gorm:"type:varchar(1000)"`
	Description string `json:"description,omitempty" gorm:"type:varchar(1000)" `
	Content     string `json:"content,omitempty" gorm:"type:text"`
	Image       string `json:"image,omitempty" gorm:"type:varchar(1000)"`
	Status      string `json:"status,omitempty" gorm:"type:varchar(1000);"`
	Author      string `json:"author,omitempty" gorm:"type:varchar(100);"`
	Uuid        string `json:"uuid,omitempty" gorm:"type:varchar(36);"`
	IsTop       int    `json:"is_top"`
	Position    int    `json:"position"`
	CountView   int    `json:"count_view,omitempty"`
	UserID      uint   `json:"user_id,omitempty"`
	AppID       uint   `json:"app_id,omitempty"`
}

func (u Article) Table() string {
	return "articles"
}

type GetImageS3Request struct {
	Data GetImageS3RequestData `json:"data"`
}
type GetImageS3RequestData struct {
	PathImage string `json:"pathImage"`
}

type GetImageS3Response struct {
	Data GetImageS3ResponseData `json:"data"`
}
type GetImageS3ResponseData struct {
	Url string `json:"url"`
}

func getArticleDetailDB(ctx context.Context, uuid string) (ArticleDTO, error) {
	var (
		article = Article{}
		err     error
		resp    = ArticleDTO{}
	)
	db, err := InitPostgres()
	if err != nil {
		return resp, err
	}
	defer func() {
		if db == nil {
			return
		}
		dbConfig, err := db.DB()
		if err == nil {
			dbConfig.Close()
		}
	}()

	err = db.Debug().Table(article.Table()).Where("uuid = ?", uuid).First(&article).Error
	if err != nil {
		return resp, err
	}
	copier.Copy(&resp, &article)
	return resp, nil

}

type Postgres struct {
	Username    string `yaml:"username" mapstructure:"username"`
	Password    string `yaml:"password" mapstructure:"password"`
	Database    string `yaml:"database" mapstructure:"database"`
	Host        string `yaml:"host" mapstructure:"host"`
	Port        int    `yaml:"port" mapstructure:"port"`
	Schema      string `yaml:"schema" mapstructure:"schema"`
	MaxIdleConn int    `yaml:"max_idle_conn" mapstructure:"max_idle_conn"`
	MaxOpenConn int    `yaml:"max_open_conn" mapstructure:"max_open_conn"`
}

// read environment of databse
func loadConfig() Postgres {
	user, _ := GetenvStr("DB_USER")
	dbpass, _ := GetenvStr("DB_PASS")
	dbhost, _ := GetenvStr("DB_HOST")
	dbport, _ := GetenvInt("DB_PORT")
	dbservice, _ := GetenvStr("DB_SERVICE")
	idleC, _ := GetenvInt("DB_MAX_IDLE_CONN")
	openC, _ := GetenvInt("DB_MAX_OPEN_CONN")
	schema, _ := GetenvStr("DB_SCHEMA")
	// password, _ := base64.StdEncoding.DecodeString(dbpass)
	return Postgres{
		Username:    user,
		Password:    dbpass,
		Database:    dbservice,
		Host:        dbhost,
		Port:        dbport,
		MaxIdleConn: idleC,
		MaxOpenConn: openC,
		Schema:      schema,
	}

}

func GetImageS3(ctx context.Context,
	req GetImageS3Request) (GetImageS3Response, error) {
	var (
		resp     = GetImageS3Response{}
		err      error
		httpResp *http.Response
	)
	chttp := NewClientHttp(0)

	httpResp, err = chttp.Post(ctx, ClientHttpRequest{
		Body: req,
		Url:  "https://qsgkg6u6n6.execute-api.ap-southeast-1.amazonaws.com/dev/image/get",
	})
	if err != nil {
		return resp, fmt.Errorf("http client post get-image-s3  err: %s", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("http client post get-image-s3  http status != 200, code: %d", httpResp.StatusCode)
	}

	bodyByte, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return resp, fmt.Errorf("http client post get-image-s3  marshal err: %s", err)
	}
	err = json.Unmarshal(bodyByte, &resp)
	if err != nil {
		return resp, fmt.Errorf("http client post get-image-s3  unmarshal err: %s", err)
	}
	return resp, nil
}

// create database postgres instance
func InitPostgres() (*gorm.DB, error) {
	log.Default().Println("connecting postgres database")
	config := loadConfig()
	fmt.Println("config postgres:", fmt.Sprintf("%#v", config))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d ", config.Host, config.Username, config.Password, config.Database, config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Default().Println("connect postgres err:", err)
		return db, err
	}
	dbc, _ := db.DB()
	dbc.SetConnMaxIdleTime(1)
	dbc.SetConnMaxLifetime(time.Duration(time.Second * 15))
	log.Default().Println("connect postgres successfully")
	return db, err
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

// get env with type float
func GetenvFloat64(key string) (float64, error) {
	s, err := GetenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// get env with type boolean
func GetenvBool(key string) (bool, error) {
	s, err := GetenvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

type ClientHttpRequest struct {
	Body        interface{}
	Method      string
	Url         string
	Header      map[string]string
	ContentType string
	Query       string
}

const (
	// contentType
	MimeJSON  = "application/json"
	URLEncode = "application/x-www-form-urlencoded"
	// timeDuration
	TimeoutHttp = 30 * time.Second
)

type ClientHttp interface {
	Post(ctx context.Context, req ClientHttpRequest) (*http.Response, error)
	// Get(ctx context.Context, req ClientHttpRequest) (*http.Response, error)
}

type clientHttp struct {
	client *http.Client
}

func timeoutHttp(timeout uint8) time.Duration {
	if timeout == 0 {
		return time.Duration(TimeoutHttp)
	}
	return time.Duration(timeout) * time.Second
}

func NewClientHttp(timeout uint8) ClientHttp {
	client := &http.Client{
		Timeout:   timeoutHttp(timeout),
		Transport: getTransport(),
	}

	return &clientHttp{
		client: client,
	}
}

func getTransport() *http.Transport {
	tr := &http.Transport{
		// MaxIdleConns:       10,
		// IdleConnTimeout:    30 * time.Second,
		// DisableCompression: true,
	}
	return tr
}

// build common header
func buildHeader(mapHeader map[string]string) (header http.Header) {
	header = make(http.Header)
	for key, value := range mapHeader {
		header.Set(key, value)
	}
	return header
}

// build body of api
func buildBody(ctx context.Context, contentType string, bodyReq interface{}) (*bytes.Reader, error) {
	var (
		body     *bytes.Reader
		err      error
		bodyByte []byte
	)
	switch contentType {
	default:
		bodyByte, err = json.Marshal(bodyReq)
	}
	if err != nil {
		return body, err
	}
	body = bytes.NewReader(bodyByte)
	return body, err
}

// build request data of http
func buildRequestHttp(ctx context.Context, req ClientHttpRequest) (*http.Request, error) {
	var (
		httpReq *http.Request
		err     error
	)
	body, err := buildBody(ctx, req.ContentType, req.Body)
	if err != nil {
		return httpReq, err
	}
	httpReq, err = http.NewRequestWithContext(ctx, req.Method, req.Url, body)
	if err != nil {
		return httpReq, err
	}
	httpReq.Header = buildHeader(req.Header)
	if req.Query != "" {
		httpReq.URL.RawQuery = req.Query
	}
	return httpReq, err
}

// post api
func (h *clientHttp) Post(ctx context.Context,
	req ClientHttpRequest) (httpResp *http.Response, err error) {
	req.Method = http.MethodPost
	reqhttp, err := buildRequestHttp(ctx, req)
	if err != nil {
		return httpResp, err
	}
	httpResp, err = h.client.Do(reqhttp)
	return httpResp, err
}
