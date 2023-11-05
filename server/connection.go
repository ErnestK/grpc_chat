package server

import (
	pb "grpc_chat/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Connect creates a message stream for a user
func (s *Server) Connect(in *pb.User, stream pb.ChatService_ConnectServer) error {
	// Add the user to the active user list
	s.messageLock.Lock()
	if _, exists := s.users[in.Username]; exists {
		s.messageLock.Unlock()
		return status.Errorf(codes.AlreadyExists, "user already connected")
	}
	s.users[in.Username] = make(chan *pb.ChatMessage, 100) // Just an example buffer size
	s.messageLock.Unlock()

	// Start a goroutine to send messages to the client
	go s.broadcastToUser(in.Username, stream)

	<-stream.Context().Done() // Block until the stream's context is done (client disconnects)

	// Clean up after disconnection
	s.disconnectUser(in.Username)

	return stream.Context().Err()
}

func (s *Server) disconnectUser(username string) {
	s.messageLock.Lock()
	delete(s.users, username)
	s.messageLock.Unlock()

	s.groupLock.Lock()
	for group, members := range s.groups {
		for i, member := range members {
			if member == username {
				s.groups[group] = append(members[:i], members[i+1:]...)
				break
			}
		}
	}
	s.groupLock.Unlock()
}
