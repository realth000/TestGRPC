package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"testgrpc/greeter_client/config"
)

const (
	defaultName = "world gRPC!"
)

// Command line flags.
var (
	flagConfigFile = kingpin.Flag("config-file", "Config file [*.toml].").String()

	flagServerUrl  = kingpin.Flag("server-ip", "Server IP.").Short('i').String()
	flagServerPort = kingpin.Flag("server-port", "Server port [0-65535].").Short('p').Uint()
	flagName       = kingpin.Flag("client-name", "Client name.").Short('n').String()

	flagSSL        = kingpin.Flag("ssl", "Use SSL in connecting with server. Use --no-ssl to disable ssl.").Default("true").Bool()
	flagMutualAuth = kingpin.Flag("mutual-auth", "Mutual authentication in SSL handshake.").Default("true").Bool()
	flagSSLCert    = kingpin.Flag("cert", "SSL credential file[*.pem] path.").String()
	flagSSLKey     = kingpin.Flag("key", "SSL private key file[*.key] path.").String()
	flagSSLCACert  = kingpin.Flag("ca-cert", "SSL CA credential file[*.pem] path.").String()

	cmdSayHello     = kingpin.Command("say-hello", "Say hello to server, used for testing.")
	cmdDownloadFile = kingpin.Command("download-file", "Download file from server.")
	flagFileName    = cmdDownloadFile.Flag("file-name", "Name of file to download.").String()
)

func loadConfigFile() {
	if *flagConfigFile == "" {
		return
	}
	var cc = config.ClientConfig{}
	err := config.LoadConfigFile(*flagConfigFile, &cc)
	if err != nil {
		log.Fatalf("can not load config file:%v", err)
	}
	fmt.Println(cc)
	//flag.Visit(func(s *flag.Flag) {
	//	switch s.Name {
	//	case "u":
	//		flagServerUrl =
	//	}
	//})
}

func checkFlag() {
	//if flagServerUrl == "" {
	//	flagServerUrl = "localhost"
	//}
	//if flagServerPort == 0 {
	//	log.Fatalln("Server port not set")
	//} else if flagServerPort > 65535 {
	//	log.Fatalf("Invalid port: %d\n", flagServerPort)
	//}
	//if actionDownloadFile && argDownloadPath == "" {
	//	log.Fatalf("Download path not set")
	//}
	//if !flagDisableSSL {
	//	if argSSLCertPath == "" {
	//		log.Fatalf("ssl enabled, but credential file[*.pem] not loaded")
	//	}
	//	if flagMutualAuth {
	//		if argSSLKeyPath == "" {
	//			log.Fatalf("ssl enabled, but private key file[*.key] not loaded")
	//		}
	//		if argSSLCACertPath == "" {
	//			log.Fatalf("mutual authentication enabled, but CA credential file[*.pem] not loaded")
	//		}
	//	}
	//}
}

func main() {
	// Setup connection.
	kingpin.Parse()
	loadConfigFile()
	checkFlag()
	var conn *grpc.ClientConn
	if !*flagSSL {
		//conn, err := grpc.Dial(address, grpc.WithInsecure()) // deprecated
		c, err := grpc.Dial(fmt.Sprintf("%s:%d", *flagServerUrl, *flagServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("can not dail %s:%d :%v", *flagServerUrl, *flagServerPort, err)
		}
		conn = c
	} else {
		if *flagMutualAuth {
			// Mutual authentication.
			cert, err := tls.LoadX509KeyPair(*flagSSLCert, *flagSSLKey)
			if err != nil {
				log.Fatalf("can not load SSL credential:%v", err)
			}

			certPool := x509.NewCertPool()
			credBytes, err := ioutil.ReadFile(*flagSSLCACert)
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
				VerifyConnection: func(cs tls.ConnectionState) error {
					// TODO: Add verify test.
					opts := x509.VerifyOptions{
						DNSName:       cs.ServerName,
						Intermediates: x509.NewCertPool(),
					}
					for _, cert := range cs.PeerCertificates[1:] {
						opts.Intermediates.AddCert(cert)
					}
					_, err := cs.PeerCertificates[0].Verify(opts)
					return err
				},
			})
			c, err := grpc.Dial(fmt.Sprintf("%s:%d", *flagServerUrl, *flagServerPort), grpc.WithTransportCredentials(cred))
			if err != nil {
				log.Fatalf("can not dail %s:%d :%v", *flagServerUrl, *flagServerPort, err)
			}
			conn = c
		} else {
			cred, err := credentials.NewClientTLSFromFile(*flagSSLCert, "")
			if err != nil {
				log.Fatalf("error ")
			}
			c, err := grpc.Dial(fmt.Sprintf("%s:%d", *flagServerUrl, *flagServerPort), grpc.WithTransportCredentials(cred))
			if err != nil {
				log.Fatalf("can not dail %s:%d :%v", *flagServerUrl, *flagServerPort, err)
			}
			conn = c
		}
	}

	defer conn.Close()

	// Contact the server and print out its response.
	var name string
	if *flagName == "" {
		name = defaultName
	} else {
		name = *flagName
	}

	switch kingpin.Parse() {
	case "say-hello":
		r, err := SayHello(conn, name)
		if err != nil {
			log.Fatalf("error greeting: %v\n", err)
		}
		log.Printf("successful greet: %s", r.Message)
	case "download-file":
		DownloadFile(conn, name, *flagFileName)
		log.Printf("download finish")
	}
}
