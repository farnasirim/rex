GO ?= go
REX := github.com/farnasirim/rex

.PHONY: rex rexd proto test coverage
all: rex rexd

rex rexd: proto
	$(GO) build $(REX)/cmd/$@

proto:
	cd proto && protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		rex.proto

clean:
	rm -f rex rexd

test:
	go test -v -race ./...

coverage:
	rm coverage.out
	go test -coverprofile=coverage.out -race ./... && go tool cover -func=coverage.out
