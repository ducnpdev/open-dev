# Folder Structure Conventions
> Folder structure options and naming conventions for software projects

## A typical top-level directory layout

    .
    ├── charts                  # định nghĩa cho helm chart, làm việc với k8s
    ├── cmd                     # command definition
        ├── server              # run server
            main.go
        ├── client              # run client
            main.go
    ├── config                  # all config per environment   
        ├── config.go 
        ├── dev.yaml            # environment develop
        ├── uat.yaml            # environment uat
        ├── prod.yaml           # environment prod
    ├── docs                    # file swagger
        ├── docs.go             # auto generation
        ├── swagger.json        # auto generation
        ├── swagger.yaml        # auto generation
    ├── internal                # handle logic api
        ├── adapter
            ├── partnerservice
                user.go
            ...
        ├── facade
            ├── payment
            ├── shipping
            user.go
            ...
        ├── api
            user.go
            product.go
            ...
        ├── usecase
            ├── payment
            ├── shipping
            user.go
            ...
        ├── repository
            ├── models
            ├── postgres
            ├── redis
            user.go
        ├── common
            const.go
            error.go
            ...
        ├── directory
            user.go
            product.go
    ├── pkg                     # reuseable
        ├── aws                 
            ├── lambda
            ├── sqs
            ...
        ├── database            # database code
            ├── mongo
                ├── connect.go
            ├── postgres
            ...
        ├── logger              # init logger, log level config
            ├── zap
            ├── logrus
            log.go
            ...
        ├── utils
    ├── protoc                  # RPC message
        ├── api.proto           # definition struct proto api, message
    ├── sql                     # migration sql or script sql
    ├── build                   # Compiled files (alternatively `dist` in nextjs)
    ├── test                    # Automated tests (alternatively `spec` or `tests`)
    ├── third_party             # Tools and utilities (alternatively `tools`)
    ├── LICENSE
    └── README.md
    └── go.mod
    └── Dockerfile
    └── Jenkinsfile
    └── Makefile

> Use short lowercase names at least for the top-level files and folders except
> `LICENSE`, `README.md`

## explain
- charts: này không bắt buộc nhưng tuỳ thuộc vào hệ thống sẽ cấu hình helm chung repo hoặc khác.
- cmd:
  - server: để chạy server bằng file main
  - client: để chạy client bằng file main
  - trong source code mình dùng cmd, vì đơn giản rất nhiều source code open source dùng nó.
  ```
  https://github.com/golang-standards/project-layout
  ```
- config: folder để load config từ file yaml hoặc env(tuỳ service)
- internal: này thấy google với một số source dùng nên dùng theo, lý do không biết.
  - api: để cấu hình một số router của api. Một số nơi sẽ dùng handler, transport,.. Tuỳ vào người dựng source.
  - dto: Data transfer object, dùng để nhận data từ client và trả data về cho client.
  - usecase: xử lý logic cũng như làm việc với repository. Một số source sẽ dùng là controller.
  - repository: là làm việc mới storage như database, cache.
  - adapter: là dùng call đến service khác. Xử lý logic của từng partner, tránh ảnh hưởng đến logic chính của usecase.
  ```
  https://refactoring.guru/design-patterns/adapter
  ```
  - facade: là một layer để handle một logic phức tạp hơn rất nhiều usecase. 
  ```
  https://refactoring.guru/design-patterns/facade
  ```
- pkg: là một folder để viết code, và có thể sử dụng lại.
- protoc: để định nghĩa message trong grpc hoặc rpc
```rpc
message MessageData {
  string status_code = 1;
  string reason_code = 2;
  string reason_message = 3;
  string jwt = 4;
}
```