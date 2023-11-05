package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	pb "grpc_chat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatClient struct {
	serviceClient pb.ChatServiceClient
	username      string
	ctx           context.Context
	conn          *grpc.ClientConn
}

func NewChatClient(address string, username string, authToken string) *ChatClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	handleFatalError("connect to server", err)

	md := metadata.New(map[string]string{"auth-header": authToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	client := pb.NewChatServiceClient(conn)
	return &ChatClient{
		serviceClient: client,
		username:      username,
		ctx:           ctx,
		conn:          conn,
	}
}

func (c *ChatClient) Connect() {
	md := metadata.New(map[string]string{"auth-header": c.username})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := c.serviceClient.Connect(ctx, &pb.User{Username: c.username})
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	log.Printf("Connected to chat as %s", c.username)

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Println("Server closed the connection.")
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a message : %v", err)
			}

			log.Printf("%s: %s\n", msg.Sender, msg.Text)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	log.Println("Disconnecting client...")
}

func (c *ChatClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *ChatClient) JoinGroupChat(channelName string) {
	ctx, cancel := c.authenticatedContext()
	defer cancel()

	_, err := c.serviceClient.JoinGroupChat(ctx, &pb.Channel{Name: channelName, Type: pb.ChannelType_GROUP})
	handleFatalError("join group chat", err)

	fmt.Println("Joined group chat:", channelName)
}

func (c *ChatClient) LeaveGroupChat(channelName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.serviceClient.LeaveGroupChat(ctx, &pb.Channel{Name: channelName, Type: pb.ChannelType_GROUP})
	handleFatalError("leave group chat", err)

	fmt.Println("Left group chat:", channelName)
}

func (c *ChatClient) CreateGroupChat(channelName string) {
	ctx, cancel := c.authenticatedContext()
	defer cancel()

	_, err := c.serviceClient.CreateGroupChat(ctx, &pb.Channel{Name: channelName, Type: pb.ChannelType_GROUP})
	handleFatalError("create group chat", err)

	fmt.Println("Created group chat:", channelName)
}

func (c *ChatClient) SendMessage(target string, targetType pb.ChannelType, text string) {
	ctx, cancel := c.authenticatedContext()
	defer cancel()

	_, err := c.serviceClient.SendMessage(ctx, &pb.Message{
		Sender:     c.username,
		Target:     target,
		TargetType: targetType,
		Text:       text,
	})
	handleFatalError("send message", err)

	fmt.Printf("Sent message to %s: %s\n", target, text)
}

func (c *ChatClient) ListChannels() {
	ctx, cancel := c.authenticatedContext()
	defer cancel()

	// Now passing the authenticated context with metadata
	resp, err := c.serviceClient.ListChannels(ctx, &emptypb.Empty{})
	handleFatalError("list channels", err)

	fmt.Println("Available channels:")
	for _, channel := range resp.Channels {
		fmt.Printf("- %s [%s]\n", channel.Name, channel.Type)
	}
}

func (c *ChatClient) authenticatedContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	md := metadata.New(map[string]string{"auth-header": c.username})
	return metadata.NewOutgoingContext(ctx, md), cancel
}

func handleFatalError(operation string, err error) {
	if err != nil {
		log.Fatalf("could not %s: %v", operation, err)
	}
}
