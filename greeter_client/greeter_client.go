package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	defaultName = "world gRPC!"
)

var (
	flagServerUrl  string
	flagServerPort uint
	flagName       string
)

func init() {
	flag.StringVar(&flagServerUrl, "u", "", "Server url")
	flag.UintVar(&flagServerPort, "p", 0, "Server port [0-65535]")
	flag.StringVar(&flagName, "n", "", "Set client name")
	flag.Parse()
	checkFlag()
}

func checkFlag() {
	if flagServerUrl == "" {
		log.Fatalln("Server url not set")
	}
	if flagServerPort == 0 {
		log.Fatalln("Server port not set")
	} else if flagServerPort > 65535 {
		log.Fatalf("Invalid port: %d\n", flagServerPort)
	}
}

func main() {
	// Setup connection.
	//conn, err := grpc.Dial(address, grpc.WithInsecure()) // deprecated
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", flagServerUrl, flagServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failed: %v\n", err)
	}
	defer conn.Close()

	// Contact the server and print out its response.
	var name string
	if flagName == "" {
		name = defaultName
	} else {
		name = flagName
	}

	r, err := SayHello(conn, name, "1/2/")
	if err != nil {
		log.Fatalf("error greeting: %v\n", err)
	}
	log.Printf("successful greet: %s", r.Message)
	DownloadFile(conn, name, "./123")
}
