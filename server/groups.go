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

// JoinGroupChat adds a user to a group
func (s *Server) JoinGroupChat(ctx context.Context, in *pb.Channel) (*emptypb.Empty, error) {
	s.groupLock.Lock()
	defer s.groupLock.Unlock()

	// Extract the username from context metadata for the current operation
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata is not provided")
	}

	// Assuming the username is passed as a custom header like "auth-header"
	usernameSlice, ok := md["auth-header"]
	if !ok || len(usernameSlice) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "username is not provided in auth-header")
	}
	username := usernameSlice[0]

	// Check if the group exists
	if _, ok := s.groups[in.Name]; !ok {
		return nil, fmt.Errorf("group %s doesn't exist", in.Name)
	}

	// Add the user to the group
	s.groups[in.Name] = append(s.groups[in.Name], username)
	return &emptypb.Empty{}, nil
}

// LeaveGroupChat removes a user from a group
func (s *Server) LeaveGroupChat(ctx context.Context, in *pb.Channel) (*emptypb.Empty, error) {
	s.groupLock.Lock()
	defer s.groupLock.Unlock()

	if members, ok := s.groups[in.Name]; ok {
		for i, member := range members {
			if member == in.Name {
				s.groups[in.Name] = append(members[:i], members[i+1:]...)
				// If the group is empty after removal, delete it
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

// CreateGroupChat creates a new group
func (s *Server) CreateGroupChat(ctx context.Context, in *pb.Channel) (*emptypb.Empty, error) {
	s.groupLock.Lock()
	defer s.groupLock.Unlock()

	if _, ok := s.groups[in.Name]; ok {
		return nil, fmt.Errorf("group %s already exists", in.Name)
	}

	s.groups[in.Name] = []string{in.Name} // The creator joins the group automatically
	return &emptypb.Empty{}, nil
}
