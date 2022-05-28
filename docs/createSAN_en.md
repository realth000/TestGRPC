# Create SSL SAN certificates

Refer from [gRPC 双向验证、自签名证书](https://juejin.cn/post/7025504258371878919)

## CA

### Create CA private key (root certificate)

~~~ shell
openssl genrsa -out ca.key 4096
~~~

### Write a config file to create *.csr

`` vim ca.conf`` fill with the content below：

> [ req ]
> default_bits = 4096
> distinguished_name = req_distinguished_name
>
> [ req_distinguished_name ]
> countryName = Country Name (2 letter code)
> countryName_default = CN
> stateOrProvinceName = State or Province Name (2 letter code)
> stateOrProvinceName_default = Beijing
> localityName = Locality Name (eg, city)
> localityName_default = Beijing
> organizationName = Organization Name (eg, company)
> organizationName_default = github.com/realth000
> commonName = Common Name (e.g. server FQDN or YOUR name)
> commonName_max = 64
> commonName_default = localhost

~~~ shell
openssl req \
-new \
-sha256 \
-out ca.csr \
-key ca.key \
-config ca.conf
~~~

### Create *.crt certificate

~~~ shell
openssl x509 \
-req \
-days 36500 \
-in ca.csr \
-signkey ca.key \
-out ca.crt
~~~

### Create *.pem for client's certificate

~~~ shell
openssl req -new -x509 -key ca.key -out ca.pem -days 36500
~~~

* This step is generating a self-signed certificate, it can be used but without guarantee as it was self-signed, which will cause failure in SSL handshake(Unknown authority).
* Skip this step if not using self-signed certificate, set ``InsecureSkipVerify=false`` and remove ``VerifyConnection`` in client.

* If using self-signed certificate, set ``InsecureSkipVerify=true`` and set ``VerifyConnection`` with a secure and reliable function when client connect server.
  * if ``InsecureSkipVerify=true``, the client will not verify server's identity.
  * if ``VerifyConnection`` set with a function, that function will be used to verify server's identity. The default function in code came from [issue on github](https://github.com/golang/go/issues/40748), and not approved to be actually safe.

## Create server certificate

### Create server private key

~~~ shell
openssl genrsa -out server.key 4096
~~~

### Write a config file to create *.csr

`` vim server.conf`` fill with the content below：

> [ req ]
> default_bits = 4096
> distinguished_name = req_distinguished_name
>
> [ req_distinguished_name ]
> countryName = Country Name (2 letter code)
> countryName_default = CN
> stateOrProvinceName = State or Province Name (2 letter code)
> stateOrProvinceName_default = Beijing
> localityName = Locality Name (eg, city)
> localityName_default = Beijing
> organizationName = Organization Name (eg, company)
> organizationName_default = github.com/realth000
> commonName = Common Name (e.g. server FQDN or YOUR name)
> commonName_max = 64
> commonName_default = localhost
> [ req_ext ]
> subjectAltName = @alt_names
> [alt_names]
> DNS.1 = localhost
> IP  = 127.0.0.1

~~~ shell
openssl req \
-new \
-sha256 \
-out server.csr \
-key server.key \
-config server.conf
~~~

### Create server.crt and server.pem

~~~ shell
openssl x509 \
-req \
-days 365 \
-CA ca.crt \
-CAkey ca.key \
-CAcreateserial \
-in server.csr \
-out server.pem \
-extensions req_ext \
-extfile server.conf
~~~

## Create client certificate

### Create client.key

~~~ shell
openssl ecparam -genkey -name secp384r1 -out client.key
~~~

### Create client.csr

~~~ shell
openssl req -new -key client.key -out client.csr -config server.conf
~~~

### Create client.pem

~~~ shell
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem -extensions req_ext -extfile server.conf
~~~

## Example(1)

* [Original](https://juejin.cn/post/7025504258371878919) example。
* Use certificates pool。
* Use mutual authentication。
* Need to modify ``InsecureSkipVerify`` and ``VerifyConnection`` as above if using self-signed certificate。
* Client uses client.pem and client.key。

### server

~~~ go
// 使用tls 进行加载  key pair
	cert, err := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	if err != nil {
		log.Println("tls 加载x509 证书失败", err)
	}
	// 创建证书池
	certPool := x509.NewCertPool()

	// 向证书池中加入证书
	cafileBytes, err := ioutil.ReadFile("cert/ca.pem")
	if err != nil {
		log.Println("读取ca.pem证书失败", err)
	}
	// 加载客户端证书
	//certPool.AddCert()

	// 加载证书从pem 文件里面
	certPool.AppendCertsFromPEM(cafileBytes)

	// 创建credentials 对象
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},        //服务端证书
		ClientAuth:   tls.RequireAndVerifyClientCert, // 需要并且验证客户端证书
		ClientCAs:    certPool,                       // 客户端证书池

	})

	lis, err := net.Listen("tcp", ":2333")
	if err != nil {
		fmt.Println("监听端口失败 err:", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	services.RegisterHelloWorldServer(s, new(services.HelloService))

	reflection.Register(s)

	err = s.Serve(lis)
	fmt.Println("gRPC starting")
	if err != nil {
		fmt.Println("启动grpc失败 err:", err)
	} else {
		fmt.Println("gRPC started")
	}
~~~

### client

~~~ go
// 和server端一样，先创建证书池
	cert, err := tls.LoadX509KeyPair("cert/client.pem","cert/client.key")
	if err!= nil{
		log.Println("加载client pem, key 失败",err)
	}

	certPool := x509.NewCertPool()
	caFile ,err :=  ioutil.ReadFile("cert/ca.pem")
	if err!= nil{
		log.Println("加载ca失败",err)
	}
	certPool.AppendCertsFromPEM(caFile)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},// 放入客户端证书
		ServerName: "localhost", //证书里面的 commonName
		RootCAs: certPool, // 证书池
	})

	conn, err := grpc.Dial(":2333", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := services.NewHelloWorldClient(conn)
	replay, err := client.HelloRPC(context.Background(), &services.SayHelloRequest{
		Name: "hello world name",
	})
	if err != nil {
		fmt.Println("rpc error:", err)
	} else {
		fmt.Println("result message :", replay.Message)
	}
~~~

## Example(2)

* Not using certificates pool。
* Disable mutual authentication。
* Client uses server.pem。

### Server

~~~ go
listener, err := net.Listen("tcp", 3080)
if err != nil {
    log.Fatalf("failed to listen: %v\n", err)
}

// gRPC server.
var s *grpc.Server
if !flagDisableSSL {
    cred, err := credentials.NewServerTLSFromFile("./server.pem", "./server.key")
    if err != nil {
        log.Fatalf("can not load TLS credentials:%v", err)
    }
    s = grpc.NewServer(grpc.Creds(cred))
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
~~~

### Client

~~~ go
cred, err := credentials.NewClientTLSFromFile(server.pem, "")
if err != nil {
    log.Fatalf("error ")
}
c, err := grpc.Dial("localhost:3080", grpc.WithTransportCredentials(cred))
if err != nil {
    log.Fatalf("can not dail %s:%d :%v", flagServerUrl, flagServerPort, err)
}
conn = c
~~~



## Something optional

### Create crt

~~~ shell
openssl req -new -x509 -key ca.key -out ca.crt -days 36500
~~~

### Check *.crt info

~~~ shell
openssl x509 -in server.crt -text -noout
~~~