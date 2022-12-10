# pipeline

## redis ZRANGE:
### usecase:
### implement
- document: https://redis.io/commands/zadd/
- code example: 

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