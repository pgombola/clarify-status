package clarifystatussvc

import (
	"context"
	"net"
	"strconv"

	"github.com/go-kit/kit/sd"
	"github.com/pgombola/gomad/client"
)

// StatusStopped means that the clarify service was not started in the scheduler.
const StatusStopped = 0
const StatusPending = 1
const StatusStarted = 2
const StatusMixed = 3

// HostStatus is a read model for the status of a Host.
type HostStatus struct {
	Host   client.Host
	Status int
}

// Service is an interface that provides the GetAllHostsStatus method.
type Service interface {
	// GetAllHostStatus returns a pointer to a collection of HostStatus.
	GetHostStatus(ctx context.Context) ([]*HostStatus, error)
}

type statusService struct {
	Nomad *client.NomadServer
}

// NewClarifyStatusService returns a new instance of the service.
func NewClarifyStatusService(nomads sd.Subscriber) Service {
	eps, _ := nomads.Endpoints()
	ep, _ := eps[0](nil, nil)
	host, portString, _ := net.SplitHostPort(ep.(string))
	port, _ := strconv.Atoi(portString)
	nomadServer := &client.NomadServer{Address: host, Port: port}
	return &statusService{Nomad: nomadServer}
}

func (s *statusService) GetHostStatus(_ context.Context) ([]*HostStatus, error) {
	hosts, _, _ := client.Hosts(s.Nomad)
	clarifyJob, _ := client.FindJob(s.Nomad, "clarify")
	hostStatus := make([]*HostStatus, len(hosts))

	for i, host := range hosts {
		alloc, err := client.FindAlloc(s.Nomad, clarifyJob, &host)
		if err != nil {
			// TODO clarify not allocated here
		}
		hostStatus[i] = &HostStatus{Host: host, Status: status(&host, clarifyJob, alloc)}
	}
	return hostStatus, nil
}

func status(host *client.Host, clarify *client.Job, alloc *client.Alloc) int {
	var status int
	if clarify.Name == "" {
		status = StatusStopped
	} else if (alloc.ClientStatus == "lost" || alloc.ClientStatus == "") && host.Drain {
		status = StatusPending
	} else if alloc.ClientStatus == "running" && !alloc.CheckTaskStates("running") {
		status = StatusMixed
	} else {
		status = StatusStarted
	}
	return status
}
