FROM golang:alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/gateway

# Install dependencies
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc git

# Build Go binary
COPY Makefile go.mod go.sum ./
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod make init && go mod download 
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,mode=0755,target=/go/pkg/mod make proto tidy build

# Deployment container
FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/gateway/gateway /gateway
ENTRYPOINT ["/gateway"]
CMD []
