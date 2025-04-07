package main

import (
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/logger"
	"github.com/dtm-labs/client/workflow"
	"github.com/dtm-labs/quick-start-sample/busi"
	"github.com/lithammer/shortuuid/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var busiCli busi.BusiClient

func main() {
	s := busi.GrpcNewServer()
	workflow.InitGrpc(busi.DtmGrpcServer, busi.BusiGrpc, s)
	busi.GrpcStartup(s)

	conn1, err := grpc.Dial(busi.BusiGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(workflow.Interceptor))
	logger.FatalIfError(err)
	busiCli = busi.NewBusiClient(conn1)

	wfName := "workflow-grpc"
	err = workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var req busi.BusiReq
		err := proto.Unmarshal(data, &req)
		logger.FatalIfError(err)
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, err := busiCli.TransOutRevert(wf.Context, &req)
			return err
		})
		_, err = busiCli.TransOut(wf.Context, &req)
		if err != nil {
			return err
		}

		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, err := busiCli.TransInRevert(wf.Context, &req)
			return err
		})
		_, err = busiCli.TransIn(wf.Context, &req)
		return err
	})
	logger.FatalIfError(err)

	req := busi.BusiReq{Amount: 30, TransInResult: "FAILURE"}
	data, err := proto.Marshal(&req)
	logger.FatalIfError(err)
	err = workflow.Execute(wfName, shortuuid.New(), data)
	logger.Infof("result of workflow.Execute is: %v", err)
	time.Sleep(3 * time.Second)
}
