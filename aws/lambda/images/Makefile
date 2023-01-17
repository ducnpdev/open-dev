.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/insertImage functions/insertImage/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/articles functions/article/list.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/articleDetail functions/articleDetail/detail.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getImage functions/getImage/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
