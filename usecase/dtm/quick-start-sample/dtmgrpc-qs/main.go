package main

import (
	"time"

	"github.com/dtm-labs/client/dtmcli/logger"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/dtm-labs/quick-start-sample/busi"
	"github.com/lithammer/shortuuid/v3"
)

func main() {
	s := busi.GrpcNewServer()
	busi.GrpcStartup(s)
	logger.Infof("grpc simple transaction begin")
	gid := shortuuid.New()
	req := &busi.BusiReq{Amount: 30}
	// req := &busi.BusiReq{Amount: 30, TransInResult: "FAILURE"}
	saga := dtmgrpc.NewSagaGrpc(busi.DtmGrpcServer, gid).
		Add(busi.BusiGrpc+"/busi.Busi/TransOut", busi.BusiGrpc+"/busi.Busi/TransOutRevert", req).
		Add(busi.BusiGrpc+"/busi.Busi/TransIn", busi.BusiGrpc+"/busi.Busi/TransInRevert", req)
	err := saga.Submit()
	if err != nil {
		panic(err)
	}
	time.Sleep(3 * time.Second)
}
