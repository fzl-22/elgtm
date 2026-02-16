install:
	@go mod download

build:
	@CGO_ENABLED=0 go build \
		-ldflags "-s -w" \
		-o ./bin/elgtm \
		./cmd/elgtm

test:
	@go test -v ./...

test-cov:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

test-html: test-cov
	@go tool cover -html=coverage.out

.PHONY: install build test test-cov test-html