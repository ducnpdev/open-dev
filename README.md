# open-dev

## Contents

- [open-dev](#open-dev)
	- [Contents](#contents)
	- [Gin Web Framework](#gin-web-framework)
		- [API Examples](#api-examples)
	- [Redis](#redis)
		- [Rate Limit](#rate-limit)
	- [Golang Usecase](#golang-usecase)
		- [Resize Image](#resize-image)
		- [Context](#context)
		- [RSA](#rsa)
		- [ecdsa algorithm](#ecdsa-algorithm)
		- [concurrency](#concurrency)
	- [Serverless Framework](#serverless-framework)
	- [terraform](#terraform)
		- [sqs](#sqs)
		- [aws-apigw](#aws-apigw)
	- [Aws](#aws)
		- [Textract](#textract)
		- [Lambda](#lambda)
	- [Performances](#performances)
		- [Standard](#standard)
	- [http](#http)
		- [reuse-http](#reuse-http)
	- [Contact:](#contact)
## Gin Web Framework
### API Examples
- source code in forder gin-web-framework

## Redis
- read article install redis: https://github.com/ducnpdev/open-dev/tree/master/redis
### Rate Limit
- read article lear more rate-limit: https://viblo.asia/p/golang-ratelimit-la-gi-su-dung-ByEZkn3qKQ0

## Golang Usecase

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
## Contact:
- facebook: https://www.facebook.com/phucducdev/
- gmail: ducnp09081998@gmail.com or phucducktpm@gmail.com
- linkedin: https://www.linkedin.com/in/phucducktpm/
- hashnode: https://hashnode.com/@OpenDev
- telegram: https://t.me/OpenDevGolang