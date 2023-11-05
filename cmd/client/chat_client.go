package main

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

// ChatClient holds the state for a client
type ChatClient struct {
	serviceClient pb.ChatServiceClient
	username      string
	ctx           context.Context  // Context with auth metadata
	conn          *grpc.ClientConn // Store the connection to close it later
}

// newAuthContext creates a context with the AUTH header using the client's username.
func (c *ChatClient) newAuthContext(parentCtx context.Context) (context.Context, context.CancelFunc) {
	md := metadata.New(map[string]string{"auth-header": c.username})
	return metadata.NewOutgoingContext(parentCtx, md), nil
}

// NewChatClient creates a new client for the chat service
func NewChatClient(address string, username string, authToken string) *ChatClient {
	// Create a connection with the insecure option (for simplicity)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	// Create a context with the auth token
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
	// Create a new auth context with the username as the auth token.
	ctx, cancel := c.newAuthContext(context.Background())
	defer cancel() // Ensure the context is canceled when this function exits

	// Attempt to create a stream with the server using the auth context.
	stream, err := c.serviceClient.Connect(ctx, &pb.User{Username: c.username})
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	log.Printf("Connected to chat as %s", c.username)

	// Start a goroutine to handle incoming messages.
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				// Server closed the stream
				log.Println("Server closed the connection.")
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a message : %v", err)
			}

			// Handle incoming message (for example, print it out)
			log.Printf("%s: %s\n", msg.Sender, msg.Text)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the client.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	// Disconnect logic if needed
	log.Println("Disconnecting client...")
	// Note: Depending on your server implementation you might need to send a disconnect message or perform other cleanup actions here.
}

func (c *ChatClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// JoinGroupChat lets a user join a group chat
func (c *ChatClient) JoinGroupChat(channelName string) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add authentication token to context metadata
	md := metadata.New(map[string]string{
		"auth-header": c.username, // using the username as auth token
	})
	// Creating a new context with the above metadata
	authCtx := metadata.NewOutgoingContext(ctx, md)

	// Now passing the authenticated context with metadata
	_, err := c.serviceClient.JoinGroupChat(authCtx, &pb.Channel{Name: channelName, Type: pb.ChannelType_GROUP})
	if err != nil {
		log.Fatalf("could not join group chat: %v", err)
	}

	fmt.Println("Joined group chat:", channelName)
}

// LeaveGroupChat lets a user leave a group chat
func (c *ChatClient) LeaveGroupChat(channelName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.serviceClient.LeaveGroupChat(ctx, &pb.Channel{Name: channelName, Type: pb.ChannelType_GROUP})
	if err != nil {
		log.Fatalf("could not leave group chat: %v", err)
	}

	fmt.Println("Left group chat:", channelName)
}

// CreateGroupChat allows a user to create a group chat
func (c *ChatClient) CreateGroupChat(channelName string) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add authentication token to context metadata
	md := metadata.New(map[string]string{
		"auth-header": c.username, // this is the username as auth token
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := c.serviceClient.CreateGroupChat(ctx, &pb.Channel{Name: channelName, Type: pb.ChannelType_GROUP})
	if err != nil {
		log.Fatalf("could not create group chat: %v", err)
	}

	fmt.Println("Created group chat:", channelName)
}

// SendMessage sends a message to a channel or a user
func (c *ChatClient) SendMessage(target string, targetType pb.ChannelType, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create metadata with the username as the auth token
	md := metadata.New(map[string]string{
		"auth-header": c.username, // using the username as auth token
	})
	// Attach the metadata to the context
	authCtx := metadata.NewOutgoingContext(ctx, md)

	// Use the authenticated context when sending the message
	_, err := c.serviceClient.SendMessage(authCtx, &pb.Message{
		Sender:     c.username,
		Target:     target,
		TargetType: targetType,
		Text:       text,
	})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	fmt.Printf("Sent message to %s: %s\n", target, text)
}

// ListChannels lists all the available channels
func (c *ChatClient) ListChannels() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add authentication token to context metadata
	md := metadata.New(map[string]string{
		"auth-header": c.username, // using the username as auth token
	})
	// Creating a new context with the above metadata
	authCtx := metadata.NewOutgoingContext(ctx, md)

	// Now passing the authenticated context with metadata
	resp, err := c.serviceClient.ListChannels(authCtx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not list channels: %v", err)
	}

	fmt.Println("Available channels:")
	for _, channel := range resp.Channels {
		fmt.Printf("- %s [%s]\n", channel.Name, channel.Type)
	}
}
