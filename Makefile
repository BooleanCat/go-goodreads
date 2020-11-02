.PHONY: test lint vet generate

test: vet lint
	go test ./... -v

vet:
	go vet ./...

lint:
ifndef SKIP_LINT
	golangci-lint run
endif

generate:
	go generate ./...