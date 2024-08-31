.PHONY: build clean zip deploy

build:
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/hello/bootstrap hello/main.go
clean:
	rm -rf ./bin
deploy: clean build zip
	sls deploy --verbose
zip:
	zip -j bin/hello.zip bin/hello/bootstrap
