# [Go] Lambda CRUD với Postgres
## Mục tiêu
* add new user
* get user by id
* update user
* delete user

## Khởi tạo project
* Ở series này mình dùng serverless để làm việc với lambda(env, create, deploy,..)
* Để biết chi tiết hơn, checkout [link](https://viblo.asia/p/golang-tao-project-lambda-bang-serverless-EoW4ob9xVml) đọc thêm.

## Script create table:
* sql create table user
  ```sql
  CREATE TABLE "users" (
    "id" bigserial,
    username character varying(50) COLLATE pg_catalog."default",
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(50) COLLATE pg_catalog."default",
    PRIMARY KEY ("id")
  );
  ```
* model database
  ```go
  type UserModel struct {
    ID       uint   `gorm:"primarykey"`
    UserName string `gorm:"column:username" bson:"username"`
    Email    string `gorm:"column:email" bson:"email"`
    Phone    string `gorm:"column:phone" bson:"phone"`
  }
  ```

## List function
### Thêm 1 user
* tạo endpoint trong file serverless:
```yml
functions:
  create:
    handler: bin/create # file binary sau khi build
    timeout: 3 # thời gian tối đa của functoin có thể xử lý 
    memorySize: 512 # resource memory được cấp phát cho function,
    description: create new user
    events:
      - http:
          path: /create # context-path
          method: post # method
```
* Code implement:
  * function connection đến postgres database:
    ```go
    // load env from os, cast to struct
    func loadConfig() Postgres {
        // những value khi get từ os lên thì đã được define trong file serverless.yaml
        user := os.Getenv("DB_USER")
        dbpass := os.Getenv("DB_PASS")
        dbhost := os.Getenv("DB_HOST")
        dbservice := os.Getenv("DB_SERVICE")
        return Postgres{
          Username: user,
          Password: dbpass,
          Database: dbservice,
          Host:     dbhost,
          Port:     5432,
        }
    }
    // create database postgres instance
    func InitPostgres() (*gorm.DB, error) {
        log.Default().Println("connecting postgres database")
        config := loadConfig()
        dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d ", config.Host, config.Username, config.Password, config.Database, config.Port)
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            log.Default().Println("connect postgres err:", err)
            return db, err
        }
        log.Default().Println("connect postgres successfully")
        return db, err
    }
    ```
  * struct nhận request client:
    ```go
    type UserDTO struct {
      Email string `json:"email"`
      User  string `json:"userName"`
      Phone string `json:"phone"`
    }
    ```
  * main function của lambda
    ```go
    func main() {
        lambda.Start(CreateUser)
    }
    ```
  * function response data
    ```go
    func ParseResponse(respBody HttpResponse) string {
      respBody.Time = time.Now().Format("2006-01-02T15:04:05.000-07:00")
      if respBody.Err != nil {
        return responseErr(respBody)
      }
      return responseOk(respBody)
    }

    func responseOk(respBody HttpResponse) string {
      var buf bytes.Buffer
      mapRes := map[string]interface{}{
        "responseId":      respBody.Uuid,
        "responseMessage": "successfully",
        "responseTime":    respBody.Time,
      }
      if respBody.Data != nil {
        mapRes["data"] = respBody.Data
      }
      body, errMarshal := json.Marshal(mapRes)
      if errMarshal != nil {
        log.Default().Println("marshal response err", errMarshal)
      }
      json.HTMLEscape(&buf, body)
      return buf.String()
    }

    func responseErr(respBody HttpResponse) string {
      var buf bytes.Buffer
      mapRes := map[string]interface{}{
        "responseId":      respBody.Uuid,
        "responseMessage": respBody.Err.Error(),
        "responseTime":    respBody.Time,
      }

      body, errMarshal := json.Marshal(mapRes)
      if errMarshal != nil {
        log.Default().Println("marshal response err", errMarshal)
      }
      json.HTMLEscape(&buf, body)
      return buf.String()
    }
    ```
  * function createUser
    ```go
    func CreateUser(ctx context.Context,
      eventReq events.APIGatewayProxyRequest) (Response, error) {
      var (
        req  = RequestBodyAPIGW{}
        resp = Response{
          StatusCode:      400,
          IsBase64Encoded: false,
          Headers: map[string]string{
            "Content-Type": "application/json",
          },
        }
      )
      err := json.Unmarshal([]byte(eventReq.Body), &req)
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{
          Uuid: req.RequestID,
          Err:  err,
        })
        return resp, nil
      }
      db, err := pkg.InitPostgres()
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{
          Uuid: req.RequestID,
          Err:  err,
        })
        resp.StatusCode = 500
        return resp, nil
      }
      // set http-code 200
      resp.StatusCode = 200
      // save new user
      err = db.Debug().Exec(`insert into users(username,email,phone) values(?,?,?)`, req.Data.User, req.Data.Email, req.Data.Phone).Error
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Err: err})
        return resp, nil
      }
      resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Data: nil})
      return resp, nil
    }
    ```
### Get user by id
* tạo endpoint trong file serverless:
```yml
functions:
  read:
    handler: bin/read
    description: get user detail by id
    events:
      - http:
          path: /read
          method: get
```
* Code implement:
  * Struct nhận request từ client
    ```go
    type UserDTO struct {
      ID string `json:"userId"`
    }
    ```
  * function get user detail by id
    ```go
    func ReadUser(ctx context.Context,
      eventReq events.APIGatewayProxyRequest) (Response, error) {
      var (
        req  = RequestBodyAPIGW{}
        resp = Response{
          StatusCode:      400,
          IsBase64Encoded: false,
          Headers: map[string]string{
            "Content-Type": "application/json",
          },
        }
      )
      // get parameter from url
      req.Data.ID = eventReq.QueryStringParameters["id"]
      // init database,
      db, err := pkg.InitPostgres()

      if err != nil {
        resp.Body = ParseResponse(HttpResponse{
          Uuid: req.RequestID,
          Err:  err,
        })
        resp.StatusCode = 500
        return resp, nil
      }
      resp.StatusCode = 200
      var user = models.UserModel{}
      err = db.Table("users").Where("id = ?", req.Data.ID).First(&user).Error
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Err: err})
        return resp, nil
      }
      resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Data: user})
      return resp, nil
    }
    ```

### Get update user
* tạo endpoint trong file serverless
  ```bash
  functions:
    update:
      handler: bin/update
      description: update user by id
      events:
        - http:
            path: /update
            method: post
  ```
* Struct nhận request từ client
  ```go
    type UserDTO struct {
      ID    int    `json:"userId"`
      Email string `json:"email"`
      Phone string `json:"phone"`
    }

  ```
* Code implement
  * function get update user
    ```go
    func UpdateUser(ctx context.Context, eventReq events. APIGatewayProxyRequest) (Response, error) {
      var (
        req  = RequestBodyAPIGW{}
        resp = Response{
          StatusCode:      400,
          IsBase64Encoded: false,
          Headers: map[string]string{
            "Content-Type": "application/json",
          },
        }
      )
      err := json.Unmarshal([]byte(eventReq.Body), &req)
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{
          Uuid: req.RequestID,
          Err:  err,
        })
        return resp, nil
      }
      db, err := pkg.InitPostgres()

      if err != nil {
        resp.Body = ParseResponse(HttpResponse{
          Uuid: req.RequestID,
          Err:  err,
        })
        resp.StatusCode = 500
        return resp, nil
      }
      resp.StatusCode = 200
      err = db.Exec(`update users set email = ?, phone = ? where id = ?`, req.Data.Email, req.Data.Phone, req.Data.ID).Error
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Err: err})
        return resp, nil
      }
      resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Data: nil})
      return resp, nil
    }
    ```

### Detele user
* tạo endpoint trong file serverless
  ```bash
  functions:
    delete:
      handler: bin/delete
      description: delete user by id
      events:
        - http:
            path: /delete
            method: post
  ```
* Struct nhận request từ client
  ```go
    type UserDTO struct {
      ID    int    `json:"userId"`
    }

  ```
* Code implement
  * function get delete user by id
    ```go
    func DeleteUser(ctx context.Context, eventReq events.APIGatewayProxyRequest) (Response, error) {
      var (
        req  = RequestBodyAPIGW{}
        resp = Response{
          StatusCode:      400,
          IsBase64Encoded: false,
          Headers: map[string]string{
            "Content-Type": "application/json",
          },
        }
      )
      err := json.Unmarshal([]byte(eventReq.Body), &req)
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{
          Uuid: req.RequestID,
          Err:  err,
        })
        return resp, nil
      }
      db, err := pkg.InitPostgres()
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{
          Uuid: req.RequestID,
          Err:  err,
        })
        resp.StatusCode = 500
        return resp, nil
      }
      resp.StatusCode = 200
      var user = models.UserModel{}
      err = db.Debug().Table("users").Delete(&user, req.Data.ID).Error
      if err != nil {
        resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Err: err})
        return resp, nil
      }
      resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Data: req.Data})
      return resp, nil
    }
    ```