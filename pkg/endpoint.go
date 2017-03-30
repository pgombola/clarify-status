package clarifystatussvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that are available for the staus service.
type Endpoints struct {
	StatusEndpoint endpoint.Endpoint
}

// GetHostStatus implements Service.
func (e Endpoints) GetHostStatus(ctx context.Context) ([]*HostStatus, error) {
	request := statusRequest{}
	response, err := e.StatusEndpoint(ctx, request)
	if err != nil {
		empty := make([]*HostStatus, 0)
		return empty, err
	}
	return response.(statusResponse).Status, response.(statusResponse).Err
}

type statusRequest struct {
	JobName *string
}

type statusResponse struct {
	Status []*HostStatus
	Err    error
}

func (r statusResponse) error() error { return r.Err }

func MakeGetHostStatusEndpoint(status Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		statusReq := request.(statusRequest)
		s, err := status.GetHostStatus(ctx, statusReq.JobName)
		return statusResponse{Status: s, Err: err}, nil
	}
}
