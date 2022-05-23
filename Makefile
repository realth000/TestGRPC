SOURCE_SERVER=$(shell find ./greeter_server/ -maxdepth 1 -type f -name "*.go")
SOURCE_CLIENT=$(shell find ./greeter_client/ -maxdepth 1 -type f -name "*.go")
TARGET_SERVER=server
TARGET_CLIENT=client
GO_CMD=go
GO_CFLAGS=-O2 -fPIE -pie -fstack-protector-all -D_FORTIFY_SOURCE=2 -s
GO_LDFLAGS="--extldflags=-Wl,-z,now,-z,relro,-z,noexecstack -s"

export GO111MODULE=on
export CGO_CFLAGS=$(GO_CFLAGS)
export CGO_CXXFLAGS=$(GO_CFLAGS)

.PHONY: all
all: protobuf server client

.PHONY: protobuf
protobuf:
	$(shell protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/greeter/greeter.proto)

.PHONY: server
server:
	$(GO_CMD) build --buildmode=exe --buildmode=pie --trimpath -ldflags $(GO_LDFLAGS) -o $(TARGET_SERVER) $(SOURCE_SERVER)

.PHONY: client
client:
	$(GO_CMD) build --buildmode=exe --buildmode=pie --trimpath -ldflags $(GO_LDFLAGS) -o $(TARGET_CLIENT) $(SOURCE_CLIENT)

.PHONY: clean
clean:
	$(RM) $(TARGET_SERVER) $(TARGET_CLIENT)
