package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	"os"
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

func (s *server) DownloadFile(req *greeter.DownloadFileRequest, stream greeter.Greeter_DownloadFileServer) error {
	file, err := os.Open(req.FilePath)
	if err != nil {
		log.Printf("error download file to client: %v\n", err)
		return err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()
	const fileChunk = 1 * (1 << 20) // 1 MB, change this to your requirement
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)
		file.Read(partBuffer)
		resp := &greeter.DownloadFileReply{
			FilePart: partBuffer,
			Process:  int32(i),
			Total:    int32(totalPartsNum),
		}

		err = stream.Send(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}
	return nil
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
