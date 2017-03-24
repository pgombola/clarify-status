package clarifystatussvc

import (
	"context"

	"github.com/pgombola/gomad/client"
)

// Status_Stopped means that the clarify service was not started in the scheduler.
const Status_Stopped = 0

// HostStatus is a read model for the status of a Host.
type HostStatus struct {
	Host   client.Host
	Status int32
}

// Service is an interface that provides the GetAllHostsStatus method.
type Service interface {
	// GetAllHostStatus returns a pointer to a collection of HostStatus.
	GetHostStatus(ctx context.Context) (*[]HostStatus, error)
}

type statusService struct{}

// NewClarifyStatusService returns a new instance of the service.
func NewClarifyStatusService() Service {
	return &statusService{}
}

func (s *statusService) GetHostStatus(_ context.Context) (*[]HostStatus, error) {
	status := make([]HostStatus, 0)
	return &status, nil
}
