.PHONY: build clean deploy

build:
	# env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/migration/x86 migration/main.go
	# env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/bootstrap migration/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/bootstrap migration/main.go
 
clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

zip:
    zip -j bin/migration/x86.zip bin/x86

buildv2:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/migration/bootstrap migration/main.go
zip:
	zip -j bin/migration.zip bin/migration/bootstrap
