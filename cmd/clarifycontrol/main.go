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
		consulAddr  = flag.String("consul.addr", ":8500", "Address of consul agent")
		instance    = flag.Int("instance", 0, "The instance count of the status service")
		serviceName = flag.String("service.name", "clarify-control", "Name of the service")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
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
		kitconsul consulsd.Client
		consul    *api.Client
		reg       *api.AgentServiceRegistration
	)
	{
		var err error
		consul, kitconsul, err = createConsulClient(consulAddr, logger)
		if err != nil {
			logger.Log("err", err)
		}
		reg, err = registerService(kitconsul, svc)
		if err != nil {
			logger.Log("err", err)
		}
	}

	// Discovery domain.
	var nomadlb lb.Balancer
	{
		nomadSubscriber := consulsd.NewSubscriber(kitconsul, nomadAddressFactory, logger, "nomad", []string{"http"}, true)
		nomadlb = lb.NewRoundRobin(nomadSubscriber)
	}

	// Mechanical domain.
	// Error handling channel.
	errc := make(chan error)
	// Signal handling channel. This allows for registry cleanup when we receive SIGTERM, SIGINT.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Business domain.
	control := clarifycontrol.NewClarifyControlService(nomadlb, consul, logger)

	// Endpoint domain.
	endpoints := clarifycontrol.MakeServerEndpoints(control)

	startHealth(errc, logger, svc)
	grpcsrv := startGRPC(errc, logger, svc, endpoints)

	for {
		select {
		case err := <-errc:
			if err != nil {
				logger.Log("err", err)
			}
		case s := <-sigs:
			logger.Log("event", "exiting", "exit_signal", s)
			grpcsrv.GracefulStop()
			kitconsul.Deregister(reg)
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

func createConsulClient(consulAddr *string, logger log.Logger) (*api.Client, consulsd.Client, error) {
	consulConfig := api.DefaultConfig()
	if len(*consulAddr) > 0 {
		consulConfig.Address = *consulAddr
	}
	consulClient, err := api.NewClient(consulConfig)
	return consulClient, consulsd.NewClient(consulClient), err
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
	logger.Log("event", "http_starting", "address", *svc.HTTPAddress)
	go func() {
		errc <- http.ListenAndServe(*svc.HTTPAddress, nil)
	}()
}

func startGRPC(errc chan error, logger log.Logger, svc *service, endpoints clarifycontrol.Endpoints) *grpc.Server {
	s := grpc.NewServer()
	go func() {
		logger.Log("event", "grpc_starting", "address", *svc.GRPCAddress)
		ln, err := net.Listen("tcp", *svc.GRPCAddress)
		if err != nil {
			errc <- err
			return
		}
		srv := clarifycontrol.MakeGRPCServer(endpoints, nil, logger)
		healthSrv := health.NewServer()
		grpc_health_v1.RegisterHealthServer(s, healthSrv)
		pb.RegisterClarifyControlServer(s, srv)
		healthSrv.SetServingStatus(*svc.Name, grpc_health_v1.HealthCheckResponse_SERVING)
		errc <- s.Serve(ln)
	}()
	return s
}

func nomadAddressFactory(instance string) (endpoint.Endpoint, io.Closer, error) {
	return func(context.Context, interface{}) (interface{}, error) {
		return instance, nil
	}, nil, nil
}
