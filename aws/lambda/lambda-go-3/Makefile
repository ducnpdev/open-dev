# .PHONY: build clean deploy

build:
	env GOARCH=amd64 CGO_ENABLED="0" GOOS=linux go build -ldflags="-s -w" -o bin/create create/main.go
	env GOARCH=amd64 CGO_ENABLED="0" GOOS=linux go build -ldflags="-s -w" -o bin/update update/main.go
	env GOARCH=amd64 CGO_ENABLED="0" GOOS=linux go build -ldflags="-s -w" -o bin/delete delete/main.go
	env GOARCH=amd64 CGO_ENABLED="0" GOOS=linux go build -ldflags="-s -w" -o bin/read read/main.go
clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
