# open-dev

## Contents

- [open-dev](#open-dev)
  - [Contents](#contents)
  - [Gin Web Framework](#gin-web-framework)
    - [API Examples](#api-examples)
  - [Redis](#redis)
    - [Rate Limit](#rate-limit)
  - [Golang Usecase](#golang-usecase)
    - [Resize Image](#resize-image)
    - [Context](#context)
    - [Thuật Toán RSA](#thuật-toán-rsa)
    - [Thuận Toán ECDSA](#thuận-toán-ecdsa)
  - [Serverless Framework](#serverless-framework)
  - [Aws](#aws)
    - [Lambda](#lambda)
  - [design pattern](#design-pattern)
    - [sso là gì](#sso-là-gì)
    - [Lưu passowrd trong database:](#lưu-passowrd-trong-database)
  - [Performances](#performances)
    - [Standard](#standard)
  - [Contact:](#contact)
## Gin Web Framework
### API Examples
- source code in demo api simples: 

## Redis
- đọc bài hướng dẫn để cài redis: https://github.com/ducnpdev/open-dev/tree/master/redis
### Rate Limit
- đọc thêm bài viết để hiểu rate-limit là gì: https://viblo.asia/p/golang-ratelimit-la-gi-su-dung-ByEZkn3qKQ0

## Golang Usecase

### Resize Image
- edit size của image từ base64: https://github.com/ducnpdev/open-dev/tree/master/usecase#resize-image

### Context
- Handle context timeout: https://github.com/ducnpdev/open-dev/blob/master/usecase/context/timeout.go

### Thuật Toán RSA
- Sử dụng thuật toán rsa để mã hoá thông tin nhạy cảm.
- code ví dụ: https://github.com/ducnpdev/open-dev/blob/master/usecase/rsa/REAME.md

### Thuận Toán ECDSA
- Sử dụng thuật toán ecdsa để làm chữ kí số, được ứng dụng rộng rãi trong blockchain.
- Code ví dụ: https://github.com/ducnpdev/open-dev/blob/master/usecase/ecdsa/ecdsa.go

## Serverless Framework
- cài đặt serverless: https://viblo.asia/p/golang-cai-dat-serverless-framework-lambada-aws-4dbZNQXkKYM
## Aws
- Cách tạo credential aws, đọc blog: https://viblo.asia/p/cau-hinh-aws-credential-zOQJwYPxVMP
### Lambda
- Create source đơn giản: https://github.com/ducnpdev/open-dev/tree/master/aws/lambda
- Series lambda function đơn giản: https://viblo.asia/s/golang-lambda-serverless-vElaB8eD5kw
- Crud with postgres: https://github.com/ducnpdev/open-dev/tree/master/aws/lambda/crud

## design pattern

### sso là gì
*   **Khái Niệm:**
    * Hiểu một cách đơn giản thì SSO (Single Sign-On) là cơ chế xác thực, nó cho phép user đăng nhập trên nhiều hệ thống khác nhau với một ID.

*   **Nó Hoạt động như thế nào:**
    * `step1`: khi user vào hệ thống Gmail, gmail sẽ kiểm tra xem có login trước đó hay không, nếu không sẽ chuyển đến trang SSO-Authen, để user nhập thông tin login.
    * `step2-3`: Server SSO authentication sẽ kiểm tra thông tin, nếu hợp lệ sẽ tạo một global session và tạo token.
    * `step4-7`: Hệ thống Gmail sẽ kiểm tra token từ SSO trả về, và gửi lại cho user.
    * `step8`: Từ Gmail, user chuyển một trang khác của hệ thống google, ví dụ như youtube.
    * `step9-10`: Youtube sẽ kiểm tra user chưa login, sẽ chuyển token đến server sso để xác thực có hợp lệ hay không, trả về token.
    * `step11-14`: Hệ thống youtube sẽ kiểm tra token từ sso, và token sẽ được đăng kí trong hệ thống youtube, cuối cùng là gửi lại token cho user đã được bảo vệ. 

### Lưu passowrd trong database:
*   **Không Nên**:
    *   Password lưu plain-text là không tốt vì với những người nắm hệ thống sẽ có thể nhìn thấy
    *   Lưu password hash là chưa đủ, vì có thể bị tấn công, ví dụ: rainbow-tables.
    *   Để giảm thiệu các rủi ro, cần thêm `salt` đến password.
*   **Vậy Salt là gì?**
    *   Theo như hướng dẫn của OWASP "salt is a unique, randomly genereted string that is added to each password as part of the hashing process"
*   **Lưu Password và Salt.**
    *   Salt ở đây không phải là `secret` nên có thể lưu plaintext trong database. Salt được sử dụng để đảm bảo rằng password hash là duy nhất trong hệ thống.
    *   Password mình sẽ lưu kiểu: hash(pass+salt)
*   **Validate Password:**
    *   User nhập password
    *   Hệ thống sẽ dựa vào user để fetch `salt` được lưu dưới database.
    *   Hệ thống sẽ hash(pass+salt) (1), pass là user nhập
    *   So sách mã hash(1) có khớp với hash được lưu dưới database không, nếu giống nhau là password hợp lệ.

## Performances
### Standard
- test 2 hàm trả error: https://opendev.hashnode.dev/golang-test-performance-function-standard-1
- so sánh thời gian xử lý hàm convert string sang int: https://opendev.hashnode.dev/golang-test-performance-function-standard-1
## Contact:
- facebook: https://www.facebook.com/phucducdev/
- gmail: ducnp09081998@gmail.com or phucducktpm@gmail.com
- linkedin: https://www.linkedin.com/in/phucducktpm/
- hashnode: https://hashnode.com/@OpenDev
- telegram: https://t.me/OpenDevGolang
