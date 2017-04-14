package clarifycontrol

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
)

type ServiceLocation struct {
	ServiceHost string
	ServicePort int64
}

type DiscoveryService interface {
	DiscoverService(ctx context.Context, serviceName, serviceTag *string, healthy bool) ([]*ServiceLocation, error)
}

type discoveryService struct {
	Logger log.Logger
	Consul consulsd.Client
}

func NewDiscoveryService(consul consulsd.Client, logger log.Logger) DiscoveryService {
	return &discoveryService{Consul: consul, Logger: logger}
}

func (s *discoveryService) DiscoverService(_ context.Context, serviceName, serviceTag *string, healthy bool) ([]*ServiceLocation, error) {
	entries, _, err := s.Consul.Service(*serviceName, *serviceTag, healthy, nil)
	locations := make([]*ServiceLocation, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error discovering service(name=%v;tag=%v;healthy=%v): %v", serviceName, serviceTag, healthy, err.Error())
		return locations, errors.New(errMsg)
	}
	for _, entry := range entries {
		sl := &ServiceLocation{ServiceHost: entry.Service.Address, ServicePort: int64(entry.Service.Port)}
		locations = append(locations, sl)
	}
	return locations, nil
}
