# dtmcli-qs
client/dtmcli的最简go使用示例

## 快速开始

### 安装运行dtm

参考[dtm安装运行](https://dtm.pub/guide/install.html)

### 启动示例

``` bash
go run main.go
```

### 输出

可以从dtmcli-qs的日志里看到执行的顺序如下：

- TransOut
- TransIn

整个saga事务执行成功

### 示例解读

``` GO
  // 具体业务微服务地址
  const qsBusi = "http://localhost:8081/api/busi_saga"
  req := &gin.H{"amount": 30} // 微服务的载荷
  // DtmServer为DTM服务的地址，是一个url
  DtmServer := "http://localhost:36789/api/dtmsvr"
  saga := dtmcli.NewSaga(DtmServer, dtmcli.MustGenGid(DtmServer)).
    // 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 补偿操作为url: qsBusi+"/TransOutCom"
    Add(qsBusi+"/TransOut", qsBusi+"/TransOutCom", req).
    // 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransIn"， 补偿操作为url: qsBusi+"/TransInCom"
    Add(qsBusi+"/TransIn", qsBusi+"/TransInCom", req)
  // 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
  err := saga.Submit()
```

### 更多示例，详见[dtm-examples](https://github.com/dtm-labs/dtm-examples)
