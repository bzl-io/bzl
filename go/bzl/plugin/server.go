package main

type Server struct {
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	//"net/http/httputil"
	"os"
	"path/filepath"
	"time"
	"sort"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/codes"
	"golang.org/x/net/context"
	_ "golang.org/x/net/trace"

	"github.com/gorilla/mux"


	api "github.com/bzlio/api"
)

const (
	ServerRunfiles = "./server.runfiles/__main__/"
	grpcPort = 6006
	httpPort = 9090
)

var (
	enableTls       = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
	tlsCertFilePath = flag.String("tls_cert_file", "../etc/localhost.crt", "Path to the CRT/PEM file.")
	tlsKeyFilePath  = flag.String("tls_key_file", "../etc/localhost.key", "Path to the private key file.")
)

func listGrpcEndpoints(grpcServer *grpc.Server) {
	endpoints := grpcweb.ListGRPCResources(grpcServer)
	sort.Strings(endpoints)
	for _, endpoint := range endpoints {
		grpclog.Println("--", endpoint)
	}
}

func main() {
	flag.Parse()
	grpclog.SetLogger(log.New(os.Stdout, "server: ", log.LstdFlags))

	grpc.EnableTracing = true
	
	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(200)}
	sopts = append(sopts, grpc.UnaryInterceptor(loggingUnaryInterceptor))
	sopts = append(sopts, grpc.StreamInterceptor(loggingStreamInterceptor))
	sopts = append(sopts, grpc.UnknownServiceHandler(unknownHandler))
	grpcServer := grpc.NewServer(sopts...)

	api.RegisterUserApiServer(grpcServer, &user.Service{})

	grpclog.Printf("Serving http at %d, gRPC at %d", httpPort, grpcPort)

	listGrpcEndpoints(grpcServer)

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

