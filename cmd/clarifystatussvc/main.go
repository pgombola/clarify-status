package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/pgombola/clarify-status/pkg"
	"github.com/pgombola/clarify-status/pkg/pb"
)

func main() {
	var (
		grpcAddr   = flag.String("grpc.addr", ":8081", "gRPC (HTTP/2) listen address")
		consulAddr = flag.String("consul.addr", "", "Address of consul agent")
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
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		if len(*consulAddr) > 0 {
			consulConfig.Address = *consulAddr
		}
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}
	reg := &api.AgentServiceRegistration{
		Name:    "clarify-status",
		Address: *grpcAddr,
		ID:      "status1",
	}
	err := client.Register(reg)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// Mechanical domain.
	// Error handling channel.
	errc := make(chan error)
	// Signal handling channel. This allows for registry cleanup when we receive SIGTERM, SIGINT.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		err := client.Deregister(reg)
		if err != nil {
			errc <- err
			return
		}
		errc <- &exitSignal{sig: sig.String()}
	}()

	// Business domain.
	var service clarifystatussvc.Service
	{
		service = clarifystatussvc.NewClarifyStatusService()
	}

	// Endpoint domain.
	var statusEp endpoint.Endpoint
	{
		statusEp = clarifystatussvc.MakeGetHostStatusEndpoint(service)
	}
	endpoints := clarifystatussvc.Endpoints{StatusEndpoint: statusEp}

	// Start gRPC.
	go func() {
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

	// Wait for exit.
	logger.Log("exit", <-errc)
}

type exitSignal struct {
	sig string
}

func (e *exitSignal) Error() string {
	return fmt.Sprintf("Received %v", e.sig)
}
