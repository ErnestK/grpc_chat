package main

import (
	"flag"
	"log"

	pb "grpc_chat/proto"
)

var (
	address     = flag.String("address", "localhost:50051", "The server address in the format of host:port")
	username    = flag.String("user", "", "Username to connect to the chat service")
	authToken   = flag.String("auth", "", "Authentication token")
	op          = flag.String("op", "", "Operation to perform: connect, join, leave, send, create, list")
	target      = flag.String("target", "", "Target user or channel")
	msg         = flag.String("msg", "", "Message to send")
	channelType = flag.String("type", "GROUP", "Channel type: GROUP or DIRECT")
)

func main() {
	flag.Parse()

	if *username == "" {
		log.Fatalf("username must be provided")
	}

	client := NewChatClient(*address, *username, *authToken)
	defer client.Close()

	switch *op {
	case "connect":
		client.Connect()
	case "join":
		if *target == "" {
			log.Fatalf("target must be provided to join a group chat")
		}
		client.JoinGroupChat(*target)
	case "leave":
		if *target == "" {
			log.Fatalf("target must be provided to leave a group chat")
		}
		client.LeaveGroupChat(*target)
	case "send":
		if *target == "" || *msg == "" {
			log.Fatalf("target and msg must be provided to send a message")
		}
		ct := pb.ChannelType(pb.ChannelType_value[*channelType])
		client.SendMessage(*target, ct, *msg)
	case "create":
		if *target == "" {
			log.Fatalf("target must be provided to create a group chat")
		}
		client.CreateGroupChat(*target)
	case "list":
		client.ListChannels()
	default:
		log.Fatalf("unknown operation: %s", *op)
	}
}
