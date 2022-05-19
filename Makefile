SOURCE_SERVER=./greeter_server/greeter_server.go
SOURCE_CLIENT=./greeter_client/greeter_client.go
TARGET_SERVER=server
TARGET_CLIENT=client
GO_CMD=go
GO_CFLAGS=-O2 -fPIE -pie -fstack-protector-all -D_FORTIFY_SOURCE=2 -s
GO_LDFLAGS="--extldflags=-Wl,-z,now,-z,relro,-z,noexecstack -s"

export GO111MODULE=on
export CGO_CFLAGS=$(GO_CFLAGS)
export CGO_CXXFLAGS=$(GO_CFLAGS)

.PHONY: all
all: server client

.PHONY: server
server:
	$(GO_CMD) build --buildmode=exe --buildmode=pie --trimpath -ldflags $(GO_LDFLAGS) -o $(TARGET_SERVER) $(SOURCE_SERVER)

.PHONY: client
client:
	$(GO_CMD) build --buildmode=exe --buildmode=pie --trimpath -ldflags $(GO_LDFLAGS) -o $(TARGET_CLIENT) $(SOURCE_CLIENT)

.PHONY: clean
clean:
	$(RM) $(TARGET_SERVER) $(TARGET_CLIENT)
