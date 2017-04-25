package clarifycontrol

import (
	"context"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pgombola/clarify-status/pkg/pb"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	status    grpctransport.Handler
	discovery grpctransport.Handler
	drain     grpctransport.Handler
	stop      grpctransport.Handler
	leader    grpctransport.Handler
}

// MakeGRPCServer makes the status service endpoints available as gRPC StatusServer.
func MakeGRPCServer(endpoints Endpoints, tracer stdopentracing.Tracer, logger log.Logger) pb.ClarifyControlServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		status: grpctransport.NewServer(
			endpoints.NodeStatusEndpoint,
			DecodeGRPCGetAllHostStatusRequest,
			EncodeGRPCGetAllHostStatusResponse,
			options...),
		discovery: grpctransport.NewServer(
			endpoints.DiscoverServiceEndpoint,
			DecodeGRPCDiscoverServiceRequest,
			EncodeGRPCDiscoverServiceResponse,
			options...),
		drain: grpctransport.NewServer(
			endpoints.DrainEndpoint,
			DecodeGRPCDrainRequest,
			EncodeGRPCDrainResponse,
			options...),
		stop: grpctransport.NewServer(
			endpoints.StopEndpoint,
			DecodeGRPCStopRequest,
			EncodeGRPCStopResponse,
			options...),
		leader: grpctransport.NewServer(
			endpoints.LeaderEndpoint,
			DecodeGRPCLeaderRequest,
			EncodeGRPCLeaderResponse,
			options...),
	}
}

func (s *grpcServer) NodeStatus(ctx oldcontext.Context, req *pb.Job) (*pb.NodeStatusReply, error) {
	_, rep, err := s.status.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.NodeStatusReply), nil
}

func DecodeGRPCGetAllHostStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	statusReq := grpcReq.(*pb.Job)
	return &statusRequest{JobName: &statusReq.JobName}, nil
}

func EncodeGRPCGetAllHostStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(statusResponse)
	if resp.Err != nil {
		return &pb.NodeStatusReply{Error: resp.Err.Error()}, nil
	}
	replyDetails := make([]*pb.NodeStatusReply_NodeDetails, 0)
	for _, s := range resp.Status {
		node := &pb.Node{Hostname: s.Host.Name}
		details := &pb.NodeStatusReply_NodeDetails{Node: node, Status: statusToPb(s.Status)}
		replyDetails = append(replyDetails, details)
	}
	return &pb.NodeStatusReply{Details: replyDetails, Error: ""}, nil
}

func statusToPb(status int) pb.ClarifyStatus {
	switch status {
	case StatusMixed:
		return pb.ClarifyStatus_NODE_ALLOC_MIXED
	case StatusPending:
		return pb.ClarifyStatus_NODE_DRAINED
	case StatusStarted:
		return pb.ClarifyStatus_NODE_ALLOC_STARTED
	case StatusStopped:
		return pb.ClarifyStatus_JOB_STOPPED
	case StatusUnallocated:
		return pb.ClarifyStatus_NODE_UNALLOCATED
	}
	return -1
}

func (s *grpcServer) ServiceLocation(ctx oldcontext.Context, req *pb.Service) (*pb.ServiceLocationReply, error) {
	_, rep, err := s.discovery.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceLocationReply), nil
}

func DecodeGRPCDiscoverServiceRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.Service)
	return &discoveryRequest{ServiceName: &req.ServiceName, ServiceTag: &req.ServiceTag, Healthy: req.Healthy}, nil
}

func EncodeGRPCDiscoverServiceResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(discoveryResponse)
	replyLocations := make([]*pb.ServiceLocationReply_ServiceLocation, 0)
	for _, l := range resp.Locations {
		replyLocation := &pb.ServiceLocationReply_ServiceLocation{ServiceHost: *l.ServiceHost, ServicePort: int64(l.ServicePort)}
		replyLocations = append(replyLocations, replyLocation)
	}
	return &pb.ServiceLocationReply{Locations: replyLocations}, nil
}

func (s *grpcServer) Drain(ctx oldcontext.Context, req *pb.DrainRequest) (*pb.DrainReply, error) {
	_, rep, err := s.drain.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DrainReply), nil
}

func DecodeGRPCDrainRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	drainreq := grpcReq.(*pb.DrainRequest)
	return &drainRequest{&drainreq.Node.Hostname, drainreq.Enabled}, nil
}

func EncodeGRPCDrainResponse(_ context.Context, resp interface{}) (interface{}, error) {
	dr := resp.(drainResponse)
	hostname := dr.Hostname
	node := &pb.Node{Hostname: *hostname}
	if dr.Drained && dr.Enabled {
		return &pb.DrainReply{Node: node, Status: pb.ClarifyStatus_NODE_DRAINED}, nil
	} else if dr.Drained && !dr.Enabled {
		return &pb.DrainReply{Node: node, Status: pb.ClarifyStatus_NODE_UNALLOCATED}, nil
	}
	return &pb.DrainReply{Node: node, Status: pb.ClarifyStatus_UNKNOWN}, nil
}

func (s *grpcServer) Stop(ctx oldcontext.Context, req *pb.Job) (*pb.StopReply, error) {
	_, rep, err := s.stop.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.StopReply), nil
}

func DecodeGRPCStopRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	job := grpcReq.(*pb.Job)
	return &stopRequest{JobName: &job.JobName}, nil
}

func EncodeGRPCStopResponse(_ context.Context, resp interface{}) (interface{}, error) {
	stopped := resp.(stopResponse).Stopped
	jobname := *resp.(stopResponse).JobName
	if stopped {
		return &pb.StopReply{
			Job:    &pb.Job{JobName: jobname},
			Status: pb.ClarifyStatus_JOB_STOPPED}, nil
	}
	return &pb.StopReply{
		Job:    &pb.Job{JobName: jobname},
		Status: pb.ClarifyStatus_UNKNOWN}, nil
}

func (s *grpcServer) Leader(ctx oldcontext.Context, req *pb.Service) (*pb.ServiceLocationReply, error) {
	_, rep, err := s.leader.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceLocationReply), err
}

func DecodeGRPCLeaderRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	servicename := &grpcReq.(*pb.Service).ServiceName
	return &discoveryRequest{ServiceName: servicename}, nil
}

func EncodeGRPCLeaderResponse(_ context.Context, resp interface{}) (interface{}, error) {
	discovered := resp.(discoveryResponse).Locations
	locations := make([]*pb.ServiceLocationReply_ServiceLocation, 0)
	if len(discovered) == 0 {
		return &pb.ServiceLocationReply{Locations: locations}, nil
	}
	locations = append(locations,
		&pb.ServiceLocationReply_ServiceLocation{
			ServiceHost: *discovered[0].ServiceHost,
			ServicePort: int64(discovered[0].ServicePort),
			ServiceName: "find"})
	return &pb.ServiceLocationReply{Locations: locations}, nil
}
