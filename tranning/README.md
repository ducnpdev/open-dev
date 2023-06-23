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
    "requestId": {{requestId}},
    "requestTime": {{timeRPC3339}},
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
- output: sẽ trả về 1 field: signatura sử dụng thuật toán sha256 hoặc hmacsha256
```
{
    "requestId": {{requestId}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "signatura": {{string}}
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
    "requestId": {{requestId}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "outEncode": {{string}},
        "outDecode": {{string}}
    }
}
```