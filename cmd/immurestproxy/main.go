package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	gw "github.com/codenotary/immudbrestproxy/pkg/api"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint",  "immud:8080", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()

	handler := cors.Default().Handler(mux)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterImmuServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", handler)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}