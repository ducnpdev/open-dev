module github.com/dtm-labs/client

go 1.16

require (
	github.com/dtm-labs/dtmdriver v0.0.6
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-resty/resty/v2 v2.7.0
	go.mongodb.org/mongo-driver v1.9.1
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/dtm-labs/logger v0.0.1
	github.com/stretchr/testify v1.8.0 // indirect
)

retract v1.18.7
