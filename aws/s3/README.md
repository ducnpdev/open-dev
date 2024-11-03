# S3
## static website
- bucket policy
```
{
    "Version":"2012-10-17",
    "Statement":[
        {
            "Sid":"PublicReadGetObject",
            "Effect":"Allow",
            "Principal": "*",
            "Action":[
                "s3:GetObject"
            ],
            "Resource":[
                "arn:aws:s3:::bucketName/*"
            ]
        }
    ]
}
```

## Code Upload & Delete File
- config:
```go
package s3

type AwsS3 struct{}
type AwsS3Config struct {
	Bucket           string
	Region           string
	Endpoint         string
	SecretID         string
	SecretKey        string
	BaseURL          string
	PathPrefix       string
	S3ForcePathStyle bool
	DisableSSL       bool
}

var configS3 = AwsS3Config{
	Bucket:           "ducnp5",
	Region:           "ap-southeast-1",
	Endpoint:         "endpoint",
	SecretID:         "xxx",
	SecretKey:        "xxx",
	BaseURL:          "opendev",
	PathPrefix:       "",
	S3ForcePathStyle: true,
	DisableSSL:       true,
}
```

- upload
```go

func (*AwsS3) UploadFile(file *multipart.FileHeader) (string, string, error) {
	session := session.Must(session.NewSession())

	uploader := s3manager.NewUploader(session)

	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	filename := configS3.PathPrefix + "/" + fileKey
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close()

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(configS3.Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		return "", "", err
	}

	return configS3.BaseURL + "/" + filename, fileKey, nil
}
```

- delete
```go

func (*AwsS3) DeleteFile(key string) error {
	session := session.Must(session.NewSession())
	svc := s3.New(session)
	filename := configS3.PathPrefix + "/" + key
	bucket := configS3.Bucket

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return errors.New("function svc.DeleteObject() failed, err:" + err.Error())
	}

	_ = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	return nil
}
```