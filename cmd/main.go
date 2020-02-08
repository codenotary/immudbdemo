package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/codenotary/immudbrestproxy/pkg/api"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"net/http"
	"os"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", os.Getenv("IMMUD_HOST")+":8080", "gRPC server endpoint")
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

	root, err := rp.GetRoot(ctx)
	log.Print(fmt.Sprintf("Starting...\nRoot index from immudb: %d \nRoot hash from immudb: %s", root.Index, hex.EncodeToString(root.Root)))
	if err != nil {
		return err
	}

	return http.ListenAndServe(":"+os.Getenv("IMMURESTPROXY_PORT"), handler)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
