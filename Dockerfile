FROM golang:1.25.3-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 10001 \
    appuser

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /bin/elgtm ./cmd/elgtm

FROM scratch AS final

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /bin/elgtm /bin/elgtm
COPY --from=builder /src/.reviewer  /etc/elgtm/defaults

ENV PROMPT_DEFAULTS=/etc/elgtm/defaults

WORKDIR /workspace

USER appuser:appuser

ENTRYPOINT [ "/bin/elgtm" ]