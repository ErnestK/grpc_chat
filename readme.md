connect in one tab
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=connect

all others in different tab
create group
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=create --target=CoolChat --type=GROUP

join to chat
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=join --target=CoolChat

send message to chat
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=send --target=CoolChat --type=GROUP --msg="Hello everyone in CoolChat!"

list available channels
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=list

leave channel
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=leave --target=CoolChat

1) Use case - get history
1.1) Run server
go run cmd/server/main.go

1.2) In tab 2 conect as Alice
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=connect

1.3) In tab 3 conect as Bob
go run cmd/client/main.go --address=localhost:50051 --user=Bob --op=connect

1.4) In tab 4
Create group as Alice
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=create --target=CoolChat --type=GROUP

1.5) In tab 4
Sent message as Alice
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=send --target=CoolChat --type=GROUP --msg="Hello everyone in Cool"

1.6) In tab 4
join group as Bob
go run cmd/client/main.go --address=localhost:50051 --user=Bob --op=join --target=CoolChat

Result: 
in tab 3 Bob get all history messages

2) Usecase - broadcasting

2.1) repeat steps from 1 usecase until 1.4( include 1.4 )

2.2) In tab 4
join group as Bob
go run cmd/client/main.go --address=localhost:50051 --user=Bob --op=join --target=CoolChat

2.3) In tab 4
Sent message as Alice
go run cmd/client/main.go --address=localhost:50051 --user=Alice --op=send --target=CoolChat --type=GROUP --msg="Hello everyone in Cool"

Result: 
in tab 2 and in tab 3 all participants( Bob and Alice ), got message.
