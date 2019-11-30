.PHONY: test lint vet generate

test: vet lint
	go test ./... -v

vet:
	go vet ./...

lint:
	golangci-lint run

generate:
	go generate ./...