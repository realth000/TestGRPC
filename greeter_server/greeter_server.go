package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"testgrpc/proto/greeter"
)

var (
	flagPort uint
)

func init() {
	flag.UintVar(&flagPort, "p", 0, "Running port [0-65535]")
	flag.Parse()
	checkFlag()
}

func checkFlag() {
	if flagPort == 0 {
		log.Fatalln("Port not set")
	} else if flagPort > 65535 {
		log.Fatalf("Invalid port: %d\n", flagPort)
	}
}

type server struct {
	greeter.UnimplementedGreeterServer // TODO: WTF IS THIS ???
}

func (s *server) SayHello(ctx context.Context, req *greeter.HelloRequest) (rsp *greeter.HelloReply, err error) {
	rsp = &greeter.HelloReply{Message: "Hello" + req.Name}
	log.Printf("Say Hello to %v\n", ctx)
	return rsp, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", flagPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	// gRPC server.
	s := grpc.NewServer()
	greeter.RegisterGreeterServer(s, &server{})

	// reflection.Register(s)
	fmt.Printf("gRPC serer running on %d\n", flagPort)
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
