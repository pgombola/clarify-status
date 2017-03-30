package clarifystatussvc

import (
	"context"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pgombola/clarify-status/pkg/pb"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	status grpctransport.Handler
}

// MakeGRPCServer makes the status service endpoints available as gRPC StatusServer.
func MakeGRPCServer(endpoints Endpoints, tracer stdopentracing.Tracer, logger log.Logger) pb.ClarifyStatusServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		status: grpctransport.NewServer(
			endpoints.StatusEndpoint,
			DecodeGRPCGetAllHostStatusRequest,
			EncodeGRPCGetAllHostStatusResponse,
			options...),
	}
}

func DecodeGRPCGetAllHostStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return &statusRequest{}, nil
}

func EncodeGRPCGetAllHostStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(statusResponse)
	if resp.Err != nil {
		return &pb.HostStatusReply{Error: resp.Err.Error()}, nil
	}
	replyHosts := make([]*pb.HostStatusReply_Host, 0)
	for _, s := range resp.Status {
		host := &pb.HostStatusReply_Host{Hostname: s.Host.Name, Status: statusToPb(s.Status)}
		replyHosts = append(replyHosts, host)
	}
	return &pb.HostStatusReply{Hosts: replyHosts, Error: ""}, nil
}

func (s *grpcServer) GetHostStatus(ctx oldcontext.Context, req *pb.HostStatusRequest) (*pb.HostStatusReply, error) {
	_, rep, err := s.status.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.HostStatusReply), nil
}

func statusToPb(status int) pb.HostStatusReply_Host_HostStatus {
	switch status {
	case StatusMixed:
		return pb.HostStatusReply_Host_MIXED
	case StatusPending:
		return pb.HostStatusReply_Host_PENDING
	case StatusStarted:
		return pb.HostStatusReply_Host_STARTED
	case StatusStopped:
		return pb.HostStatusReply_Host_STOPPED
	}
	return -1
}
