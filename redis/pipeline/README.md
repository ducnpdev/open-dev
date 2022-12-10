# pipeline

## redis ZRANGE:
### usecase:
- requirement: handle logic attempt access for period
### implement
- document: https://redis.io/commands/zadd/
- code example: https://github.com/ducnpdev/open-dev/blob/master/redis/pipeline/range.go

### benchmark
- execute only one
```
go test -bench=. 
```
- execute with count
```
go test -bench=. -count 5
```
- output:
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