package main

import (
	"context"
	"google.golang.org/grpc"
	"testgrpc/proto/greeter"
	"time"
)

func SayHello(conn *grpc.ClientConn, name string) (*greeter.HelloReply, error) {
	c := greeter.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := c.SayHello(ctx, &greeter.HelloRequest{Name: name})
	return r, err
}
