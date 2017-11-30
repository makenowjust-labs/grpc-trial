.PHONY: all server client
all: server client
server: bin/server
client: bin/client

bin/server: echo/echo.pb.go
	go build -o ./bin/server ./cmd/server

bin/client: echo/echo.pb.go
	go build -o ./bin/client ./cmd/client

echo/echo.pb.go: echo/echo.proto
	protoc -I echo echo/echo.proto --go_out=plugins=grpc:echo
