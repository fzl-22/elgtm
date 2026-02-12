install:
	@go mod download

build:
	@CGO_ENABLED=0 go build \
		-ldflags "-s -w" \
		-o ./bin/elgtm \
		./cmd/elgtm

.PHONY: install build run