package main

import (
	"log"
	"net"

	pb "grpc_chat/proto"
	"grpc_chat/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listening on port 50051")

	// Create a new gRPC server with the interceptor attached
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(server.AuthInterceptor), // Register the interceptor
	}
	s := grpc.NewServer(opts...)

	// Instantiate your server using NewServer
	chatServer := server.NewServer()

	// Register your service with the server instance you created
	pb.RegisterChatServiceServer(s, chatServer)
	log.Println("Chat service registered")

	// Start serving
	log.Println("Server is ready to accept connections")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
