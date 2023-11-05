package server

import (
	"context"
	"fmt"

	pb "grpc_chat/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) JoinGroupChat(ctx context.Context, in *pb.Channel) (*emptypb.Empty, error) {
	s.groupLock.Lock()
	defer s.groupLock.Unlock()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata is not provided")
	}

	usernameSlice, ok := md["auth-header"]
	if !ok || len(usernameSlice) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "username is not provided in auth-header")
	}
	username := usernameSlice[0]

	if _, ok := s.groups[in.Name]; !ok {
		return nil, fmt.Errorf("group %s doesn't exist", in.Name)
	}

	s.groups[in.Name] = append(s.groups[in.Name], username)

	go func() {
		userChan := s.users[username]
		s.historyLock.RLock()
		history := s.messageHistory[in.Name]
		s.historyLock.RUnlock()

		for _, msg := range history {
			userChan <- msg
		}
	}()
	return &emptypb.Empty{}, nil
}

func (s *Server) LeaveGroupChat(ctx context.Context, in *pb.Channel) (*emptypb.Empty, error) {
	s.groupLock.Lock()
	defer s.groupLock.Unlock()

	if members, ok := s.groups[in.Name]; ok {
		for i, member := range members {
			if member == in.Name {
				s.groups[in.Name] = append(members[:i], members[i+1:]...)
				if len(s.groups[in.Name]) == 0 {
					delete(s.groups, in.Name)
				}
				break
			}
		}
	} else {
		return nil, fmt.Errorf("group %s doesn't exist", in.Name)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) CreateGroupChat(ctx context.Context, in *pb.Channel) (*emptypb.Empty, error) {
	s.groupLock.Lock()
	defer s.groupLock.Unlock()

	if _, ok := s.groups[in.Name]; ok {
		return nil, fmt.Errorf("group %s already exists", in.Name)
	}

	s.groups[in.Name] = []string{in.Name}
	return &emptypb.Empty{}, nil
}
