package main

import (
	"context"
	"fmt"
	"github.com/tv2169145/golang-project/grpc/03-client/echo"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := echo.NewEchoServerClient(conn)
	response, err := c.Echo(ctx, &echo.EchoRequest{
		Message:"hello jimmy!",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("get from grpc server :", response.Response)

}
