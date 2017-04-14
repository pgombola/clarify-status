package clarifycontrol

import "github.com/go-kit/kit/endpoint"

// Endpoints collects all of the endpoints that are available for the staus service.
type Endpoints struct {
	StatusEndpoint           endpoint.Endpoint
	ServiceDiscoveryEndpoint endpoint.Endpoint
}
