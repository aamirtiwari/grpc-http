// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	pb "test/helloworld"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
)

func run() {
	//run grpc server
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	lis, err := net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.Serve(lis)
}

func main() {
	flag.Parse()
	go run()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx := context.Background()
	err := pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	http.ListenAndServe(":8081", mux)

}
