package clarifycontrol

import (
	"context"

	"github.com/pgombola/clarify-status/pkg/pb"
	oldcontext "golang.org/x/net/context"
)

func DecodeGRPCGetAllHostStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	statusReq := grpcReq.(*pb.HostStatusRequest)
	return &statusRequest{JobName: &statusReq.JobName}, nil
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
	case StatusUnallocated:
		return pb.HostStatusReply_Host_UNALLOCATED
	}
	return -1
}
