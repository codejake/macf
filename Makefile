.DEFAULT_GOAL := build

.PHONY: fmt test vet build clean

fmt:
	go fmt ./...

test: fmt
	go test ./...

vet: test
	go vet ./...

build: vet
	go build -o macf

clean:
	go clean
