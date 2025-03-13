.PHONY: test fmt vet

vet:
	@go vet ./...

fmt: vet
	@go fmt ./...

test: vet
	@go test ./...

build: test
	@go build -o bin/$(BINARY_NAME) -v

install: build
	@go install -v
