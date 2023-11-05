package server

import (
	"context"
	"fmt"
	pb "grpc_chat/proto"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) broadcastToUser(username string, stream pb.ChatService_ConnectServer) {
	for {
		select {
		case msg := <-s.users[username]:
			if err := stream.Send(msg); err != nil {
				// Handle send error, which could mean the client has disconnected
				log.Printf("Error sending to user %v: %v", username, err)
				return
			}
		case <-stream.Context().Done():
			// Handle disconnection
			return
		}
	}
}

func (s *Server) SendMessage(ctx context.Context, in *pb.Message) (*emptypb.Empty, error) {
	s.messageLock.Lock()
	defer s.messageLock.Unlock()

	// If it's a group message, deliver it to all users in the group
	if in.TargetType == pb.ChannelType_GROUP {
		s.groupLock.RLock()
		users, ok := s.groups[in.Target]
		s.groupLock.RUnlock()

		if !ok {
			return nil, fmt.Errorf("group %s doesn't exist", in.Target)
		}

		// Record the message in the group's history before broadcasting
		s.historyLock.Lock() // Ensure thread safety when accessing the history
		s.messageHistory[in.Target] = append(s.messageHistory[in.Target], &pb.ChatMessage{Sender: in.Sender, Text: in.Text})
		s.historyLock.Unlock()

		// Broadcast the message to all users in the group
		for _, username := range users {
			if userChan, ok := s.users[username]; ok {
				userChan <- &pb.ChatMessage{Sender: in.Sender, Text: in.Text}
			}
		}
	} else { // It's a direct message to a user
		if userChan, ok := s.users[in.Target]; ok {
			userChan <- &pb.ChatMessage{Sender: in.Sender, Text: in.Text}
		} else {
			return nil, fmt.Errorf("user %s doesn't exist", in.Target)
		}
	}

	return &emptypb.Empty{}, nil
}
