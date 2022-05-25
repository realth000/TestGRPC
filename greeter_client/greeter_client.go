package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
)

const (
	defaultName = "world gRPC!"
)

var (
	flagServerUrl      string
	flagServerPort     uint
	flagName           string
	flagDisableSSL     bool
	flagMutualAuth     bool
	actionSayHello     bool
	actionDownloadFile bool
	argDownloadPath    string
	argSSLCertPath     string
	argSSLKeyPath      string
	argSSLCACertPath   string
)

func init() {
	flag.StringVar(&flagServerUrl, "u", "", "Server url")
	flag.UintVar(&flagServerPort, "p", 0, "Server port [0-65535]")
	flag.StringVar(&flagName, "n", "", "Set client name")
	flag.BoolVar(&flagDisableSSL, "disablessl", false, "disable ssl(use http instead of https)")
	flag.BoolVar(&flagMutualAuth, "mutualAuth", true, "use mutual authentication in SSL handshake")
	flag.BoolVar(&actionSayHello, "sayhello", false, "Say hello to server")
	flag.BoolVar(&actionDownloadFile, "downloadfile", false, "Download file from server")
	flag.StringVar(&argDownloadPath, "downloadpath", "", "Specify download path")
	flag.StringVar(&argSSLCertPath, "cert", "", "SSL credential file[*.pem] path")
	flag.StringVar(&argSSLKeyPath, "key", "", "SSL private key file[*.key] path")
	flag.StringVar(&argSSLCACertPath, "cacert", "", "CA credential file[*.pem] path")
	flag.Parse()
	checkFlag()
}

func checkFlag() {
	if flagServerUrl == "" {
		flagServerUrl = "localhost"
	}
	if flagServerPort == 0 {
		log.Fatalln("Server port not set")
	} else if flagServerPort > 65535 {
		log.Fatalf("Invalid port: %d\n", flagServerPort)
	}
	if actionDownloadFile && argDownloadPath == "" {
		log.Fatalf("Download path not set")
	}
	if !flagDisableSSL {
		if argSSLCertPath == "" {
			log.Fatalf("ssl enabled, but credential file[*.pem] not loaded")
		}
		if flagMutualAuth {
			if argSSLKeyPath == "" {
				log.Fatalf("ssl enabled, but private key file[*.key] not loaded")
			}
			if argSSLCACertPath == "" {
				log.Fatalf("mutual authentication enabled, but CA credential file[*.pem] not loaded")
			}
		}
	}
}

func main() {
	// Setup connection.
	var conn *grpc.ClientConn
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
				//ServerName: "",
				RootCAs: certPool,
				// InsecureSkipVerify should be true to pass self-signed certificate.
				InsecureSkipVerify: true,
				// FIXME: Use some custom VerifyConnection to ensure handshake if using self-signed certificate.
				//VerifyConnection
			})
			c, err := grpc.Dial(fmt.Sprintf("%s:%d", flagServerUrl, flagServerPort), grpc.WithTransportCredentials(cred))
			if err != nil {
				log.Fatalf("can not dail %s:%d :%v", flagServerUrl, flagServerPort, err)
			}
			conn = c
		} else {
			cred, err := credentials.NewClientTLSFromFile(argSSLCertPath, "")
			if err != nil {
				log.Fatalf("error ")
			}
			c, err := grpc.Dial(fmt.Sprintf("%s:%d", flagServerUrl, flagServerPort), grpc.WithTransportCredentials(cred))
			if err != nil {
				log.Fatalf("can not dail %s:%d :%v", flagServerUrl, flagServerPort, err)
			}
			conn = c
		}
	} else {
		//conn, err := grpc.Dial(address, grpc.WithInsecure()) // deprecated
		c, err := grpc.Dial(fmt.Sprintf("%s:%d", flagServerUrl, flagServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("can not dail %s:%d :%v", flagServerUrl, flagServerPort, err)
		}
		conn = c
	}

	defer conn.Close()

	// Contact the server and print out its response.
	var name string
	if flagName == "" {
		name = defaultName
	} else {
		name = flagName
	}

	if actionSayHello {
		r, err := SayHello(conn, name)
		if err != nil {
			log.Fatalf("error greeting: %v\n", err)
		}
		log.Printf("successful greet: %s", r.Message)
	}

	if actionDownloadFile {
		DownloadFile(conn, name, argDownloadPath)
		log.Printf("download finish")
	}
}
