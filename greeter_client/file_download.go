package main

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testgrpc/proto/greeter"
	"time"
)

func DownloadFile(conn *grpc.ClientConn, name string, filePath string) {
	c := greeter.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	request := greeter.DownloadFileRequest{
		ClientName: name,
		FileName:   filepath.Base(filePath),
		FilePath:   filePath,
	}
	r, err := c.DownloadFile(ctx, &request)
	if err != nil {
		log.Fatalf("error downloading file:%v\n", err)
	}

	b := new(bytes.Buffer)
	for {
		size, err := r.Recv()

		if err == io.EOF {
			log.Println("receive finish")
			ioutil.WriteFile("new_"+request.FileName, b.Bytes(), 0755|os.ModeAppend)
			break
		}
		if err != nil {
			log.Fatalf("error receving file:%v\n", err)
			break
		}
		fmt.Println(100 * size.Process / size.Total)
		b.Write(size.FilePart)
	}
}
