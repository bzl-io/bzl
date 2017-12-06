package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/codes"
	"golang.org/x/net/context"
	_ "golang.org/x/net/trace"

	api "github.com/bzl-io/bzl/api"
)

const (
	grpcPort = 6006
)

func main() {
	flag.Parse()
	grpclog.SetLogger(log.New(os.Stdout, "server: ", log.LstdFlags))

	grpc.EnableTracing = true
	
	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(200)}
	sopts = append(sopts, grpc.UnaryInterceptor(loggingUnaryInterceptor))
	sopts = append(sopts, grpc.StreamInterceptor(loggingStreamInterceptor))
	sopts = append(sopts, grpc.UnknownServiceHandler(unknownHandler))
	grpcServer := grpc.NewServer(sopts...)

	api.RegisterPluginApiServer(grpcServer, &PluginService{})

	grpclog.Printf("Serving bzl plugins at http2://localhost:%d", grpcPort)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}	

	grpcServer.Serve(lis)
}

func unknownHandler(srv interface{}, stream grpc.ServerStream) error {
	grpclog.Printf("Unknown %s", srv)
	return grpc.Errorf(codes.Unauthenticated, "user unauthenticated")
}

func loggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	grpclog.Printf("Unary Request %s", info.FullMethod)
	return handler(ctx, req)
}

func loggingStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	grpclog.Printf("Streaming Request %s", info.FullMethod)
	return handler(srv, ss)
}

func startTcpServer(grpcServer *grpc.Server) {
}

type PluginService struct {
}

func (p *PluginService) GetPlugin(ctx context.Context, req *api.PluginRequest) (*api.PluginResponse, error) {

	grpclog.Printf("Getting plugin %s/%s", req.Organization, req.Name);

	res := &api.PluginResponse{
		Name: req.Name,
		Organization: req.Organization,
		Length: 0,
	}
	
	return res, nil
}
