# pipeline

## redis ZRANGE:
### usecase:
- thường được sử dụng trong việc check số lần truy cập trong 1 khoảng thời gian.
- Yêu Cầu:
    -   ở Bank A có một yêu cầu. Trên application khi thực hiện xem số cvv của thẻ tín dụng.
    -   Trong vòng 5phút user chỉ được request hiển thị không quá 3 lần. 
- Giải quyết:
    -   Với cách thông thường thì mỗi lần request sẽ insert vào database, sau đó thực hiện cộng hoặc trừ để biết trong 1 giai đoạn có bao nhiêu lần request.
    -   Một cách khác đó là sử dụng redis pipeline. Tại sao dùng redis để thay thế đơn giản nó nhanh.
### implement
- document: https://redis.io/commands/zadd/
- code mẫu: 

### benchmark
- chỉ thực hiện 1 lân.
```
go test -bench=. 
```
- thực hiện với số lần bằng 5.
```
go test -bench=. -count 5
```
- kết quả:
```
goos: darwin
goarch: arm64
pkg: open-dev/redis/pipeline
BenchmarkPrimeNumbers-8             1303            848796 ns/op
BenchmarkPrimeNumbers-8             1336            850908 ns/op
BenchmarkPrimeNumbers-8             1425            812590 ns/op
PASS
ok      open-dev/redis/pipeline 5.469s
```