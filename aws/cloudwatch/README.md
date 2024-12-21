# cloudwatch

# Log Group
- code golang example: https://github.com/ducnpdev/open-dev/tree/master/aws/cloudwatch/logGroup/README.md

## log insight
- thống kê lượt call theo api key id
```json
fields @timestamp, @message
| filter @message like /API Key ID/
| parse @message "API Key ID: *" as @apikeyid
| stats count(*) by @apikeyid
```
- thống kê lượt call theo ip
```json
fields @timestamp, @message
| filter @message like /X-Forwarded-For/
| parse @message "X-Forwarded-For=*," as @apikeyid
| stats count(*) as tmp by @apikeyid
| sort tmp desc
```

- số lượng call theo path
```json
fields @timestamp, @message
| filter @message like /Resource Path/
| parse @message "Resource Path: *" as @apipath
| stats count(*) as tmp by @apipath
| sort tmp desc
```

- số lượng call vào domain
```json
fields @timestamp, @message
| filter @message like /Host/
| parse @message "Host=*," as @apipath
| stats count(*) as tmp by @apipath
| sort tmp desc
```

- sort latency
```json
fields @timestamp, @message, @logStream, @log, latency
| sort latency desc
| limit 20
```
- max, min, avg latency
```json
fields @timestamp, @message, @logStream, @log, latency
| max(latency) as max, min(latency) as min, avg(latency)
```
- filter latency 
```json
fields @timestamp, @message, @logStream, @log, latency
| filter latency > 1000
```

- filter status, http code
```json
fields @timestamp, @message, @logStream, @duration 
| filter @message like 'status":"500"'
| filter @message like 'pathRouter'
| sort @timestamp desc
```

## lambda
- memory usage
```json
filter @type = "REPORT"
| stats max(@memorySize / 1000 / 1000) as provisionedMemoryMB,
  min(@maxMemoryUsed / 1000 / 1000) as smallestMemoryRequestMB,
  avg(@maxMemoryUsed / 1000 / 1000) as avgMemoryUsedMB,
  max(@maxMemoryUsed / 1000 / 1000) as maxMemoryUsedMB,
  provisionedMemoryMB - maxMemoryUsedMB as overProvisionedMB
```
- latency lambda
```json
filter @type = "REPORT"
| fields @requestId, @billedDuration, @logStream, @duration
| sort by @billedDuration desc
---
REPORT RequestId: 56ca8bba-aa7d-4585-b4d6-5aa53651349a Duration: 93.16 ms Billed Duration: 94 ms Memory Size: 512 MB Max Memory Used: 60 MB

- time max process json field
```json
fields duration
| filter @message like /time_process_action/
| sort duration desc
| limit 10000
```

- sort time duration call external service
```json
fields duration,@timestamp, @message, @logStream, @log
| filter @logStream like /kafka/
| filter @message like /time_process_action/
| sort duration desc
```