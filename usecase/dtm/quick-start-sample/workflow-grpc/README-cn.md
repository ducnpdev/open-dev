# workflow-grpc
client/workflow的最简grpc使用示例

## 快速开始

### 安装运行dtm

参考[dtm安装运行](https://dtm.pub/guide/install.html)

### 启动示例

``` bash
go run main.go
```

### 输出

可以从workflow-grpc的日志里看到执行的顺序如下：

- TransOut
- TransIn

整个workflow事务执行成功

## 步骤
完整的使用例子会包括以下步骤

### 初始化 Workflow
``` Go
	s := busi.GrpcNewServer()
	workflow.InitGrpc(busi.DtmGrpcServer, busi.BusiGrpc, s)
```

### 添加Workflow 的拦截器
Workflow 会自动拦截gRPC的请求，并记录进度
``` Go
	conn1, err := grpc.Dial(busi.BusiGrpc, nossl, grpc.WithUnaryInterceptor(workflow.Interceptor))
	// check err
	busiCli = busi.NewBusiClient(conn1)
```

### 注册一个workflow
``` Go
	wfName := "workflow-grpc"
	err = workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
	// ...
		_, err = busiCli.TransOut(wf.Context, &req)
	// grpc call should use clients with workflow's interceptor, and workflow's Context
	}
```

### 执行一个 Workflow
``` Go
	err = workflow.Execute(wfName, shortuuid.New(), data)
```

### 更多示例，详见[dtm-examples](https://github.com/dtm-labs/dtm-examples)
