package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"testgrpc/greeter_server/check_permission"
	"testgrpc/proto/greeter"
)

var (
	flagPort          uint
	flagDisableSSL    bool
	flagMutualAuth    bool
	argPermitConfFile string
	argSSLCertPath    string
	argSSLKeyPath     string
	argSSLCACertPath  string
)

func init() {
	flag.UintVar(&flagPort, "p", 0, "Running port [0-65535]")
	flag.BoolVar(&flagDisableSSL, "disablessl", false, "disable ssl(use http instead of https)")
	flag.BoolVar(&flagMutualAuth, "mutualAuth", true, "use mutual authentication in SSL handshake")
	flag.StringVar(&argPermitConfFile, "permitfile", "", "--permitfile [filepath] Load permission(permit) config file")
	flag.StringVar(&argSSLCertPath, "cert", "", "SSL credential file[*.pem] path")
	flag.StringVar(&argSSLKeyPath, "key", "", "SSL private key file[*.key] path")
	flag.StringVar(&argSSLCACertPath, "cacert", "", "CA credential file[*.pem] path")
	flag.Parse()
	checkFlag()
}

func checkFlag() {
	if flagPort == 0 {
		log.Fatalln("Port not set")
	} else if flagPort > 65535 {
		log.Fatalf("Invalid port: %d\n", flagPort)
	}
	if !flagDisableSSL {
		if argSSLCertPath == "" {
			log.Fatalf("ssl enabled, but credential file[*.pem] not loaded")
		}
		if argSSLKeyPath == "" {
			log.Fatalf("ssl enabled, but private key file[*.key] not loaded")
		}
		if flagMutualAuth && argSSLCACertPath == "" {
			log.Fatalf("mutual authentication enabled, but CA credential file[*.pem] not loaded")
		}
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
	if !check_permission.CheckPathPermission(req.FilePath) {
		e := status.Error(codes.PermissionDenied, "denied to access this file")
		log.Printf("client=%s, %s", req.ClientName, e)
		return e
	}

	file, err := os.Open(req.FilePath)
	if err != nil {
		e := status.Error(codes.NotFound, "can not open this file"+err.Error())
		log.Printf("client=%s, %s", req.ClientName, e)
		return e
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

	err = check_permission.LoadPermission(argPermitConfFile)
	if err != nil {
		log.Fatalf("error loading permission:%v\n", err)
	}
	// gRPC server.
	var s *grpc.Server
	if !flagDisableSSL {
		if flagMutualAuth {
			// Mutual authentication.
			cert, err := tls.LoadX509KeyPair(argSSLCertPath, argSSLKeyPath)
			if err != nil {
				log.Fatalf("can not load SSL credential:%v", err)
			}
			certPool := x509.NewCertPool()
			credBytes, err := ioutil.ReadFile(argSSLCACertPath)
			if err != nil {
				log.Fatalf("can not load CA credential:%v", err)
			}
			certPool.AppendCertsFromPEM(credBytes)
			cred := credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    certPool,
			})
			s = grpc.NewServer(grpc.Creds(cred))
		} else {
			cred, err := credentials.NewServerTLSFromFile(argSSLCertPath, argSSLKeyPath)
			if err != nil {
				log.Fatalf("can not load SSL credential:%v", err)
			}
			s = grpc.NewServer(grpc.Creds(cred))
		}

	} else {
		s = grpc.NewServer()
	}
	greeter.RegisterGreeterServer(s, &server{})

	// reflection.Register(s)
	fmt.Printf("gRPC serer running on %d\n", flagPort)
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
