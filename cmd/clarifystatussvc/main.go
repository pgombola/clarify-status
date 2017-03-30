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

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
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
		serviceName = flag.String("service.name", "clarify-status", "Name of the service")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Registry domain.
	var (
		client          consulsd.Client
		reg             *api.AgentServiceRegistration
		nomadSubscriber sd.Subscriber
	)
	{
		var err error
		client, err = createConsulClient(consulAddr, logger)
		if err != nil {
			logger.Log("err", err)
		}
		reg, err = registerService(client, httpAddr, grpcAddr, *serviceName, instance)
		if err != nil {
			logger.Log("err", err)
		}
		nomadSubscriber = consulsd.NewSubscriber(client, nomadAddressFactory, logger, "nomad", []string{"http"}, true)
	}

	//Discovery doman.

	// Mechanical domain.
	// Error handling channel.
	errc := make(chan error)
	// Signal handling channel. This allows for registry cleanup when we receive SIGTERM, SIGINT.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Business domain.
	var service clarifystatussvc.Service
	{
		service = clarifystatussvc.NewClarifyStatusService(nomadSubscriber)
	}

	// Endpoint domain.
	var statusEp endpoint.Endpoint
	{
		statusEp = clarifystatussvc.MakeGetHostStatusEndpoint(service)
	}
	endpoints := clarifystatussvc.Endpoints{StatusEndpoint: statusEp}

	startHealth(errc, logger, httpAddr)
	startGRPC(errc, logger, grpcAddr, endpoints)

	for {
		select {
		case err := <-errc:
			if err != nil {
				logger.Log(err)
			}
		case s := <-sigs:
			logger.Log("err", fmt.Sprintf("Captured %v. Exiting...", s))
			client.Deregister(reg)
			os.Exit(0)
		}
	}
}

func createConsulClient(consulAddr *string, logger log.Logger) (consulsd.Client, error) {
	consulConfig := api.DefaultConfig()
	if len(*consulAddr) > 0 {
		consulConfig.Address = *consulAddr
	}
	consulClient, err := api.NewClient(consulConfig)
	return consulsd.NewClient(consulClient), err
}

func registerService(client consulsd.Client, httpAddr *string, grpcAddr *string, serviceName string, instance *int) (*api.AgentServiceRegistration, error) {
	check := &api.AgentServiceCheck{
		HTTP:     fmt.Sprintf("http://%v/health", *httpAddr),
		Interval: "10s",
		Timeout:  "1s",
	}
	host, strPort, _ := net.SplitHostPort(*grpcAddr)
	port, _ := strconv.Atoi(strPort)
	reg := &api.AgentServiceRegistration{
		Name:    serviceName,
		Address: host,
		Port:    port,
		ID:      serviceName + "-" + strconv.Itoa(*instance),
		Tags:    []string{"urlprefix-/status"},
		Check:   check,
	}
	err := client.Register(reg)
	return reg, err
}

func startHealth(errc chan error, logger log.Logger, httpAddr *string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	logger.Log("info", fmt.Sprintf("Starting /health at %v", *httpAddr))
	go func() {
		errc <- http.ListenAndServe(*httpAddr, nil)
	}()
}

func startGRPC(errc chan error, logger log.Logger, grpcAddr *string, endpoints clarifystatussvc.Endpoints) {
	go func() {
		logger.Log("info", fmt.Sprintf("Starting gRPC at %v", *grpcAddr))
		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}
		srv := clarifystatussvc.MakeGRPCServer(endpoints, nil, nil)
		s := grpc.NewServer()
		pb.RegisterClarifyStatusServer(s, srv)
		errc <- s.Serve(ln)
	}()
}

func nomadAddressFactory(instance string) (endpoint.Endpoint, io.Closer, error) {
	return func(context.Context, interface{}) (interface{}, error) {
		return instance, nil
	}, nil, nil
}

type exitSignal struct {
	sig string
}

func (e *exitSignal) Error() string {
	return fmt.Sprintf("Received %v", e.sig)
}
