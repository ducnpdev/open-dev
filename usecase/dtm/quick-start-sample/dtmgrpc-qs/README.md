English | [简体中文](./README-cn.md)

# dtmgrpc-qs
Minimal usage example for client/dtmgrpc

## Quick start

### Installing and running dtm

Refer to [dtm installation and running](https://en.dtm.pub/guide/install.html)

### Startup example

``` bash
go run main.go
```

### Output

The order of execution can be seen in the log of dtmgrpc-qs as follows.

- TransOut
- TransIn

The entire saga transaction was executed successfully

### Example interpretation

``` GO
	gid := shortuuid.New()
	req := &busi.BusiReq{Amount: 30} // load of the microservice

	saga := dtmgrpc.NewSagaGrpc(busi.DtmGrpcServer, gid).
    // Add a subtransaction of TransOut, the forward action is grpc url: busi.BusiGrpc+"/busi.Busi/TransOut"， and the compensating action is similar
		Add(busi.BusiGrpc+"/busi.Busi/TransOut", busi.BusiGrpc+"/busi.Busi/TransOutRevert", req).
    // Add a subtransaction of TransIn, with grpc url: busi.BusiGrpc+"/busi.Busi/TransIn"， and the compensating action is similar
		Add(busi.BusiGrpc+"/busi.Busi/TransIn", busi.BusiGrpc+"/busi.Busi/TransInRevert", req)
  // commit saga transaction, dtm will complete all subtransactions/rollback all subtransactions
	err := saga.Submit()
```

### For more examples, see [dtm-examples](https://github.com/dtm-labs/dtm-examples)
