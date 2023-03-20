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

# test runs all of the unit tests locally. Each test is run 5 times to minimize
# flakiness.
.PHONY: tests
tests:
	$(GC) clean -testcache && go test -v -count=5 -failfast ./...

# fmt runs the formatter.
.PHONY: fmt
fmt:
	gofumpt -l -w .

# lint runs the linter.
.PHONY: lint
lint:
	golangci-lint run --fix
	golangci-lint run --config .golangci.yml

# add-license adds the license to all the top of all the .go files.
.PHONY: add-license
add-license:
	chmod +x ./scripts/add-license.sh
	./scripts/add-license.sh
