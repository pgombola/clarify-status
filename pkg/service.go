package clarifystatussvc

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/sd/lb"
	"github.com/pgombola/gomad/client"
)

// StatusStopped means that the clarify service was not started in the scheduler.
const StatusStopped = 0

// StatusPending means that a host has been drained.
const StatusPending = 1

// StatusStarted means that specified job has been allocated and is running on the given host.
const StatusStarted = 2

// StatusMixed means that one or more tasks on the allocation are not in the "running" state.
const StatusMixed = 3

// StatusUnallocated means that an allocation hasn't been placed on a host
// or an allocation has been removed through adjusting the count
const StatusUnallocated = 4

// HostStatus is a read model for the status of a Host.
type HostStatus struct {
	Host   client.Host
	Status int
}

// Service is an interface that provides the GetAllHostsStatus method.
type Service interface {
	// GetAllHostStatus returns a pointer to a collection of HostStatus.
	GetHostStatus(ctx context.Context, jobName *string) ([]*HostStatus, error)
}

type statusService struct {
	Nomad  lb.Balancer
	Logger log.Logger
}

// NewClarifyStatusService returns a new instance of the service.
func NewClarifyStatusService(nomad lb.Balancer, logger log.Logger) Service {
	return &statusService{Nomad: nomad, Logger: logger}
}

func (s *statusService) GetHostStatus(_ context.Context, jobName *string) ([]*HostStatus, error) {
	nomad := s.getNomadServer()
	hosts, _, _ := client.Hosts(nomad)
	job, _ := client.FindJob(nomad, *jobName)
	hostStatus := make([]*HostStatus, len(hosts))

	for i, host := range hosts {
		alloc, _ := client.FindAlloc(nomad, job, &host)
		hostStatus[i] = &HostStatus{Host: host, Status: status(&host, job, alloc)}
	}
	return hostStatus, nil
}

func (s *statusService) getNomadServer() *client.NomadServer {
	endpoint, err := s.Nomad.Endpoint()
	if err != nil {
		s.Logger.Log("err", "Unable to locate nomad server.")
	}
	ep, _ := endpoint(nil, nil)
	host, portString, _ := net.SplitHostPort(ep.(string))
	port, _ := strconv.Atoi(portString)
	s.Logger.Log("info", fmt.Sprintf("Discovered nomad server @ %v:%v", host, port))
	return &client.NomadServer{Address: host, Port: port}
}

func status(host *client.Host, job *client.Job, alloc *client.Alloc) int {
	var status int
	if job.Name == "" {
		status = StatusStopped
	} else if (alloc.ClientStatus == "lost" || alloc.ClientStatus == "") && host.Drain {
		status = StatusPending
	} else if alloc.ClientStatus == "running" && !alloc.CheckTaskStates("running") {
		status = StatusMixed
	} else if alloc.ClientStatus == "complete" || (alloc.ClientStatus == "" && !host.Drain) {
		status = StatusUnallocated
	} else {
		status = StatusStarted
	}
	return status
}
