FROM golang:alpine AS builder
LABEL maintainer="Ali Khedmati"
ENV CGO_ENABLED 0 \
    GOPATH /go \
    GOCACHE /go-build \
    GOOS linux \
    GOARCH=amd64 \
    GO111MODULE=on 
WORKDIR /cdn-app
COPY . /cdn-app
# Install Dependencies
ENV GOPROXY=http://nexus-repo.omaxplatform.com/repository/go-proxy,direct
RUN go env GOPROXY
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go mod download

# Build the binary
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go build -o cdn-app main.go

# Start a new stage from scratch
FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk*
WORKDIR /usr/local/bin/cdn-app
COPY --from=builder /cdn-app/.env .
COPY --from=builder /cdn-app/cdn-app .
CMD ["./cdn-app", "app" ,"bootstrap"]
