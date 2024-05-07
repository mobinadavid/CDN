FROM golang:alpine AS builder
LABEL maintainer="Ali Khedmati"
# ENVs
ENV CGO_ENABLED 0 \
    GOPATH /go \
    GOCACHE /go-build \
    GOOS linux \
    GOARCH=amd64 \
    GO111MODULE=on
# Copy
WORKDIR /cdn
COPY . /cdn
# Install Dependencies
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go mod download

# Build the binary
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go build -o cdn main.go

# Start a new stage from scratch
FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk*
WORKDIR /usr/local/bin/cdn
COPY --from=builder /cdn/.env .
COPY --from=builder /cdn/cdn .
CMD ["./cdn", "app" ,"bootstrap"]
