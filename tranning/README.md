# golang courses
## Buoi 1
- Deploy lambda với golang, expose ra 2 restapi
- tất cả các api điều có mặc định 2 field: requestId, requestTime
### api-1
- input đầu vào có 2 field: value1, value2
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "value1": {{number}},
        "value2": {{number}}
    }
}
```
- output: sẽ trả về giá trị của value1+value2
```
{
   "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "data": {
        "sum": {{value1+value2}}
    }
}
```
### api-2
- input đầu vào có 2 field: plaintText, secretKey
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "plaintText": {{string}},
        "secretKey": {{string}}
    }
}
```
- output: sẽ trả về 1 field: signature sử dụng thuật toán sha256 hoặc hmacsha256
```
{
   "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "data": {
        "signature": {{string}}
    }
}
```
### api-3
- là dùng base64, input có 2 filed: needEncode, needDecode
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "needEncode": {{string}},
        "needDecode": {{string}}
    }
}
```
- output: sẽ trả về 2 field: outEncode là output của base64 field needEncode, outDecode là output của field needDecode
```
{
   "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "data": {
        "outEncode": {{string}},
        "outDecode": {{string}}
    }
}
```

## Buoi 2
- Viết 4 api, gọi vào database postgres, oracle or mysql
  - craete user
  - update user
  - delete user by username
  - get user detail by username
- Trong tất cả các api điều phải validate username, name, phone không được rỗng, KHÔNG được sài pointer
- Thiết kế api path cũng như method hợp lý.
- Script create table user
```sql
CREATE TABLE "users" (
    "id" bigserial,
    username character varying(50) COLLATE pg_catalog."default",
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(50) COLLATE pg_catalog."default",
    PRIMARY KEY ("id")
);
```
### api 1
- create user với username, name, phone. Phải check username tồn tại duy nhất trong table, không sử dụng unique của database
- input:
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "username": {{string}},
        "name": {{string}},
        "phone": {{string}}
    }
}
```
- output:
```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}}
}
```

### api 2
- update user by username. Thông tin update là name và phone.
- input:
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "username": {{string}},
        "name": {{string}},
        "phone": {{string}}
    }
}
```

- output:
```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}},
}
```

### api 3
- delete user by username
- input:
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "username": {{string}}
    }
}
```

- output:
```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}},
}
```

### api 4
- get user detail by username

- output:
```
{
    "responseId": {{uuid}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}},
    "data": {
        "username": {{string}},
        "name": {{string}},
        "phone": {{string}}
    }
}
```

## Buoi 3
- Thực hiện call api của service khác
- Dùng http client
- Trong buổi 2 có api tạo user, mình sẽ dùng api đó để trước khi insert user vào database. Cần call qua service khác bằng restapi để check số phone có hợp lệ hay không. Khi api trả về `00` thì phone hợp lệ, `!00` là phone không hợp lệ. Thông tin api(postman) sẽ gửi riêng.
- Trong api create user phải verify signature, sử dụng sha256 đã học buổi 1.
- input:
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "signature": {{sha256(requestId+phone+username+secretKey)}}
    // vidu: 
    // requestId=162b757c-2ab6-4ee4-9201-6b670afca615 
    // phone=0335287777
    // username=ducnp5
    // secretKey=golang
    // ==> signature=sha256(162b757c-2ab6-4ee4-9201-6b670afca6150335287777ducnp5golang)
    "data": {
        "username": {{string}},
        "name": {{string}},
        "phone": {{string}}
    }
}
```
- output:
```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}}
}
```