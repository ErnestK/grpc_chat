package server

import (
	"context"

	pb "grpc_chat/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ListChannels(ctx context.Context, in *emptypb.Empty) (*pb.ListChannelsResponse, error) {
	s.groupLock.RLock()
	defer s.groupLock.RUnlock()

	response := &pb.ListChannelsResponse{}

	for user := range s.users {
		response.Channels = append(response.Channels, &pb.Channel{Name: user, Type: pb.ChannelType_USER})
	}
	for group := range s.groups {
		response.Channels = append(response.Channels, &pb.Channel{Name: group, Type: pb.ChannelType_GROUP})
	}
	return response, nil
}
