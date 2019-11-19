.PHONY: test lint vet test-unit

ginkgo := go run github.com/onsi/ginkgo/ginkgo --race --randomizeAllSpecs -r
lint := go run github.com/golangci/golangci-lint/cmd/golangci-lint

test: vet lint test-unit

vet:
	go vet ./...

lint:
	$(lint) run

test-unit:
	$(ginkgo) --skipPackage acceptance