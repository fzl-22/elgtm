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

test-int:
	@go test -tags=integration -v ./...

test-int-cov:
	@go test -tags=integration -coverprofile=coverage-int.out ./...
	@go tool cover -func=coverage-int.out

test-int-html: test-int-cov
	@go tool cover -html=coverage-int.out

.PHONY: install build test test-cov test-html test-int test-int-cov test-int-html