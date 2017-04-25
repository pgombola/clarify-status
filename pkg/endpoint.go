package clarifycontrol

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that are available for the staus service.
type Endpoints struct {
	NodeStatusEndpoint      endpoint.Endpoint
	DiscoverServiceEndpoint endpoint.Endpoint
	DrainEndpoint           endpoint.Endpoint
	StopEndpoint            endpoint.Endpoint
	LeaderEndpoint          endpoint.Endpoint
}

func MakeServerEndpoints(control ControlService) Endpoints {
	return Endpoints{
		NodeStatusEndpoint:      MakeNodeStatusEndpoint(control),
		DiscoverServiceEndpoint: MakeDiscoverServiceEndpoint(control),
		DrainEndpoint:           MakeDrainEndpoint(control),
		StopEndpoint:            MakeStopEndpoint(control),
		LeaderEndpoint:          MakeLeaderEndpoint(control),
	}
}

type statusRequest struct {
	JobName *string
}

type statusResponse struct {
	Status []*NodeDetails
	Err    error
}

func (r statusResponse) error() error { return r.Err }

func MakeNodeStatusEndpoint(control ControlService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		statusReq := request.(*statusRequest)
		s, err := control.NodeStatus(ctx, statusReq.JobName)
		return statusResponse{Status: s, Err: err}, nil
	}
}

type discoveryRequest struct {
	ServiceName *string
	ServiceTag  *string
	Healthy     bool
}

type discoveryResponse struct {
	Locations []*ServiceLocation
}

func MakeDiscoverServiceEndpoint(control ControlService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*discoveryRequest)
		s, err := control.DiscoverService(ctx, req.ServiceName, req.ServiceTag, req.Healthy)
		if err != nil {
			return discoveryResponse{}, err
		}
		return discoveryResponse{Locations: s}, nil
	}
}

type drainRequest struct {
	Hostname *string
}

type drainResponse struct {
	Hostname *string
	Drained  bool
}

func MakeDrainEndpoint(control ControlService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*drainRequest)
		drained, err := control.DrainNode(ctx, req.Hostname)
		if err != nil {
			return drainResponse{}, err
		}
		return drainResponse{Drained: drained, Hostname: req.Hostname}, nil
	}
}

type stopRequest struct {
	JobName *string
}

type stopResponse struct {
	JobName *string
	Stopped bool
}

func MakeStopEndpoint(control ControlService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*stopRequest)
		stopped, err := control.StopJob(ctx, req.JobName)
		if err != nil {
			return stopResponse{}, err
		}
		return stopResponse{Stopped: stopped, JobName: req.JobName}, nil
	}
}

func MakeLeaderEndpoint(control ControlService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*discoveryRequest)
		loc, err := control.Leader(ctx, req.ServiceName)
		if err != nil {
			return discoveryResponse{}, err
		}
		locs := make([]*ServiceLocation, 1)
		locs[0] = loc
		return discoveryResponse{Locations: locs}, nil
	}
}
