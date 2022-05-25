# 生成SSL SAN 证书

参考自[gRPC 双向验证、自签名证书](https://juejin.cn/post/7025504258371878919)

## CA

### 创建一个CA私钥（根证书）

~~~ shell
openssl genrsa -out ca.key 4096
~~~

### 创建一个conf 用来生成csr（请求签名证书文件）

`` vim ca.conf``填入以下内容：

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

### 生成证书 crt文件

~~~ shell
openssl x509 \
-req \
-days 36500 \
-in ca.csr \
-signkey ca.key \
-out ca.crt
~~~

### 生成pem文件（用于client证书）

~~~ shell
openssl req -new -x509 -key ca.key -out ca.pem -days 36500
~~~

原文没有这一步，生成client.pem时会出错

## 生成server证书

### 创建server私钥 server.key

~~~ shell
openssl genrsa -out server.key 4096
~~~

### 创建一个conf 用来生成csr

`` vim server.conf``填入以下内容：

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

### 生成server.crt和server.pem

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

## 创建client证书

### 生成client.key

~~~ shell
openssl ecparam -genkey -name secp384r1 -out client.key
~~~

### 生成client.csr

~~~ shell
openssl req -new -key client.key -out client.csr -config server.conf
~~~

### 生成client.pem

~~~ shell
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem -extensions req_ext -extfile server.conf
~~~

## 使用例

* 原文使用例。
* 使用证书池。
* client使用client.pem和client.key。

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

## 使用例

* 不使用证书池。
* client使用server.pem。

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



## 可选

### 生成crt

~~~ shell
openssl req -new -x509 -key ca.key -out ca.crt -days 36500
~~~

### 查看crt信息

~~~ shell
openssl x509 -in server.crt -text -noout
~~~