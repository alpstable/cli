# GC is the go compiler.
GC = go

export GO111MODULE=on

default:
	$(GC) build -o gidari-test cmd/gidari/main.go

# proto is a phony target that will generate the protobuf files.
.PHONY: proto
proto:
	protoc \
		--proto_path=api/proto \
		--go-grpc_out=api/service/web \
		--go_out=api/service/web api/proto/web.proto

