# dtmgrpc-qs
client/dtmgrpc的最简go使用示例

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
	gid := shortuuid.New() // 生成gid
	req := &busi.BusiReq{Amount: 30} // 微服务的载荷

	saga := dtmgrpc.NewSagaGrpc(busi.DtmGrpcServer, gid).
    // 添加一个TransOut的子事务，正向操作为grpc的url: busi.BusiGrpc+"/busi.Busi/TransOut"， 补偿操作类似
		Add(busi.BusiGrpc+"/busi.Busi/TransOut", busi.BusiGrpc+"/busi.Busi/TransOutRevert", req).
    // 添加一个TransIn的子事务，正向操作为grpc的url: busi.BusiGrpc+"/busi.Busi/TransIn"， 补偿操作类似
		Add(busi.BusiGrpc+"/busi.Busi/TransIn", busi.BusiGrpc+"/busi.Busi/TransInRevert", req)
  // 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	err := saga.Submit()
```

### 更多示例，详见[dtm-examples](https://github.com/dtm-labs/dtm-examples)
