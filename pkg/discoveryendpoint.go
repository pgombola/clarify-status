package clarifycontrol

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type discoveryRequest struct {
	ServiceName *string
	ServiceTag  *string
	Healthy     bool
}

type discoveryResponse struct {
	Locations []*ServiceLocation
}

func MakeServiceDiscoveryEndpoint(discovery DiscoveryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		discoverReq := request.(*discoveryRequest)
		s, err := discovery.DiscoverService(
			ctx,
			discoverReq.ServiceName,
			discoverReq.ServiceTag,
			discoverReq.Healthy)
		if err != nil {
			return &discoveryResponse{}, err
		}
		return discoveryResponse{Locations: s}, nil
	}
}
