# .PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/push push/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/receive receive/main.go
clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
