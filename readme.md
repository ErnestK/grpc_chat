connect in one tab
go run cmd/client/main.go cmd/client/chat_client.go --address=localhost:50051 --user=Alice --op=connect

all others in different tab
create group
go run cmd/client/main.go cmd/client/chat_client.go --address=localhost:50051 --user=Alice --op=create --target=CoolChat --type=GROUP

join to chat
go run cmd/client/main.go cmd/client/chat_client.go --address=localhost:50051 --user=Alice --op=join --target=CoolChat

send message to chat
go run cmd/client/main.go cmd/client/chat_client.go --address=localhost:50051 --user=Alice --op=send --target=CoolChat --type=GROUP --msg="Hello everyone in CoolChat!"

list available channels
go run cmd/client/main.go cmd/client/chat_client.go --address=localhost:50051 --user=Alice --op=list

leave channel
go run cmd/client/main.go cmd/client/chat_client.go --address=localhost:50051 --user=Alice --op=leave --target=CoolChat
