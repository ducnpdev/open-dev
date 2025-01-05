English | [简体中文](./README-cn.md)

# workflow-grpc
Minimal usage example for client/workflow using grpc protocol

## Quick start

### Installing and running dtm

Refer to [dtm installation and running](https://en.dtm.pub/guide/install.html)

### Startup example

``` bash
go run main.go
```

### Output

The order of execution can be seen in the log of workflow-grpc as follows.

- TransOut
- TransIn

The entire workflow transaction was executed successfully

## Steps
A complete example includes following steps:

### Init Workflow
``` Go
	s := busi.GrpcNewServer()
	workflow.InitGrpc(busi.DtmGrpcServer, busi.BusiGrpc, s)
```

### Add Workflow Interceptor
Workflow will automaticly intercept gRPC call, and record the progresses.
``` Go
	conn1, err := grpc.Dial(busi.BusiGrpc, nossl, grpc.WithUnaryInterceptor(workflow.Interceptor))
	// check err
	busiCli = busi.NewBusiClient(conn1)
```

### Register a Workflow
``` Go
	wfName := "workflow-grpc"
	err = workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
	// ...
		_, err = busiCli.TransOut(wf.Context, &req)
	// grpc call should use clients with workflow's interceptor, and workflow's Context
	}
```

### Execute a Workflow
``` Go
	err = workflow.Execute(wfName, shortuuid.New(), data)
```

### For more examples, see [dtm-examples](https://github.com/dtm-labs/dtm-examples)
