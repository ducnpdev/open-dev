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
