package main

import (
	"flag"
	"fmt"
	"github.com/MakeNowJust-Labo/grpc-trial/echo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type echoServer struct{}

func (es *echoServer) Connect(c echo.Echo_ConnectServer) error {
	for {
		msg, err := c.Recv()
		if err != nil {
			return err
		}

		log.Print(msg)
		c.Send(msg)
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 4567, "port to listen gRPC server")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	echo.RegisterEchoServer(grpcServer, &echoServer{})

	grpcServer.Serve(lis)
}
