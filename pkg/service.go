package clarifycontrol

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	kitlog "github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	consul "github.com/hashicorp/consul/api"
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

type ControlService interface {
	NodeStatus(ctx context.Context, jobName *string) ([]*NodeDetails, error)
	DiscoverService(ctx context.Context, serviceName, serviceTag *string, healthy bool) ([]*ServiceLocation, error)
	DrainNode(ctx context.Context, hostname *string, enabled bool) (bool, error)
	StopJob(ctx context.Context, jobName *string) (bool, error)
	Leader(ctx context.Context, serviceName *string) (*ServiceLocation, error)
}

type controlService struct {
	Nomad        lb.Balancer
	ConsulClient *consul.Client
	Discovery    consulsd.Client
	Logger       kitlog.Logger
}

type ServiceLocation struct {
	ServiceName *string
	ServiceHost *string
	ServicePort int
}

type NodeDetails struct {
	Host   client.Host
	Status int
}

func NewClarifyControlService(nomad lb.Balancer, consul *consul.Client, logger kitlog.Logger) ControlService {
	return &controlService{
		Nomad:        nomad,
		ConsulClient: consul,
		Discovery:    consulsd.NewClient(consul),
		Logger:       logger}
}

func (s *controlService) NodeStatus(ctx context.Context, jobName *string) ([]*NodeDetails, error) {
	nomad := s.getNomadServer()
	hosts, _, _ := client.Hosts(nomad)
	job, _ := client.FindJob(nomad, *jobName)
	hostStatus := make([]*NodeDetails, len(hosts))

	for i, host := range hosts {
		alloc, _ := client.FindAlloc(nomad, job, &host)
		hostStatus[i] = &NodeDetails{Host: host, Status: jobstatus(&host, job, alloc)}
	}
	return hostStatus, nil
}

func (s *controlService) DiscoverService(ctx context.Context, serviceName, serviceTag *string, healthy bool) ([]*ServiceLocation, error) {
	entries, _, err := s.Discovery.Service(*serviceName, *serviceTag, healthy, nil)
	locations := make([]*ServiceLocation, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error discovering service(name=%v;tag=%v;healthy=%v): %v", serviceName, serviceTag, healthy, err.Error())
		s.Logger.Log("err", errMsg)
		return locations, errors.New(errMsg)
	}
	for _, entry := range entries {
		sl := &ServiceLocation{ServiceHost: &entry.Service.Address, ServicePort: entry.Service.Port}
		locations = append(locations, sl)
	}
	return locations, nil
}

func (s *controlService) DrainNode(ctx context.Context, hostname *string, enabled bool) (bool, error) {
	nomad := s.getNomadServer()
	host, err := client.HostID(nomad, hostname)
	if err != nil {
		return false, err
	}
	status, err := client.Drain(nomad, host.ID, enabled)
	if err != nil {
		return false, err
	}
	if status == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func (s *controlService) StopJob(ctx context.Context, jobName *string) (bool, error) {
	nomad := s.getNomadServer()
	job := &client.Job{Name: *jobName}
	status, err := client.StopJob(nomad, job)
	if err != nil {
		return false, err
	}
	if status == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func (s *controlService) Leader(ctx context.Context, serviceName *string) (*ServiceLocation, error) {
	key := fmt.Sprintf("/service/%s/leader", *serviceName)
	kv, _, err := s.ConsulClient.KV().Get(key, nil)
	if err != nil {
		return &ServiceLocation{}, errors.New("kv_error")
	}
	if kv != nil {
		if address := string(kv.Value); len(address) > 0 {
			host, portstr, err := net.SplitHostPort(address)
			if err != nil {
				return &ServiceLocation{}, errors.New("invalid_address")
			}
			port, _ := strconv.Atoi(portstr)
			return &ServiceLocation{
				ServiceName: serviceName,
				ServiceHost: &host,
				ServicePort: port}, nil
		}
	}
	return &ServiceLocation{}, errors.New("no_leader")
}

func (s *controlService) getNomadServer() *client.NomadServer {
	endpoint, err := s.Nomad.Endpoint()
	if err != nil {
		s.Logger.Log("err", "Unable to locate nomad server.")
	}
	ep, _ := endpoint(nil, nil)
	host, portString, _ := net.SplitHostPort(ep.(string))
	port, _ := strconv.Atoi(portString)
	s.Logger.Log("nomad", fmt.Sprintf("%v:%v", host, port))
	return &client.NomadServer{Address: host, Port: port}
}

func jobstatus(host *client.Host, job *client.Job, alloc *client.Alloc) int {
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
