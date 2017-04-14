package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/hashicorp/consul/api"
	"github.com/pgombola/clarify-status/pkg"
	"github.com/pgombola/clarify-status/pkg/pb"
)

func main() {
	var (
		grpcAddr    = flag.String("grpc.addr", ":8081", "gRPC (HTTP/2) listen address")
		httpAddr    = flag.String("http.addr", ":8082", "HTTP listen address (/health)")
		consulAddr  = flag.String("consul.addr", "localhost:8500", "Address of consul agent")
		instance    = flag.Int("instance", 0, "The instance count of the status service")
		serviceName = flag.String("service.name", "clarify-control", "Name of the service")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	svc := &service{
		HTTPAddress: httpAddr,
		GRPCAddress: grpcAddr,
		Name:        serviceName,
		Instance:    *instance}

	// Registry domain.
	var (
		client consulsd.Client
		reg    *api.AgentServiceRegistration
	)
	{
		var err error
		client, err = createConsulClient(consulAddr, logger)
		if err != nil {
			logger.Log("err", err)
		}
		reg, err = registerService(client, svc)
		if err != nil {
			logger.Log("err", err)
		}
	}

	// Discovery domain.
	var (
		balancer lb.Balancer
	)
	{
		nomadSubscriber := consulsd.NewSubscriber(client, nomadAddressFactory, logger, "nomad", []string{"http"}, true)
		balancer = lb.NewRoundRobin(nomadSubscriber)
	}

	// Mechanical domain.
	// Error handling channel.
	errc := make(chan error)
	// Signal handling channel. This allows for registry cleanup when we receive SIGTERM, SIGINT.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Business domain.
	var statusService clarifycontrol.StatusService
	var discoveryService clarifycontrol.DiscoveryService
	{
		statusService = clarifycontrol.NewClarifyStatusService(balancer, logger)
		discoveryService = clarifycontrol.NewDiscoveryService(client, logger)
	}

	// Endpoint domain.
	var statusEp endpoint.Endpoint
	var discoveryEp endpoint.Endpoint
	{
		statusEp = clarifycontrol.MakeGetHostStatusEndpoint(statusService)
		discoveryEp = clarifycontrol.MakeServiceDiscoveryEndpoint(discoveryService)
	}
	endpoints := clarifycontrol.Endpoints{StatusEndpoint: statusEp, ServiceDiscoveryEndpoint: discoveryEp}

	startHealth(errc, logger, svc)
	startGRPC(errc, logger, svc, endpoints)

	for {
		select {
		case err := <-errc:
			if err != nil {
				logger.Log(err)
			}
		case s := <-sigs:
			logger.Log("error", fmt.Sprintf("Captured %v. Exiting...", s))
			client.Deregister(reg)
			os.Exit(0)
		}
	}
}

type service struct {
	GRPCAddress *string
	HTTPAddress *string
	Instance    int
	Name        *string
}

func createConsulClient(consulAddr *string, logger log.Logger) (consulsd.Client, error) {
	consulConfig := api.DefaultConfig()
	if len(*consulAddr) > 0 {
		consulConfig.Address = *consulAddr
	}
	consulClient, err := api.NewClient(consulConfig)
	return consulsd.NewClient(consulClient), err
}

func registerService(client consulsd.Client, svc *service) (*api.AgentServiceRegistration, error) {
	check := &api.AgentServiceCheck{
		HTTP:     fmt.Sprintf("http://%v/health", *svc.HTTPAddress),
		Interval: "10s",
		Timeout:  "3s",
	}
	host, strPort, _ := net.SplitHostPort(*svc.GRPCAddress)
	port, _ := strconv.Atoi(strPort)
	reg := &api.AgentServiceRegistration{
		Name:    *svc.Name,
		Address: host,
		Port:    port,
		ID:      *svc.Name + "-" + strconv.Itoa(svc.Instance),
		Tags:    []string{"grpc"},
		Check:   check,
	}
	err := client.Register(reg)
	return reg, err
}

func startHealth(errc chan error, logger log.Logger, svc *service) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		conn, err := grpc.Dial(*svc.GRPCAddress, grpc.WithInsecure())
		if err != nil {
			errc <- err
		}
		defer conn.Close()
		client := grpc_health_v1.NewHealthClient(conn)
		health, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: *svc.Name})
		if err != nil {
			errc <- err
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(health.String()))
	})
	logger.Log("info", fmt.Sprintf("Starting /health at %v", *svc.HTTPAddress))
	go func() {
		errc <- http.ListenAndServe(*svc.HTTPAddress, nil)
	}()
}

func startGRPC(errc chan error, logger log.Logger, svc *service, endpoints clarifystatussvc.Endpoints) {
	go func() {
		logger.Log("info", fmt.Sprintf("Starting gRPC at %v", *svc.GRPCAddress))
		ln, err := net.Listen("tcp", *svc.GRPCAddress)
		if err != nil {
			errc <- err
			return
		}
		srv := clarifystatussvc.MakeGRPCServer(endpoints, nil, nil)
		healthSrv := health.NewServer()
		s := grpc.NewServer()
		grpc_health_v1.RegisterHealthServer(s, healthSrv)
		pb.RegisterClarifyStatusServer(s, srv)
		healthSrv.SetServingStatus(*svc.Name, grpc_health_v1.HealthCheckResponse_SERVING)
		errc <- s.Serve(ln)
	}()
}

func nomadAddressFactory(instance string) (endpoint.Endpoint, io.Closer, error) {
	return func(context.Context, interface{}) (interface{}, error) {
		return instance, nil
	}, nil, nil
}
