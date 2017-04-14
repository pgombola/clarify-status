package clarifycontrol

import (
	"context"

	"github.com/pgombola/clarify-status/pkg/pb"
	oldcontext "golang.org/x/net/context"
)

func DecodeGRPCDiscoverServiceRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ServiceLocationRequest)
	return &discoveryRequest{ServiceName: &req.ServiceName, ServiceTag: &req.ServiceTag, Healthy: req.Healthy}, nil
}

func EncodeGRPCDiscoverServiceResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(discoveryResponse)
	replyLocations := make([]*pb.ServiceLocationReply_ServiceLocation, 0)
	for _, l := range resp.Locations {
		replyLocation := &pb.ServiceLocationReply_ServiceLocation{ServiceHost: l.ServiceHost, ServicePort: l.ServicePort}
		replyLocations = append(replyLocations, replyLocation)
	}
	return &pb.ServiceLocationReply{Locations: replyLocations}, nil
}

func (s *grpcServer) GetServiceLocation(ctx oldcontext.Context, req *pb.ServiceLocationRequest) (*pb.ServiceLocationReply, error) {
	_, rep, err := s.discovery.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceLocationReply), nil
}
