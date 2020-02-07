package main

import (
	"context"
	"flag"
	"github.com/codenotary/immudbrestproxy/pkg/api"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net/http"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "127.0.0.1:8083", "gRPC server endpoint")
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

	err := api.RegisterImmuServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)

	conn, err := grpc.Dial(*grpcServerEndpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", grpcServerEndpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", grpcServerEndpoint, cerr)
			}
		}()
	}()

	client := api.NewImmuServiceClient(conn)
	rp := api.NewRootprovider(client)

	rp.GetRoot(ctx)

	if err != nil {
		return err
	}

	return http.ListenAndServe(":8081", handler)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
