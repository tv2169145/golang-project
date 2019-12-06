package main

import (
	"context"
	"fmt"
	"github.com/tv2169145/golang-project/grpc/02-server/echo"
	"google.golang.org/grpc"
	"net"
)

type EchoServer struct {}

func (e *EchoServer) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Response: "My Echo :" + req.Message,
	}, nil
}

func main() {
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	srv := &EchoServer{}
	echo.RegisterEchoServerServer(s, srv)
	fmt.Println("Now serving at port 8080")
	err = s.Serve(lst)
	if err != nil {
		panic(err)
	}

}
