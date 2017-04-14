package clarifycontrol

import (
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pgombola/clarify-status/pkg/pb"
)

type grpcServer struct {
	status    grpctransport.Handler
	discovery grpctransport.Handler
}

// MakeGRPCServer makes the status service endpoints available as gRPC StatusServer.
func MakeGRPCServer(endpoints Endpoints, tracer stdopentracing.Tracer, logger log.Logger) pb.ClarifyControlServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		status: grpctransport.NewServer(
			endpoints.StatusEndpoint,
			DecodeGRPCGetAllHostStatusRequest,
			EncodeGRPCGetAllHostStatusResponse,
			options...),
		discovery: grpctransport.NewServer(
			endpoints.ServiceDiscoveryEndpoint,
			DecodeGRPCDiscoverServiceRequest,
			EncodeGRPCDiscoverServiceResponse,
			options...),
	}
}
