FROM golang:1.25.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /bin/elgtm ./cmd/elgtm

FROM alpine:latest AS final

RUN apk add --no-cache ca-certificates tzdata

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 10001 \
    appuser

COPY --from=builder /bin/elgtm /bin/elgtm
COPY --from=builder /src/.reviewer  /etc/elgtm/defaults

WORKDIR /workspace

RUN chown appuser:appuser /workspace

ENV PROMPT_DEFAULTS=/etc/elgtm/defaults

USER appuser

ENTRYPOINT [ "/bin/elgtm" ]