# golang migration runtime go1.x -> provided.al2

## runtime go1.x
- make file
```
env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o x86 migration/main.go
```
- serverless
```
service: golang-migration
frameworkVersion: '3'

provider:
  name: aws
  runtime: go1.x
  region: ap-southeast-1

functions:
  hello:
    handler: x86
    events:
      - httpApi:
          path: /x86
          method: get

```

## runtime provided.al2
- make file
```
env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/bootstrap migration/main.go
```
- serverless
```
service: golang-migration
frameworkVersion: '3'

provider:
  name: aws
  runtime: provided.al2 # <- change from go1.x to provided.al2
  architecture: arm64   # <- change from x86_64 to arm64
  region: ap-southeast-1

functions:
  hello:
    # handler: bin/x86 # old
    handler: bootstrap # new
    events:
      - httpApi:
          # path: /x86 # old
          path: /arm64 # new
          method: get
```