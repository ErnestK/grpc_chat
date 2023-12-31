package server

import (
	"context"
	"sync"

	pb "grpc_chat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedChatServiceServer
	users          map[string]chan *pb.ChatMessage // Streams for each user
	groups         map[string][]string             // Group names with user lists
	groupLock      sync.RWMutex                    // To control concurrent access to groups
	messageLock    sync.Mutex                      // To control concurrent access to message streams
	messageHistory map[string][]*pb.ChatMessage    // Group names with message lists
	historyLock    sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		users:          make(map[string]chan *pb.ChatMessage),
		groups:         make(map[string][]string),
		messageHistory: make(map[string][]*pb.ChatMessage),
	}
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing credentials")
	}

	// Assuming the username is passed as a custom header like "auth-header"
	usernameSlice, ok := md["auth-header"]
	if !ok || len(usernameSlice) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing auth token")
	}
	username := usernameSlice[0]
	if username == "" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
	}

	return handler(ctx, req)
}
