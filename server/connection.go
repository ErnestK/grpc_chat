package server

import (
	pb "grpc_chat/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Connect(in *pb.User, stream pb.ChatService_ConnectServer) error {
	s.messageLock.Lock()
	if _, exists := s.users[in.Username]; exists {
		s.messageLock.Unlock()
		return status.Errorf(codes.AlreadyExists, "user already connected")
	}
	s.users[in.Username] = make(chan *pb.ChatMessage, 100)
	s.messageLock.Unlock()

	go s.broadcastToUser(in.Username, stream)

	<-stream.Context().Done()

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
