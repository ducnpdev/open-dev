package busi

import (
	context "context"
	"fmt"
	"net"
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/dtm-labs/client/dtmgrpc/dtmgimp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// DtmGrpcServer dtm grpc service address
const DtmGrpcServer = "localhost:36790"

// BusiGrpcPort 1
const BusiGrpcPort = 50589

// BusiGrpc busi service grpc address
var BusiGrpc string = fmt.Sprintf("localhost:%d", BusiGrpcPort)

func handleGrpcBusiness(in *BusiReq, result1 string, busi string) error {
	res := dtmimp.OrString(result1, dtmcli.ResultSuccess)
	dtmimp.Logf("grpc busi %s result: %s", busi, res)
	if res == dtmcli.ResultSuccess {
		return nil
	} else if res == dtmcli.ResultFailure {
		return status.New(codes.Aborted, "FAILURE").Err()
	}
	return status.New(codes.Internal, fmt.Sprintf("unknow result %s", res)).Err()
}

// busiServer is used to implement helloworld.GreeterServer.
type busiServer struct {
	UnimplementedBusiServer
}

// GrpcNewServer new a Server
func GrpcNewServer() *grpc.Server {
	return grpc.NewServer(grpc.UnaryInterceptor(dtmgimp.GrpcServerLog))
}

// GrpcStartup for grpc
func GrpcStartup(s *grpc.Server) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", BusiGrpcPort))
	dtmimp.FatalIfError(err)
	s.RegisterService(&Busi_ServiceDesc, &busiServer{})
	go func() {
		dtmimp.Logf("busi grpc listening at %v", lis.Addr())
		err := s.Serve(lis)
		dtmimp.FatalIfError(err)
	}()
	time.Sleep(100 * time.Millisecond)
}

func (s *busiServer) TransInRevert(ctx context.Context, in *BusiReq) (*BusiReply, error) {
	return &BusiReply{}, handleGrpcBusiness(in, "", dtmimp.GetFuncName())
}

func (s *busiServer) TransOutRevert(ctx context.Context, in *BusiReq) (*BusiReply, error) {
	return &BusiReply{}, handleGrpcBusiness(in, "", dtmimp.GetFuncName())
}

func (s *busiServer) TransIn(ctx context.Context, in *BusiReq) (*BusiReply, error) {
	return &BusiReply{}, handleGrpcBusiness(in, in.TransInResult, dtmimp.GetFuncName())
}

func (s *busiServer) TransOut(ctx context.Context, in *BusiReq) (*BusiReply, error) {
	return &BusiReply{}, handleGrpcBusiness(in, in.TransOutResult, dtmimp.GetFuncName())
}
