# open-dev

## Contents

- [open-dev](#open-dev)
	- [Contents](#contents)
	- [Let's Go](#lets-go)
	- [Gin Web Framework](#gin-web-framework)
		- [API Examples](#api-examples)
	- [Queue](#queue)
		- [Kafka](#kafka)
			- [Consumer Advanced](#consumer-advanced)
			- [Producer Advanced](#producer-advanced)
			- [Example](#example)
		- [RabbitMQ](#rabbitmq)
			- [Code Example](#code-example)
	- [Redis](#redis)
		- [Rate Limit](#rate-limit)
		- [redis pub-sub](#redis-pub-sub)
	- [Golang Usecase](#golang-usecase)
		- [Strategy Management Connect database](#strategy-management-connect-database)
		- [Resize Image](#resize-image)
		- [Context](#context)
		- [RSA](#rsa)
		- [ecdsa algorithm](#ecdsa-algorithm)
		- [concurrency](#concurrency)
			- [pattern](#pattern)
	- [Serverless Framework](#serverless-framework)
	- [terraform](#terraform)
		- [sqs](#sqs)
		- [aws-apigw](#aws-apigw)
	- [Aws](#aws)
		- [Textract](#textract)
		- [Lambda](#lambda)
		- [S3](#s3)
	- [Performances](#performances)
		- [Standard](#standard)
	- [http](#http)
		- [reuse-http](#reuse-http)
	- [tls](#tls)
		- [tls-base64](#tls-base64)
	- [cryto](#cryto)
		- [most commonly](#most-commonly)
	- [Contact:](#contact)
	- [Give a Star! ⭐](#give-a-star-)
	- [Buy me a coffee](#buy-me-a-coffee)

## Let's Go  
- build web application golang
- link: https://github.com/ducnpdev/open-dev/tree/master/letgo

## Gin Web Framework
### API Examples
- source code in forder gin-web-framework

## Queue
### Kafka
#### Consumer Advanced
- read document vietnamese: https://github.com/ducnpdev/open-dev/blob/master/kafka/README_VI.md#consumer-advanced
#### Producer Advanced
- read document vietnamese: https://github.com/ducnpdev/open-dev/blob/master/kafka/README_VI.md#topic-advanced
#### Example
- read article install kafka: https://github.com/ducnpdev/open-dev/tree/master/kafka
### RabbitMQ
#### Code Example
- read source: https://github.com/ducnpdev/golang-rabbitmq

## Redis
- read article install redis: https://github.com/ducnpdev/open-dev/tree/master/redis
### Rate Limit
- read article lear more rate-limit: https://viblo.asia/p/golang-ratelimit-la-gi-su-dung-ByEZkn3qKQ0

### redis pub-sub
- pub-sub code example: https://github.com/ducnpdev/open-dev/tree/master/redis/pubsub

## Golang Usecase

### Strategy Management Connect database
- link: https://viblo.asia/p/golang-chien-luoc-quan-ly-ket-noi-database-thong-qua-gorm-zXRJ8r6dVGq

### Resize Image
- edit size of image from base64: https://github.com/ducnpdev/open-dev/tree/master/usecase#resize-image


### Context
- Handle context timeout: https://github.com/ducnpdev/open-dev/blob/master/usecase/context/timeout.go

### RSA
- Encrypt Decrypt Data, code example: https://github.com/ducnpdev/open-dev/blob/master/usecase/rsa/REAME.md

### ecdsa algorithm
- This is code example: https://github.com/ducnpdev/open-dev/blob/master/usecase/ecdsa/ecdsa.go

### concurrency
- source-code example:
#### pattern
- 7 pattern thiết yếu trong golang: https://github.com/ducnpdev/open-dev/tree/master/concurrency/patterns

## Serverless Framework
## terraform
### sqs
- create sqs simple: https://github.com/ducnpdev/open-dev/tree/master/terraform/apps/sqs
### aws-apigw
- this is code: https://github.com/ducnpdev/open-dev/tree/master/terraform/apps/apigw
## Aws

### Textract 
- code example use aws textract identity: https://docs.aws.amazon.com/textract/index.html
- code: https://github.com/ducnpdev/open-dev/blob/master/aws/textract/main.go
### Lambda
- Demo simple lambda function: https://github.com/ducnpdev/open-dev/tree/master/aws/lambda
- Crud with postgres, link source: https://github.com/ducnpdev/open-dev/tree/master/aws/lambda/crud
- build docker images: https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/

### S3
- link code examle: https://github.com/ducnpdev/open-dev/tree/master/aws/s3

## Performances
### Standard
- Test simple of function return error: https://opendev.hashnode.dev/golang-test-performance-function-standard-1
- Test convert string to int of 3 function: https://opendev.hashnode.dev/golang-test-performance-function-standard-1

## http
### reuse-http
```go
package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
)

func main() {
	Reuse()
	//
	NonReuse()
}

// NonReuse, not reuse http
func NonReuse() {
	// client trace to log whether the request's underlying tcp connection was re-used
	clientTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf("conn was reused: %t", info.Reused)
		},
	}
	traceCtx := httptrace.WithClientTrace(context.Background(), clientTrace)

	// 1st request
	req, err := http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// 2nd request
	req, err = http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

// Reuse, reuse http client
func Reuse() {
	var (
		err error
	)
	// client trace to log whether the request's underlying tcp connection was re-used
	clientTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf("conn was reused: %t", info.Reused)
		},
	}
	traceCtx := httptrace.WithClientTrace(context.Background(), clientTrace)

	// 1st request
	req, err := http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	// 2nd request
	req, err = http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
```

## tls
### tls-base64
- convert file key and file pem to base64
- link: https://github.com/ducnpdev/open-dev/tree/master/tls/base64parse
  
## cryto
### most commonly
- link: https://github.com/ducnpdev/open-dev/tree/master/usecase/crypto-demo

## Contact:
- facebook: https://www.facebook.com/phucducdev/
- gmail: ducnp09081998@gmail.com or phucducktpm@gmail.com
- linkedin: https://www.linkedin.com/in/phucducktpm/
- hashnode: https://hashnode.com/@OpenDev
- telegram: https://t.me/OpenDevGolang

## Give a Star! ⭐

If you like or are using this project to learn or start your solution, please give it a star. Thanks!

## Buy me a coffee