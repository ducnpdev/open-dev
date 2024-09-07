# http client retry

- library: `github.com/hashicorp/go-retryablehttp`

## test call retry api:
- config number retry:
```go
retryClient.RetryMax = 3
```

- config time retry:
```go
retryClient.RetryWaitMin = 1 * time.Second // Minimum wait time before retry
retryClient.RetryWaitMax = 3 * time.Second // Maximum wait time before retry
```

- output:
```console
2024/09/07 17:13:15 [DEBUG] GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello
2024/09/07 17:13:15 [ERR] GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello request failed: Get "https://***.execute-api.ap-southeast-1.amazonaws.com/hello": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

2024/09/07 17:13:15 [DEBUG] GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello: retrying in 1s (3 left)
2024/09/07 17:13:16 [ERR] GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello request failed: Get "https://***.execute-api.ap-southeast-1.amazonaws.com/hello": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

2024/09/07 17:13:16 [DEBUG] GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello: retrying in 2s (2 left)
2024/09/07 17:13:18 [ERR] GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello request failed: Get "https://***.execute-api.ap-southeast-1.amazonaws.com/hello": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

2024/09/07 17:13:18 [DEBUG] GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello: retrying in 3s (1 left)
2024/09/07 17:13:21 [ERR] GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello request failed: Get "https://***.execute-api.ap-southeast-1.amazonaws.com/hello": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

2024/09/07 17:13:21 Error creating request: GET https://***.execute-api.ap-southeast-1.amazonaws.com/hello giving up after 4 attempt(s): Get "https://***.execute-api.ap-southeast-1.amazonaws.com/hello": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
exit status 1
```