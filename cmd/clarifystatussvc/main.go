package main

import (
	"flag"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/pgombola/clarify-status/pkg"
	"github.com/pgombola/clarify-status/pkg/pb"
)

func main() {
	var (
		grpcAddr = flag.String("grpc.addr", ":8081", "gRPC (HTTP/2) listen address")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Mechanical domain.
	errc := make(chan error)
	// ctx := context.Background()

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

	logger.Log("exit", <-errc)
}
